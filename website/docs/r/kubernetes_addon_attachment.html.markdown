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

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.16.0.0/16"
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = "10.31.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "keep"
  cluster_desc            = "test cluster desc"
  cluster_version         = "1.20.6"
  cluster_max_service_num = 32

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"
}

resource "tencentcloud_kubernetes_addon_attachment" "addon_cbs" {
  cluster_id = "cls-xxxxxxxx"
  name       = "cbs"
  version    = "1.0.0"
  values = [
    "rootdir=/var/lib/kubelet"
  ]
}

resource "tencentcloud_kubernetes_addon_attachment" "addon_tcr" {
  name    = "tcr"
  version = "1.0.0"
  values = [
    # imagePullSecretsCrs is an array which can configure image pull
    "global.imagePullSecretsCrs[0].name=sample-vpc",
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
  ]
}
```

directly

```hcl
resource "tencentcloud_kubernetes_addon_attachment" "addon_cbs" {
  cluster_id = "cls-xxxxxxxx"
  req_body = { \ "spec\":{\"chart\":{\"chartName\":\"cbs\",\"chartVersion\":\"1.0.0\"},\"values\":{\"rawValuesType\":\"yaml\",\"values\":[]}}}
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) ID of cluster.
* `name` - (Required, ForceNew) Name of addon.
* `request_body` - (Optional) Serialized json string as request body of addon spec. If set, will ignore `version` and `values`.
* `values` - (Optional) Values the addon passthroughs. Conflict with `request_body`.
* `version` - (Optional) Addon version, default latest version. Conflict with `request_body`.

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

