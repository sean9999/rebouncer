# rebouncer
A debouncer in Go with map-reduce mojo

Rebouncer melds the concept of a debouncer with the concepts of map-reduce to provide a flexible solution the problem of needing to take many events that occur over a short span of time, and reduce them to fewer events over a long span of time.

It was primarily written for my dev server, fasthak, as a way to deal with the enormous amount of inotify events that some IDEs generate.  

Rebouncer, though, aims to be a generic solution. []ff
     sdfsdfsdfsdf