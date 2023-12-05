Provides a resource to create a ckafka datahub_topic

Example Usage

```hcl
data "tencentcloud_user_info" "user" {}

resource "tencentcloud_ckafka_datahub_topic" "datahub_topic" {
  name = format("%s-tf", data.tencentcloud_user_info.user.app_id)
  partition_num = 20
  retention_ms = 60000
  note = "for test"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

ckafka datahub_topic can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_datahub_topic.datahub_topic datahub_topic_name
```