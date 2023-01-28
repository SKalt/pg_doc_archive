.PHONY: all clean
all: ./data/site.tar.gz
clean:
	rm -rf ./data/{site,url}* # it's safe to hold onto .cache, though

./data/site.tar.gz: ./data/url_list.tsv ./bin/archive
	./bin/archive <./data/url_list.tsv > ./data/site.tar.gz

./bin/spider_urls: ./scripts/spider_urls/scraper.go go.mod go.sum
	go build -o bin/spider_urls ./scripts/spider_urls

./bin/archive: ./scripts/archive_ensemble/archive.go go.mod go.sum
	go build -o bin/archive ./scripts/archive_ensemble

./data/url_list.tsv: ./bin/spider_urls
	./bin/spider_urls | sort -u > ./data/url_list.tsv
