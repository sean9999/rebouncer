# rebouncer

> A debouncer in Go with map-reduce mojo


[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)

Rebouncer melds the concept of a debouncer with the concepts of map-reduce to provide a flexible solution the problem of needing to take many events that occur over a short span of time, and reduce them to fewer events over a long span of time.

It was primarily written for my dev server, fasthak, as a way to deal with the enormous amount of inotify events that some IDEs generate, many of which should be ignored.

