package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lateralusd/gdylib"
	"github.com/spf13/cobra"
	"github.com/ulikunitz/xz"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	dylibName = "FridaGadget.dylib"
	dylibPath = "@executable_path/FridaGadget.dylib"
)

var patchCmd = &cobra.Command{
	Use:   "patch [ipa file] [CFBundleExecutable]",
	Short: "Patch application or binary with FridaGadget",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("missing ipa file and CFBundleExecutable")
		}
		pth := args[0]
		exe := args[1]
		if !strings.HasSuffix(pth, ".ipa") {
			return errors.New("file does not end with .ipa")
		}
		if err := downloadGadget(); err != nil {
			return err
		}

		tdir, execPath, err := extractIPA(pth, exe)
		if err != nil {
			return err
		}

		if err := copyGadgetToIPA(tdir, execPath); err != nil {
			return err
		}

		if err := addLoad(tdir, execPath); err != nil {
			return err
		}

		if err := createIPA(tdir, pth); err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("Removing %s", dylibName))
		logger.Info("Removing temp directory")
		os.RemoveAll(tdir)
		return os.Remove(dylibName)
	},
}

func init() {
	patchCmd.Flags().BoolP("remove-codesig", "c", true, "remove code signature")
	patchCmd.Flags().BoolP("rpath", "r", true, "add LC_RPATH after adding dylib")
	rootCmd.AddCommand(patchCmd)
}

type data struct {
	Version string  `json:"tag_name"`
	Assets  []asset `json:"assets"`
}

type asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}

func downloadGadget() error {
	u := "https://api.github.com/repos/frida/frida/releases/latest"

	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := c.Get(u)
	if err != nil {
		return err
	}

	var dt data

	if err := json.NewDecoder(resp.Body).Decode(&dt); err != nil {
		resp.Body.Close()
		return err
	}

	logger.Info(fmt.Sprintf("Downloading frida version %s", dt.Version))

	assetName := fmt.Sprintf("frida-gadget-%s-ios-universal.dylib.xz", dt.Version)
	var ur string

	for _, asset := range dt.Assets {
		if asset.Name == assetName {
			ur = asset.DownloadURL
		}
	}

	resp, err = c.Get(ur)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	r, err := xz.NewReader(resp.Body)
	if err != nil {
		return err
	}

	f, err := os.Create(dylibName)
	if err != nil {
		return err
	}
	defer f.Close()

	io.Copy(f, r)

	return nil
}

func extractIPA(path, exe string) (string, string, error) {
	logger.Info(fmt.Sprintf("Extracting %s", path))

	dir := os.TempDir()
	tdir := filepath.Join(dir, "ihelper")
	os.Mkdir(tdir, os.ModePerm)

	r, err := zip.OpenReader(path)
	if err != nil {
		return "", "", err
	}

	var execPath string

	for _, file := range r.File {
		if file.FileInfo().IsDir() {
			fileDir := filepath.Join(tdir, file.Name)
			os.Mkdir(fileDir, os.ModePerm)
			continue
		}

		if filepath.Base(file.Name) == exe {
			execPath = file.Name
		}

		dst, err := os.Create(filepath.Join(tdir, file.Name))
		if err != nil {
			return "", "", err
		}

		f, err := file.Open()
		if err != nil {
			return "", "", err
		}

		io.Copy(dst, f)

		dst.Close()
		f.Close()
	}

	return tdir, execPath, nil
}

func copyGadgetToIPA(tdir, execPath string) error {
	dir := filepath.Dir(execPath)
	logger.Info(fmt.Sprintf("Copying %s to extracted IPA", dylibName))
	pth := filepath.Join(tdir, dir, dylibName)
	f, err := os.Create(pth)
	if err != nil {
		return err
	}
	defer f.Close()

	dylib, err := os.Open(dylibName)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, dylib)
	return err
}

func addLoad(tdir, execPath string) error {
	logger.Info("Adding LC_LOAD_DYLIB")
	buff := new(bytes.Buffer)

	fullPath := filepath.Join(tdir, execPath)
	r, err := gdylib.Run(fullPath, dylibPath,
		gdylib.WithLoadType(gdylib.DYLIB),
		gdylib.WithRemoveCodeSig(true))
	if err != nil {
		return err
	}

	if _, err := io.Copy(buff, r); err != nil {
		return err
	}

	tempfile, err := os.CreateTemp("", "ihelperpatched")
	if err != nil {
		return err
	}
	defer os.Remove(tempfile.Name())

	tempfile.Write(buff.Bytes())

	logger.Info("Adding LC_RPATH")

	withRpath := new(bytes.Buffer)

	r, err = gdylib.Run(tempfile.Name(), dylibPath,
		gdylib.WithLoadType(gdylib.RPATH))
	if err != nil {
		return err
	}

	if _, err := io.Copy(withRpath, r); err != nil {
		return err
	}

	old, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer old.Close()

	_, err = old.Write(withRpath.Bytes())
	return err
}

func createIPA(tdir, originalIPA string) error {
	newIPAName := strings.TrimSuffix(originalIPA, filepath.Ext(originalIPA))
	newIPAName += "_patched.ipa"
	logger.Info(fmt.Sprintf("Creating new %s file", newIPAName))
	file, err := os.Create(newIPAName)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		strippedPath := strings.TrimPrefix(path, tdir)
		strippedPath = strings.TrimLeft(strippedPath, "/")
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		f, err := w.Create(strippedPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	return filepath.Walk(tdir, walker)
}
