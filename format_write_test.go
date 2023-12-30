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

type tcW struct {
	Version  *version.Version
	Expected string
	RaiseErr bool
}

func testW(t *testing.T, format string, cases []tcW) {
	for _, c := range cases {
		s, err := version.Format(format, c.Version)
		if c.RaiseErr != (err != nil) {
			t.Errorf("error expectation failed, expected: %v, actual: %+v; format: %+v, input: %+v", c.RaiseErr, err, format, c.Version)
		}
		if !c.RaiseErr && err == nil && s != c.Expected {
			if !c.RaiseErr {
				t.Errorf("version expectation failed, expected: %+v, actual: %+v; format: %+v, input: %+v", c.Expected, s, format, c.Version)
			}
		}
	}
}

func TestWMajorMinor(t *testing.T) {
	format := "5.4"

	cases := []tcW{
		{
			&version.Version{
				Major: 1,
			},
			"1.0",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			"1.1",
			false,
		},
		{
			&version.Version{
				Major: 3,
				Minor: 5,
			},
			"3.5",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			"1.1",
			false,
		},
	}

	testW(t, format, cases)
}

func TestWMajorMinorPatch(t *testing.T) {
	format := "5.4.3"

	cases := []tcW{
		{
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			"1.1.0",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			"1.1.1",
			false,
		},
		{
			&version.Version{
				Major: 5,
				Minor: 7,
				Patch: 8,
			},
			"5.7.8",
			false,
		},
		{
			&version.Version{
				Major: 2,
				Minor: 3,
			},
			"2.3.0",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: 4,
			},
			"1.2.3",
			false,
		},
	}

	testW(t, format, cases)
}

func TestWReleaseAbbrBuild(t *testing.T) {
	format := "5.4.3b-1"

	cases := []tcW{
		{
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			"1.1.00", // ill-formed when tag=release
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
			},
			"1.1.1rc-0",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
				Build:  1,
			},
			"1.1.1rc-1",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			"1.1.1a-1",
			false,
		},
		{
			&version.Version{
				Major:  2,
				Minor:  3,
				Patch:  4,
				PreRel: version.Beta,
				Build:  5,
			},
			"2.3.4b-5",
			false,
		},
	}

	testW(t, format, cases)
}

func TestWReleaseDashBuild(t *testing.T) {
	format := "5.4.3.beta-1"

	cases := []tcW{
		{
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			"1.1.0.0",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
			},
			"1.1.1.rc-0",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
				Build:  1,
			},
			"1.1.1.rc-1",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			"1.1.1.alpha-1",
			false,
		},
		{
			&version.Version{
				Major:  2,
				Minor:  3,
				Patch:  4,
				PreRel: version.Beta,
				Build:  5,
			},
			"2.3.4.beta-5",
			false,
		},
		{
			&version.Version{
				Major:  5,
				Minor:  6,
				Patch:  7,
				PreRel: version.Release,
				Build:  8,
			},
			"5.6.7.8",
			false,
		},
	}

	testW(t, format, cases)
}

func TestWDashReleaseBuild(t *testing.T) {
	format := "5.4.3-beta.1"

	cases := []tcW{
		{
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			"1.1.0.0",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
			},
			"1.1.1-rc.0",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
				Build:  1,
			},
			"1.1.1-rc.1",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			"1.1.1-alpha.1",
			false,
		},
		{
			&version.Version{
				Major:  2,
				Minor:  3,
				Patch:  4,
				PreRel: version.Beta,
				Build:  5,
			},
			"2.3.4-beta.5",
			false,
		},
		{
			&version.Version{
				Major:  5,
				Minor:  6,
				Patch:  7,
				PreRel: version.Release,
				Build:  8,
			},
			"5.6.7.8",
			false,
		},
	}

	testW(t, format, cases)
}

func TestWDashReleaseDashBuild(t *testing.T) {
	format := "5.4.3-beta-1"

	cases := []tcW{
		{
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			"1.1.00", // ill-formed when tag=release
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
			},
			"1.1.1-rc-0",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
				Build:  1,
			},
			"1.1.1-rc-1",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.Alpha,
				Build:  1,
			},
			"1.1.1-alpha-1",
			false,
		},
		{
			&version.Version{
				Major:  2,
				Minor:  3,
				Patch:  4,
				PreRel: version.Beta,
				Build:  5,
			},
			"2.3.4-beta-5",
			false,
		},
		{
			&version.Version{
				Major:  5,
				Minor:  6,
				Patch:  7,
				PreRel: version.Release,
				Build:  8,
			},
			"5.6.78", // ill-formed when tag=release
			false,
		},
	}

	testW(t, format, cases)
}

