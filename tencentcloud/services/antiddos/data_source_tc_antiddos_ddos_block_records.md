Use this data source to query AntiDDoS DDoS block/unblock records and unblock quota info.

Example Usage

```hcl
data "tencentcloud_antiddos_ddos_block_records" "example" {
  start_time = "2026-02-04T11:30:00+08:00"
  end_time   = "2026-03-04T11:30:00+08:00"

  filters {
    name   = "Status"
    values = ["Blocked"]
  }
}
```