# GWLB Target Group Port Field Drift Fix Test

This directory contains test configurations to verify the fix for the `port` field drift issue in `tencentcloud_gwlb_target_group` resource.

## Test Cases

### Test Case 1: Without Port Field
Resource: `tencentcloud_gwlb_target_group.without_port`
- Does NOT specify the `port` field
- API should return a default port value
- Terraform should accept the default without detecting drift

**Expected Behavior:**
1. `terraform apply` creates the resource successfully
2. `terraform plan` shows no changes (no drift)
3. Output `without_port_value` shows the API-provided default port

### Test Case 2: With Explicit Port Value
Resource: `tencentcloud_gwlb_target_group.with_port`
- Explicitly sets `port = 6081`
- API should use and return the specified value
- Terraform should not detect drift

**Expected Behavior:**
1. `terraform apply` creates the resource with port 6081
2. `terraform plan` shows no changes
3. Output `with_port_value` shows 6081

## How to Test

```bash
# Initialize Terraform
terraform init

# Plan and verify no errors
terraform plan

# Apply the configuration
terraform apply

# Verify no drift - should show "No changes"
terraform plan

# Clean up
terraform destroy
```

## Success Criteria

✅ Both resources create successfully  
✅ No drift detected after initial apply  
✅ Port values are correctly reflected in outputs  
✅ Second `terraform plan` shows "No changes. Your infrastructure matches the configuration."

## Fix Details

**Problem**: The `port` field caused drift because:
- API returns a default value even when not specified by user
- Field lacked `Computed: true` attribute
- Terraform treated API default as configuration drift

**Solution**: Added `Computed: true` attribute to the schema definition:
```go
"port": {
    Type:        schema.TypeInt,
    Optional:    true,
    Computed:    true,  // Added this line
    Description: "Default port of the target group...",
},
```

This allows Terraform to accept API-provided defaults while still supporting user-specified values.
