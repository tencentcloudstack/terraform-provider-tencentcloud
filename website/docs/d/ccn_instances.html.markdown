---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_instances"
sidebar_current: "docs-tencentcloud-datasource-ccn_instances"
description: |-
  Use this data source to query detailed information of CCN instances.
---

# tencentcloud_ccn_instances

Use this data source to query detailed information of CCN instances.

## Example Usage

```hcl
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

data "tencentcloud_ccn_instances" "id_instances" {
  ccn_id = tencentcloud_ccn.main.id
}

data "tencentcloud_ccn_instances" "name_instances" {
  name = tencentcloud_ccn.main.name
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Optional) ID of the CCN to be queried.
* `name` - (Optional) Name of the CCN to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - Information list of CCN.
  * `attachment_list` - Information list of instance is attached.
    * `attached_time` - Time of attaching.
    * `cidr_block` - A network address block of the instance that is attached.
    * `instance_id` - ID of instance is attached.
    * `instance_region` - The region that the instance locates at.
    * `instance_type` - Type of attached instance network, and available values include VPC, DIRECTCONNECT, BMVPC and VPNGW.
    * `state` - States of instance is attached, and available values include PENDING, ACTIVE, EXPIRED, REJECTED, DELETED, FAILED(asynchronous forced disassociation after 2 hours), ATTACHING, DETACHING and DETACHFAILED(asynchronous forced disassociation after 2 hours).
  * `ccn_id` - ID of the CCN.
  * `create_time` - Creation time of resource.
  * `description` - Description of the CCN.
  * `name` - Name of the CCN.
  * `qos` - Service quality of CCN, and the available value include 'PT', 'AU', 'AG'. The default is 'AU'.
  * `state` - States of instance. The available value include 'ISOLATED'(arrears) and 'AVAILABLE'.


