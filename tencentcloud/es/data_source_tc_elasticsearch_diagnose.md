Use this data source to query detailed information of elasticsearch diagnose

Example Usage

```hcl
data "tencentcloud_elasticsearch_diagnose" "diagnose" {
  instance_id = "es-xxxxxx"
  date = "20231030"
  limit = 1
}
```