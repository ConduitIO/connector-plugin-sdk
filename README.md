# Conduit Connector Plugin SDK

[![License](https://img.shields.io/badge/license-Apache%202-blue)](https://github.com/ConduitIO/connector-plugin-sdk/blob/main/LICENSE.md)
[![Build](https://github.com/ConduitIO/connector-plugin-sdk/actions/workflows/build.yml/badge.svg)](https://github.com/ConduitIO/connector-plugin-sdk/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/conduitio/connector-plugin-sdk)](https://goreportcard.com/report/github.com/conduitio/connector-plugin-sdk)
[![Go Reference](https://pkg.go.dev/badge/github.com/conduitio/connector-plugin-sdk.svg)](https://pkg.go.dev/github.com/conduitio/connector-plugin-sdk)

This repository contains the Go software development kit for implementing a connector plugin for
[Conduit](https://github.com/conduitio/conduit). If you want to implement a connector plugin in another language please
have a look at the [connector plugin protocol](https://github.com/conduitio/connector-plugin).

## Quickstart

Create a new folder and initialize a fresh go module:
```
go mod init example.com/conduit-plugin-demo
```

Add the connector plugin SDK dependency:

```
go get github.com/conduitio/connector-plugin-sdk
```

With this you can start implementing the connector plugin. To implement a source (a plugin that reads from a 3rd party
resource and sends data to Conduit) create a struct that implements
[`sdk.Source`](https://pkg.go.dev/github.com/conduitio/connector-plugin-sdk#Source). To implement a destination (a
plugin that receives data from Conduit and writes it to a 3rd party resource) create a struct that implements
[`sdk.Destination`](https://pkg.go.dev/github.com/conduitio/connector-plugin-sdk#Destination). You can implement both to
make a plugin that can be used both as a source or a destination.

Apart from the source and/or destination you should create constructor functions that return a `sdk.Source`,
`sdk.Destination` and `sdk.Specification` respectively.

The last part is the entrypoint, it needs to call `sdk.Serve` and pass in the constructor functions mentioned before. If
the plugin does not implement a source or destination you should pass in `nil` instead.

```go
package main

import (
	demo "example.com/conduit-plugin-demo"
	sdk "github.com/conduitio/connector-plugin-sdk"
)

func main() {
	sdk.Serve(
		demo.Specification,  // func Specification() sdk.Specification { ... }
		demo.NewSource,      // func NewSource() sdk.Source { ... }
		demo.NewDestination, // func NewDestination() sdk.Destination { ... }
	)
}
```

Now you can build the standalone plugin:

```
go build path/to/main.go
```

You will get a compiled binary which Conduit can use as a plugin. To run your plugin as part of a Conduit pipeline you
can create a connector using the connectors API and specify the path to the compiled plugin binary in the field `plugin`.

Here is an example request to `POST /v1/connectors` (find more about the [Conduit API](https://github.com/conduitio/conduit#api)):

```json
{
  "type": "TYPE_SOURCE",
  "plugin": "/path/to/compiled/plugin/binary",
  "pipelineId": "...",
  "config": {
    "name": "my-plugin",
    "settings": {
      "my-key": "my-value"
    }
  }
}
```

## Examples

For examples of simple plugins you can look at existing plugins like
[conduit-plugin-generator](https://github.com/ConduitIO/conduit-plugin-generator) or
[conduit-plugin-file](https://github.com/ConduitIO/conduit-plugin-file).
