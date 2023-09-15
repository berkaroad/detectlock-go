go-bench:
	GO111MODULE=on GOPROXY=https://goproxy.cn,direct go test -run=none -count=1 -benchtime=10000x -benchmem -bench=. ./... | grep Benchmark
