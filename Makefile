all: clean linux macos macm1

clean:
	rm -f build/*

linux:
	GOOS=linux GOARCH=amd64 go build -o build/notion-offliner-linux-amd64

macos:
	GOOS=darwin GOARCH=amd64 go build -o build/notion-offliner-macos

macm1:
	GOOS=darwin GOARCH=arm64 go build -o build/notion-offliner-macos-m1
