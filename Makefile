BINARY=wikiracer

bench:
	go test -bench -v

test:
	go test -v

build:
	go get github.com/PurpureGecko/go-lfc
	go get github.com/patrickmn/go-cache
	go build -o ${BINARY} main.go worker.go parser.go

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

