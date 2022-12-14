# Mgmt
# Copyright (C) 2013-2022+ James Shubin and the project contributors
# Written by James Shubin <james@shubin.ca> and the project contributors
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

SHELL = /usr/bin/env bash
.PHONY: all build clean fuzz

all: build

build: parser/lexer.nn.go parser/y.go interpolate/parse.generated.go
	@# recursively run make in child dir named types
	@$(MAKE) --quiet -C types

clean:
	$(MAKE) --quiet -C types clean
	@rm -f parser/lexer.nn.go parser/y.go parser/y.output interpolate/parse.generated.go || true

parser/lexer.nn.go: parser/lexer.nex
	@echo "Generating: lexer..."
	nex -e -o $@ $<
	@ROOT="$$( cd "$$( dirname "$${BASH_SOURCE[0]}" )" && cd .. && pwd )" && $$ROOT/misc/header.sh 'parser/lexer.nn.go'

parser/y.go: parser/parser.y
	@echo "Generating: parser..."
	goyacc -v parser/y.output -o $@ $<
	@ROOT="$$( cd "$$( dirname "$${BASH_SOURCE[0]}" )" && cd .. && pwd )" && $$ROOT/misc/header.sh 'parser/y.go'

interpolate/parse.generated.go: interpolate/parse.rl
	@echo "Generating: interpolation..."
	ragel -Z -G2 -o $@ $<
	# remove the ragel header so our header test passes
	sed -i -e 1,2d $@
	gofmt -s -w $@

fuzz:
	@$(MAKE) --quiet -C fuzz
