---
subcategory: "TencentCloud ServiceMesh(TCM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcm_tracing_config"
sidebar_current: "docs-tencentcloud-resource-tcm_tracing_config"
description: |-
  Provides a resource to create a tcm tracing_config
---

# tencentcloud_tcm_tracing_config

Provides a resource to create a tcm tracing_config

~> **NOTE:** If you use the config attribute tracing in tencentcloud_tcm_mesh, do not use tencentcloud_tcm_tracing_config

## Example Usage

```hcl
resource "tencentcloud_tcm_tracing_config" "tracing_config" {
  mesh_id = "mesh-xxxxxxxx"
  enable  = true
  apm {
    enable      = true
    region      = "ap-guangzhou"
    instance_id = "apm-xxx"
  }
  sampling =
  zipkin {
    address = "10.10.10.10:9411"
  }
}

resource "tencentcloud_tcm_tracing_config" "delete_config" {
  mesh_id = "mesh-rofjmxxx"
  enable  = true
  apm {
    enable = false
    # region = "ap-guangzhou"
    # instance_id = "apm-xxx"
  }
  sampling = 0
  zipkin {
    address = ""
  }
}
```

## Argument Reference

The following arguments are supported:

* `mesh_id` - (Required, String) Mesh ID.
* `apm` - (Optional, List) APM config.
* `enable` - (Optional, Bool) Whether enable tracing.
* `sampling` - (Optional, Float64) Tracing sampling, 0.0-1.0.
* `zipkin` - (Optional, List) Third party zipkin config.

The `apm` object supports the following:

* `enable` - (Optional, Bool) Whether enable APM.
* `instance_id` - (Optional, String) Instance id of the APM.
* `region` - (Optional, String) Region.

The `zipkin` object supports the following:

* `address` - (Required, String) Zipkin address.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tcm tracing_config can be imported using the mesh_id, e.g.
```
$ terraform import tencentcloud_tcm_tracing_config.tracing_config mesh-rofjmxxx
```

