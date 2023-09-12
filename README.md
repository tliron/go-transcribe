Transcribe for Go
=================

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/tliron/go-transcribe.svg)](https://pkg.go.dev/github.com/tliron/go-transcribe)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/go-transcribe)](https://goreportcard.com/report/github.com/tliron/go-transcribe)

Go programs often need to output data in a structured representation format, such as
JSON or YAML. But why not provide wider compatibility and support all common formats,
letting the user choose? This library provides a unified API over conversion to several
formats.

It also supports "pretty" printing to terminals with semantic colorization (including
"it just works" support for colorizing Windows terminals, which defy common standards).

Supported formats:

* [YAML](https://yaml.org/)
* [JSON](https://www.json.org/), including a convention for extending JSON to support
  additional type differentiation
* [XML](https://www.w3.org/XML/) via a conventional schema
* [CBOR](https://cbor.io/)
* [MessagePack](https://msgpack.org/)

The binary formats (CBOR, MessagePack) can be output as is (incompatible with terminals)
or textualized into [Base64](https://datatracker.ietf.org/doc/html/rfc4648#section-4).
