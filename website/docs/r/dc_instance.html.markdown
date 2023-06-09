---
subcategory: "Direct Connect(DC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_instance"
sidebar_current: "docs-tencentcloud-resource-dc_instance"
description: |-
  Provides a resource to create a dc instance
---

# tencentcloud_dc_instance

Provides a resource to create a dc instance

## Example Usage

```hcl
resource "tencentcloud_dc_instance" "instance" {
  access_point_id         = "ap-shenzhen-b-ft"
  bandwidth               = 10
  customer_contact_number = "0"
  direct_connect_name     = "terraform-for-test"
  line_operator           = "In-houseWiring"
  port_type               = "10GBase-LR"
  sign_law                = true
  vlan                    = -1
}
```

## Argument Reference

The following arguments are supported:

* `access_point_id` - (Required, String) Access point of connection.You can call `DescribeAccessPoints` to get the region ID. The selected access point must exist and be available.
* `direct_connect_name` - (Required, String) Connection name.
* `line_operator` - (Required, String) ISP that provides connections. Valid values: ChinaTelecom (China Telecom), ChinaMobile (China Mobile), ChinaUnicom (China Unicom), In-houseWiring (in-house wiring), ChinaOther (other Chinese ISPs), InternationalOperator (international ISPs).
* `port_type` - (Required, String) Port type of connection. Valid values: 100Base-T (100-Megabit electrical Ethernet interface), 1000Base-T (1-Gigabit electrical Ethernet interface), 1000Base-LX (1-Gigabit single-module optical Ethernet interface; 10 KM), 10GBase-T (10-Gigabit electrical Ethernet interface), 10GBase-LR (10-Gigabit single-module optical Ethernet interface; 10 KM). Default value: 1000Base-LX.
* `bandwidth` - (Optional, Int) Connection port bandwidth in Mbps. Value range: [2,10240]. Default value: 1000.
* `circuit_code` - (Optional, String) Circuit code of a connection, which is provided by the ISP or connection provider.
* `customer_address` - (Optional, String) User-side IP address for connection debugging, which is automatically assigned by default.
* `customer_contact_mail` - (Optional, String) Email address of connection applicant, which is obtained from the account system by default.
* `customer_contact_number` - (Optional, String) Contact number of connection applicant, which is obtained from the account system by default.
* `customer_name` - (Optional, String) Name of connection applicant, which is obtained from the account system by default.
* `fault_report_contact_number` - (Optional, String) Fault reporting contact number.
* `fault_report_contact_person` - (Optional, String) Fault reporting contact person.
* `location` - (Optional, String) Local IDC location.
* `redundant_direct_connect_id` - (Optional, String) ID of redundant connection.
* `sign_law` - (Optional, Bool) Whether the connection applicant has signed the service agreement. Default value: true.
* `tencent_address` - (Optional, String) Tencent-side IP address for connection debugging, which is automatically assigned by default.
* `vlan` - (Optional, Int) VLAN for connection debugging, which is enabled and automatically assigned by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dc instance can be imported using the id, e.g.

```
terraform import tencentcloud_dc_instance.instance dc_id
```

