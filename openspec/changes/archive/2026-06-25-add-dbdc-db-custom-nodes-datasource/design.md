## Context

The Terraform Provider for TencentCloud currently supports the dbdc (Database Dedicated Cluster) product but lacks a datasource for querying DB Custom nodes. Users who manage DB Custom nodes need a way to discover and reference existing nodes in their Terraform configurations without hardcoding values. The `DescribeDBCustomNodes` cloud API (dbdc v20201029) provides the capability to query nodes by IDs, filters, or tags, and returns detailed node attributes including CPU, memory, disk info, network info, status, and billing info.

The existing provider pattern for datasource resources follows a consistent structure: schema definition in `data_source_tc_<service>_<name>.go`, service layer method in `service_tencentcloud_<service>.go`, registration in `provider.go`, and documentation in `.md` files. The reference implementation `tencentcloud_igtm_instance_list` demonstrates the standard datasource pattern with filters and list output.

## Goals / Non-Goals

**Goals:**
- Add a RESOURCE_KIND_DATASOURCE `tencentcloud_dbdc_db_custom_nodes` that queries DB Custom nodes via `DescribeDBCustomNodes` API
- Support optional filter inputs: `node_ids`, `filters`, and `tags`
- Expose full node attribute details in `node_set` computed output (flattened, no nested wrapper)
- Implement internal pagination to fetch all results without exposing `limit`/`offset` to users
- Follow the established datasource pattern (retry, error handling, nil checks)

**Non-Goals:**
- Creating CRUD resources for DB Custom nodes (only datasource needed)
- Modifying any existing dbdc resources or datasources
- Adding write operations (create, update, delete) for DB Custom nodes

## Decisions

1. **Pagination approach**: Internal pagination using `Offset`/`Limit` parameters in the API call, with `Limit` set to 100 (max per API docs). Users will NOT see `limit`/`offset` in the schema. The service layer will loop until all results are collected.

2. **Schema flattening**: The `node_set` output will be a `TypeList` of `TypeMap` (schema.Resource) where each element's fields are flattened at the top level. No `xxx_set`/`xxx_list` wrapper around all fields. Nested structures (`SystemDisk`, `DataDisks`, `Tags`) will be represented as `TypeList`/`TypeMap` sub-blocks within each node element.

3. **Filter input format**: `filters` uses `TypeList` with `name` (Required, string) and `values` (Required, TypeSet of string) sub-fields, matching the standard Filter pattern. `tags` uses `TypeList` with `key` (Required, string) and `value` (Required, string) sub-fields. `node_ids` uses `TypeList` of string.

4. **Result ID**: Use `helper.BuildToken()` as the datasource ID (standard pattern for list datasources).

5. **Retry and error handling**: Use `tccommon.ReadRetryTimeout` with `resource.Retry()` for the Read operation. If the API returns empty results (nil response/empty NodeSet), return `NonRetryableError` instead of silently clearing the ID, per the datasource guidelines.

6. **Nil field checks**: Before calling `d.Set()` for any response field, check if the corresponding API response field is nil. Skip setting if nil.

7. **Service layer**: Add a `DescribeDbdcDbCustomNodesByFilter` method in `service_tencentcloud_dbdc.go` that handles pagination, request construction, and response parsing.

## Risks / Trade-offs

- **[API pagination limits]** The API allows max Limit=100 per request. If a user has more than 100 nodes, multiple API calls are needed internally. → Mitigation: Implement pagination loop in service layer.
- **[Nil response fields]** Some fields like `SystemDisk`, `DataDisks`, and `Tags` in `DBCustomNode` may return nil per API documentation. → Mitigation: Check nil before setting each field.
- **[No result validation]** As a datasource, there's no strong guarantee the API will return data for every query combination. → Mitigation: Return `NonRetryableError` on empty results to let retry mechanism handle transient API issues.
