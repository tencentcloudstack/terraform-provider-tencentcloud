---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_webhook_trigger"
sidebar_current: "docs-tencentcloud-resource-tcr_webhook_trigger"
description: |-
  Provides a resource to create a tcr webhook trigger
---

# tencentcloud_tcr_webhook_trigger

Provides a resource to create a tcr webhook trigger

## Example Usage

### Create a tcr webhook trigger instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "basic"
  delete_bucket = true

  tags = {
    test = "test"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf_example_ns_retention"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

data "tencentcloud_tcr_namespaces" "example" {
  instance_id = tencentcloud_tcr_namespace.example.instance_id
}

locals {
  ns_id = data.tencentcloud_tcr_namespaces.example.namespace_list.0.id
}

resource "tencentcloud_tcr_webhook_trigger" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  namespace   = tencentcloud_tcr_namespace.example.name
  trigger {
    name = "trigger-example"
    targets {
      address = "http://example.org/post"
      headers {
        key    = "X-Custom-Header"
        values = ["a"]
      }
    }
    event_types  = ["pushImage"]
    condition    = ".*"
    enabled      = true
    description  = "example for trigger description"
    namespace_id = local.ns_id

  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String) namespace name.
* `registry_id` - (Required, String) instance Id.
* `trigger` - (Required, List) trigger parameters.
* `tags` - (Optional, Map) Tag description list.

The `headers` object of `targets` supports the following:

* `key` - (Required, String) Header Key.
* `values` - (Required, Set) Header Values.

The `targets` object of `trigger` supports the following:

* `address` - (Required, String) target address.
* `headers` - (Optional, List) custom Headers.

The `trigger` object supports the following:

* `condition` - (Required, String) trigger rule.
* `enabled` - (Required, Bool) enable trigger.
* `event_types` - (Required, Set) trigger action.
* `name` - (Required, String) trigger name.
* `targets` - (Required, List) trigger target.
* `description` - (Optional, String) trigger description.
* `namespace_id` - (Optional, Int) the namespace Id to which the trigger belongs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tcr webhook_trigger can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_webhook_trigger.example webhook_trigger_id
```

