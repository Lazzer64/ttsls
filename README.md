# ttsls - Tabletop Simulator Language Server
`ttsls` is an implementation of the [Language Server Protocol](https://microsoft.github.io/language-server-protocol/) integrating extensions provided by the [Tabletop Simulator External Editor API](https://api.tabletopsimulator.com/externaleditorapi/), and a suite of other tooling, directly with your editor of choice.

## Installation
Precompiled binaries are not avilable at this time, `ttsls` can be compiled from source using the following steps:
```shell
$ git clone https://github.com/Lazzer64/ttsls
$ cd ttsls
$ go build -o ttsls cmd/ttsls.go
```

For integration with your faviorite editor see a list of LSP clients [here](https://microsoft.github.io/language-server-protocol/implementors/tools/).

## Currently Supported Features
- Support for inlining external files with the `#include` directive
- Go-to-definition navigation for `#include`'s
- Code completion suggestions for Tabletop Simulator's custom lua functions/objects
- Hover and signature help support for Tabletop Simulator's custom lua functions/objects
- Direct script execution in Tabletop Simulator via the `tts.exec` command
