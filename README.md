[![License][License-Image]][License-URL] [![ReportCard][ReportCard-Image]][ReportCard-URL] [![Build][Build-Status-Image]][Build-Status-URL] [![Chat][Chat-Image]][Chat-URL]
# visago
Visual AI Aggregator

```
Usage:
  visago <files/urls> [flags]

Flags:
  -h, --help              help for visago
  -j, --json              provide JSON output
  -l, --list-plugins      list supported plugins
  -t, --tag-score float   minimum tag score
  -v, --verbose           verbose mode
      --version           display version
```

## Install

```
go get -v github.com/zquestz/visago
cd $GOPATH/src/github.com/zquestz/visago
make
make install
```

## Integration

The `visagoapi` package is available for developers who want to integrate visual AI results in their software.

```
pluginConfig := &visagoapi.PluginConfig{
	URLs: []string{"http://example.com/image.png"},
	Files: []string{"filename"},
}

output, _ := visagoapi.RunPlugins(pluginConfig, true)

fmt.Printf("%#v\n", output)
```

There is also an example integration in `/example/main.go`.

## Plugins

* Clarifai - [https://www.clarifai.com/](https://www.clarifai.com/)
* Google Vision - [https://cloud.google.com/vision/](https://cloud.google.com/vision/)
* Imagga - [https://imagga.com/](https://imagga.com/)

## Configuration

To setup your own default configuration just create `~/.visago/config`. The configuration file is in UCL format. JSON is also fully supported as UCL can parse JSON files.

For more information about UCL visit:
[https://github.com/vstakhov/libucl](https://github.com/vstakhov/libucl)

The following keys are supported:

* blacklist (array of plugins to exclude)
# tag_score (minimum tag score)
* verbose (verbose mode)
* whitelist (array of plugins to include)

## Contributors

* [Josh Ellithorpe (zquestz)](https://github.com/zquestz/)

## License

visago is released under the MIT license.

[License-URL]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg
[ReportCard-URL]: http://goreportcard.com/report/zquestz/visago
[ReportCard-Image]: https://goreportcard.com/badge/github.com/zquestz/visago
[Build-Status-URL]: http://travis-ci.org/zquestz/visago
[Build-Status-Image]: https://travis-ci.org/zquestz/visago.svg?branch=master
[Chat-Image]: https://badges.gitter.im/zquestz/visago.svg
[Chat-URL]: https://gitter.im/zquestz/visago?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge
