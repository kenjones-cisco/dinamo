# Dynamic Generator

Lightweight command-line utility for generating file(s) from using go templates.

## Status

[![GoDoc][1]][2]
[![GoCard][3]][4]
[![Coverage Status][5]][6]
[![Build Status][7]][8]
[![Maintainability][11]][12]
[![MIT License][9]][10]

[1]: https://godoc.org/github.com/kenjones-cisco/dinamo?status.svg
[2]: https://godoc.org/github.com/kenjones-cisco/dinamo
[3]: https://goreportcard.com/badge/kenjones-cisco/dinamo
[4]: https://goreportcard.com/report/github.com/kenjones-cisco/dinamo
[5]: https://coveralls.io/repos/github/kenjones-cisco/dinamo/badge.svg?branch=master
[6]: https://coveralls.io/github/kenjones-cisco/dinamo?branch=master
[7]: https://travis-ci.org/kenjones-cisco/dinamo.svg?branch=master
[8]: https://travis-ci.org/kenjones-cisco/dinamo
[9]: http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square
[10]: https://github.com/kenjones-cisco/dinamo/blob/master/LICENSE
[11]: https://api.codeclimate.com/v1/badges/f26bf4e7607a7940d26b/maintainability
[12]: https://codeclimate.com/github/kenjones-cisco/dinamo/maintainability


## Install

- Download pre-built binaries using [Published Releases](https://github.com/kenjones-cisco/dinamo/releases).

- Alternatively install using `go get`:
```
go get github.com/kenjones-cisco/dinamo/cmd/dinamo
```

## Features

- Generate file(s) using a template
- Supports multiple data sources, key-value arguments, environment variables, YAML files, and JSON files
- Supports [sprig](http://masterminds.github.io/sprig) template functions


## Usage

See detailed [Usage](docs/usage/dinamo.md)

```bash
# create output.txt from config.tmpl using the key-value pairs
dinamo gen -t config.tmpl -f output.txt key1=value1 key2=value2

# create output.txt from config.tmpl using the JSON data file source.json
dinamo gen -t config.tmpl -f output.txt -d source.json

# create output.txt from config.tmpl using the YAML data file source.yaml
dinamo gen -t config.tmpl -f output.txt -d source.yaml

# create output.txt from config.tmpl using environment variables
dinamo gen -t config.tmpl -f output.txt -e

# create output.txt from config.tmpl using the key-value pairs and the YAML data file source.yml
dinamo gen -t config.tmpl -f output.txt -d source.yml key1=value1 key2=value2

# create output.txt from config.tmpl using the key-value pairs and the JSON data file source.json
dinamo gen -t config.tmpl -f output.txt -d source.json key1=value1 key2=value2

# create output.txt from config.tmpl using the key-value pairs and environment variables
dinamo gen -t config.tmpl -f output.txt -e key1=value1 key2=value2

# create output.txt from config.tmpl using the key-value pairs, environment variables, and the YAML data file source.yml
dinamo gen -t config.tmpl -f output.txt -e -d source.yml key1=value1 key2=value2
```

## Contributions

Contributing guidelines are in [CONTRIBUTING.md](CONTRIBUTING.md).
