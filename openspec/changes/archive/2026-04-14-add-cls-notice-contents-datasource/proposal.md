# Add tencentcloud_cls_notice_contents Data Source

## What

Add a new Terraform data source `tencentcloud_cls_notice_contents` to query CLS notification content template lists via `DescribeNoticeContents`.

## Why

The existing `tencentcloud_cls_notice_content` resource handles CRUD for individual templates, but there is no data source to list and filter templates. Users need to enumerate notice content templates for referencing in alarm policies.

## APIs Used

| Data Source | API | Pagination |
|---|---|---|
| `tencentcloud_cls_notice_contents` | `DescribeNoticeContents` | Limit(max=100) + Offset + TotalCount |
