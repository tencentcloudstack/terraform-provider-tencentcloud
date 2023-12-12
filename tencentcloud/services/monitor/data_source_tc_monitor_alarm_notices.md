Use this data source to Interlude notification list.

Example Usage

```hcl
data "tencentcloud_monitor_alarm_notices" "notices" {
    order = "DESC"
    owner_uid = 1
    name = ""
    receiver_type = ""
    user_ids = []
    group_ids = []
    notice_ids = []
}
```