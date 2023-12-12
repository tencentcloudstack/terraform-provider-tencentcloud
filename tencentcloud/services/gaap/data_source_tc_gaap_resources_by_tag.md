Use this data source to query detailed information of gaap resources by tag

Example Usage

```hcl
data "tencentcloud_gaap_resources_by_tag" "resources_by_tag" {
  tag_key = "tagKey"
  tag_value = "tagValue"
}
```