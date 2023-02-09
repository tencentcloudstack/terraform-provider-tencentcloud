---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_function_targets_attachment"
sidebar_current: "docs-tencentcloud-resource-clb_function_targets_attachment"
description: |-
  Provides a resource to create a clb function_targets_attachment
---

# tencentcloud_clb_function_targets_attachment

Provides a resource to create a clb function_targets_attachment

## Example Usage

```hcl
resource "tencentcloud_clb_function_targets_attachment" "function_targets" {
  domain           = "xxx.com"
  listener_id      = "lbl-nonkgvc2"
  load_balancer_id = "lb-5dnrkgry"
  url              = "/"

  function_targets {
    weight = 10

    function {
      function_name           = "keep-tf-test-1675954233"
      function_namespace      = "default"
      function_qualifier      = "$LATEST"
      function_qualifier_type = "VERSION"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `function_targets` - (Required, List, ForceNew) List of cloud functions to be bound.
* `listener_id` - (Required, String, ForceNew) Load Balancer Listener ID.
* `load_balancer_id` - (Required, String, ForceNew) Load Balancer Instance ID.
* `domain` - (Optional, String, ForceNew) The domain name of the target forwarding rule. If the LocationId parameter has been entered, this parameter will not take effect.
* `location_id` - (Optional, String, ForceNew) The ID of the target forwarding rule. When binding the cloud function to a layer-7 forwarding rule, this parameter or the Domain+Url parameter must be entered.
* `url` - (Optional, String, ForceNew) The URL of the target forwarding rule. If the LocationId parameter has been entered, this parameter will not take effect.

The `function_targets` object supports the following:

* `function` - (Required, List) Information about cloud functions.&quot;Note: This field may return null, indicating that no valid value can be obtained.
* `weight` - (Optional, Int) Weight.

The `function` object supports the following:

* `function_name` - (Required, String) The name of function.
* `function_namespace` - (Required, String) The namespace of function.
* `function_qualifier` - (Required, String) The version name or alias of the function.
* `function_qualifier_type` - (Optional, String) Identifies the type of FunctionQualifier parameter, possible values: VERSION, ALIAS.Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

clb function_targets_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_clb_function_targets_attachment.function_targets loadBalancerId#listenerId#locationId or loadBalancerId#listenerId#domain#rule
```

