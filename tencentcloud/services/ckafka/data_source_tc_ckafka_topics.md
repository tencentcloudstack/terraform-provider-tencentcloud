Use this data source to query detailed information of ckafka topic.

Example Usage

```hcl
data "tencentcloud_ckafka_topics" "example" {
  instance_id = "ckafka-vv7wp5nx"
  topic_name  = "tf_example"
}
```