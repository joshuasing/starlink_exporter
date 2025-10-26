#!/bin/sh
#
# Copyright (c) 2025 Joshua Sing <joshua@joshuasing.dev>
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
#

set -eu

CHANGELOG_FILE="CHANGELOG.md"

fatal() {
	echo "$*" 1>&2
	exit 1
}

run() {
	version="$1"
	changelog=""

	# Check if the changelog file exists.
	if [ ! -f "$CHANGELOG_FILE" ]; then
		fatal "error: changelog file not found: $CHANGELOG_FILE"
	fi

	# Read the changelog file line by line.
	found_version=false
	while IFS= read -r line; do
		if ! $found_version; then
			# Find the version line.
			if echo "$line" | grep -Eq "^## \[$version\]"; then
				found_version=true
			fi
			continue
		fi

		# Stop when we reach another version or a horizontal separator with 5 or
		# more dashes (as used in the footer).
		echo "$line" | grep -Eq "^(-----+|## \[[0-9]+\.[0-9]+\.[0-9]+\])" && break
		changelog="${changelog}${line}\n"
	done < "$CHANGELOG_FILE"

	# If the version was not found, print an error.
	if ! $found_version; then
		fatal "error: version was not found in changelog: $version" 1>&2
	fi

	# Tidy up the changelog for displaying on GitHub.
	changelog=$(echo "$changelog" | sed -e '1{/^$/d;}' \
		-e 's/^### /## /' -e 's/^#### /### /' -e 's/^##### /#### /')

	# Fixup usernames
	changelog=$(echo "$changelog" | \
		sed -E 's/\[@([A-Za-z0-9_-]+)\]\(https:\/\/github\.com\/\1\)/@\1/g')

	# Print the changelog for the version.
	echo "$changelog"
}

# Check the version argument was provided.
if [ "$#" -ne 1 ]; then
    fatal "usage: $0 <version>"
fi

run "$1"
