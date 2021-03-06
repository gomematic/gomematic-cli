# Gomematic: CLI client

[![Build Status](https://cloud.drone.io/api/badges/gomematic/gomematic-cli/status.svg)](https://cloud.drone.io/gomematic/gomematic-cli)
[![Join the Matrix chat at https://matrix.to/#/#gomematic:matrix.org](https://img.shields.io/badge/matrix-%23gomematic-7bc9a4.svg)](https://matrix.to/#/#gomematic:matrix.org)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6bbad9ecac6c42d6b0a9722d97979398)](https://www.codacy.com/app/gomematic/gomematic-cli?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=gomematic/gomematic-cli&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/gomematic/gomematic-cli?status.svg)](http://godoc.org/github.com/gomematic/gomematic-cli)
[![Go Report](http://goreportcard.com/badge/github.com/gomematic/gomematic-cli)](http://goreportcard.com/report/github.com/gomematic/gomematic-cli)
[![](https://images.microbadger.com/badges/image/gomematic/gomematic-cli.svg)](http://microbadger.com/images/gomematic/gomematic-cli "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

Within this repository we are building the command-line client to interact with the [Gomematic API](https://github.com/gomematic/gomematic-api) server, for further information take a look at our [documentation](https://gomematic.tech).


## Install

You can download prebuilt binaries from the GitHub releases or from our [download site](http://dl.gomematic.tech/cli). You are a Mac user? Just take a look at our [homebrew formula](https://github.com/gomematic/homebrew-gomematic).


## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.11.

```bash
git clone https://github.com/gomematic/gomematic-cli.git
cd gomematic-cli

make sync generate build

./bin/gomematic-cli -h
```


## Security

If you find a security issue please contact gomematic@webhippie.de first.


## Contributing

Fork -> Patch -> Push -> Pull Request


## Authors

* [Thomas Boerger](https://github.com/tboerger)


## License

Apache-2.0


## Copyright

```
Copyright (c) 2018 Thomas Boerger <thomas@webhippie.de>
```
