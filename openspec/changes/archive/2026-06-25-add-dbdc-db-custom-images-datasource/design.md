## Context

The TencentCloud Terraform Provider already supports DBDC (Database Dedicated Cluster) via the `tencentcloud_dbdc_db_custom_clusters` data source, which queries `DescribeDBCustomClusters`. Users creating DB Custom clusters or adding nodes need to discover valid OS image IDs via `DescribeDBCustomImages`, which is currently not exposed as a Terraform data source.

The `DescribeDBCustomImages` API is a simple list query with only pagination parameters (Offset, Limit). It returns `ImageSet` containing `DBCustomImage` objects with four fields: `ImageId`, `OsName`, `ImageType`, `Architecture`. This is a straightforward read-only data source with no filtering beyond pagination.

## Goals / Non-Goals

**Goals:**
- Add `tencentcloud_dbdc_db_custom_images` data source enabling Terraform users to query available DB Custom OS images
- Follow existing dbdc data source patterns (matching `tencentcloud_dbdc_db_custom_clusters`)
- Flatten the `ImageSet` list into top-level computed fields per the project's requirement that datasource list results should be expanded, not nested under a wrapper list
- Implement automatic pagination in the service layer (not exposing Offset/Limit to users)
- Provide gomonkey-based unit tests for the data source

**Non-Goals:**
- Adding image selection or filtering capabilities beyond what the cloud API provides (DescribeDBCustomImages has no filter parameters)
- Modifying the existing `tencentcloud_dbdc_db_custom_clusters` data source
- Adding CRUD resources for DB Custom images (images are read-only metadata)

## Decisions

### 1. Schema Structure: Flatten ImageSet fields at top level
**Decision**: The `image_set` result field uses `TypeList` with `Elem: &schema.Resource{}` containing `image_id`, `os_name`, `image_type`, `architecture` as individual `Computed: true` fields.
**Rationale**: Following the project rule that datasource list results must be expanded/flattened, not wrapped under a redundant nesting layer. Matches the pattern used by `tencentcloud_dbdc_db_custom_clusters` (`cluster_set` with flattened cluster fields).

### 2. No filter/input parameters for users
**Decision**: The data source has no input parameters besides `result_output_file`. The `DescribeDBCustomImages` API only accepts Offset and Limit for pagination, with no filtering capability.
**Rationale**: Since the cloud API has no filter parameters, exposing Offset/Limit to users would violate the project's rule that datasource pagination should be handled internally. The data source simply returns all available images.

### 3. Service layer pagination
**Decision**: The service layer method `DescribeDBCustomImagesByFilter` will implement internal pagination with Limit=100 (the API max) and accumulate all results across pages.
**Rationale**: Consistent with existing dbdc service pattern (`DescribeDBCustomClustersByFilter` uses Limit=100 pagination). The API documentation states Limit range is [1, 100].

### 4. Unit test approach: gomonkey mock
**Decision**: Use gomonkey to mock the `DescribeDBCustomImages` client method, testing the Read function logic directly.
**Rationale**: Per project rules, new Terraform resources use gomonkey-based mock tests rather than Terraform acceptance test suite. This follows the pattern established in the existing `data_source_tc_dbdc_db_custom_clusters_test.go`.

### 5. Error handling for empty response
**Decision**: In the service layer, if `DescribeDBCustomImages` returns nil/empty response, return `NonRetryableError` rather than silently clearing the data source ID.
**Rationale**: Per project rules for RESOURCE_KIND_DATASOURCE, returning empty should not directly `d.SetId("")` but should return `NonRetryableError` to let the retry mechanism handle transient API failures, with a `log.Printf("[DATASOURCE] read empty, skip SetId")` message for troubleshooting.

## Risks / Trade-offs

- **[API has no filtering]** → The data source returns all available images with no way to narrow results. Users must filter in Terraform logic. This is a limitation of the cloud API, not the Terraform implementation.
- **[Limited image metadata]** → DBCustomImage only has 4 fields (ImageId, OsName, ImageType, Architecture). If the cloud API adds fields later, the schema will need updating. Mitigation: Follow standard pattern of nil-checking all fields before setting, making future additions straightforward.
