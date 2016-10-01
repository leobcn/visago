[![License][License-Image]][License-URL] [![ReportCard][ReportCard-Image]][ReportCard-URL] 
# visago
Visual AI Aggregator

```
Usage:
  visago <files> [flags]

Flags:
  -h, --help           help for visago
  -l, --list-plugins   list supported plugins
  -v, --verbose        verbose mode
      --version        display version

```

## Install

```
go get -v github.com/zquestz/visago
cd $GOPATH/src/github.com/zquestz/visago
make
make install
```

## Configuration

To setup your own default configuration just create `~/.visago/config`. The configuration file is in UCL format. JSON is also fully supported as UCL can parse JSON files.

For more information about UCL visit:
[https://github.com/vstakhov/libucl](https://github.com/vstakhov/libucl)

The following keys are supported:

* blacklist (array of plugins to exclude)
* verbose (verbose mode)
* whitelist (array of plugins to include)

## Contributors

* [Josh Ellithorpe (zquestz)](https://github.com/zquestz/)

## License

visago is released under the MIT license.

[License-URL]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg
[ReportCard-URL]: http://goreportcard.com/report/zquestz/s
[ReportCard-Image]: https://goreportcard.com/badge/github.com/zquestz/visago
