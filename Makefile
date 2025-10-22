# Copyright (c) 2024-2025 Joshua Sing <joshua@joshuasing.dev>
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

.PHONY: all
all: lint build

.PHONY: all-deps
all-deps: deps lint-deps gen-deps

.PHONY: deps
deps:
	go mod download
	go mod verify

define LICENSE_HEADER
Copyright (c) {{.year}} {{.author}}

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
endef
export LICENSE_HEADER

.PHONY: lint
lint:
	golangci-lint fmt ./...
	golangci-lint run --fix ./...
	golicenser -tmpl="$$LICENSE_HEADER" -author="Joshua Sing <joshua@Joshuasing.dev>" -year-mode=git-range -fix ./...

.PHONY: lint-deps
lint-deps:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5
	go install github.com/joshuasing/golicenser/cmd/golicenser@v0.3

.PHONY: build
build:
	go build -trimpath -o ./bin/starlink_exporter ./cmd/starlink_exporter

.PHONY: gen
gen:
	./scripts/proto-gen.sh

.PHONY: gen-deps
gen-deps:
	@echo "Using $$(protoc --version)"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
