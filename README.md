# ihelper

ihelper is a tool containing small tools useful in penetration testing iOS applications.

# Installation

You first need to have [frida-go](https://github.com/frida/frida-go) setup correctly, meaning you need to have frida devkit downloaded.

```bash
$ go install github.com/lateralusd/ihelper@latest
```

# Usage

```bash
$ ihelper --help
iOS penetration testing helpers

Usage:
  ihelper [command]

Available Commands:
  dl          Download file/binary from the application
  help        Help about any command
  patch       Patch application or binary with FridaGadget

Flags:
  -h, --help   help for ihelper

Use "ihelper [command] --help" for more information about a command.
```