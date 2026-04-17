# Design: CAM Policy Detail Data Source

## File Layout

| File | Action |
|---|---|
| `tencentcloud/services/cam/data_source_tc_cam_policy_detail.go` | New |
| `tencentcloud/services/cam/data_source_tc_cam_policy_detail.md` | New |
| `tencentcloud/services/cam/data_source_tc_cam_policy_detail_test.go` | New |
| `tencentcloud/services/cam/service_tencentcloud_cam.go` | Modified (append method) |
| `tencentcloud/provider.go` | Modified (register data source) |

## API & Pagination

`GetPolicy` takes a single required `PolicyId` (uint64) and returns one policy record. No pagination.

The service layer method is wrapped in `resource.Retry(ReadRetryTimeout, ...)`.

## Schema: tencentcloud_cam_policy_detail

**Required input:**
- `policy_id` (Int) — Policy ID.

**Computed outputs (in `policy_info` TypeList MaxItems=1):**
- `policy_name` (String)
- `description` (String)
- `type` (Int) — 1=custom, 2=preset
- `add_time` (String)
- `update_time` (String)
- `policy_document` (String)
- `preset_alias` (String)
- `is_service_linked_role_policy` (Int)
- `tags` (List) — key/value

**SetId:** string(policy_id)
