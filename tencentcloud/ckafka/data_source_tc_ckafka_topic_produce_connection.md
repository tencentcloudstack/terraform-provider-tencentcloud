Use this data source to query detailed information of ckafka topic_produce_connection

Example Usage

```hcl
data "tencentcloud_ckafka_topic_produce_connection" "topic_produce_connection" {
  instance_id = "ckafka-xxxxxx"
  topic_name = "topic-xxxxxx"
}
```