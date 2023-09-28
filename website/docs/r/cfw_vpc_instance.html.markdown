---
subcategory: "Cfw"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_vpc_instance"
sidebar_current: "docs-tencentcloud-resource-cfw_vpc_instance"
description: |-
  Provides a resource to create a cfw vpc_instance
---

# tencentcloud_cfw_vpc_instance

Provides a resource to create a cfw vpc_instance

## Example Usage

### If mode is 0

```hcl
resource "tencentcloud_cfw_vpc_instance" "example" {
  name = "tf_example"
  mode = 0

  vpc_fw_instances {
    name = "fw_ins_example"
    vpc_ids = [
      "vpc-291vnoeu",
      "vpc-39ixq9ci"
    ]
    fw_deploy {
      deploy_region = "ap-guangzhou"
      width         = 1024
      cross_a_zone  = 1
      zone_set = [
        "ap-guangzhou-6",
        "ap-guangzhou-7"
      ]
    }
  }

  switch_mode = 1
  fw_vpc_cidr = "auto"
}
```

### If mode is 1

```hcl
resource "tencentcloud_cfw_vpc_instance" "example" {
  name = "tf_example"
  mode = 1

  vpc_fw_instances {
    name = "fw_ins_example"
    fw_deploy {
      deploy_region = "ap-guangzhou"
      width         = 1024
      cross_a_zone  = 0
      zone_set = [
        "ap-guangzhou-6"
      ]
    }
  }

  ccn_id      = "ccn-peihfqo7"
  switch_mode = 1
  fw_vpc_cidr = "auto"
}
```

## Argument Reference

The following arguments are supported:

* `mode` - (Required, Int) Mode 0: private network mode; 1: CCN cloud networking mode.
* `name` - (Required, String) VPC firewall (group) name.
* `switch_mode` - (Required, Int) Switch mode of firewall instance. 1: Single point intercommunication; 2: Multi-point communication; 4: Custom Routing.
* `vpc_fw_instances` - (Required, List) List of firewall instances under firewall (group).
* `ccn_id` - (Optional, String) Cloud networking id, suitable for cloud networking mode.
* `fw_vpc_cidr` - (Optional, String) auto Automatically select the firewall network segment; 10.10.10.0/24 The firewall network segment entered by the user.

The `fw_deploy` object supports the following:

* `deploy_region` - (Required, String) Firewall Deployment Region.
* `width` - (Required, Int) Bandwidth, unit: Mbps.
* `zone_set` - (Required, Set) Zone list.
* `cross_a_zone` - (Optional, Int) Off-site disaster recovery 1: use off-site disaster recovery; 0: do not use off-site disaster recovery; if it is empty, off-site disaster recovery will not be used by default.

The `vpc_fw_instances` object supports the following:

* `fw_deploy` - (Required, List) Deploy regional information.
* `name` - (Required, String) Firewall instance name.
* `vpc_ids` - (Optional, Set) List of VpcIds accessed in private network mode; only used in private network mode.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cfw vpc_instance can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_vpc_instance.example cfwg-4ee69507
```

