Provides a resource to bind a kubeconfig credential to an existing BH (Bastion Host) container account.

~> **NOTE:** The HCL field `account_id` corresponds to the SDK request field `Id` (the container account Id). It is renamed in HCL because `id` is reserved by the Terraform Plugin SDK as the resource's internal identifier.

~> **NOTE — Read is a no-op:** The BH service does not currently provide a query API for the kubeconfig binding. Read does not call any SDK; state is authoritative. External drift (e.g. credential rotation in the web console) is invisible to Terraform — re-apply will overwrite the backend with the HCL value.

~> **NOTE — Delete is a no-op:** The BH service does not provide an unbind API. `terraform destroy` removes the resource from local state but does NOT remove the kubeconfig binding on the backend. Manage with care.

Example Usage

```hcl
resource "tencentcloud_bh_bind_device_account_kubeconfig" "example" {
  account_id       = 12345
  kubeconfig       = file("${path.module}/kubeconfig.yaml")
  manage_dimension = 1
}
```
