## Context

Terraform Provider for TencentCloud currently has no datasource for DB Custom regions (dbdc product). The dbdc SDK package (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`) is available in the vendor directory, and the `DescribeDBCustomRegions` API supports querying the DB Custom supported region list. Unlike other dbdc describe APIs (e.g., `DescribeDBCustomClusters`, `DescribeDBCustomImages`), `DescribeDBCustomRegions` has **no input parameters** and its response has **no pagination fields** (no `Limit`, `Offset`, `TotalCount`) — it returns the full `RegionSet` list in a single call.

The provider already has a dbdc service directory (`tencentcloud/services/dbdc/`) with an existing `service_tencentcloud_dbdc.go` file and multiple datasources (`tencentcloud_dbdc_db_custom_clusters`, `tencentcloud_dbdc_db_custom_images`, etc.). The new datasource will follow the same code pattern as the existing `tencentcloud_dbdc_db_custom_images` datasource (which is the closest analog: a simple no-filter query datasource).

## Goals / Non-Goals

**Goals:**
- Add `tencentcloud_dbdc_db_custom_regions` datasource that queries DB Custom supported region list via `DescribeDBCustomRegions` API
- Return region details including `region` (region name) and `region_state` (sale status: SELL / SOLD_OUT)
- Add the `DescribeDBCustomRegions` service method to the existing `service_tencentcloud_dbdc.go` file
- Register the datasource in `provider.go` and `provider.md`
- Add documentation and unit tests

**Non-Goals:**
- Not creating a full resource (CRUD) for DB Custom regions — only a datasource (Read) is needed
- Not implementing pagination logic — the `DescribeDBCustomRegions` API has no pagination fields and returns all regions in one call
- Not exposing any filter parameters to users — the API has no input parameters
- Not adding support for other dbdc APIs

## Decisions

### Decision 1: Schema design — flatten region_set fields at top level

**Choice**: Flatten the `RegionInfo` fields directly into the `region_set` TypeList element schema, following the project convention for datasource list-type responses.

**Rationale**: The project rule states that in Describe interfaces (returning resource lists), parameters should be expanded (flattened) — no extra nesting layer like `xxx_set`/`xxx_list` wrapping all fields again. Each element in `region_set` contains the flat fields from `RegionInfo`: `region` and `region_state`.

**Alternatives considered**:
- Nesting all fields under a sub-object — rejected per project rule #13
- Only returning a subset of fields — rejected; users need both region name and sale status

### Decision 2: Service layer method — DescribeDBCustomRegions (no pagination)

**Choice**: Add a service method `DescribeDBCustomRegions` to the existing `service_tencentcloud_dbdc.go` file that builds the API request (no parameters), calls the API inside a `resource.Retry(tccommon.ReadRetryTimeout, ...)` block, and returns the `RegionSet` list directly. **No pagination loop** is needed because the API response has no `Limit`/`Offset`/`TotalCount` fields.

**Rationale**: Unlike `DescribeDBCustomClustersByFilter` and `DescribeDBCustomImagesByFilter` which support pagination via `Limit`/`Offset` parameters and `TotalCount` response, the `DescribeDBCustomRegions` API has no input parameters and no pagination fields. The full region list is returned in a single API call. The service method still wraps the API call in `resource.Retry` to handle transient failures, consistent with the established service layer pattern.

**Alternatives considered**:
- Implementing pagination like other dbdc service methods — rejected; the API has no pagination support and no pagination fields in the response
- Calling SDK API directly in the datasource Read — rejected; service layer pattern provides better separation of concerns and reusability

### Decision 3: Retry and error handling pattern

**Choice**: Use `resource.Retry` with `tccommon.ReadRetryTimeout` in both the service layer (wrapping the API call) and the datasource Read function (wrapping the service method call). In the service layer retry block, if the API returns empty results (response is nil, Response is nil, or RegionSet is nil), return `resource.NonRetryableError` instead of silently clearing the state.

**Rationale**: Per project rule #14 for RESOURCE_KIND_DATASOURCE, empty API responses should NOT clear the state id. Returning `NonRetryableError` ensures the retry mechanism properly handles temporary API fluctuations without data loss. A `log.Printf("[DATASOURCE] read empty, skip SetId")` should be kept in the service layer retry failure path.

**Note on nested retry**: Per project rule, do not nest retry inside retry. Since the service layer method already wraps the API call in `resource.Retry`, the datasource Read function will call the service method inside its own `resource.Retry`. To avoid double-retry, the service method will perform the retry internally and return the final result/error, while the datasource Read's `resource.Retry` will simply propagate any error from the service method (the service method's internal retry is the authoritative one for API errors). This mirrors the existing `DescribeDBCustomImagesByFilter` pattern where the service method has its own retry and the Read function also has a retry wrapper.

### Decision 4: ID generation for datasource

**Choice**: Use `helper.BuildToken()` to generate the datasource ID, consistent with other list-type datasources.

**Rationale**: List-type datasources don't have a natural single ID; `BuildToken()` is the standard pattern.

## Risks / Trade-offs

- **[Risk] API rate limiting** → Mitigation: Using `tccommon.ReadRetryTimeout` for retry in the service layer
- **[Risk] Nil pointer fields in RegionInfo** → Mitigation: `RegionSet` field in API response is noted as "may return null"; nil checks are applied before calling set methods per project rule #8
- **[Trade-off] No pagination support** → The `DescribeDBCustomRegions` API returns all regions in one call with no pagination; this is acceptable because the region list is small and bounded
- **[Trade-off] No filter parameters** → Users cannot filter the region list, but the API itself has no input parameters, so this is a cloud API limitation, not a provider limitation
