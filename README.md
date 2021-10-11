# pg_doc_archive

A scraper to archive the postgresql documentation site.

## Why scrape the docs?

To make it easier to build on the postgres docs. For example, the tarball this scraper spits out could be used to ship the postgres docs as a part of a postgres language server or to automatically compare a specific language feature between postgres versions.

The postgres [docs download page](https://www.postgresql.org/docs/) offers PDF downloads, but doesn't offer collections of web pages.
There are other archives that might offer collections of the HTML postgres docs (for example, [archive.org](https://web.archive.org/web/*/www.postgres.org) and [devdocs.io](https://devdocs.io/postgresql~13/)), but none seem to offer a service as easy as downloading a tarball from GitHub.

### How to consume

Go to the github page's releases tab and download the latest tarball.

You should be able to serve the unpacked tarball as a local copy of the postgres documentation site.
There may be some broken links, but most of the relevant CSS and images should be preserved.
You could then point a scraper such as [`colly`](https://github.com/gocolly/colly/blob/master/_examples/local_files/local_files.go) ([archive](https://github.com/gocolly/colly/blob/19b3ce62c774973a898eea3063c962aea206ec19/_examples/local_files/local_files.go)) at either the filesystem or your local webserver.

## Licensing

The scraper itself [MIT licensed](./licenses/mit.license.md).
The postgres docs are licensed under [the Postgres license](./licenses/postgres.license.md).
Per the postgres license, a copy of the postgres license is included in all copies of the scraped documentation.
