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

package version_test

import (
	"testing"

	"github.com/gsxab/go-version"
)

func TestVersion(t *testing.T) {
	v1 := &version.Version{
		Major:  2,
		Minor:  3,
		Patch:  4,
		PreRel: version.Beta,
	}
	v2 := &version.Version{
		Major:  1,
		Minor:  4,
		Patch:  5,
		PreRel: version.Release,
	}
	v3 := &version.Version{
		Major:  3,
		Minor:  4,
		Patch:  5,
		PreRel: version.Release,
	}
	v4 := &version.Version{
		Major:  2,
		Minor:  2,
		Patch:  4,
		PreRel: version.Release,
	}
	v5 := &version.Version{
		Major:  2,
		Minor:  4,
		Patch:  4,
		PreRel: version.Release,
	}
	v6 := &version.Version{
		Major:  2,
		Minor:  3,
		Patch:  3,
		PreRel: version.Release,
	}
	v7 := &version.Version{
		Major:  2,
		Minor:  3,
		Patch:  5,
		PreRel: version.Release,
	}
	v8 := &version.Version{
		Major:  2,
		Minor:  3,
		Patch:  4,
		PreRel: version.Alpha,
	}
	v9 := &version.Version{
		Major:  2,
		Minor:  3,
		Patch:  4,
		PreRel: version.ReleaseCandidate,
	}
	v10 := &version.Version{
		Major:  2,
		Minor:  3,
		Patch:  4,
		PreRel: version.Beta,
	}
	v11 := &version.Version{
		Major:  2,
		Minor:  3,
		Patch:  4,
		PreRel: version.Beta,
		Other:  "123",
	}
	v12 := &version.Version{
		Major:  2,
		Minor:  3,
		Patch:  4,
		PreRel: version.Beta,
		Build:  123,
	}

	expectedLess := []*version.Version{
		v2, v4, v6, v8,
	}
	expectedGreater := []*version.Version{
		v3, v5, v7, v9, v12,
	}
	expectedEqual := []*version.Version{
		v10, v11,
	}

	for _, v := range expectedLess {
		if v.EQ(v1) {
			t.Errorf("EQ returns true when expecting less, lhs=%v, rhs=%v", v1, v)
		}
		if !v.LT(v1) {
			t.Errorf("LT returns false when expecting less, lhs=%v, rhs=%v", v1, v)
		}
		if !v.LE(v1) {
			t.Errorf("LE returns true when expecting less, lhs=%v, rhs=%v", v1, v)
		}
	}

	for _, v := range expectedGreater {
		if v.EQ(v1) {
			t.Errorf("EQ returns true when expecting greater, lhs=%v, rhs=%v", v1, v)
		}
		if v.LT(v1) {
			t.Errorf("LT returns true when expecting greater, lhs=%v, rhs=%v", v1, v)
		}
		if v.LE(v1) {
			t.Errorf("LE returns true when expecting greater, lhs=%v, rhs=%v", v1, v)
		}
	}

	for _, v := range expectedEqual {
		if !v.EQ(v1) {
			t.Errorf("EQ returns false when expecting equal, lhs=%v, rhs=%v", v1, v)
		}
		if v.LT(v1) {
			t.Errorf("LT returns true when expecting equal, lhs=%v, rhs=%v", v1, v)
		}
		if !v.LE(v1) {
			t.Errorf("LE returns false when expecting equal, lhs=%v, rhs=%v", v1, v)
		}
	}
}
