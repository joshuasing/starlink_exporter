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

#
# Usage: ./scripts/contributors.sh
# Creates a list of contributors since the last tag.
#

GH_REPO="joshuasing/starlink_exporter"

# Fetch previous tag
previous_tag=$(gh release view --repo "$GH_REPO" --json tagName | jq -r '.tagName // empty')

# Retrieve all contributors
contributors=$(gh api -H "Accept: application/vnd.github+json" -H "X-GitHub-Api-Version: 2022-11-28" \
	"/repos/${GH_REPO}/compare/$previous_tag...HEAD" --jq '.commits[].author.login' | sed '/^null$/d' | sort -u)

# Print in Markdown format with link to profile.
echo "$contributors" | while IFS= read -r user; do
	if [ -z "$user" ] || (echo "$user" | grep -qE "\[bot]"); then
		# Ignore empty lines or bot users
		continue
	fi
	echo "- [@$user](https://github.com/$user)"
done
