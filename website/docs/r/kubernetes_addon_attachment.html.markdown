---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_addon_attachment"
sidebar_current: "docs-tencentcloud-resource-kubernetes_addon_attachment"
description: |-
  Provide a resource to configure kubernetes cluster app addons.
---

# tencentcloud_kubernetes_addon_attachment

Provide a resource to configure kubernetes cluster app addons.

~> **NOTE**: Avoid to using legacy "1.0.0" version, leave the versions empty so we can fetch the latest while creating.

## Example Usage

```hcl
resource "tencentcloud_kubernetes_addon_attachment" "addon_cbs" {
  cluster_id = "cls-xxxxxxxx"
  name       = "cbs"
  # version = "1.0.5"
  values = [
    "rootdir=/var/lib/kubelet"
  ]
}

resource "tencentcloud_kubernetes_addon_attachment" "addon_tcr" {
  cluster_id = "cls-xxxxxxxx"
  name       = "tcr"
  # version = "1.0.0"
  values = [
    # imagePullSecretsCrs is an array which can configure image pull
    "global.imagePullSecretsCrs[0].name=unique-sample-vpc",
    "global.imagePullSecretsCrs[0].namespaces=tcr-assistant-system",
    "global.imagePullSecretsCrs[0].serviceAccounts=*",
    "global.imagePullSecretsCrs[0].type=docker",
    "global.imagePullSecretsCrs[0].dockerUsername=100012345678",
    "global.imagePullSecretsCrs[0].dockerPassword=a.b.tcr-token",
    "global.imagePullSecretsCrs[0].dockerServer=xxxx.tencentcloudcr.com",
    "global.imagePullSecretsCrs[1].name=sample-public",
    "global.imagePullSecretsCrs[1].namespaces=*",
    "global.imagePullSecretsCrs[1].serviceAccounts=*",
    "global.imagePullSecretsCrs[1].type=docker",
    "global.imagePullSecretsCrs[1].dockerUsername=100012345678",
    "global.imagePullSecretsCrs[1].dockerPassword=a.b.tcr-token",
    "global.imagePullSecretsCrs[1].dockerServer=sample",
    # Specify global hosts
    "global.hosts[0].domain=sample-vpc.tencentcloudcr.com",
    "global.hosts[0].ip=10.16.0.49",
    "global.hosts[0].disabled=false",
  ]
}
```

Install new addon by passing spec json to req_body directly

```hcl
resource "tencentcloud_kubernetes_addon_attachment" "addon_cbs" {
  cluster_id   = "cls-xxxxxxxx"
  request_body = <<EOF
  {
    "spec":{
        "chart":{
            "chartName":"cbs",
            "chartVersion":"1.0.5"
        },
        "values":{
            "rawValuesType":"yaml",
            "values":[
              "rootdir=/var/lib/kubelet"
            ]
        }
    }
  }
EOF
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) ID of cluster.
* `name` - (Required, String, ForceNew) Name of addon.
* `request_body` - (Optional, String) Serialized json string as request body of addon spec. If set, will ignore `version` and `values`.
* `values` - (Optional, List: [`String`]) Values the addon passthroughs. Conflict with `request_body`.
* `version` - (Optional, String) Addon version, default latest version. Conflict with `request_body`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `response_body` - Addon response body.
* `status` - Addon current status.


## Import

Addon can be imported by using cluster_id#addon_name
```
$ terraform import tencentcloud_kubernetes_addon_attachment.addon_cos cls-xxxxxxxx#cos
```

