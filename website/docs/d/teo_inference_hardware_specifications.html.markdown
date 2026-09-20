---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_hardware_specifications"
sidebar_current: "docs-tencentcloud-datasource-teo_inference_hardware_specifications"
description: |-
  Use this data source to query TEO inference hardware specifications by zone ID, returning the list of available hardware specifications and their resource configuration (CPU, memory, GPU, etc.) for creating inference services.
---

# tencentcloud_teo_inference_hardware_specifications

Use this data source to query TEO inference hardware specifications by zone ID, returning the list of available hardware specifications and their resource configuration (CPU, memory, GPU, etc.) for creating inference services.

## Example Usage

```hcl
data "tencentcloud_teo_inference_hardware_specifications" "specifications" {
  zone_id = "zone-2qtuhspy7cr6"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Zone ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `hardware_specifications` - Inference hardware specification list.
  * `allowed_gpu_nums` - List of GPU card counts currently supported by the specification.
  * `cpu_num` - Number of CPU cores allocated by default for the specification.
  * `disk_size` - Disk size allocated by default for the specification. Unit: MB.
  * `gpu_mem_size` - GPU memory size allocated by default for the specification. Unit: MB.
  * `gpu_num` - Number of GPU cards allocated by default for the specification.
  * `hardware_spec_id` - Specification unique identifier ID.
  * `mem_size` - Memory size allocated by default for the specification. Unit: MB.
  * `name` - Specification name.
  * `spec` - Specification identifier. Deprecated, refer to `hardware_spec_id`.


