Provides a resource to create Redis instance log delivery land set its attributes.

~> **NOTE:** When you use an existing cls logset and topic to enable logging, there is no need to set parameters such
as `period`, `create_index`, `log_region`, etc.

Example Usage

Use cls logset and topic which existed

```hcl
resource "tencentcloud_redis_log_delivery" "delivery" {
  instance_id = "crs-dmjj8en7"
  logset_id   = "cc31d9d6-74c0-4888-8b2f-b8148c3bcc5c"
  topic_id    = "5c2333e9-0bab-41fd-9f75-c602b3f9545f"
}
```

Use exist cls logset and create new topic

```hcl
resource "tencentcloud_redis_log_delivery" "delivery" {
  instance_id = "crs-dmjj8en7"
  logset_id   = "cc31d9d6-74c0-4888-8b2f-b8148c3bcc5c"
  topic_name   = "test13"
  period       = 20
  create_index = true
}
```

Create new cls logset and topic

```hcl
resource "tencentcloud_redis_log_delivery" "delivery" {
  instance_id  = "crs-dmjj8en7"
  log_region   = "ap-guangzhou"
  logset_name  = "test"
  topic_name   = "test"
  period       = 20
  create_index = true
}
```

Import

Redis log delivery can be imported, e.g.

```
$ terraform import tencentcloud_redis_log_delivery.delivery crs-dmjj8en7
```