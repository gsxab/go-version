# Go-Version

It is a golang library to parse formatted version tags,
designed to resemble the formatting and parsing API of `time` in the standard library.

## Terms

To describe components in version strings,
we generally use terms like semver,
naming the common scheme of four separated parts *major.minor.patch.build* .
However, to make it generalized to other common versioning schemes, we have to tweak the definitions a little.

<dl>
    <dt><strong>Counter</strong></dt>
    <dd>What “count”s, i.e., to increment by one each time.</dd>
    <dt><strong>Numeric Counter</strong></dt>
    <dd>A numeric counter counts as <em>1</em>, <em>2</em>, <em>3</em>, <em>4</em>, etc.</dd>
    <dt><strong>Alphabetic Counter</strong></dt>
    <dd>An alphabetic counter counts as <em>a</em>, <em>b</em>, <em>c</em>, <em>d</em>, etc., and it counts as <em>aa</em>, <em>ab</em>, <em>ac</em> after <em>z</em>. <br/> If capitals are used, i.e. <em>A</em>, <em>B</em>, <em>C</em>, and <em>Z</em>, <em>AA</em>, <em>AB</em>, we say the counter is in capitals. </dd>
    <dt>Roman Counter</dt>
    <dd>An roman counter counts as <em>i</em>, <em>ii</em>, <em>iii</em>, <em>iv</em>, etc. <br/> If capitals are used, i.e. <em>I</em>, <em>II</em>, <em>III</em>, <em>IV</em>, we say the counter is in capitals. <br/> Roman counters are only reserved and not planned to be implemented now. </dd>
    <dt><strong>Major</strong></dt>
    <dd>The first counter in the versioning scheme.</dd>
    <dt><strong>Minor</strong></dt>
    <dd>The counter following <strong>major</strong> in the versioning scheme.</dd>
    <dt><strong>Patch</strong></dt>
    <dd>The counter following <strong>minor</strong> in the versioning scheme. But never after a <strong>pre-release tag</strong>.</dd>
    <dt><strong>Pre-Release Tag</strong> (Pre-Rel Tag)</dt>
    <dd>The tag showing pre-release stages, i.e. <em>alpha</em>, <em>beta</em> and <em>rc</em> for release candidate. A version without a pre-rel tag is considered to be tagged as release.</dd>
    <dt><strong>Build</strong></dt>
    <dd>The counter following <strong>patch</strong> and/or <strong>pre-release tag</strong> .</dd>
    <dt><strong>Zero</strong></dt>
    <dd>A version <em>1.1</em> is considered the same as <em>1.1.0</em>, where the <strong>patch</strong> is zero. <br/> <em>0</em> is considered zero for numeric counters, an empty string for alphabetic and roman counters, and release for pre-rel tags.</dd>
</dl>

## Placeholders

As in the `time` library, we use a specific special version as the format pattern. It is:

```text
5.4.3-beta.1
```

The 5<sup>th</sup> major version, the 4<sup>th</sup> minor version, the 3<sup>rd</sup> patch, the 2<sup>nd</sup> pre-release tag, and the 1<sup>st</sup> build.
Quite easy to remember, right?

Note that we perform full match on the pattern to ensure expected results,
matching pattern `5.4.3` against either `1.1.1.1` or `1.1` gives an error.

To allow a scheme to end earlier, for example, `1.2.3` for a common version but `1.2` for `1.2.0`,
we introduced a `$` sign for “allowing string end” to mark it `5.4$.3`.
In this case, we allow an end-of-string directly following the minor, while still rejecting an ill-formed end in input like `1.2.`.
When writing, the sign `$` finishes the string if all tokens following it are defined as omittable, e.g. `.0.0` in `1.2.0.0`.

Besides, all literals in version strings, like `v` and `.`, are optional when reading.

Consult the following table for details about all tokens in pattern strings.

