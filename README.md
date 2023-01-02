# Golang web analytics for Snowplow

[![actively-maintained]][tracker-classificiation] [![Build Status][travis-image]][travis] [![Coveralls][coveralls-image]][coveralls] [![Goreport][goreport-image]][goreport] [![Release][release-image]][releases] [![GoDoc][godoc-image]][godoc] [![License][license-image]][license]

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

| Technical Docs                 | Setup Guide              | Roadmap                | Contributing                     |
|--------------------------------|--------------------------|------------------------|----------------------------------|
| ![i1][techdocs-image]          | ![i2][setup-image]       | ![i3][roadmap-image]   | ![i4][contributing-image]        |
| **[Technical Docs][techdocs]** | **[Setup Guide][setup]** | **[Roadmap][roadmap]** | **[Contributing][contributing]** |

## Copyright and license

The Snowplow Golang Tracker is copyright 2016-2023 Snowplow Analytics Ltd.

Licensed under the **[Apache License, Version 2.0][license]** (the "License");
you may not use this software except in compliance with the License.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[travis-image]: https://travis-ci.org/snowplow/snowplow-golang-tracker.svg?branch=master
[travis]: https://travis-ci.org/snowplow/snowplow-golang-tracker

[release-image]: http://img.shields.io/badge/release-3.0.0-6ad7e5.svg?style=flat
[releases]: https://github.com/snowplow/snowplow-golang-tracker/releases

[license-image]: http://img.shields.io/badge/license-Apache--2-blue.svg?style=flat
[license]: http://www.apache.org/licenses/LICENSE-2.0

[coveralls-image]: https://coveralls.io/repos/github/snowplow/snowplow-golang-tracker/badge.svg?branch=master
[coveralls]: https://coveralls.io/github/snowplow/snowplow-golang-tracker?branch=master

[godoc-image]: https://godoc.org/gopkg.in/snowplow/snowplow-golang-tracker.v2/tracker?status.svg
[godoc]: https://godoc.org/gopkg.in/snowplow/snowplow-golang-tracker.v2/tracker

[goreport-image]: https://goreportcard.com/badge/github.com/snowplow/snowplow-golang-tracker
[goreport]: https://goreportcard.com/report/github.com/snowplow/snowplow-golang-tracker

[techdocs-image]: https://d3i6fms1cm1j0i.cloudfront.net/github/images/techdocs.png
[setup-image]: https://d3i6fms1cm1j0i.cloudfront.net/github/images/setup.png
[roadmap-image]: https://d3i6fms1cm1j0i.cloudfront.net/github/images/roadmap.png
[contributing-image]: https://d3i6fms1cm1j0i.cloudfront.net/github/images/contributing.png

[techdocs]: https://github.com/snowplow/snowplow/wiki/Golang-Tracker
[setup]: https://github.com/snowplow/snowplow/wiki/Golang-Tracker-Setup
[roadmap]: https://github.com/snowplow/snowplow/wiki/Product-roadmap
[contributing]: https://github.com/snowplow/snowplow/wiki/Contributing

[tracker-classificiation]: https://github.com/snowplow/snowplow/wiki/Tracker-Maintenance-Classification
[actively-maintained]: https://img.shields.io/static/v1?style=flat&label=Snowplow&message=Actively%20Maintained&color=6638b8&labelColor=9ba0aa&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAMAAAAoLQ9TAAAAeFBMVEVMaXGXANeYANeXANZbAJmXANeUANSQAM+XANeMAMpaAJhZAJeZANiXANaXANaOAM2WANVnAKWXANZ9ALtmAKVaAJmXANZaAJlXAJZdAJxaAJlZAJdbAJlbAJmQAM+UANKZANhhAJ+EAL+BAL9oAKZnAKVjAKF1ALNBd8J1AAAAKHRSTlMAa1hWXyteBTQJIEwRgUh2JjJon21wcBgNfmc+JlOBQjwezWF2l5dXzkW3/wAAAHpJREFUeNokhQOCA1EAxTL85hi7dXv/E5YPCYBq5DeN4pcqV1XbtW/xTVMIMAZE0cBHEaZhBmIQwCFofeprPUHqjmD/+7peztd62dWQRkvrQayXkn01f/gWp2CrxfjY7rcZ5V7DEMDQgmEozFpZqLUYDsNwOqbnMLwPAJEwCopZxKttAAAAAElFTkSuQmCC
