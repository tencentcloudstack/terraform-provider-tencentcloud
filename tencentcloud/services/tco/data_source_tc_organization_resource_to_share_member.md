Use this data source to query detailed information of Organization resource to share member

Example Usage

```hcl
data "tencentcloud_organization_resource_to_share_member" "example" {
  area                 = "ap-guangzhou"
  search_key           = "tf-example"
  type                 = "CVM"
  product_resource_ids = ["ins-69hg2ze0", "ins-0cxjwrog"]
}
```
