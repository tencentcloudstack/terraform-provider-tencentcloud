Provides a resource to create a tsf bind_api_group

Example Usage

```hcl
resource "tencentcloud_tsf_bind_api_group" "bind_api_group" {
  gateway_deploy_group_id = "group-vzd97zpy"
  group_id = "grp-qp0rj3zi"
}
```

Import

tsf bind_api_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_bind_api_group.bind_api_group bind_api_group_id
```