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

type PreRelTag int64

const (
	Alpha PreRelTag = iota - 3
	Beta
	ReleaseCandidate
	Release
)

type Version struct {
	Major  int64
	Minor  int64
	Patch  int64
	PreRel PreRelTag
	Build  int64
	Other  string
}

func (v *Version) EQ(v2 *Version) bool {
	return v.Major == v2.Major && v.Minor == v2.Minor && v.Patch == v2.Patch &&
		v.PreRel == v2.PreRel && v.Build == v2.Build
}

func (v *Version) LT(v2 *Version) bool {
	return v.Major < v2.Major || v.Major == v2.Major &&
		(v.Minor < v2.Minor || v.Minor == v2.Minor &&
			(v.Patch < v2.Patch || v.Patch == v2.Patch &&
				(v.PreRel < v2.PreRel || v.PreRel == v2.PreRel &&
					(v.Build < v2.Build))))
}

func (v *Version) LE(v2 *Version) bool {
	return v.Major < v2.Major || v.Major == v2.Major &&
		(v.Minor < v2.Minor || v.Minor == v2.Minor &&
			(v.Patch < v2.Patch || v.Patch == v2.Patch &&
				(v.PreRel < v2.PreRel || v.PreRel == v2.PreRel &&
					(v.Build <= v2.Build))))
}
