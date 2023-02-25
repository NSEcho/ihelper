package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/frida/frida-go/frida"
	"os"
	"path/filepath"
	"time"
)

func download(target, directory, filename string, isFile bool) (string, int, error) {
	d := frida.USBDevice()
	if d == nil {
		return "", 0, errors.New("could not attach to USB device")
	}
	session, err := d.Attach(target, nil)
	if err != nil {
		return "", 0, err
	}

	script, err := session.CreateScript(scriptJS)
	if err != nil {
		return "", 0, err
	}

	var name string
	var length int
	done := make(chan struct{})

	script.On("message", func(message string, data []byte) {
		if len(data) > 0 {
			unmarshalled := make(map[string]string)
			json.Unmarshal([]byte(message), &unmarshalled)

			name = filepath.Base(unmarshalled["payload"])
			length = len(data)
			if err := os.WriteFile(name, data, os.ModePerm); err != nil {
				logger.Fatal("could not save %s: %v", name, err)
			}
			go func() {
				done <- struct{}{}
			}()
		}
	})
	script.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if isFile {
		err := script.ExportsCallWithContext(ctx, "download_file", directory, filename)
		if err != nil {
			return "", 0, err.(error)
		}
	} else {
		err := script.ExportsCallWithContext(ctx, "download_bin")
		if err != nil {
			return "", 0, err.(error)
		}
	}
	<-done
	return name, length, nil
}
