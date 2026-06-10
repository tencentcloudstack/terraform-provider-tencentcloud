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
resource "tencentcloud_cfw_cluster_fw_bypass_config" "vpc_fw_example" {
  fw_type = "VPC_FW"
  ccn_id  = "ccn-xxxxxxxx"
  enable  = false
}
```

### NAT_FW type

```hcl
resource "tencentcloud_cfw_cluster_fw_bypass_config" "nat_fw_example" {
  fw_type    = "NAT_FW"
  ccn_id     = "ccn-xxxxxxxx"
  nat_ins_id = "nat-xxxxxxxx"
  enable     = false
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
terraform import tencentcloud_cfw_cluster_fw_bypass_config.vpc_fw_example VPC_FW#ccn-xxxxxxxx
```

For NAT_FW type, the format is `{fw_type}#{ccn_id}#{nat_ins_id}`:

```
terraform import tencentcloud_cfw_cluster_fw_bypass_config.nat_fw_example NAT_FW#ccn-xxxxxxxx#nat-xxxxxxxx
```

