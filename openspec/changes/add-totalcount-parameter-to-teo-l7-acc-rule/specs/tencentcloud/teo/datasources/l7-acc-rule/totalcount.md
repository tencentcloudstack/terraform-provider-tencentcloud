# TotalCount Parameter Specification for tencentcloud_teo_l7_acc_rule Data Source

## Overview

The `tencentcloud_teo_l7_acc_rule` data source retrieves Layer 7 acceleration rules from Tencent Cloud EdgeOne (TEO). This specification documents the addition of the `TotalCount` parameter to the data source, which exposes the total number of rules returned by the underlying API.

## Purpose

The `TotalCount` parameter provides users with the total count of records available for the queried criteria. This enables:

- Pagination validation and UI enhancements
- Data completeness checks for automation scripts
- Better understanding of the data volume for the given query parameters
- Integration with monitoring and logging systems

## API Details

The parameter is sourced from the `DescribeL7AccRules` API operation, which returns a `TotalCount` field in its response structure.

**API Response Structure:**
```json
{
  "Response": {
    "RequestId": "string",
    "TotalCount": 1234,
    "Rules": [
      // ... rule items ...
    ]
  }
}
```

## Data Source Specification

**Name:** `tencentcloud_teo_l7_acc_rule`

**Schema Extension:**

| Attribute | Type | Computed | Description |
|-----------|------|----------|-------------|
| `total_count` | Number | Yes | Total number of Layer 7 acceleration rules matching the query criteria |

**Usage Example:**

```hcl
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-123456789"
}

output "rule_count" {
  description = "Total number of rules in the zone"
  value       = data.tencentcloud_teo_l7_acc_rule.example.total_count
}
```

## Behavioral Expectations

1. **Read-Only Property**: The `total_count` field is computed only and cannot be set by users.

2. **Availability**: The field will always be available when the data source is queried successfully, even if the result set is empty (in which case the value will be `0`).

3. **Consistency**: The value is sourced directly from the API response and reflects the API's perspective on the total count, which may differ from the actual number of items returned due to API-side pagination or filtering.

4. **Type Safety**: The value is a non-negative integer (0 or positive). Negative values should be considered invalid.

## Error Handling

- **API Failure**: If the `DescribeL7AccRules` API call fails, the data source query will fail, and no `total_count` value will be returned.
- **Null/Undefined**: The field should not be null or undefined in a successful API response. However, if the API changes in the future, the provider should handle such cases gracefully (e.g., by defaulting to `0` or omitting the field).

## Backward Compatibility

This change is fully backward compatible. Existing configurations that do not reference the `total_count` field will continue to function exactly as before. The addition of this computed field does not affect the data source's query parameters or the set of rules returned.

## Migration Notes

No migration is required for existing Terraform configurations. Users can opt to use the new field in their configurations at their convenience.

## Testing Requirements

1. **Unit Tests**: Verify that the `total_count` field is correctly extracted from the API response and exposed in the data source schema.

2. **Integration Tests**: Ensure that the field is properly set in real API responses and is accessible in Terraform configurations.

3. **Edge Cases**: Test scenarios including:
   - Empty result sets (total_count = 0)
   - Very large total counts (to ensure proper integer handling)
   - Concurrent queries to validate consistency

## Documentation Updates

The following documentation should be updated to reflect this change:

1. **Data Source Documentation**: Add the `total_count` attribute to the schema documentation for `tencentcloud_teo_l7_acc_rule`.

2. **Examples**: Include usage examples demonstrating how to reference the `total_count` field in outputs and conditional logic.

3. **Changelog**: Document the addition of this attribute in the provider's changelog.
