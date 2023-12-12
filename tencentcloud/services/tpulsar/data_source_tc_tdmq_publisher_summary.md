Use this data source to query detailed information of tdmq publisher_summary

Example Usage

```hcl
data "tencentcloud_tdmq_publisher_summary" "publisher_summary" {
  cluster_id = "pulsar-9n95ax58b9vn"
  namespace  = "keep-ns"
  topic      = "keep-topic"
}
```