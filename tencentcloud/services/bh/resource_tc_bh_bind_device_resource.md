Provides a resource to bind devices to a BH (Bastion Host) service instance.

~> **NOTE:** This resource must exclusive in one bh resource, do not declare additional device id resources of this device elsewhere.

Example Usage

```hcl
resource "tencentcloud_bh_bind_device_resource" "example" {
  device_id_set = [4173, 4175]
  resource_id   = "bh-saas-4ikvobas"
  domain_id     = "net-telc7g8p"
}
```

K8S cluster managed scenario

```hcl
resource "tencentcloud_bh_bind_device_resource" "example" {
  device_id_set     = [3434]
  resource_id       = "bh-saas-sk8eyhcn"
  domain_id         = "net-89sng6ha"
  manage_dimension  = 1
  manage_account_id = 3970
}
```

Import

BH bind device resource can be imported using the resource_id, e.g.

```
terraform import tencentcloud_bh_bind_device_resource.example bh-saas-4ikvobas
```
