Use this resource to create dayu eip rule

Example Usage

```hcl
resource "tencentcloud_dayu_eip" "test" {
  resource_id = "bgpip-000004xg"
  eip = "162.62.163.50"
  bind_resource_id = "ins-4m0jvxic"
  bind_resource_region = "hk"
  bind_resource_type = "cvm"
}
```