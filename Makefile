BINARY_NAME=netmap
VERSION=v0.1.3
DIST_DIR=dist

build:
	mkdir -p ${DIST_DIR}
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -v -o ${DIST_DIR}/${BINARY_NAME}-darwin-amd64-${VERSION} .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -v -o ${DIST_DIR}/${BINARY_NAME}-darwin-arm64-${VERSION} .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"  -v -o ${DIST_DIR}/${BINARY_NAME}-linux-amd64-${VERSION} .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -v -o ${DIST_DIR}/${BINARY_NAME}-linux-arm64-${VERSION} .
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -v -o ${DIST_DIR}/${BINARY_NAME}-windows-amd64-${VERSION}.exe .

clean:
	go clean
	rm -rf ${DIST_DIR}