| Pattern | Parsing / Reading | Formatting / Writing | Omittable |
| :-: | --- | --- | --- |
| `v` | Reads an optional `v`. | Writes a `v`. | Always. |
| `V` | Reads an optional `V`. | Writes a `V`. | Always. |
| `.` | Reads an optional `.`. | Writes a `.`. | Always. |
| `$` | Finishes reading if the end of the string is met. | Does not write any suffix if every token in the suffix is omittable. | - |
| `5` | Reads a numeric major. | Writes a numeric major. | If zero. |
| `4` | Reads a numeric minor. | Writes a numeric minor. | If zero. |
| `3` | Reads a numeric patch. | Writes a numeric patch. | If zero. |
| `y` | Reads an alphabetic patch. | Writes an alphabetic patch. | If zero. |
| `Y` | Reads an alphabetic patch in capitals. | Writes an alphabetic patch in capitals. | If zero. |
| (reserved, not implemented yet) <br/> `i` | Reads an roman patch. | Writes an roman patch. | If zero. |
| (reserved, not implemented yet) <br/> `I` | Reads an roman patch in capitals. | Writes an roman patch in capitals. | If zero. |
| `b` | Reads a pre-rel tag, `a` for alpha, `b` for beta, `rc` for release candidate or nothing for release. | Writes a pre-rel tag, `a` for alpha, `b` for beta, `rc` for release candidate or nothing or release. | If zero. |
| `B` | Like `b`, but reads capitals instead. | Like `b`, but writes in capitals instead. | If zero. |
| `beta` | Reads a pre-rel tag, `alpha` for alpha, `beta` for beta, `rc` for release candidate or nothing or release. | Writes a pre-rel tag, `alpha` for alpha, `beta` for beta, `rc` for release candidate or nothing or release. | If zero. |
| `Beta` | Like `beta`, but reads in title case instead, i.e. `Alpha`, `Beta`, `RC`. | Like `beta`, but writes in title case instead, i.e. `Alpha`, `Beta`, `RC`. | If zero. |
| `BETA` | Like `beta`, but reads in all caps instead, i.e. `ALPHA`, `BETA`, `RC`. | Like `beta`, but writes in all caps instead, i.e. `ALPHA`, `BETA`, `RC`. | If zero. |
| `-b`, `-beta`, etc. | Like `b`, `beta`, etc., and reads an optional hythen before the tag. | Like `b`, `beta`, etc., and writes a hythen before the tag unless it is a release. | If zero. |
| `b-`, `beta-`, etc. | Like `b`, `beta`, etc., and reads an optional hythen after the tag. | Like `b`, `beta`, etc., and writes a hythen after the tag unless it is a release. | If zero. |
| `-b-`, `-beta-`, etc. | Like `b`, `beta`, etc., and reads optional hythens both before and after the tag. | Like `b`, `beta`, etc., and writes hythens both before and after the tag unless it is a release. | If zero. |
| (not implemented yet) <br/> `b?`, `beta?`, `-b-?` | Like `b`, `beta`, `-b-`, etc., but reads a `.` as release instead. | Like `b`, `beta`, `-b-`, but writes a `.` if it is a release. | If zero. |
| `1` | Reads a numeric build. | Writes a numeric build. | If zero. |
| `z` | Reads an alphabetic build. | Writes an alphabetic build. | If zero. |
| `Z` | Reads an alphabetic build in capitals. | Writes an alphabetic build in capitals. | If zero. |
| (reserved, not implemented yet) <br/> `j` | Reads an roman build. | Writes an roman build. | If zero. |
| (reserved, not implemented yet) <br/> `J` | Reads an roman build in capitals. | Writes an roman build in capitals. | If zero. |
| `o` | Reads all remaining text as “other”, which not considered to be part of the version number. | Writes the stored“other”. | Always. |
| (for robustness only) <br/> other | Reads the character optionally. | Writes the character. | Always. |

The token `i`/`I` and `j`/`J` are only reserved for roman numbers.
Anyway, I do not expected any of their presence in versioning though.

Any characters out of the tokens list are used to read and write to keep compability,
but it is *highly NOT recommended* because this behavior may change for some characters in the future,
if they are appended to the list as new tokens.
This forward-compatibility behavior may be removed in later versions.
