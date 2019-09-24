---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnats"
sidebar_current: "docs-tencentcloud-datasource-dnats"
description: |-
  Use this data source to query detailed information of DNATs.
---

# tencentcloud_dnats

Use this data source to query detailed information of DNATs.

## Example Usage

```hcl
data "tencentcloud_dnats" "foo" {
  name         = "main"
  vpc_id       = "vpc-xfqag"
  nat_id       = "nat-xfaq1"
  elastic_ip   = "123.207.115.136"
  elastic_port = "80"
  private_ip   = "172.16.0.88"
  private_port = "9001"
  description  = "test"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the nat forward.
* `elastic_ip` - (Optional) Network address of the eip.
* `elastic_port` - (Optional) Port of the eip.
* `nat_id` - (Optional) ID of the nat.
* `private_ip` - (Optional) Network address of the backend service.
* `private_port` - (Optional) Port of intranet.
* `result_output_file` - (Optional) Used to save results.
* `vpc_id` - (Optional) ID of the vpc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dnat_list` - Information list of the dnats.
  * `elastic_ip` - Network address of the eip.
  * `elastic_port` - Port of the eip.
  * `nat_id` - ID of the nat.
  * `private_ip` - Network address of the backend service.
  * `private_port` - Port of intranet.
  * `protocol` - Type of the network protocol, the available values include： TCP and UDP.
  * `vpc_id` - ID of the vpc.


