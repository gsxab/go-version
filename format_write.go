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
	"strconv"
	"strings"
)

// format

func (field Field) FormatField(v *Version, layout string) string {
	switch field {
	case build:
		return formatInt(v.Build)
	case alphabetic_build:
		return formatAlpha(v.Build)
	case preRelTag:
		return formatTag(layout, v.PreRel)
	case patch:
		return formatInt(v.Patch)
	case alphabetic_patch:
		return formatAlpha(v.Patch)
	case minor:
		return formatInt(v.Minor)
	case major:
		return formatInt(v.Major)
	case other:
		return v.Other
	case fixed:
		return layout
	default:
		panic("unexpected field to set")
	}
}

func formatInt(val int64) string {
	return strconv.FormatInt(val, 10)
}

func formatAlpha(val int64) string {
	str := make([]byte, 0)
	for val > 26 {
		str = append(str, byte(val%26)+'a')
		val /= 26
	}
	return string(str)
}

func formatTag(layout string, val PreRelTag) string {
	if val == Release {
		return ""
	}

	parts := make([]string, 3)

	if layout[0] == '-' {
		parts[0] = "-"
		layout = layout[1:]
	}
	if layout[len(layout)-1] == '-' {
		parts[2] = "-"
		layout = layout[:len(layout)-1]
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
	parts[1] = tagSet[val]

	return strings.Join(parts, "")
}

func Format(layout string, version *Version) (string, error) {
	// layout example: 5.4.3-beta.1(.other)
	parts := make([]string, 0)
	for len(layout) > 0 {
		fieldFmt, field, suffix, err := nextChunk(layout)
		if err != nil {
			return "", err
		}
		parts = append(parts, field.FormatField(version, fieldFmt))
		layout = suffix
	}
	return strings.Join(parts, ""), nil
}
