Use this data source to query detailed information of ckafka topic_subscribe_group

Example Usage

```hcl
data "tencentcloud_ckafka_topic_subscribe_group" "topic_subscribe_group" {
  instance_id = "ckafka-xxxxxx"
  topic_name = "xxxxxx"
}
```