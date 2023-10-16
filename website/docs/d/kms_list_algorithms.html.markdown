---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_list_algorithms"
sidebar_current: "docs-tencentcloud-datasource-kms_list_algorithms"
description: |-
  Use this data source to query detailed information of kms list_algorithms
---

# tencentcloud_kms_list_algorithms

Use this data source to query detailed information of kms list_algorithms

## Example Usage

```hcl
data "tencentcloud_kms_list_algorithms" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `asymmetric_algorithms` - Asymmetric encryption algorithms supported in this region.
  * `algorithm` - Algorithm.
  * `key_usage` - Key usage.
* `asymmetric_sign_verify_algorithms` - Asymmetric signature verification algorithms supported in this region.
  * `algorithm` - Algorithm.
  * `key_usage` - Key usage.
* `symmetric_algorithms` - Symmetric encryption algorithms supported in this region.
  * `algorithm` - Algorithm.
  * `key_usage` - Key usage.


