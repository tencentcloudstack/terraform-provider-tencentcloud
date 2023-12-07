Provides a resource to create a gaap proxy group

Example Usage

```hcl
resource "tencentcloud_gaap_proxy_group" "proxy_group" {
  project_id = 0
  group_name = "tf-test-update"
  real_server_region = "Beijing"
  ip_address_version = "IPv4"
  package_type = "Thunder"
}
```

Import

gaap proxy_group can be imported using the id, e.g.

```
terraform import tencentcloud_gaap_proxy_group.proxy_group proxy_group_id
```