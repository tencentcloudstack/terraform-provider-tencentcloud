---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_cluster_fw_bypass_config"
sidebar_current: "docs-tencentcloud-resource-cfw_cluster_fw_bypass_config"
description: |-
  Provides a resource to manage CFW (Cloud Firewall) cluster firewall bypass configuration.
---

# tencentcloud_cfw_cluster_fw_bypass_config

Provides a resource to manage CFW (Cloud Firewall) cluster firewall bypass configuration.

## Example Usage

### VPC_FW type

```hcl
resource "tencentcloud_cfw_cluster_fw_bypass_config" "example" {
  fw_type = "VPC_FW"
  ccn_id  = "ccn-p3mlp0tj"
  enable  = false
}
```

### NAT_FW type

```hcl
resource "tencentcloud_cfw_cluster_fw_bypass_config" "example" {
  fw_type    = "NAT_FW"
  ccn_id     = "ccn-p3mlp0tj"
  nat_ins_id = "nat-h1i1mf4n"
  enable     = true
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String, ForceNew) CCN instance ID.
* `enable` - (Required, Bool) Bypass switch. `true` - enable Bypass (traffic bypasses firewall), `false` - disable Bypass (traffic goes through firewall).
* `fw_type` - (Required, String, ForceNew) Firewall type. `VPC_FW` - VPC firewall, `NAT_FW` - NAT firewall.
* `nat_ins_id` - (Optional, String, ForceNew) NAT firewall instance ID. Required when fw_type is `NAT_FW`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CFW cluster firewall bypass config can be imported using the composite ID.

For VPC_FW type, the format is `{fw_type}#{ccn_id}`:

```
terraform import tencentcloud_cfw_cluster_fw_bypass_config.example VPC_FW#ccn-p3mlp0tj
```

For NAT_FW type, the format is `{fw_type}#{nat_ins_id},{ccn_id}`:

```
terraform import tencentcloud_cfw_cluster_fw_bypass_config.example NAT_FW#nat-h1i1mf4n,ccn-p3mlp0tj
```

