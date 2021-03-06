[![License][License-Image]][License-URL] [![ReportCard][ReportCard-Image]][ReportCard-URL] [![Build][Build-Status-Image]][Build-Status-URL] [![Release][Release-Image]][Release-URL] [![Chat][Chat-Image]][Chat-URL]
# visago
Visual AI Aggregator.

```
Usage:
  visago <files/urls> [flags]

Flags:
  -c, --colors            display colors
  -f, --faces             display faces
  -j, --json              provide JSON output
  -l, --list-plugins      list supported plugins
  -s, --tag-score float   minimum tag score
  -t, --tags              display tags
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

If you have issues building visago, you can vendor the dependencies by using [gvt](https://github.com/FiloSottile/gvt):

```
go get -u github.com/FiloSottile/gvt
cd $GOPATH/src/github.com/zquestz/visago
gvt restore
```

## Examples

Get metadata about any set of URLs/Files. This requests all currently supported features (Tags/Colors/Faces).
```
visago --json \
  landscape.jpg \
  http://digital-photography-school.com/wp-content/uploads/flickr/2746960560_8711acfc60_o.jpg \
  http://cdn.wegotthiscovered.com/wp-content/uploads/keanureeves.jpg
```

The lengthy response is contained in this [gist](https://gist.github.com/zquestz/08712a847d0b0da1700338f6711d89c8).

To filter returned tags by tag score, use the `-s` flag. This value is between 0 and 1.
```
visago -s .5 landscape.jpg
```

You can also enable specific features if you do not require all supported functionality.
That way you can be much more efficient on API calls.

To only fetch tags pass the `-t` flag.
```
visago -t mountain.png
```

To only fetch facial data pass the `-f` flag.
```
visago -f bio.jpg
```

To only fetch color data pass the `-c` flag.
```
visago -c elmo.jpg
```

## Integration

The `visagoapi` package is available for developers who want to integrate visual AI results in their software.

```
pluginConfig := &visagoapi.PluginConfig{
	URLs: []string{"http://example.com/image.png"},
	Files: []string{"filename"},
	// To only enable select features, set them below.
	// By default all features are enabled.
	// Features: []string{visagoapi.TagsFeature, visagoapi.ColorsFeature, visagoapi.FacesFeature},
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

* blacklist - []string (plugins to exclude)
* colors - bool (display colors)
* faces - bool (display faces)
* json_output - bool (output JSON)
* tag_score - float64 (minimum tag score)
* tags - bool (display tags)
* verbose - bool (verbose mode)
* whitelist - []string (plugins to include)

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
[Release-URL]: https://github.com/zquestz/visago/releases/tag/v0.3.2
[Release-Image]: http://img.shields.io/badge/release-v0.3.2-1eb0fc.svg
[Chat-Image]: https://badges.gitter.im/zquestz/visago.svg
[Chat-URL]: https://gitter.im/zquestz/visago?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge
