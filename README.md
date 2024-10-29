# izu
izu is a unified hotkey config that's made to parse its own config and generate it into any other hotkey config available.

This can thus be used to manage multiple hotkey daemons on different hosts.

The primary reason for this is switching display protocols or window managers on Linux (using NixOS managed config files).

It's inspired by [sxhkd](https://github.com/baskerville/sxhkd) and shares part of the config syntax.

## Usage

```
NAME:
   izu - A unified hotkey config based on sxhkd.

USAGE:
   izu [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value     Path to the configuration file
   --formatter value, -f value  Path to the formatter lua file
   --version, -v                Print the version (default: false)
   --verbose, -V                Print verbose output (default: false)
   --silent, -S                 Silent output, does not output any logs or errors unless when panicking (default: false)
   --string value, -s value     String to parse
   --help, -h                   show help
```

Example:
```
izu --config ./configfile --formatter sway
```
## Supported formatters
 - sxhkd (done)
 - hyprland (needs improvement)
 - sway (needs improvement)

## Examples
For configuration examples look in `./example/`

For formatter examples look in `./pkg/izu/formatters/`

## License
MIT
