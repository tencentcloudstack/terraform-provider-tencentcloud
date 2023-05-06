---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_access_address"
sidebar_current: "docs-tencentcloud-datasource-tse_access_address"
description: |-
  Use this data source to query detailed information of tse access_address
---

# tencentcloud_tse_access_address

Use this data source to query detailed information of tse access_address

## Example Usage

```hcl
data "tencentcloud_tse_access_address" "access_address" {
  instance_id = "ins-7eb7eea7"
  # vpc_id = "vpc-xxxxxx"
  # subnet_id = "subnet-xxxxxx"
  # workload = "pushgateway"
  engine_region = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) engine instance Id.
* `engine_region` - (Optional, String) Deploy region.
* `result_output_file` - (Optional, String) Used to save results.
* `subnet_id` - (Optional, String) Subnet ID, Zookeeper does not need to pass vpcid and subnetid; nacos and Polaris need to pass vpcid and subnetid.
* `vpc_id` - (Optional, String) VPC ID, Zookeeper does not need to pass vpcid and subnetid; nacos and Polaris need to pass vpcid and subnetid.
* `workload` - (Optional, String) Name of other engine components(pushgateway, polaris-limiter).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `console_internet_address` - Console public network access addressNote: This field may return null, indicating that a valid value is not available.
* `console_internet_band_width` - Console public network bandwidthNote: This field may return null, indicating that a valid value is not available.
* `console_intranet_address` - Console Intranet access addressNote: This field may return null, indicating that a valid value is not available.
* `env_address_infos` - Apollo Multi-environment public ip address.
  * `config_internet_service_ip` - config public network ip.
  * `config_intranet_address` - config Intranet access addressNote: This field may return null, indicating that a valid value is not available.
  * `enable_config_internet` - Whether to enable the config public network.
  * `enable_config_intranet` - Whether to enable the config Intranet clbNote: This field may return null, indicating that a valid value is not available.
  * `env_name` - env name.
  * `internet_band_width` - Client public network bandwidthNote: This field may return null, indicating that a valid value is not available.
* `internet_address` - Public access address.
* `internet_band_width` - Client public network bandwidthNote: This field may return null, indicating that a valid value is not available.
* `intranet_address` - Intranet access address.
* `limiter_address_infos` - Access IP address of the Polaris traffic limiting server nodeNote: This field may return null, indicating that a valid value is not available.
  * `intranet_address` - VPC access IP address listNote: This field may return null, indicating that a valid value is not available.


