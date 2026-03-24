Use this data source to query detailed information of cls alarm notices

Example Usage

Query all cls alarm notices

```hcl
data "tencentcloud_cls_alarm_notices" "example" {}
```

Query by filters

```hcl
data "tencentcloud_cls_alarm_notices" "example" {
  filters {
    key    = "name"
    values = ["tf-example"]
  }

  filters {
    key    = "alarmNoticeId"
    values = ["notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101"]
  }
}
```
