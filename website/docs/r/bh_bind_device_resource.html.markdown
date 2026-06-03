---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_bind_device_resource"
sidebar_current: "docs-tencentcloud-resource-bh_bind_device_resource"
description: |-
  Provides a resource to bind devices to a BH (Bastion Host) service instance.
---

# tencentcloud_bh_bind_device_resource

Provides a resource to bind devices to a BH (Bastion Host) service instance.

## Example Usage

```hcl
resource "tencentcloud_bh_bind_device_resource" "example" {
  device_id_set = [4173, 4175]
  resource_id   = "bh-saas-4ikvobas"
  domain_id     = "net-telc7g8p"
}
```

### K8S cluster managed scenario

```hcl
resource "tencentcloud_bh_bind_device_resource" "example" {
  device_id_set     = [3434]
  resource_id       = "bh-saas-sk8eyhcn"
  domain_id         = "net-89sng6ha"
  manage_dimension  = 1
  manage_account_id = 3970
}
```

## Argument Reference

The following arguments are supported:

* `device_id_set` - (Required, Set: [`Int`]) Device ID set.
* `resource_id` - (Required, String, ForceNew) Bindable bastion host service ID.
* `domain_id` - (Optional, String, ForceNew) Network domain ID.
* `manage_account_id` - (Optional, Int, ForceNew) K8S cluster managed account ID.
* `manage_account` - (Optional, String, ForceNew) K8S cluster managed account name.
* `manage_dimension` - (Optional, Int, ForceNew) K8S cluster managed account dimension. 1-cluster, 2-namespace, 3-workload.
* `manage_kubeconfig` - (Optional, String, ForceNew) K8S cluster managed account kubeconfig credential.
* `namespace` - (Optional, String, ForceNew) K8S cluster managed namespace.
* `workload` - (Optional, String, ForceNew) K8S cluster managed workload.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

BH bind device resource can be imported using the resource_id, e.g.

```
terraform import tencentcloud_bh_bind_device_resource.example bh-saas-4ikvobas
```

