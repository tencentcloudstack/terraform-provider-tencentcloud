---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_customized_config_attachment"
sidebar_current: "docs-tencentcloud-resource-clb_customized_config_attachment"
description: |-
  Provides a resource to create a CLB customized config attachment.
---

# tencentcloud_clb_customized_config_attachment

Provides a resource to create a CLB customized config attachment.

~> **NOTE:** This resource must exclusive in one CLB customized config attachment, do not declare additional rule resources of this CLB customized config attachment elsewhere.

## Example Usage

### If config_type is SERVER

```hcl
resource "tencentcloud_clb_customized_config_v2" "example" {
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  config_name    = "tf-example"
  config_type    = "SERVER"
}

resource "tencentcloud_clb_customized_config_attachment" "example" {
  config_id = tencentcloud_clb_customized_config_v2.example.config_id
  bind_list {
    load_balancer_id = "lb-g1miv1ok"
    listener_id      = "lbl-9bsa90io"
    domain           = "demo1.com"
  }

  bind_list {
    load_balancer_id = "lb-g1miv1ok"
    listener_id      = "lbl-qfljudr4"
    domain           = "demo2.com"
  }
}
```

### If config_type is LOCATION

```hcl
resource "tencentcloud_clb_customized_config_v2" "example" {
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  config_name    = "tf-example"
  config_type    = "LOCATION"
}

resource "tencentcloud_clb_customized_config_attachment" "example" {
  config_id = tencentcloud_clb_customized_config_v2.example.config_id
  bind_list {
    load_balancer_id = "lb-g1miv1ok"
    listener_id      = "lbl-9bsa90io"
    domain           = "demo1.com"
    location_id      = "loc-5he3og2u"
  }

  bind_list {
    load_balancer_id = "lb-g1miv1ok"
    listener_id      = "lbl-qfljudr4"
    domain           = "demo2.com"
    location_id      = "loc-0oxl4lfw"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bind_list` - (Required, Set) Associated server or location.
* `config_id` - (Required, String, ForceNew) ID of Customized Config.

The `bind_list` object supports the following:

* `domain` - (Required, String) Domain.
* `listener_id` - (Required, String) Listener ID.
* `load_balancer_id` - (Required, String) Clb ID.
* `location_id` - (Optional, String) Location ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CLB customized config attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_customized_config_attachment.example pz-ivj39268
```

