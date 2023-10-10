---
subcategory: "Cwp"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cwp_license_bind_attachment"
sidebar_current: "docs-tencentcloud-resource-cwp_license_bind_attachment"
description: |-
  Provides a resource to create a cwp license_bind_attachment
---

# tencentcloud_cwp_license_bind_attachment

Provides a resource to create a cwp license_bind_attachment

## Example Usage

```hcl
data "tencentcloud_cwp_machines_simple" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"
  project_ids    = [1210293]

  filters {
    name        = "Version"
    values      = ["BASIC_VERSION"]
    exact_match = true
  }
}

resource "tencentcloud_cwp_license_order" "example" {
  alias        = "tf_example"
  license_type = 0
  license_num  = 1
  region_id    = 1
  project_id   = 0
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_cwp_license_bind_attachment" "example" {
  resource_id  = tencentcloud_cwp_license_order.example.resource_id
  license_id   = tencentcloud_cwp_license_order.example.license_id
  license_type = 0
  quuid        = data.tencentcloud_cwp_machines_simple.example.machines[0].quuid
}
```

## Argument Reference

The following arguments are supported:

* `license_id` - (Required, Int, ForceNew) License ID.
* `license_type` - (Required, Int, ForceNew) LicenseType, 0 CWP Pro - Pay as you go, 1 CWP Pro - Monthly subscription, 2 CWP Ultimate - Monthly subscription. Default is 0.
* `quuid` - (Required, String, ForceNew) Machine quota that needs to be bound.
* `resource_id` - (Required, String, ForceNew) Resource ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `agent_status` - agent status.
* `is_switch_bind` - Is it allowed to change the binding, false is not allowed to change the binding.
* `is_unbind` - Allow unbinding, false does not allow unbinding.
* `machine_ip` - machine ip.
* `machine_name` - machine name.
* `machine_wan_ip` - machine wan ip.
* `uuid` - uuid.


## Import

cwp license_bind_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cwp_license_bind_attachment.example cwplic-ab3edffa#44#2c7e5cce-1cec-4456-8d18-018f160dd987#0
```

