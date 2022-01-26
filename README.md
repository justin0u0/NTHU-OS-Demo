# NTHU-OS-Demo

A CLI tool for OS demo.

Powered by [spf13/cobra](https://github.com/spf13/cobra), [pterm/pterm](https://github.com/pterm/pterm), [martinlindhe/imgcat](https://github.com/martinlindhe/imgcat) and [AlecAivazis/survey](https://github.com/AlecAivazis/survey).

Note that imgcat only support on iterm2.

## Installation

If you have golang version greater than v1.17 installed. Run the following command to install the CLI:

```bash
go install github.com/justin0u0/NTHU-OS-Demo/cmd/osdemo@latest
```

If you get error of `osdemo` command not found, add `export PATH=$(go env GOPATH)/bin:$PATH` to your bash config file.

## Features

All with a single command.

- Generate random questions from a single JSON config file. Example: [Link](question/assets/example.json)
- Generate a form to fill. Record the form into a JSON file.
- Customize the form to fill with a single JSON config file. Example: [Link](record/assets/example.json)
- Export .csv file with the recorded result. Customize the exporter with a single JSON config file.

## Demo

## TODOs

- Add feature to export summary.
