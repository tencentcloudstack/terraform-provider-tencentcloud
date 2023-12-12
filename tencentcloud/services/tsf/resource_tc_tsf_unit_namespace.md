Provides a resource to create a tsf unit_namespace

Example Usage

```hcl
resource "tencentcloud_tsf_unit_namespace" "unit_namespace" {
  gateway_instance_id = "gw-ins-lvdypq5k"
  namespace_id = "namespace-vwgo38wy"
  namespace_name = "keep-terraform-cls"
}
```

Import

tsf unit_namespace can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_unit_namespace.unit_namespace gw-ins-lvdypq5k#namespace-vwgo38wy
```