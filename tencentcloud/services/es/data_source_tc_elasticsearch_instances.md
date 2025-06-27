Use this data source to query Elasticsearch(ES) instances.

Example Usage

Query ES instances by filters

```hcl
data "tencentcloud_elasticsearch_instances" "example" {
  instance_id   = "es-bxffils7"
  instance_name = "tf-example"
  tags = {
    createBy = "Terraform"
  }
}
```
