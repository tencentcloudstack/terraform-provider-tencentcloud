# Fix SCF Function image_uri Drift

## What

In `tencentcloud_scf_function`, the `image_config.image_uri` field accepts three equivalent formats:

- **Format A** (tag only): `registry/repo:tag`
- **Format B** (digest only): `registry/repo@sha256:digest`
- **Format C** (tag + digest): `registry/repo:tag@sha256:digest`

The API always returns Format C. When users supply Format A or Format B, the Read handler writes the Format C value verbatim to state, causing a perpetual plan diff on every subsequent apply.

## Why

Users cannot suppress the false diff without changing their `.tf` file to Format C, which defeats the purpose of supporting multiple formats.

## Fix

In the Read handler, after receiving `imageConfigResp.ImageUri` (always Format C), apply format-matching logic to align the value stored in state with the user-supplied format:

1. If the user-supplied value is Format C → store the API response as-is.
2. If the user-supplied value is Format A (`repo:tag`, no digest) → verify that `repo:tag` portion of the API response matches, then strip `@sha256:…` before storing.
3. If the user-supplied value is Format B (`repo@sha256:digest`, no tag) → verify that `repo` and `sha256:digest` portions match, then strip `:tag` before storing.
4. If none of the above match (genuinely different image) → store the API response as-is so the real diff is surfaced.
