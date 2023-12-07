Use this data source to query detailed information of ckafka topic_flow_ranking

Example Usage

```hcl
data "tencentcloud_ckafka_topic_flow_ranking" "topic_flow_ranking" {
  instance_id = "ckafka-xxxxxx"
  ranking_type = "PRO"
  begin_date = "2023-05-29T00:00:00+08:00"
  end_date = "2021-05-29T23:59:59+08:00"
}
```