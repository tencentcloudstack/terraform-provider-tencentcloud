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
  device_id_set = [123, 456]
  resource_id   = "bh-saas-abc123"
  domain_id     = "dm-domain01"
}
```

### K8S cluster managed scenario

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

## Argument Reference

The following arguments are supported:

* `device_id_set` - (Required, List: [`Int`], ForceNew) Device ID set.
* `resource_id` - (Required, String) Bindable bastion host service ID.
* `domain_id` - (Optional, String) Network domain ID.
* `manage_account_id` - (Optional, Int) K8S cluster managed account ID.
* `manage_account` - (Optional, String) K8S cluster managed account name.
* `manage_dimension` - (Optional, Int) K8S cluster managed account dimension. 1-cluster, 2-namespace, 3-workload.
* `manage_kubeconfig` - (Optional, String) K8S cluster managed account kubeconfig credential.
* `namespace` - (Optional, String) K8S cluster managed namespace.
* `workload` - (Optional, String) K8S cluster managed workload.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

BH bind device resource can be imported using the composite ID `device_ids_comma_separated#resource_id`, e.g.

```
terraform import tencentcloud_bh_bind_device_resource.example 123,456#bh-saas-abc123
```

