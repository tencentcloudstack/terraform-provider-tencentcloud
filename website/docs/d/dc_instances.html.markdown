---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_instances"
sidebar_current: "docs-tencentcloud-datasource-dc_instances"
description: |-
  Use this data source to query detailed information of DC instances.
---

# tencentcloud_dc_instances

Use this data source to query detailed information of DC instances.

## Example Usage

```hcl
data "tencentcloud_dc_instances" "name_select"{
    name = "t"
}

data "tencentcloud_dc_instances"  "id" {
    dcx_id = "dc-kax48sg7"
}

```

## Argument Reference

The following arguments are supported:

* `dc_id` - (Optional, ForceNew) ID of the DC to be queried.
* `name` - (Optional, ForceNew) Name of the DC to be queried.
* `result_output_file` - (Optional, ForceNew) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - Information list of the DC.
  * `access_point_id` - Access point ID of tne DC.
  * `bandwidth` - Bandwidth of the DC.
  * `circuit_code` - The circuit code provided by the operator for the DC.
  * `create_time` - Creation time of resource.
  * `customer_address` - Interconnect IP of the DC within client. Note: This field may return null, indicating that no valid values are taken.
  * `customer_email` - Applicant email of the DC, the default is obtained from the account. Note: This field may return null, indicating that no valid values are taken.
  * `customer_name` - Applicant name of the DC, the default is obtained from the account. Note: This field may return null, indicating that no valid values are taken.
  * `customer_phone` - Applicant phone number of the DC, the default is obtained from the account. Note: This field may return null, indicating that no valid values are taken.
  * `dc_id` - ID of the DC.
  * `enabled_time` - Enable time of resource.
  * `expired_time` - Expire date of resource.
  * `fault_report_contact_person` - Contact of reporting a faulty. Note: This field may return null, indicating that no valid values are taken.
  * `fault_report_contact_phone` - Phone number of reporting a faulty. Note: This field may return null, indicating that no valid values are taken.
  * `line_operator` - Operator of the DC, and available values include ChinaTelecom, ChinaMobile, ChinaUnicom, In-houseWiring, ChinaOther and InternationalOperator.
  * `location` - The DC location where the connection is located.
  * `name` - Name of the DC.
  * `port_type` - Port type of the DC in client, and available values include 100Base-T, 1000Base-T, 1000Base-LX, 10GBase-T and 10GBase-LR. The default value is 1000Base-LX.
  * `redundant_dc_id` - ID of the redundant DC.
  * `state` - State of the DC, and available values include REJECTED, TOPAY, PAID, ALLOCATED, AVAILABLE, DELETING and DELETED.
  * `tencent_address` - Interconnect IP of the DC within Tencent. Note: This field may return null, indicating that no valid values are taken.