func TestWDollarSign(t *testing.T) {
	format := "5.4$.3.1"

	cases := []tcW{
		{
			&version.Version{
				Major: 1,
			},
			"1.0",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			"1.1",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			"1.1.1.0",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
			},
			"1.1.1.0",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
				Build:  0,
			},
			"1.1.1.0",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
				Build: 1,
			},
			"1.1.1.1",
			false,
		},
		{
			&version.Version{
				Major: 2,
				Minor: 3,
				Patch: 4,
				Build: 5,
			},
			"2.3.4.5",
			false,
		},
		{
			&version.Version{
				Major: 2,
				Minor: 3,
				Build: 5,
			},
			"2.3.0.5",
			false,
		},
	}

	testW(t, format, cases)
}

func TestWDollarSign2(t *testing.T) {
	format := "5.4$.3$.1"

	cases := []tcW{
		{
			&version.Version{
				Major: 1,
			},
			"1.0",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			"1.1",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			"1.1.1",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
			},
			"1.1.1",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
				Build:  0,
			},
			"1.1.1",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
				Build: 1,
			},
			"1.1.1.1",
			false,
		},
		{
			&version.Version{
				Major: 2,
				Minor: 3,
				Patch: 4,
				Build: 5,
			},
			"2.3.4.5",
			false,
		},
		{
			&version.Version{
				Major: 2,
				Minor: 3,
				Build: 5,
			},
			"2.3.0.5",
			false,
		},
	}

	testW(t, format, cases)
}

func TestWVMajorMinor(t *testing.T) {
	format := "v5.4"

	cases := []tcW{
		{
			&version.Version{
				Major: 1,
			},
			"v1.0",
			false,
		},
		{
			&version.Version{
				Minor: 1,
			},
			"v0.1",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			"v1.1",
			false,
		},
		{
			&version.Version{
				Major: 3,
				Minor: 5,
			},
			"v3.5",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
				Build: 1,
			},
			"v1.1",
			false,
		},
	}

	testW(t, format, cases)
}

func TestWAlphabeticPatch(t *testing.T) {
	format := "5.4y"

	cases := []tcW{
		{
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			"1.1",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			"1.1a",
			false,
		},
		{
			&version.Version{
				Major: 5,
				Minor: 7,
				Patch: 8,
			},
			"5.7h",
			false,
		},
		{
			&version.Version{
				Major: 2,
				Minor: 3,
			},
			"2.3",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: 4,
			},
			"1.2c",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3*26 + 4,
				Build: 4,
			},
			"1.2cd",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 2,
				Patch: 2*26*26 + 3*26 + 4,
				Build: 4,
			},
			"1.2bcd",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 2,
				Patch: 26*26 + 26,
				Build: 4,
			},
			"1.2zz",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 2,
				Patch: 27*26 + 1,
				Build: 4,
			},
			"1.2aaa",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 2,
				Patch: 26*26*26 + 26*26 + 26,
				Build: 4,
			},
			"1.2zzz",
			false,
		},
	}

	testW(t, format, cases)
}

func TestWAlphabeticBuild(t *testing.T) {
	format := "5.4.3z"

	cases := []tcW{
		{
			&version.Version{
				Major: 1,
			},
			"1.0.0",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
			},
			"1.1.0",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			"1.1.1",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
			},
			"1.1.1",
			false,
		},
		{
			&version.Version{
				Major:  1,
				Minor:  1,
				Patch:  1,
				PreRel: version.ReleaseCandidate,
				Build:  0,
			},
			"1.1.1",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
				Build: 1,
			},
			"1.1.1a",
			false,
		},
		{
			&version.Version{
				Major: 2,
				Minor: 3,
				Patch: 4,
				Build: 5,
			},
			"2.3.4e",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
				Build: 27,
			},
			"1.1.1aa",
			false,
		},
		{
			&version.Version{
				Major: 3,
				Minor: 5,
				Patch: 8,
				Build: 2*26 + 3,
			},
			"3.5.8bc",
			false,
		},
		{
			&version.Version{
				Major: 3,
				Minor: 5,
				Patch: 8,
				Build: 26*26 + 2*26 + 3,
			},
			"3.5.8abc",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: 26*26 + 26,
			},
			"1.2.3zz",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: 27*26 + 1,
			},
			"1.2.3aaa",
			false,
		},
		{
			&version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: 26*26*26 + 26*26 + 26,
			},
			"1.2.3zzz",
			false,
		},
	}

	testW(t, format, cases)
}
