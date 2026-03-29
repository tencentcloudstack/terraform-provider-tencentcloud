# Proposal: Add TotalCount Parameter to tencentcloud_teo_l7_acc_rule

## Why

The `DescribeL7AccRules` API returns a `TotalCount` field that indicates the total number of records, which is essential for pagination scenarios and data validation. Currently, this field is not exposed in the Terraform provider for the `tencentcloud_teo_l7_acc_rule` resource and data source. By exposing this parameter, users can better understand the total number of records available and implement more robust data handling in their Terraform configurations.

## What Changes

- **New Data Source Parameter**: Add `TotalCount` output field to the `tencentcloud_teo_l7_acc_rule` data source to expose the total number of records returned by the `DescribeL7AccRules` API
- **Schema Update**: Extend the schema to include the `TotalCount` attribute as a computed (read-only) field
- **Implementation**: Update the data source read function to capture and set the `TotalCount` value from the API response

## Capabilities

### New Capabilities
None - This is extending an existing data source's output capabilities

### Modified Capabilities
- **teo_l7_acc_rule**: Existing data source capability being extended with the TotalCount field

## Impact

- **Code Changes**: Modifications to `tencentcloud_teo_l7_acc_rule` data source implementation
- **API Calls**: No changes to API calls - `DescribeL7AccRules` already returns TotalCount
- **Backward Compatibility**: Fully backward compatible - adding a computed field does not break existing configurations
- **Documentation**: Update the data source documentation to include the TotalCount field
