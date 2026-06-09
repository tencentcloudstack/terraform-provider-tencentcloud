---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_boot_script"
sidebar_current: "docs-tencentcloud-resource-emr_boot_script"
description: |-
  Provide a resource to create an EMR boot script.
---

# tencentcloud_emr_boot_script

Provide a resource to create an EMR boot script.

## Example Usage

```hcl
resource "tencentcloud_emr_boot_script" "example" {
  instance_id = "emr-qe336v2e"
  boot_type   = "resourceAfter"
  pre_executed_file_settings {
    path          = "demo.py"
    bucket        = "tf-1309115522"
    cos_file_name = "demo"
    region        = "ap-guangzhou"
  }
}
```

## Argument Reference

The following arguments are supported:

* `boot_type` - (Required, String, ForceNew) Boot type. Valid values: `resourceAfter`, `clusterAfter`, `clusterBefore`.
* `instance_id` - (Required, String, ForceNew) EMR instance ID.
* `pre_executed_file_settings` - (Optional, List) List of pre-execution script settings.

The `pre_executed_file_settings` object supports the following:

* `app_id` - (Optional, String) COS AppId.
* `args` - (Optional, String) Script execution parameters.
* `bucket` - (Optional, String) COS bucket name.
* `cos_file_name` - (Optional, String) Script file name.
* `cos_file_uri` - (Optional, String) Script COS address.
* `cos_secret_id` - (Optional, String) COS SecretId.
* `cos_secret_key` - (Optional, String) COS SecretKey.
* `domain` - (Optional, String) COS domain data.
* `path` - (Optional, String) Script path on COS.
* `region` - (Optional, String) COS region name.
* `remark` - (Optional, String) Remark.
* `run_order` - (Optional, Int) Execution order.
* `when_run` - (Optional, String) Execution timing. Valid values: `resourceAfter`, `clusterAfter`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

EMR boot script can be imported using the id (`{instance_id}#{boot_type}`), e.g.

```
terraform import tencentcloud_emr_boot_script.example emr-qe336v2e#resourceAfter
```

