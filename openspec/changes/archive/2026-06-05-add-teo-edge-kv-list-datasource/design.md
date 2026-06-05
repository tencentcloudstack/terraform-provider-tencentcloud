## Context

TencentCloud EdgeOne (TEO) provides Edge KV storage for serverless edge functions. The `EdgeKVList` API allows listing key names within a namespace with optional prefix filtering and cursor-based pagination. The Terraform provider already has TEO service infrastructure in `tencentcloud/services/teo/`, and the SDK is vendored at `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`.

## Goals / Non-Goals

**Goals:**
- Provide a `tencentcloud_teo_edge_k_v_list` data source that queries KV key names from a TEO zone namespace.
- Support optional prefix filtering and cursor-based traversal.
- Follow existing data source patterns (reference: `tencentcloud_igtm_instance_list`).
- Include unit tests using gomonkey mocking.

**Non-Goals:**
- Automatic pagination (cursor-based traversal across all pages) — the user controls cursor input/output explicitly.
- CRUD operations on KV entries (this is read-only data source).
- Exposing the `Limit` parameter to users — it will be hardcoded to the maximum value (1000) internally.

## Decisions

1. **Single-page query with cursor passthrough**: The `EdgeKVList` API uses cursor-based pagination rather than offset/limit. Since the user may want to control pagination themselves, we expose `cursor` as both an Optional input and a Computed output. The internal `Limit` is set to 1000 (maximum) to return as many keys as possible per call.

2. **File placement**: Following the existing TEO service structure, the data source goes in `tencentcloud/services/teo/data_source_tc_teo_edge_k_v_list.go`. The TEO service layer (`service_tencentcloud_teo.go`) already exists and will be extended with a helper method.

3. **ID generation**: Since this is a query-only data source with no unique resource identity, use `helper.BuildToken()` to generate a random ID (consistent with other data sources like `tencentcloud_igtm_instance_list`).

4. **Retry pattern**: Use `resource.Retry` with `tccommon.ReadRetryTimeout` and `tccommon.RetryError` for transient failures, matching the standard provider pattern.

5. **Unit testing with gomonkey**: Mock the TEO SDK client's `EdgeKVListWithContext` method to test the Read logic without real API calls.

## Risks / Trade-offs

- [Risk] Cursor semantics may confuse users who expect automatic full-list retrieval → Mitigation: Document that cursor is used for pagination and provide example usage in the .md doc.
- [Risk] API returns empty cursor when all data is traversed → Mitigation: Set cursor as Computed so it's always available in state for inspection.
