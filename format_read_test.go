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
	"version"
)

type tc struct {
	VersionString string
	Expected      *version.Version
	RaiseErr      bool
}

func test(t *testing.T, format string, cases []tc) {
	for _, c := range cases {
		v, err := version.Parse(format, c.VersionString)
		if c.RaiseErr != (err != nil) {
			t.Errorf("error expectation failed, expected: %v, actual: %+v; format: %+v, input: %+v", c.RaiseErr, err, format, c.VersionString)
		}
		if !c.RaiseErr && err == nil && !v.EQ(c.Expected) {
			if !c.RaiseErr {
				t.Errorf("version expectation failed, expected: %+v, actual: %+v; format: %+v, input: %+v", c.Expected, v, format, c.VersionString)
			}
		}
	}
}

func TestMajorMinor(t *testing.T) {
	format := "5.4"

	cases := []tc{
		{
			"1",
			nil,
			true,
		},
		{
			"a.b",
			nil,
			true,
		},
		{
			"1.1",
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			false,
		},
		{
			"3.5",
			&version.Version{
				Major: 3,
				Minor: 5,
			},
			false,
		},
		{
			"1.1.1",
			nil,
			true,
		},
	}

	test(t, format, cases)
}

func TestMajorMinorPatch(t *testing.T) {
	format := "5.4.3"

	cases := []tc{
		{
			"1.1",
			nil,
			true,
		},
		{
			"1.1.1",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			false,
		},
		{
			"5.7.8",
			&version.Version{
				Major: 5,
				Minor: 7,
				Patch: 8,
			},
			false,
		},
		{
			"2.3",
			nil,
			true,
		},
		{
			"1.1.1a",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			true,
		},
	}

	test(t, format, cases)
}

func TestReleaseAbbrBuild(t *testing.T) {
	format := "5.4.3.b1"

	cases := []tc{
		{
			"1.1",
			nil,
			true,
		},
		{
			"1.1.1",
			nil,
			true,
		},
		{
			"1.1.1rc1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
				Build:  1,
			},
			false,
		},
		{
			"1.1.1.rc",
			nil,
			true,
		},
		{
			"1.1.1.a1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			false,
		},
		{
			"2.3.4.b5",
			&version.Version{
				Major:  2,
				Minor:  3,
				Patch:  4,
				PreRel: version.Beta,
				Build:  5,
			},
			false,
		},
		{
			"5.6.7.8",
			&version.Version{
				Major:  5,
				Minor:  6,
				Patch:  7,
				PreRel: version.Release,
				Build:  8,
			},
			false,
		},
	}

	test(t, format, cases)
}

func TestReleaseDashBuild(t *testing.T) {
	format := "5.4.3.beta-1"

	cases := []tc{
		{
			"1.1",
			nil,
			true,
		},
		{
			"1.1.1",
			nil,
			true,
		},
		{
			"1.1.1.alpha-1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			false,
		},
		{
			"1.1.1alpha-1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			false,
		},
		{
			"1.1.1alpha1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			false,
		},
		{
			"2.3.4.beta-5",
			&version.Version{
				Major:  2,
				Minor:  3,
				Patch:  4,
				PreRel: version.Beta,
				Build:  5,
			},
			false,
		},
		{
			"5.6.7.rc-8",
			&version.Version{
				Major:  5,
				Minor:  6,
				Patch:  7,
				PreRel: version.ReleaseCandidate,
				Build:  8,
			},
			false,
		},
		{
			"5.6.7.8",
			&version.Version{
				Major:  5,
				Minor:  6,
				Patch:  7,
				PreRel: version.Release,
				Build:  8,
			},
			false,
		},
	}

	test(t, format, cases)
}

func TestDashReleaseBuild(t *testing.T) {
	format := "5.4.3-beta.1"

	cases := []tc{
		{
			"1.1",
			nil,
			true,
		},
		{
			"1.1.1",
			nil,
			true,
		},
		{
			"1.1.1-alpha.1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			false,
		},
		{
			"1.1.1alpha.1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			false,
		},
		{
			"1.1.1alpha1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			false,
		},
		{
			"2.3.4-beta.5",
			&version.Version{
				Major:  2,
				Minor:  3,
				Patch:  4,
				PreRel: version.Beta,
				Build:  5,
			},
			false,
		},
		{
			"5.6.7-rc.8",
			&version.Version{
				Major:  5,
				Minor:  6,
				Patch:  7,
				PreRel: version.ReleaseCandidate,
				Build:  8,
			},
			false,
		},
		{
			"5.6.7.8",
			&version.Version{
				Major:  5,
				Minor:  6,
				Patch:  7,
				PreRel: version.Release,
				Build:  8,
			},
			false,
		},
	}

	test(t, format, cases)
}

