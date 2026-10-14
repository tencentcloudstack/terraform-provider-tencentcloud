---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_reserved_concurrency_config"
sidebar_current: "docs-tencentcloud-resource-scf_reserved_concurrency_config"
description: |-
  Provides a resource to create a scf reserved_concurrency_config
---

# tencentcloud_scf_reserved_concurrency_config

Provides a resource to create a scf reserved_concurrency_config

## Example Usage

```hcl
resource "tencentcloud_scf_reserved_concurrency_config" "reserved_concurrency_config" {
  function_name            = "keep-1676351130"
  reserved_concurrency_mem = 128000
  namespace                = "default"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String, ForceNew) Specifies the function of which you want to configure the reserved quota.
* `reserved_concurrency_mem` - (Required, Int, ForceNew) Reserved memory quota of the function. Note: the upper limit for the total reserved quota of the function is the user's total concurrency memory minus 12800.
* `namespace` - (Optional, String, ForceNew) Function namespace. Default value: default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [PutReservedConcurrencyConfig](https://www.tencentcloud.com/document/product/583/39224) - This API is used to configure the reserved quota of a function.
* Read: [GetReservedConcurrencyConfig](https://www.tencentcloud.com/document/product/583/39226) - This API is used to obtain the reserved quota details of a function.
* Delete: [DeleteReservedConcurrencyConfig](https://www.tencentcloud.com/document/product/583/39228) - This API is used to delete the configuration of reserved quota.

## Import

scf reserved_concurrency_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_reserved_concurrency_config.reserved_concurrency_config reserved_concurrency_config_id
```

