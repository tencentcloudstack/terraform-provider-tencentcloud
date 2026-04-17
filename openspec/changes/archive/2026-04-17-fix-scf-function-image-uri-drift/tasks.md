## 1. Resource Fix

- [x] 1.1 Add `normalizeImageUri(apiValue, userValue string) string` helper function to `resource_tc_scf_function.go`
- [x] 1.2 In the Read handler, replace `"image_uri": imageConfigResp.ImageUri` with a call to `normalizeImageUri`, passing the API value and the current state value `d.Get("image_config.0.image_uri").(string)`
- [x] 1.3 Switch `image_uri` from Read-time normalization to `DiffSuppressFunc` approach: revert the Read handler to write the raw API value, add `DiffSuppressFunc` on the schema field that calls `normalizeImageUri` to compare config vs state
