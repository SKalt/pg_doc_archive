bin/spider_urls: scripts/spider_urls/scraper.go
	go build -o bin/spider_urls ./scripts/spider_urls

bin/archive: scripts/archive_ensemble/archive.go
	go build -o bin/archive ./scripts/archive_ensemble

data/url_list.tsv: bin/spider_urls
	./bin/spider_urls | sort -u > ./data/url_list.tsv

data/site.tar.gz: data/url_list.tsv bin/archive
	./bin/archive <./data/url_list.tsv 2>/tmp/archive.log | gzip -9 > ./data/site.tar.gz
