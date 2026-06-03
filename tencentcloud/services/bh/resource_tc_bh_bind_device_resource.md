Provides a resource to bind devices to a BH (Bastion Host) service instance.

Example Usage

```hcl
resource "tencentcloud_bh_bind_device_resource" "example" {
  device_id_set = [123, 456]
  resource_id   = "bh-saas-abc123"
  domain_id     = "dm-domain01"
}
```

K8S cluster managed scenario

```hcl
resource "tencentcloud_bh_bind_device_resource" "k8s_example" {
  device_id_set    = [789]
  resource_id      = "bh-saas-abc123"
  manage_dimension = 1
  manage_account   = "admin"
  namespace        = "default"
  workload         = "deployment/nginx"
}
```

Import

BH bind device resource can be imported using the composite ID `device_ids_comma_separated#resource_id`, e.g.

```
terraform import tencentcloud_bh_bind_device_resource.example 123,456#bh-saas-abc123
```
