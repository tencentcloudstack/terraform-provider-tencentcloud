## Context

The `tencentcloud_teo_customize_error_page` resource manages TEO custom error pages. The resource currently exposes individual fields from a single `CustomErrorPage` object (`page_id`, `name`, `content_type`, `description`, `content`), but does not expose the `References` field from the API response. The `CustomErrorPage` struct in the SDK contains a `References` field (of type `[]*ErrorPageReference`) that indicates which business rules reference the error page.

The `DescribeCustomErrorPages` API accepts `ZoneId` and `Filters` as input parameters and returns `ErrorPages` (a list of `CustomErrorPage` objects) in the response. The existing service method `DescribeTeoCustomizeErrorPageById` already uses this API with `Filters` internally to query by `page-id`.

## Goals / Non-Goals

**Goals:**
- Add the `references` computed attribute at the top level of the resource schema
- Expose the `References` sub-field from `CustomErrorPage` as a flat list of business ID strings
- Maintain backward compatibility with existing Terraform configurations and state

**Non-Goals:**
- Adding an `error_pages` nested list attribute (the existing top-level fields already cover page_id, name, content_type, description, content)
- Changing the existing schema fields or their behavior
- Creating a new data source for custom error pages

## Decisions

1. **Add `references` as a top-level computed TypeList of TypeString**: The `ErrorPageReference` struct only contains a single `BusinessId` field. Rather than creating a nested object list, `references` is a flat `TypeList` of `TypeString` containing the `BusinessId` values directly. This avoids unnecessary nesting since the outer-level fields (`page_id`, `name`, `content_type`, `description`, `content`) already exist at the top level.

2. **Populate `references` in the read function**: The `resourceTencentCloudTeoCustomizeErrorPageRead` function already calls `DescribeTeoCustomizeErrorPageById` which uses the `DescribeCustomErrorPages` API. The response data is already available; we just need to flatten the `References` field and set the `references` attribute.

## Risks / Trade-offs

- **[Risk] Adding a computed attribute may affect state** → Mitigation: The `references` field is computed-only, so any changes are reflected on the next read without requiring schema migration.
- **[Risk] The `References` field may change in future API versions** → Mitigation: Use computed-only fields so any changes are reflected automatically.
