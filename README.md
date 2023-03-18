# Rebouncer

<!--
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
-->

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/sean9999/rebouncer/graphs/commit-activity)

<!--
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/sean9999/go/rebouncer)
-->

[![Go Report Card](https://goreportcard.com/badge/github.com/sean9999/rebouncer)](https://goreportcard.com/report/github.com/sean9999/rebouncer)

[![Go version](https://img.shields.io/github/go-mod/go-version/sean9999/rebouncer.svg)](https://github.com/sean9999/rebouncer)

## A debouncer on steroids

![Flower Of Life](flower_of_life.webp)

![Hand](hand.jpg)

![Sacred](sacred_geo.webp)

![Moon Phases](moon_phases.avif)


Rebouncer melds the concept of a debouncer with the concepts of map-reduce to provide a flexible solution the problem of needing to take many events that occur over a short span of time, and reduce them to fewer events over a long span of time.

It was primarily written for my dev server, fasthak, as a way to deal with the enormous amount of inotify events that some IDEs generate, many of which should be ignored.

