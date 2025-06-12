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

~> **NOTE**: This resource has been deprecated in Terraform TencentCloud provider version 1.81.200. Please use `tencentcloud_kubernetes_addon` instead.

~> **NOTE**: Avoid to using legacy "1.0.0" version, leave the versions empty so we can fetch the latest while creating.

## Example Usage

### Install cbs addon by passing values

```hcl
resource "tencentcloud_kubernetes_addon_attachment" "example" {
  cluster_id = "cls-fdy7hm1q"
  name       = "cbs"
  version    = "1.1.7"
  values = [
    "rootdir=/var/lib/kubelet"
  ]
}
```

### Install tcr addon by passing values

```hcl
resource "tencentcloud_kubernetes_addon_attachment" "example" {
  cluster_id = "cls-fdy7hm1q"
  name       = "tcr"
  version    = "1.0.0"
  values = [
    # imagePullSecretsCrs is an array which can configure image pull
    "global.imagePullSecretsCrs[0].name=${local.tcr_id}-vpc",                              #specify a unique name, invalid format as: `${tcrId}-vpc`
    "global.imagePullSecretsCrs[0].namespaces=${local.ns_name}",                           #input the specified namespaces of the cluster, or input `*` for all.
    "global.imagePullSecretsCrs[0].serviceAccounts=*",                                     #input the specified service account of the cluster, or input `*` for all.
    "global.imagePullSecretsCrs[0].type=docker",                                           #only support docker now
    "global.imagePullSecretsCrs[0].dockerUsername=${local.user_name}",                     #input the access username, or you can create it from `tencentcloud_tcr_token`
    "global.imagePullSecretsCrs[0].dockerPassword=${local.token}",                         #input the access token, or you can create it from `tencentcloud_tcr_token`
    "global.imagePullSecretsCrs[0].dockerServer=${local.tcr_name}-vpc.tencentcloudcr.com", #invalid format as: `${tcr_name}-vpc.tencentcloudcr.com`
    "global.imagePullSecretsCrs[1].name=${local.tcr_id}-public",                           #specify a unique name, invalid format as: `${tcr_id}-public`
    "global.imagePullSecretsCrs[1].namespaces=${local.ns_name}",
    "global.imagePullSecretsCrs[1].serviceAccounts=*",
    "global.imagePullSecretsCrs[1].type=docker",
    "global.imagePullSecretsCrs[1].dockerUsername=${local.user_name}",                 #refer to previous description
    "global.imagePullSecretsCrs[1].dockerPassword=${local.token}",                     #refer to previous description
    "global.imagePullSecretsCrs[1].dockerServer=${local.tcr_name}.tencentcloudcr.com", #invalid format as: `${tcr_name}.tencentcloudcr.com`
    "global.cluster.region=gz",
    "global.cluster.longregion=ap-guangzhou",

    # Specify global hosts(optional), the numbers of hosts must be matched with the numbers of imagePullSecretsCrs
    "global.hosts[0].domain=${local.tcr_name}-vpc.tencentcloudcr.com", #Corresponds to the dockerServer in the imagePullSecretsCrs above
    "global.hosts[0].ip=${local.end_point}",                           #input InternalEndpoint of tcr instance, you can get it from data source `tencentcloud_tcr_instances`
    "global.hosts[0].disabled=false",                                  #disabled this host config or not
    "global.hosts[1].domain=${local.tcr_name}.tencentcloudcr.com",
    "global.hosts[1].ip=${local.end_point}",
    "global.hosts[1].disabled=false",
  ]
}

locals {
  tcr_id    = tencentcloud_tcr_instance.example.id
  tcr_name  = tencentcloud_tcr_instance.example.name
  ns_name   = tencentcloud_tcr_namespace.example.name
  user_name = tencentcloud_tcr_token.example.user_name
  token     = tencentcloud_tcr_token.example.token
  end_point = data.tencentcloud_tcr_instances.example.instance_list.0.internal_end_point
}

resource "tencentcloud_tcr_token" "example" {
  instance_id = local.tcr_id
  description = "tcr token"
}

data "tencentcloud_tcr_instances" "example" {
  instance_id = local.tcr_id
}

resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example"
  instance_type = "basic"
  delete_bucket = true

  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = local.tcr_id
  name           = "tf-example"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}
```

### Install new addon by passing spec json to req_body directly

```hcl
resource "tencentcloud_kubernetes_addon_attachment" "example" {
  cluster_id   = "cls-fdy7hm1q"
  name         = "cbs"
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
* `raw_values_type` - (Optional, String) The type of raw Values. Required with `raw_values`.
* `raw_values` - (Optional, String) Raw Values. Conflict with `request_body`. Required with `raw_values_type`.
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
$ terraform import tencentcloud_kubernetes_addon_attachment.example cls-fdy7hm1q#cbs
```

