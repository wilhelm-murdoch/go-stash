BIN=go
BUILD_DATE=$(shell date +'%d-%m-%Y')
BUILD_REF=v0.0.0
BUILD_STAGE=development
BUILD_COMMIT=$(shell git rev-parse HEAD)

reset:
	rm dist/*.json dist/tags/*.json

clean:
	rm stash *.out *.test 

build:
	CGO_ENABLED=0 ${BIN} build -a -installsuffix cgo -ldflags="-s -w -X main.Date=${BUILD_DATE} -X main.Version=${BUILD_REF} -X main.Stage=${BUILD_STAGE} -X main.Commit=${BUILD_COMMIT}" ./cmd/stash

test:
	${BIN} test -race -v .

lint:
	staticcheck -f stylish

bench:
	${BIN} test -run . -bench . -benchtime 5s -count 10 -benchmem -cpuprofile cpu.out -memprofile mem.out -trace trace.out

pprof-cpu:
	${BIN} tool pprof -http :8800 cpu.out

pprof-mem:
	${BIN} tool pprof -http :8900 mem.out

trace:
	${BIN} tool trace trace.out

coverage:
	${BIN} test -v -coverprofile cover.out .