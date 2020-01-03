# Golang web analytics for Snowplow

[![Build Status][travis-image]][travis] [![Coveralls][coveralls-image]][coveralls] [![Goreport][goreport-image]][goreport] [![Release][release-image]][releases] [![GoDoc][godoc-image]][godoc] [![License][license-image]][license]

## Overview

Snowplow event tracker for Golang. Add analytics to your Go apps and servers.

## Developer Quickstart

### Building

Assuming git, **[Vagrant][vagrant-install]** and **[VirtualBox][virtualbox-install]** installed:

```bash
 host> git clone https://github.com/snowplow/snowplow-golang-tracker
 host> cd snowplow-golang-tracker
 host> vagrant up && vagrant ssh
guest> cd /opt/gopath/src/github.com/snowplow/snowplow-golang-tracker
guest> make
guest> make test
```

## Find out more

| Technical Docs                 | Setup Guide              | Roadmap                | Contributing                     |
|--------------------------------|--------------------------|------------------------|----------------------------------|
| ![i1][techdocs-image]          | ![i2][setup-image]       | ![i3][roadmap-image]   | ![i4][contributing-image]        |
| **[Technical Docs][techdocs]** | **[Setup Guide][setup]** | **[Roadmap][roadmap]** | **[Contributing][contributing]** |

## Copyright and license

The Snowplow Golang Tracker is copyright 2016-2019 Snowplow Analytics Ltd.

Licensed under the **[Apache License, Version 2.0][license]** (the "License");
you may not use this software except in compliance with the License.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[travis-image]: https://travis-ci.org/snowplow/snowplow-golang-tracker.svg?branch=master
[travis]: https://travis-ci.org/snowplow/snowplow-golang-tracker

[release-image]: http://img.shields.io/badge/release-2.2.1-6ad7e5.svg?style=flat
[releases]: https://github.com/snowplow/snowplow-golang-tracker/releases

[license-image]: http://img.shields.io/badge/license-Apache--2-blue.svg?style=flat
[license]: http://www.apache.org/licenses/LICENSE-2.0

[coveralls-image]: https://coveralls.io/repos/github/snowplow/snowplow-golang-tracker/badge.svg?branch=master
[coveralls]: https://coveralls.io/github/snowplow/snowplow-golang-tracker?branch=master

[godoc-image]: https://godoc.org/gopkg.in/snowplow/snowplow-golang-tracker.v2/tracker?status.svg
[godoc]: https://godoc.org/gopkg.in/snowplow/snowplow-golang-tracker.v2/tracker

[goreport-image]: https://goreportcard.com/badge/github.com/snowplow/snowplow-golang-tracker
[goreport]: https://goreportcard.com/report/github.com/snowplow/snowplow-golang-tracker

[vagrant-install]: http://docs.vagrantup.com/v2/installation/index.html
[virtualbox-install]: https://www.virtualbox.org/wiki/Downloads

[techdocs-image]: https://d3i6fms1cm1j0i.cloudfront.net/github/images/techdocs.png
[setup-image]: https://d3i6fms1cm1j0i.cloudfront.net/github/images/setup.png
[roadmap-image]: https://d3i6fms1cm1j0i.cloudfront.net/github/images/roadmap.png
[contributing-image]: https://d3i6fms1cm1j0i.cloudfront.net/github/images/contributing.png

[techdocs]: https://github.com/snowplow/snowplow/wiki/Golang-Tracker
[setup]: https://github.com/snowplow/snowplow/wiki/Golang-Tracker-Setup
[roadmap]: https://github.com/snowplow/snowplow/wiki/Product-roadmap
[contributing]: https://github.com/snowplow/snowplow/wiki/Contributing
