---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_attachment"
sidebar_current: "docs-tencentcloud-resource-clb_attachment"
description: |-
  Provides a resource to create a CLB attachment.
---

# tencentcloud_clb_attachment

Provides a resource to create a CLB attachment.

~> **NOTE:** This resource is designed to manage the entire set of binding relationships associated with a particular CLB (Cloud Load Balancer). As such, it does not allow the simultaneous use of this resource for the same CLB across different contexts or environments.

## Example Usage

### Bind a Cvm instance by using rule_id

```hcl
resource "tencentcloud_clb_attachment" "example" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9"
  rule_id     = "loc-4xxr2cy7"

  targets {
    instance_id = "ins-1flbqyp8"
    port        = 80
    weight      = 10
  }
}
```

### Bind a Cvm instance by using domian and url

```hcl
resource "tencentcloud_clb_attachment" "example" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9"
  domain      = "test.com"
  url         = "/"

  targets {
    instance_id = "ins-1flbqyp8"
    port        = 80
    weight      = 10
  }
}
```

### Bind multiple Cvm instances by using rule_id

```hcl
resource "tencentcloud_clb_attachment" "example" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9"
  rule_id     = "loc-4xxr2cy7"

  targets {
    instance_id = "ins-1flbqyp8"
    port        = 80
    weight      = 10
  }

  targets {
    instance_id = "ins-ekloqpa1"
    port        = 81
    weight      = 10
  }
}
```

### Bind multiple Cvm instances by using domian and url

```hcl
resource "tencentcloud_clb_attachment" "example" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9"
  domain      = "test.com"
  url         = "/"

  targets {
    instance_id = "ins-1flbqyp8"
    port        = 80
    weight      = 10
  }

  targets {
    instance_id = "ins-ekloqpa1"
    port        = 81
    weight      = 10
  }
}
```

### Bind backend target is ENI by using rule_id

```hcl
resource "tencentcloud_clb_attachment" "example" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9"
  rule_id     = "loc-4xxr2cy7"

  targets {
    eni_ip = "172.16.16.52"
    port   = 8090
    weight = 50
  }
}
```

### Bind backend target is ENI by using domian and url

```hcl
resource "tencentcloud_clb_attachment" "example" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9"
  domain      = "test.com"
  url         = "/path"

  targets {
    eni_ip = "172.16.16.52"
    port   = 8090
    weight = 50
  }
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, String, ForceNew) ID of the CLB.
* `listener_id` - (Required, String, ForceNew) ID of the CLB listener.
* `targets` - (Required, Set) Information of the backends to be attached.
* `domain` - (Optional, String, ForceNew) Domain of the target forwarding rule. Does not take effect when parameter `rule_id` is provided.
* `rule_id` - (Optional, String, ForceNew) ID of the CLB listener rule. Only supports listeners of `HTTPS` and `HTTP` protocol.
* `url` - (Optional, String, ForceNew) URL of the target forwarding rule. Does not take effect when parameter `rule_id` is provided.

The `targets` object supports the following:

* `port` - (Required, Int) Port of the backend server. Valid value ranges: (0~65535).
* `eni_ip` - (Optional, String) Eni IP address of the backend server, conflict with `instance_id` but must specify one of them.
* `instance_id` - (Optional, String) CVM Instance Id of the backend server, conflict with `eni_ip` but must specify one of them.
* `weight` - (Optional, Int) Forwarding weight of the backend service. Valid value ranges: (0~100). defaults to `10`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `protocol_type` - Type of protocol within the listener.


## Import

CLB attachment can be imported using the id, e.g.

If use rule_id

```
$ terraform import tencentcloud_clb_attachment.example loc-4xxr2cy7#lbl-hh141sn9#lb-7a0t6zqb
```

If use domain & url

```
$ terraform import tencentcloud_clb_attachment.example test.com,/path#lbl-hh141sn9#lb-7a0t6zqb
```

Of if use layer-4 forwarding rule

```
$ terraform import tencentcloud_clb_attachment.example ""#lbl-hh141sn9#lb-7a0t6zqb
```

