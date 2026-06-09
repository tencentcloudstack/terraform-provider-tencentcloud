---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_bind_device_account_kubeconfig"
sidebar_current: "docs-tencentcloud-resource-bh_bind_device_account_kubeconfig"
description: |-
  Provides a resource to bind a kubeconfig credential to an existing BH (Bastion Host) container account.
---

# tencentcloud_bh_bind_device_account_kubeconfig

Provides a resource to bind a kubeconfig credential to an existing BH (Bastion Host) container account.

~> **NOTE:** The HCL field `account_id` corresponds to the SDK request field `Id` (the container account Id). It is renamed in HCL because `id` is reserved by the Terraform Plugin SDK as the resource's internal identifier.

~> **NOTE — Read is a no-op:** The BH service does not currently provide a query API for the kubeconfig binding. Read does not call any SDK; state is authoritative. External drift (e.g. credential rotation in the web console) is invisible to Terraform — re-apply will overwrite the backend with the HCL value.

~> **NOTE — Delete is a no-op:** The BH service does not provide an unbind API. `terraform destroy` removes the resource from local state but does NOT remove the kubeconfig binding on the backend. Manage with care.

## Example Usage

```hcl
resource "tencentcloud_bh_bind_device_account_kubeconfig" "example" {
  account_id       = 12345
  kubeconfig       = file("${path.module}/kubeconfig.yaml")
  manage_dimension = 1
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required, Int, ForceNew) Container account Id. Maps to the SDK request field `Id`. Renamed in HCL because `id` is reserved by the Terraform Plugin SDK as the resource's internal identifier.
* `kubeconfig` - (Required, String) Container account kubeconfig credential.
* `manage_dimension` - (Optional, Int) Manage dimension. 1 means cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



