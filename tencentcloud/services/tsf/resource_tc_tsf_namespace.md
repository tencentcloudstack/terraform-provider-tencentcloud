Provides a resource to create a tsf namespace

Example Usage

```hcl
resource "tencentcloud_tsf_namespace" "namespace" {
  namespace_name = "namespace-name"
  # cluster_id = "cls-xxxx"
  namespace_desc = "namespace desc"
  # namespace_resource_type = ""
  namespace_type = "DEF"
  # namespace_id = ""
  is_ha_enable = "0"
  # program_id = ""
}
```