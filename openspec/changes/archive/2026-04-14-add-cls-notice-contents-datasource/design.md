# Design: tencentcloud_cls_notice_contents Data Source

## File Layout

| File | Action |
|---|---|
| `data_source_tc_cls_notice_contents.go` | New |
| `data_source_tc_cls_notice_contents.md` | New |
| `data_source_tc_cls_notice_contents_test.go` | New |
| `service_tencentcloud_cls.go` | Modified — append `DescribeClsNoticeContentsByFilter` |
| `provider.go` | Modified — register data source |

## Pagination

Limit=100 (max), Offset+=Limit until len(items) >= TotalCount. Each page call in `resource.Retry`.

## Schema

**Optional filters:** `filters` (list of name/values)

**Computed:** `notice_content_list` (List of `NoticeContentTemplate`)
- `notice_content_id`, `name`, `type` (Int), `flag` (Int), `uin` (Int), `sub_uin` (Int), `create_time` (Int), `update_time` (Int)
- `notice_contents` (List) → `type`, `trigger_content`/`recovery_content` (nested: `title`, `content`, `headers`)
