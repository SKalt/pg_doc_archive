```
./bin
|- spider_urls # scrapes a list of all postgres doc-site urls
|- archive     # reads each url from stdin and tar.gz's it to stdout
`- README.md
```

See [the makefile](../Makefile) for how to produce these files.
Git ignores the other files in this directory apart from this README.
