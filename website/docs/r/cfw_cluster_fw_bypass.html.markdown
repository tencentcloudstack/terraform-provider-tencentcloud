---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_cluster_fw_bypass"
sidebar_current: "docs-tencentcloud-resource-cfw_cluster_fw_bypass"
description: |-
  Provides a resource to manage CFW cluster firewall bypass configuration.
---

# tencentcloud_cfw_cluster_fw_bypass

Provides a resource to manage CFW cluster firewall bypass configuration.

## Example Usage

### VPC_FW type

```hcl
resource "tencentcloud_cfw_cluster_fw_bypass" "vpc_fw_example" {
  fw_type = "VPC_FW"
  ccn_id  = "ccn-xxxxxxxx"
  enable  = false
}
```

### NAT_FW type

```hcl
resource "tencentcloud_cfw_cluster_fw_bypass" "nat_fw_example" {
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
* `filters` - (Optional, List) Filter condition list. Supports filtering by Common (general search), InsObj (instance ID), ObjName (instance name), etc.
* `nat_ins_id` - (Optional, String, ForceNew) NAT firewall instance ID. Required when fw_type is `NAT_FW`.
* `nat_type` - (Optional, String) NAT firewall type filter. `nat` - VPC internal protection type, `nat_ccn` - CCN cluster mode type. If not specified, both types are queried.

The `filters` object supports the following:

* `name` - (Required, String) Filter key.
* `operator_type` - (Required, Int) Operator type. 1: equal, 2: greater than, 3: less than, 4: greater than or equal, 5: less than or equal, 6: not equal, 8: not in, 9: fuzzy match.
* `values` - (Required, List) Filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `data` - NAT firewall switch detail list.
  * `asset_type` - Asset type. `nat` - VPC internal protection, `nat_ccn` - CCN cluster mode.
  * `attach_id` - Associated ID. For nat_ccn asset type, this is the CCN instance ID.
  * `attach_name` - Associated ID instance name. For nat_ccn asset type, this is the CCN name.
  * `bypass` - Bypass status. 0: disabled, 1: enabled.
  * `fw_type` - Firewall type.
  * `ins_obj` - NAT instance ID.
  * `obj_name` - Instance name.
  * `region` - Region.
  * `routing_mode` - Traffic routing method. 0: multi-route table, 1: policy routing.
  * `status` - Switch status. 0: not enabled, 1: enabled, 2: enabling, 3: disabling.
  * `switch_mode` - Switch access mode. 1: automatic, 2: manual.
* `region_list` - Region list.
  * `text` - Display name.
  * `value` - Actual value.
* `total` - Total number of records matching the conditions.


## Import

CFW cluster firewall bypass config can be imported using the composite ID.

For VPC_FW type, the format is `{fw_type}#{ccn_id}`:

```
terraform import tencentcloud_cfw_cluster_fw_bypass.vpc_fw_example VPC_FW#ccn-xxxxxxxx
```

For NAT_FW type, the format is `{fw_type}#{ccn_id}#{nat_ins_id}`:

```
terraform import tencentcloud_cfw_cluster_fw_bypass.nat_fw_example NAT_FW#ccn-xxxxxxxx#nat-xxxxxxxx
```

