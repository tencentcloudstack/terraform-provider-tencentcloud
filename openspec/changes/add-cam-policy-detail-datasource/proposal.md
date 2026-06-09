# Add CAM Policy Detail Data Source

## What

Add a new Terraform data source `tencentcloud_cam_policy_detail` that retrieves the full detail of a single CAM policy by its ID using the `GetPolicy` API.

## Why

Users managing CAM policies via Terraform need to inspect a specific policy's content, type, description, and associated tags without importing it as a managed resource.

## APIs Used

| Data Source | API | Pagination |
|---|---|---|
| `tencentcloud_cam_policy_detail` | `GetPolicy` | None (single-record query by PolicyId) |
