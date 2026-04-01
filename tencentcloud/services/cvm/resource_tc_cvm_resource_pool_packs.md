Provides a resource to create CVM resource pool packs for advance compute resource pool management.

Resource pool packs allow users to purchase and manage compute resource pools in advance. This resource supports standard Terraform operations: create, read, and delete. Update operations are not available from CVM API, all field changes will trigger a destroy and recreate.

Example Usage

Basic Usage

```hcl
resource "tencentcloud_cvm_resource_pool_packs" "pool_packs" {
  zone                              = "ap-guangzhou-7"
  instance_type                     = "SA9.96XLARGE1152"
  period                            = 12
  resource_pool_pack_type           = "EXCLUSIVE"
  auto_placement                    = true
  dedicated_resource_pool_pack_name = "my-resource-pool"
  renew_flag                        = "NOTIFY_AND_MANUAL_RENEW"
}
```

With Auto Renewal

```hcl
resource "tencentcloud_cvm_resource_pool_packs" "pool_packs_auto_renew" {
  zone                              = "ap-guangzhou-7"
  instance_type                     = "SA9.96XLARGE1152"
  period                            = 6
  resource_pool_pack_type           = "EXCLUSIVE"
  auto_placement                    = true
  dedicated_resource_pool_pack_name = "auto-renew-pool"
  renew_flag                        = "NOTIFY_AND_AUTO_RENEW"
}
```

Import

CVM resource pool packs can be imported using the resource pool pack ID, e.g.

```
terraform import tencentcloud_cvm_resource_pool_packs.pool_packs rpp-xxxxx
```
