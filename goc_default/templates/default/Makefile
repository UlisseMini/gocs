# Copyright {{.Year}} {{.Author}}. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

all: clean test release

clean:
	@rm -rf "bin"

release:
	@mkdir bin
	CGO_ENABLED=0 GOOS=windows go build -o bin/{{.Project}}_win64.exe
	CGO_ENABLED=0 GOOS=linux go build -o bin/{{.Project}}_linux64
	CGO_ENABLED=0 GOOS=darwin go build -o bin/{{.Project}}_mac64

test:
	@go test -cover ./...
