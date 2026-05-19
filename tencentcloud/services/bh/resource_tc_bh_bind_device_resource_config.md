Provides a resource to manage BH bind device resource config.

Example Usage

```hcl
resource "tencentcloud_bh_bind_device_resource_config" "example" {
  resource_id = "bh-saas-jn2p3"
  device_id_set = [5186, 5187]
  domain_id = "net-4sovwr11w7"
}
```

With K8S managed account

```hcl
resource "tencentcloud_bh_bind_device_resource_config" "example_k8s" {
  resource_id       = "bh-saas-jn2p3"
  device_id_set     = [5200]
  manage_dimension  = 2
  manage_account_id = 100
  manage_account    = "k8s-admin"
  namespace         = "default"
  workload          = "my-deployment"
}
```

Import

BH bind device resource config can be imported using the id (resource_id), e.g.

```
terraform import tencentcloud_bh_bind_device_resource_config.example bh-saas-jn2p3
```
