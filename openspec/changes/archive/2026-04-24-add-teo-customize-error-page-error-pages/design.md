## Context

The `tencentcloud_teo_customize_error_page` resource manages TEO custom error pages. The resource currently exposes individual fields from a single `CustomErrorPage` object (`page_id`, `name`, `content_type`, `description`, `content`), but does not expose the `error_pages` computed attribute from the `DescribeCustomErrorPages` API response. The `CustomErrorPage` struct in the SDK also contains a `References` field (of type `[]*ErrorPageReference`) that indicates which business rules reference the error page, which is not available in the current Terraform schema.

The `DescribeCustomErrorPages` API accepts `ZoneId` and `Filters` as input parameters and returns `ErrorPages` (a list of `CustomErrorPage` objects) in the response. The existing service method `DescribeTeoCustomizeErrorPageById` already uses this API with `Filters` internally to query by `page-id`.

## Goals / Non-Goals

**Goals:**
- Add the `error_pages` computed attribute to the `tencentcloud_teo_customize_error_page` resource schema
- Expose the `References` sub-field from each `CustomErrorPage` in the `error_pages` attribute
- Maintain backward compatibility with existing Terraform configurations and state

**Non-Goals:**
- Adding the `filters` parameter as a user-facing schema attribute (it is used internally in the service layer)
- Changing the existing schema fields or their behavior
- Creating a new data source for custom error pages

## Decisions

1. **Add `error_pages` as a computed TypeList attribute**: The `error_pages` field will be a `TypeList` of `schema.TypeMap` (or nested object blocks) containing all fields from `CustomErrorPage` including the `References` field. This follows the existing pattern in the codebase for exposing API response lists.

2. **Flatten `References` as a list of strings**: The `ErrorPageReference` struct only contains a single `BusinessId` field. Rather than creating a deeply nested schema, the `references` sub-field within each error page will be a `TypeList` of `TypeString` containing the `BusinessId` values.

3. **Populate `error_pages` in the read function**: The `resourceTencentCloudTeoCustomizeErrorPageRead` function already calls `DescribeTeoCustomizeErrorPageById` which uses the `DescribeCustomErrorPages` API. The response data is already available; we just need to flatten the `References` field and set the `error_pages` attribute.

4. **Schema structure for `error_pages`**: Use nested `schema.Resource` blocks for the list elements to maintain type safety and follow existing patterns in the codebase (similar to how other resources handle list-of-object computed attributes).

## Risks / Trade-offs

- **[Risk] Adding a computed attribute may increase state size** → Mitigation: The `error_pages` list typically contains only one element for a single resource read, so the impact is minimal.
- **[Risk] The `References` field may change in future API versions** → Mitigation: Use computed-only fields so any changes are reflected on the next read without requiring schema migration.
