---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_addon"
sidebar_current: "docs-tencentcloud-resource-kubernetes_addon"
description: |-
  Provide a resource to configure kubernetes cluster app addons.
---

# tencentcloud_kubernetes_addon

Provide a resource to configure kubernetes cluster app addons.

## Example Usage

### Install tcr addon

```hcl
resource "tencentcloud_kubernetes_addon" "example" {
  cluster_id = "cls-k2o1ws9g"
  addon_name = "tcr"
  raw_values = jsonencode({
    global = {
      imagePullSecretsCrs = [
        {
          name            = "tcr-h3ff76s9"
          namespaces      = "*"
          serviceAccounts = "*"
          type            = "docker"
          dockerUsername  = "100038911322"
          dockerPassword  = "eyJhbGciOiJSUzI1NiIsImtpZCI6************"
          dockerServer    = "testcd.tencentcloudcr.com"
        }
      ]
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `addon_name` - (Required, String, ForceNew) Name of addon.
* `cluster_id` - (Required, String, ForceNew) ID of cluster.
* `addon_version` - (Optional, String) Version of addon. If no set, the latest version will be installed by default.
* `raw_values` - (Optional, String) Params of addon, base64 encoded json format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `phase` - Status of addon.
* `reason` - Reason of addon failed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `3m`) Used when creating the resource.
* `update` - (Defaults to `3m`) Used when updating the resource.
* `delete` - (Defaults to `3m`) Used when deleting the resource.

## Import

kubernetes cluster app addons can be imported using the id, e.g.
```
$ terraform import tencentcloud_kubernetes_addon.example cls-k2o1ws9g#tcr
```