func TestDashReleaseDashBuild(t *testing.T) {
	format := "5.4.3-beta-1"

	cases := []tc{
		{
			"1.1",
			nil,
			true,
		},
		{
			"1.1.1",
			nil,
			true,
		},
		{
			"1.1.1-alpha-1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			false,
		},
		{
			"1.1.1-alpha1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			false,
		},
		{
			"1.1.1alpha-1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			false,
		},
		{
			"1.1.1alpha1",
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			false,
		},
		{
			"2.3.4-beta-5",
			&version.Version{
				Major:  2,
				Minor:  3,
				Patch:  4,
				PreRel: version.Beta,
				Build:  5,
			},
			false,
		},
		{
			"5.6.7-rc-8",
			&version.Version{
				Major:  5,
				Minor:  6,
				Patch:  7,
				PreRel: version.ReleaseCandidate,
				Build:  8,
			},
			false,
		},
		{
			"5.6.7.8",
			nil,
			true,
		},
	}

	test(t, format, cases)
}

func TestDollarSign(t *testing.T) {
	format := "5.4$.3$.1"

	cases := []tc{
		{
			"1",
			nil,
			true,
		},
		{
			"1.1",
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			false,
		},
		{
			"1.1.1",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			false,
		},
		{
			"1.1.1rc",
			nil,
			true,
		},
		{
			"1.1.1rc",
			nil,
			true,
		},
		{
			"1.1.1-alpha1",
			nil,
			true,
		},
		{
			"1.1.1-1",
			nil,
			true,
		},
		{
			"1.1.1.1",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
				Build: 1,
			},
			false,
		},
		{
			"2.3.4.5",
			&version.Version{
				Major: 2,
				Minor: 3,
				Patch: 4,
				Build: 5,
			},
			false,
		},
	}

	test(t, format, cases)
}

func TestVMajorMinor(t *testing.T) {
	format := "v5.4"

	cases := []tc{
		{
			"v1",
			nil,
			true,
		},
		{
			"1",
			nil,
			true,
		},
		{
			"v.1",
			nil,
			true,
		},
		{
			"1.1",
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			false,
		},
		{
			"v1.1",
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			false,
		},
		{
			"v3.5",
			&version.Version{
				Major: 3,
				Minor: 5,
			},
			false,
		},
		{
			"v1.1.1",
			nil,
			true,
		},
	}

	test(t, format, cases)
}

func TestAlphabeticPatch(t *testing.T) {
	format := "5.4y"

	cases := []tc{
		{
			"v1",
			nil,
			true,
		},
		{
			"1",
			nil,
			true,
		},
		{
			"1.1",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 0,
			},
			false,
		},
		{
			"1.1a",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			false,
		},
		{
			"1.1z",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 26,
			},
			false,
		},
		{
			"1.1aa",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 27,
			},
			false,
		},
		{
			"3.5bc",
			&version.Version{
				Major: 3,
				Minor: 5,
				Patch: 2*26 + 3,
			},
			false,
		},
		{
			"3.5b1",
			nil,
			true,
		},
	}

	test(t, format, cases)
}

func TestAlphabeticBuild(t *testing.T) {
	format := "5.4.3z"

	cases := []tc{
		{
			"v1",
			nil,
			true,
		},
		{
			"1",
			nil,
			true,
		},
		{
			"1.1",
			nil,
			true,
		},
		{
			"1.1.1",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
				Build: 0,
			},
			false,
		},
		{
			"1.1.1a",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
				Build: 1,
			},
			false,
		},
		{
			"1.1.1z",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
				Build: 26,
			},
			false,
		},
		{
			"1.1.1aa",
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
				Build: 27,
			},
			false,
		},
		{
			"3.5.8bc",
			&version.Version{
				Major: 3,
				Minor: 5,
				Patch: 8,
				Build: 2*26 + 3,
			},
			false,
		},
		{
			"3.5b.1",
			nil,
			true,
		},
	}

	test(t, format, cases)
}
