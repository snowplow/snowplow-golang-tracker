# Golang web analytics for Snowplow

[![actively-maintained]][tracker-classification] [![Build Status][gh-actions-image]][gh-actions] [![Coveralls][coveralls-image]][coveralls] [![Goreport][goreport-image]][goreport] [![Release][release-image]][releases] [![GoDoc][godoc-image]][godoc] [![License][license-image]][license]

## Overview

Snowplow event tracker for Golang. Add analytics to your Go apps and servers.

## Developer Quickstart

### Building

Assuming git is installed:

```bash
 host> git clone https://github.com/snowplow/snowplow-golang-tracker
 host> cd snowplow-golang-tracker
 host> make test
 host> make
```

## Find out more

| Snowplow Docs                 | Contributing                        |
|-------------------------------|-------------------------------------|
| ![i1][techdocs-image]         | ![i2][contributing-image]           |
| **[Snowplow Docs][techdocs]** | **[Contributing](Contributing.md)** |

## Copyright and license

The Snowplow Golang Tracker is copyright 2016-2023 Snowplow Analytics Ltd.

Licensed under the **[Apache License, Version 2.0][license]** (the "License");
you may not use this software except in compliance with the License.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[gh-actions]: https://github.com/snowplow/snowplow-golang-tracker/actions
[gh-actions-image]: https://github.com/snowplow/snowplow-golang-tracker/workflows/Build/badge.svg

[release-image]: http://img.shields.io/badge/release-3.0.0-6ad7e5.svg?style=flat
[releases]: https://github.com/snowplow/snowplow-golang-tracker/releases

[license-image]: http://img.shields.io/badge/license-Apache--2-blue.svg?style=flat
[license]: http://www.apache.org/licenses/LICENSE-2.0

[coveralls-image]: https://coveralls.io/repos/github/snowplow/snowplow-golang-tracker/badge.svg?branch=master
[coveralls]: https://coveralls.io/github/snowplow/snowplow-golang-tracker?branch=master

[godoc-image]: https://pkg.go.dev/badge/github.com/snowplow/snowplow-golang-tracker/v3
[godoc]: https://pkg.go.dev/github.com/snowplow/snowplow-golang-tracker/v3

[goreport-image]: https://goreportcard.com/badge/github.com/snowplow/snowplow-golang-tracker
[goreport]: https://goreportcard.com/report/github.com/snowplow/snowplow-golang-tracker

[techdocs-image]: https://d3i6fms1cm1j0i.cloudfront.net/github/images/techdocs.png
[contributing-image]: https://d3i6fms1cm1j0i.cloudfront.net/github/images/contributing.png
[techdocs]: https://docs.snowplow.io/docs/collecting-data/collecting-from-own-applications/golang-tracker/

[tracker-classification]: https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/tracker-maintenance-classification/
[actively-maintained]: https://img.shields.io/static/v1?style=flat&label=Snowplow&message=Actively%20Maintained&color=6638b8&labelColor=9ba0aa&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAMAAAAoLQ9TAAAAeFBMVEVMaXGXANeYANeXANZbAJmXANeUANSQAM+XANeMAMpaAJhZAJeZANiXANaXANaOAM2WANVnAKWXANZ9ALtmAKVaAJmXANZaAJlXAJZdAJxaAJlZAJdbAJlbAJmQAM+UANKZANhhAJ+EAL+BAL9oAKZnAKVjAKF1ALNBd8J1AAAAKHRSTlMAa1hWXyteBTQJIEwRgUh2JjJon21wcBgNfmc+JlOBQjwezWF2l5dXzkW3/wAAAHpJREFUeNokhQOCA1EAxTL85hi7dXv/E5YPCYBq5DeN4pcqV1XbtW/xTVMIMAZE0cBHEaZhBmIQwCFofeprPUHqjmD/+7peztd62dWQRkvrQayXkn01f/gWp2CrxfjY7rcZ5V7DEMDQgmEozFpZqLUYDsNwOqbnMLwPAJEwCopZxKttAAAAAElFTkSuQmCC
