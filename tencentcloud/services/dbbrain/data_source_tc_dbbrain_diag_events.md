Use this data source to query detailed information of DBbrain diag events

Example Usage

Query events only by time

```hcl
data "tencentcloud_dbbrain_diag_events" "example" {
  start_time = "2025-01-01T00:00:00+08:00"
  end_time   = "2026-12-31T00:00:00+08:00"
}
```

Or add another filters

```hcl
data "tencentcloud_dbbrain_diag_events" "example" {
  start_time = "2026-01-01T00:00:00+08:00"
  end_time   = "2026-12-31T00:00:00+08:00"
  instance_ids = [
    "crs-kpyy0txj"
  ]

  product    = "redis"
  severities = [1, 2, 3, 4, 5]
}
```
