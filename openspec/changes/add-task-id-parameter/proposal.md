# Change: Add TaskId Parameter to MPS Operation Resources

## Why

The MPS (Media Processing Service) operation resources (`tencentcloud_mps_process_media_operation`, `tencentcloud_mps_edit_media_operation`, `tencentcloud_mps_process_live_stream_operation`, and others) call Tencent Cloud APIs that return a `TaskId` in the response. However, these resources do not expose this `TaskId` as a computed field in their Terraform schema.

Users need access to the task ID to:
- Track the status of media processing tasks
- Correlate Terraform operations with task status in the console or other tools
- Debug and monitor long-running media processing operations

Currently, while the resources use `TaskId` as the Terraform resource ID internally, it is not available as a readable attribute in the schema, limiting users' ability to reference task IDs in other resources or for operational visibility.

## What Changes

- Add `task_id` computed field to MPS operation resources that return TaskId from cloud API responses
- Update resource schemas to include the new field with proper documentation
- Ensure backward compatibility - all new fields are optional computed fields

### Affected Resources

The following operation resources return TaskId from their respective APIs and need the task_id field added:

1. `tencentcloud_mps_process_media_operation` - ProcessMedia API returns TaskId
2. `tencentcloud_mps_edit_media_operation` - EditMedia API returns TaskId
3. `tencentcloud_mps_process_live_stream_operation` - ProcessLiveStream API returns TaskId
4. `tencentcloud_mps_start_flow_operation` - StartFlow API returns TaskId (if applicable)
5. `tencentcloud_mps_withdraws_watermark_operation` - WithdrawsWatermark API returns TaskId (if applicable)

### New Computed Fields

For each affected resource, add:
```go
"task_id": {
    Type:        schema.TypeString,
    Computed:    true,
    Description:  "Task ID returned by the API, used to track the media processing task status.",
}
```

## Capabilities

### New Capabilities
- `mps-operation-task-id`: Adds task_id computed field to MPS operation resources to expose API-returned task IDs for user visibility and reference

### Modified Capabilities
None - this is an additive change only, no spec-level behavior changes

## Impact

### Affected Specs
- `resource-mps-process-media-operation` (spec for tencentcloud_mps_process_media_operation)
- `resource-mps-edit-media-operation` (spec for tencentcloud_mps_edit_media_operation)
- `resource-mps-process-live-stream-operation` (spec for tencentcloud_mps_process_live_stream_operation)
- Additional MPS operation resources as applicable

### Affected Code
- `tencentcloud/services/mps/resource_tc_mps_process_media_operation.go` - Add task_id to schema and set in Read function
- `tencentcloud/services/mps/resource_tc_mps_edit_media_operation.go` - Add task_id to schema and set in Read function
- `tencentcloud/services/mps/resource_tc_mps_process_live_stream_operation.go` - Add task_id to schema and set in Read function
- `tencentcloud/services/mps/resource_tc_mps_start_flow_operation.go` - Add task_id if API returns it
- `tencentcloud/services/mps/resource_tc_mps_withdraws_watermark_operation.go` - Add task_id if API returns it
- `tencentcloud/services/mps/resource_tc_mps_*.md` - Update documentation files with new field
- `website/docs/r/mps_*.html.markdown` - Update public documentation

### Breaking Changes
None - all changes are additive (new computed fields only)

### Dependencies
None - uses existing API calls that already return TaskId

### Testing Requirements
- Verify task_id is correctly set from API response in Create function
- Verify task_id is correctly read and set in Read function
- Test with each affected resource type
- Verify backward compatibility with existing configurations
