## Context

Terraform Provider for TencentCloud currently has no datasource for DB Custom Cluster (dbdc product). The dbdc SDK package (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`) is available in the vendor directory, and the `DescribeDBCustomClusters` API supports querying DB Custom cluster lists with filtering by cluster IDs, filter conditions, and tags.

The provider follows the standard pattern where datasources are organized under `tencentcloud/services/<service>/` directories. The dbdc service directory does not exist yet and needs to be created. The datasource will follow the same code pattern as existing list-type datasources like `tencentcloud_igtm_instance_list`.

## Goals / Non-Goals

**Goals:**
- Add `tencentcloud_dbdc_db_custom_clusters` datasource that queries DB Custom cluster list via `DescribeDBCustomClusters` API
- Support filtering by `cluster_ids`, `filters` (cluster-name, cluster-status), and `tags`
- Return full cluster details including cluster_id, cluster_name, region, cluster_level, cluster_status, cluster_version, cluster_node_num, cluster_description, created_time, and tags
- Create the dbdc service layer with `DescribeDBCustomClustersByFilter` method
- Register the datasource in `provider.go` and `provider.md`
- Add documentation and unit tests

**Non-Goals:**
- Not creating a full resource (CRUD) for DB Custom clusters — only a datasource (Read) is needed
- Not exposing pagination parameters (limit/offset) to users — internal auto-pagination will fetch all results
- Not adding support for other dbdc APIs (e.g., DescribeDBCustomNodes, DescribeDBCustomImages)

## Decisions

### Decision 1: Schema design — flatten cluster_set fields at top level

**Choice**: Flatten the DBCustomCluster fields directly into the `cluster_set` TypeList element schema, following the project convention for datasource list-type responses.

**Rationale**: The project rule states that in Describe interfaces (returning resource lists), parameters should be expanded (flattened) — no extra nesting layer like `xxx_set`/`xxx_list` wrapping all fields again. Each element in `cluster_set` contains the flat fields from `DBCustomCluster`.

**Alternatives considered**:
- Nesting all fields under a sub-object — rejected per project rule #13
- Only returning a subset of fields — rejected; users need full cluster info for resource planning

### Decision 2: Service layer method — DescribeDBCustomClustersByFilter

**Choice**: Create a service method `DescribeDBCustomClustersByFilter` in a new `service_tencentcloud_dbdc.go` file that accepts a parameter map, builds the API request, handles pagination internally, and returns the full cluster list.

**Rationale**: Following the established pattern (e.g., `IgtmService.DescribeIgtmInstanceListByFilter`), the service layer encapsulates API calls and pagination. The datasource Read function calls this service method inside a `resource.Retry` block.

**Alternatives considered**:
- Calling SDK API directly in the datasource Read — rejected; service layer pattern provides better separation of concerns and reusability
- Using an existing service method — rejected; no dbdc service exists yet

### Decision 3: Auto-pagination — internal limit/offset handling

**Choice**: The service method will internally handle pagination by setting `Limit` to the API maximum (100) and incrementing `Offset` until all results are fetched. The `Limit` and `Offset` parameters are NOT exposed in the Terraform schema.

**Rationale**: Per project convention, datasources should not expose limit/offset to users. All data should be fetched internally via auto-pagination.

### Decision 4: Retry and error handling pattern

**Choice**: Use `resource.Retry` with `tccommon.ReadRetryTimeout` in the datasource Read function, wrapping the service method call. In the retry block, if the API returns empty results (response is nil, Response is nil, or ClusterSet length is 0), return `resource.NonRetryableError` instead of `d.SetId("")`.

**Rationale**: Per project rule #14 for RESOURCE_KIND_DATASOURCE, empty API responses should NOT clear the state id. Returning `NonRetryableError` ensures the retry mechanism properly handles temporary API fluctuations without data loss. A `log.Printf("[DATASOURCE] read empty, skip SetId")` should be kept in the retry failure path.

### Decision 5: ID generation for datasource

**Choice**: Use `helper.BuildToken()` to generate the datasource ID, consistent with other list-type datasources.

**Rationale**: List-type datasources don't have a natural single ID; `BuildToken()` is the standard pattern.

## Risks / Trade-offs

- **[Risk] API rate limiting** → Mitigation: Using `tccommon.ReadRetryTimeout` for retry; Limit set to 100 (API maximum) to minimize pagination rounds
- **[Risk] Large result sets could cause slow Terraform plans** → Mitigation: Users can filter by cluster_ids/filters/tags to narrow results; auto-pagination still fetches all matching results
- **[Risk] Nil pointer fields in DBCustomCluster.Tags** → Mitigation: Tags field in API response is noted as "may return null"; nil checks are applied before calling set methods per project rule #8
- **[Trade-off] Not exposing pagination params** → Users cannot limit the number of results for very large accounts, but this follows the established project convention
