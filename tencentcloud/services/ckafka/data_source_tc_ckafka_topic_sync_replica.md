Use this data source to query detailed information of ckafka topic_sync_replica

Example Usage

```hcl
data "tencentcloud_ckafka_topic_sync_replica" "topic_sync_replica" {
  instance_id = "ckafka-xxxxxx"
  topic_name = "xxxxxx"
}
```