/*
 * SPDX-License-Identifier: Apache-2.0
 *
 * Copyright (c) 2023 Gsxab
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package version

type Field int

const (
	build     Field = 1
	preRelTag Field = 2
	patch     Field = 3
	minor     Field = 4
	major     Field = 5
	other     Field = 1 + iota
	fixed
	allowEnd
	alphabetic_build
	alphabetic_patch
)

// format tokenizer

func nextChunk(layout string) (string, Field, string, error) {
	// number field
	if isAsciiNum(layout[0]) {
		value, offset, err := readInt(layout)
		if err != nil {
			return "", 0, "", err
		}
		return layout[:offset], Field(value), layout[offset:], nil
	}
	// alphabetic field
	if layout[0] == 'z' || layout[0] == 'Z' {
		return layout[:1], alphabetic_build, layout[1:], nil
	}
	if layout[0] == 'y' || layout[0] == 'Y' {
		return layout[:1], alphabetic_patch, layout[1:], nil
	}
	// allow end
	if layout[0] == '$' {
		return layout[:1], allowEnd, layout[1:], nil
	}
	// dot
	if layout[0] == '.' {
		return layout[:1], fixed, layout[1:], nil
	}
	// v
	if layout[0] == 'v' || layout[0] == 'V' {
		return layout[:1], fixed, layout[1:], nil
	}
	// other
	if layout[0] == 'o' {
		return layout, other, "", nil
	}
	// pre-release tag
	index := 0
	if layout[0] == '-' {
		index = 1
	}
	if len(layout) <= index || (layout[index] != 'b' && layout[index] != 'B') {
		return layout[:1], fixed, layout[1:], nil
	}
	// next: -?b(eta)?-?
	if len(layout) >= index+4 && (layout[index:index+4] == "beta" || layout[index:index+4] == "Beta") {
		index += 4
	} else {
		index++
	}
	if len(layout) > index && layout[index] == '-' {
		index++
	}
	return layout[:index], preRelTag, layout[index:], nil
}
