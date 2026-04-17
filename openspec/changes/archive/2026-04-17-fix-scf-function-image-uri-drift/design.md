# Design: SCF Function image_uri Drift Fix

## Affected File

`tencentcloud/services/scf/resource_tc_scf_function.go`

## Change Location

Read handler — inside the `if resp.ImageConfig != nil` block, replace:

```go
"image_uri":  imageConfigResp.ImageUri,
```

with a call to a new helper function.

## Helper Function: `normalizeImageUri`

```go
// normalizeImageUri aligns the API-returned image URI (always Format C: repo:tag@sha256:digest)
// with the user-supplied format to prevent false plan diffs.
//
// Three valid input formats:
//   A: registry/repo:tag
//   B: registry/repo@sha256:digest
//   C: registry/repo:tag@sha256:digest   ← what the API always returns
func normalizeImageUri(apiValue, userValue string) string {
    // If userValue already has both tag and digest (Format C), no normalization needed.
    if strings.Contains(userValue, ":") && strings.Contains(userValue, "@") {
        return apiValue
    }

    // Parse the API value (Format C): split on "@" first, then on ":"
    // apiValue = "registry/repo:tag@sha256:digest"
    atIdx := strings.LastIndex(apiValue, "@")
    if atIdx == -1 {
        // API returned something unexpected; return as-is
        return apiValue
    }
    apiRepoTag := apiValue[:atIdx]   // "registry/repo:tag"
    apiDigest := apiValue[atIdx+1:]  // "sha256:digest"

    // Locate the tag separator within apiRepoTag
    colonIdx := strings.LastIndex(apiRepoTag, ":")
    apiRepoOnly := apiRepoTag
    if colonIdx != -1 {
        apiRepoOnly = apiRepoTag[:colonIdx]  // "registry/repo"
    }

    // Format A: userValue has a tag but no digest → strip "@sha256:…"
    if !strings.Contains(userValue, "@") && strings.Contains(userValue, ":") {
        // Verify repo:tag match before suppressing
        if apiRepoTag == userValue {
            return userValue
        }
        return apiValue
    }

    // Format B: userValue has a digest but no tag → strip ":tag"
    if strings.Contains(userValue, "@") && !strings.Contains(userValue[:strings.Index(userValue, "@")], ":") {
        // Parse userValue: "registry/repo@sha256:digest"
        userAtIdx := strings.Index(userValue, "@")
        userRepo := userValue[:userAtIdx]
        userDigest := userValue[userAtIdx+1:]
        if apiRepoOnly == userRepo && apiDigest == userDigest {
            return userRepo + "@" + apiDigest
        }
        return apiValue
    }

    return apiValue
}
```

## Integration Point

```go
// Before (line ~990)
"image_uri": imageConfigResp.ImageUri,

// After
"image_uri": normalizeImageUri(
    helper.PString(imageConfigResp.ImageUri),
    d.Get("image_config.0.image_uri").(string),
),
```

Since `imageConfigResp.ImageUri` is `*string`, dereference safely with nil check.

## Edge Cases

| Scenario | Action |
|---------|--------|
| API returns nil ImageUri | Skip normalization, leave field unchanged |
| User value is empty string | Return API value as-is |
| API value has no `@` (unexpected format) | Return API value as-is |
| repo or digest portion does not match | Return API value as-is → real diff surfaces |
