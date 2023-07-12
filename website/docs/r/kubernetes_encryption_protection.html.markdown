---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_encryption_protection"
sidebar_current: "docs-tencentcloud-resource-kubernetes_encryption_protection"
description: |-
  Provides a resource to create a tke encryption_protection
---

# tencentcloud_kubernetes_encryption_protection

Provides a resource to create a tke encryption_protection

## Example Usage

### Enable tke encryption protection

```hcl
variable "example_region" {
  default = "ap-guangzhou"
}

variable "example_cluster_cidr" {
  default = "10.32.0.0/16"
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = var.example_cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf_example_cluster"
  cluster_desc            = "a tf example cluster for the kms test"
  cluster_max_service_num = 32
  cluster_deploy_type     = "MANAGED_CLUSTER"
}

resource "tencentcloud_kms_key" "example" {
  alias       = "tf-example-kms-key"
  description = "example of kms key instance"
  key_usage   = "ENCRYPT_DECRYPT"
  is_enabled  = true
}

resource "tencentcloud_kubernetes_encryption_protection" "example" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  kms_configuration {
    key_id     = tencentcloud_kms_key.example.id
    kms_region = var.example_region
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) cluster id.
* `kms_configuration` - (Required, List, ForceNew) kms encryption configuration.

The `kms_configuration` object supports the following:

* `key_id` - (Optional, String) kms id.
* `kms_region` - (Optional, String) kms region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - kms encryption status.


