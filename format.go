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

import (
	"fmt"
	"strconv"
	"strings"
)

type Field int

const (
	other Field = iota
	build
	preRelTag
	patch
	minor
	major
	fixed
	allowEnd
	alphabetic_build
	alphabetic_patch
)

func (field Field) SetField(v *Version, val int64) {
	switch field {
	case build, alphabetic_build:
		v.Build = val
	case preRelTag:
		v.PreRel = PreRelTag(val)
	case patch, alphabetic_patch:
		v.Patch = val
	case minor:
		v.Minor = val
	case major:
		v.Major = val
	default:
		panic("unexpected field to set")
	}
}

func (field Field) Read(v *Version, layout string, source string) (int, error) {
	switch field {
	case build, patch, minor, major:
		val, offset, err := readInt(source)
		if err != nil {
			return 0, err
		}
		field.SetField(v, val)
		return offset, nil
	case alphabetic_build, alphabetic_patch:
		val, offset, err := readAlpha(source)
		if err != nil {
			return 0, err
		}
		field.SetField(v, val)
		return offset, nil
	case preRelTag:
		tag, offset, err := readTag(layout, source)
		if err != nil {
			return 0, err
		}
		v.PreRel = tag
		return offset, nil
	case fixed:
		if strings.HasPrefix(source, layout) {
			return len(layout), nil
		} else {
			return 0, nil
		}
	case other:
		v.Other = source
		return len(layout), nil
	case allowEnd:
		if source == "" {
			return 1, nil
		}
		return 0, nil
	default:
		panic("unexpected field to set")
	}
}

func readInt(source string) (int64, int, error) {
	var i int
	for i = 0; i < len(source); i++ {
		if !isAsciiNum(source[i]) {
			break
		}
	}
	val, err := strconv.ParseInt(source[:i], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return val, i, nil
}

func readAlpha(source string) (int64, int, error) {
	var i int
	for i = 0; i < len(source); i++ {
		if !isAsciiAlpha(source[i]) {
			break
		}
	}
	if i == 0 {
		return 0, 0, nil
	}
	str := source[:i]
	val := int64(0)
	mod := int64(1)
	for j := len(str) - 1; j >= 0; j-- {
		val = val + alphaToNumOrd(str[j])*mod
		mod = mod * 26
	}
	return val, i, nil
}

func readTag(layout string, source string) (tag PreRelTag, offset int, err error) {
	tag = Release
	if len(source) == 0 {
		return
	}

	var prefixDash, suffixDash bool
	if layout[0] == '-' {
		prefixDash = true
		layout = layout[1:]
	}
	if layout[len(layout)-1] == '-' {
		suffixDash = true
		layout = layout[:len(layout)-1]
	}

	if prefixDash {
		if len(source) > 0 && source[0] == '-' {
			source = source[1:]
			offset++
		}
	}
	tagSet := make(map[PreRelTag]string)
	switch layout {
	case "b":
		tagSet = map[PreRelTag]string{
			Alpha:            "a",
			Beta:             "b",
			ReleaseCandidate: "rc",
		}
	case "B":
		tagSet = map[PreRelTag]string{
			Alpha:            "A",
			Beta:             "B",
			ReleaseCandidate: "RC",
		}
	case "beta":
		tagSet = map[PreRelTag]string{
			Alpha:            "alpha",
			Beta:             "beta",
			ReleaseCandidate: "rc",
		}
	case "Beta":
		tagSet = map[PreRelTag]string{
			Alpha:            "Alpha",
			Beta:             "Beta",
			ReleaseCandidate: "RC",
		}
	}
	for tagVal, tagName := range tagSet {
		if strings.HasPrefix(source, tagName) {
			tag = tagVal
			source = source[len(tagName):]
			offset += len(tagName)
			break
		}
	}
	if suffixDash && len(source) > 0 && source[0] == '-' {
		source = source[1:]
		offset++
	}

	_ = source // suppress unused value caused by advancing in source
	return
}

func Parse(layout string, versionString string) (*Version, error) {
	// layout example: 5.4.3-beta.1(.other)
	v := &Version{
		Major:  0,
		Minor:  0,
		Patch:  0,
		PreRel: 0,
		Build:  0,
	}
	for len(layout) > 0 {
		fieldFmt, field, suffix, err := nextChunk(layout)
		if err != nil {
			return nil, err
		}
		// fmt.Printf("fieldFmt=%+v, field=%+v\n", fieldFmt, field)
		advance, err := field.Read(v, fieldFmt, versionString)
		if err != nil {
			return nil, err
		}
		if field == allowEnd && advance == 1 {
			break // allow end, and meets end of versionString
		}
		layout = suffix
		// fmt.Printf("advance=%+v, taken=%+v, remaining=%+v\n", advance, versionString[:advance], versionString[advance:])
		versionString = versionString[advance:]
	}
	if len(versionString) > 0 {
		return nil, fmt.Errorf("version string not ended, left: %s", versionString)
	}
	return v, nil
}

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
