# Tasks: Add decode_values Field and Optimize ID for tencentcloud_kubernetes_addons

## Preparation
- [ ] Review current implementation of `data_source_tc_kubernetes_addons.go`
- [ ] Verify base64 encoding/decoding library availability in Go
- [ ] Confirm `cluster_id` is always available and non-empty in the read function
- [ ] Backup current code state

## Implementation

### Task 1: Add decode_values Field to Schema
- [ ] Open file `tencentcloud/services/tke/data_source_tc_kubernetes_addons.go`
- [ ] Locate the `addons` schema definition (around line 31-64)
- [ ] Add new `decode_values` field after `raw_values` (after line 51):
  ```go
  "decode_values": {
      Type:        schema.TypeString,
      Computed:    true,
      Description: "Add-on parameters in JSON format (decoded from raw_values). Note: This field may return empty string if raw_values is not set.",
  },
  ```
- [ ] Verify field is properly positioned within the schema map
- [ ] Save file

### Task 2: Add Base64 Import
- [ ] Check if `encoding/base64` is already imported at the top of file
- [ ] If not present, add to import block (around line 4-13):
  ```go
  "encoding/base64"
  ```
- [ ] Ensure imports are properly formatted
- [ ] Save file

### Task 3: Implement decode_values Logic
- [ ] Locate the data processing section in `dataSourceTencentCloudKubernetesAddonsRead` function
- [ ] Find where `raw_values` is set (around line 123-125)
- [ ] Add decoding logic immediately after setting `raw_values`:
  ```go
  if addons.RawValues != nil {
      addonsMap["raw_values"] = addons.RawValues
      
      // Decode base64 to get decode_values
      if decodedBytes, err := base64.StdEncoding.DecodeString(*addons.RawValues); err == nil {
          addonsMap["decode_values"] = string(decodedBytes)
      } else {
          // If decode fails, set empty string
          addonsMap["decode_values"] = ""
      }
  }
  ```
- [ ] Verify error handling is in place
- [ ] Save file

### Task 4: Update Resource ID to Use cluster_id
- [ ] Locate the ID setting code (line 142)
- [ ] Replace `d.SetId(helper.DataResourceIdsHash(ids))` with:
  ```go
  d.SetId(clusterId)
  ```
- [ ] Verify `clusterId` variable is available in scope (should be defined around line 85-90)
- [ ] Remove or comment out the `ids` variable construction (line 109, 135) if no longer needed
  - **Note**: Keep `ids` construction for now as it may be used elsewhere, verify first
- [ ] Save file

### Task 5: Code Cleanup (Optional)
- [ ] Review if `ids` variable (line 109) is still needed
- [ ] If only used for the removed hash ID, consider removing:
  ```go
  // Line 109: ids := make([]string, 0, len(respData))
  // Line 135: ids = append(ids, strings.Join([]string{clusterId, *addons.AddonName}, tccommon.FILED_SP))
  ```
- [ ] Verify no other code depends on `ids` variable
- [ ] Save file if changes made

### Task 6: Format Code
- [ ] Run `go fmt` on the modified file:
  ```bash
  go fmt tencentcloud/services/tke/data_source_tc_kubernetes_addons.go
  ```
- [ ] Verify no formatting errors
- [ ] Review the formatted output

## Testing

### Task 7: Compile Verification
- [ ] Run compilation check:
  ```bash
  go build ./tencentcloud/services/tke/
  ```
- [ ] Fix any compilation errors if present
- [ ] Verify no type errors or undefined references

### Task 8: Manual Testing (Recommended)
- [ ] Create a test Terraform configuration:
  ```hcl
  data "tencentcloud_kubernetes_addons" "test" {
    cluster_id = "cls-xxxxxx"  # Use a real cluster ID
    addon_name = "nginx-ingress"
  }
  
  output "raw_values" {
    value = data.tencentcloud_kubernetes_addons.test.addons[0].raw_values
  }
  
  output "decode_values" {
    value = data.tencentcloud_kubernetes_addons.test.addons[0].decode_values
  }
  
  output "id" {
    value = data.tencentcloud_kubernetes_addons.test.id
  }
  ```
- [ ] Run `terraform init`
- [ ] Run `terraform plan`
- [ ] Verify `decode_values` contains decoded JSON
- [ ] Verify ID matches `cluster_id`
- [ ] Test with addon that has no `raw_values` (should handle gracefully)

### Task 9: Edge Case Testing
- [ ] Test with cluster that has no addons
- [ ] Test with addon that has empty `raw_values`
- [ ] Test with invalid base64 in `raw_values` (verify empty string is set)
- [ ] Test with multiple addons in the same cluster

### Task 10: State Migration Testing (Optional but Recommended)
- [ ] Use existing Terraform state with old provider version
- [ ] Upgrade to new provider version
- [ ] Run `terraform refresh`
- [ ] Verify state is updated correctly
- [ ] Run `terraform plan` and verify no unexpected changes

## Documentation

### Task 11: Update Documentation (If Exists)
- [ ] Check if documentation file exists for this data source
- [ ] If yes, update with new `decode_values` field description
- [ ] Add example showing both `raw_values` and `decode_values` usage
- [ ] Note the ID change in migration notes

## Validation

### Task 12: Final Review
- [ ] Verify all code changes are saved
- [ ] Confirm `go fmt` was run
- [ ] Check for any linter warnings
- [ ] Review all modified sections
- [ ] Ensure error handling is proper
- [ ] Confirm backward compatibility

### Task 13: Code Quality Checks
- [ ] Review for any hardcoded values
- [ ] Check for potential nil pointer dereferences
- [ ] Verify all error cases are handled
- [ ] Ensure code follows project conventions

## Cleanup
- [ ] Remove any test Terraform configurations created
- [ ] Clear any temporary test outputs
- [ ] Run final `go fmt` before committing

## Notes
- The `decode_values` field provides a convenience method for users to access decoded addon configurations
- Using `cluster_id` as the resource ID is more intuitive than a hash-based ID
- Both changes are backward compatible - existing configurations will continue to work
- Error handling ensures the data source remains functional even with invalid base64 data
- After upgrade, users may need to run `terraform refresh` once to update the state ID

## Verification Checklist
- [ ] `decode_values` field is added to schema
- [ ] `encoding/base64` import is present
- [ ] Base64 decoding logic is implemented with error handling
- [ ] Resource ID uses `cluster_id` instead of hash
- [ ] Code is formatted with `go fmt`
- [ ] Compilation succeeds without errors
- [ ] Manual testing confirms functionality
- [ ] Documentation is updated (if applicable)
