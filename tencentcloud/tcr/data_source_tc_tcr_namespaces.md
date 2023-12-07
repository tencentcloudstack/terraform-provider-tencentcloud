Use this data source to query detailed information of TCR namespaces.

Example Usage

```hcl
data "tencentcloud_tcr_namespaces" "name" {
  instance_id 			= "cls-satg5125"
  namespace_name       = "test"
}
```