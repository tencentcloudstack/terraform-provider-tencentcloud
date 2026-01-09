---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_log_config"
sidebar_current: "docs-tencentcloud-resource-kubernetes_log_config"
description: |-
  Provides a resource to create a tke kubernetes_log_config
---

# tencentcloud_kubernetes_log_config

Provides a resource to create a tke kubernetes_log_config

## Example Usage

### Create a cls log config

```hcl
locals {
  first_vpc_id    = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  zone_id         = data.tencentcloud_availability_zones_by_product.gz.zones.0.id
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

data "tencentcloud_vpc_subnets" "vpc_one" {
  is_default        = true
  availability_zone = "ap-guangzhou-3"
}

data "tencentcloud_availability_zones_by_product" "gz" {
  name    = "ap-guangzhou-3"
  product = "ckafka"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = local.first_vpc_id
  cluster_cidr            = var.example_cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf_example_cluster"
  cluster_desc            = "example for tke cluster"
  cluster_max_service_num = 32
  cluster_internet        = false # (can be ignored) open it after the nodes added
  cluster_version         = "1.22.5"
  cluster_os              = "tlinux2.2(tkernel3)x86_64"
  cluster_deploy_type     = "MANAGED_CLUSTER"
  log_agent {
    enabled = true
  }
  # without any worker config
}

resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-test-example"
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_kubernetes_log_config" "kubernetes_log_config_cls" {
  log_config_name = "tf-test-cls"
  cluster_id      = tencentcloud_kubernetes_cluster.example.id
  logset_id       = tencentcloud_cls_logset.logset.id
  log_config = jsonencode({
    "apiVersion" : "cls.cloud.tencent.com/v1",
    "kind" : "LogConfig",
    "metadata" : {
      "name" : "tf-test-cls"
    },
    "spec" : {
      "clsDetail" : {
        "extractRule" : {
          "backtracking" : "0",
          "isGBK" : "false",
          "jsonStandard" : "false",
          "unMatchUpload" : "false"
        },
        "indexs" : [
          {
            "indexName" : "namespace"
          },
          {
            "indexName" : "pod_name"
          },
          {
            "indexName" : "container_name"
          }
        ],
        "logFormat" : "default",
        "logType" : "minimalist_log",
        "maxSplitPartitions" : 0,
        "region" : "ap-guangzhou",
        "storageType" : "",
        #   "topicId" : "c26b66bd-617e-4923-bea0-test"
      },
      "inputDetail" : {
        "containerStdout" : {
          "metadataContainer" : [
            "namespace",
            "pod_name",
            "pod_ip",
            "pod_uid",
            "container_id",
            "container_name",
            "image_name",
            "cluster_id"
          ],
          "nsLabelSelector" : "",
          "workloads" : [
            {
              "kind" : "deployment",
              "name" : "testlog1",
              "namespace" : "default"
            }
          ]
        },
        "type" : "container_stdout"
      }
    }
  })
}
```

### Create a ckafka log config

```hcl
locals {
  ckafka_topic = tencentcloud_ckafka_topic.example.topic_name
}

resource "tencentcloud_ckafka_instance" "example" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = local.zone_id
  vpc_id             = local.first_vpc_id
  subnet_id          = local.first_subnet_id
  msg_retention_time = 1300
  kafka_version      = "1.1.1"
  disk_size          = 500
  band_width         = 20
  disk_type          = "CLOUD_BASIC"
  partition          = 400
  charge_type        = "POSTPAID_BY_HOUR"

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }
}

resource "tencentcloud_ckafka_topic" "example" {
  instance_id                    = tencentcloud_ckafka_instance.example.id
  topic_name                     = "tmp"
  note                           = "topic note"
  replica_num                    = 2
  partition_num                  = 1
  clean_up_policy                = "delete"
  sync_replica_min_num           = 1
  unclean_leader_election_enable = false
  retention                      = 60000
}

resource "tencentcloud_kubernetes_log_config" "kubernetes_log_config_ckafka" {
  log_config_name = "tf-test-ckafka"
  cluster_id      = tencentcloud_kubernetes_cluster.example.id
  logset_id       = tencentcloud_cls_logset.logset.id
  log_config = jsonencode({
    "apiVersion" : "cls.cloud.tencent.com/v1",
    "kind" : "LogConfig",
    "metadata" : {
      "name" : "tf-test-ckafka"
    },
    "spec" : {
      "inputDetail" : {
        "containerStdout" : {
          "allContainers" : true,
          "namespace" : "default",
          "nsLabelSelector" : ""
        },
        "type" : "container_stdout"
      },
      "kafkaDetail" : {
        "brokers" : "172.16.0.30:9092", # your ckafka brokers
        "extractRule" : {},
        "instanceId" : "",
        "kafkaType" : "SelfBuildKafka",
        "logType" : "minimalist_log",
        "messageKey" : {
          "value" : "",
          "valueFrom" : {
            "fieldRef" : {
              "fieldPath" : ""
            }
          }
        },
        "metadata" : {},
        "timestampFormat" : "double",
        "timestampKey" : "",
        "topic" : local.ckafka_topic
      }
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `log_config_name` - (Required, String, ForceNew) Log config name.
* `log_config` - (Required, String, ForceNew) JSON expression of log collection configuration. For more details, please refer to the guide: https://www.tencentcloud.com/zh/document/product/457/64846.
* `cluster_type` - (Optional, String, ForceNew) The current cluster type supports tke and eks, default is tke.
* `logset_id` - (Optional, String, ForceNew) CLS log set ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



