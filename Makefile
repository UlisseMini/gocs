# Copyright 2019 Ulisse Mini. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

all: clean test release

LDFLAGS=-ldflags="-s -w"

clean:
	@rm -rf "bin"

release:
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=windows packr build $(LDFLAGS) -o bin/goc_win64.exe
	CGO_ENABLED=0 GOOS=linux packr   build $(LDFLAGS) -o bin/goc_linux64
	CGO_ENABLED=0 GOOS=darwin packr  build $(LDFLAGS) -o bin/goc_mac64

test:
	@go test

upx:
	upx --best bin/*
