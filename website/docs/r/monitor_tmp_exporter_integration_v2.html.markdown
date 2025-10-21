---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_exporter_integration_v2"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_exporter_integration_v2"
description: |-
  Provides a resource to create a monitor Exporter Integration
---

# tencentcloud_monitor_tmp_exporter_integration_v2

Provides a resource to create a monitor Exporter Integration

~> **NOTE:** If you only want to upgrade the exporter version with same config, you can set `version` under `instanceSpec` with any value to trigger the change.

## Example Usage

### Use qcloud-exporter

```hcl
resource "tencentcloud_monitor_tmp_exporter_integration_v2" "example" {
  instance_id = "prom-gzg3f1em"
  kind        = "qcloud-exporter"
  content     = "{\"name\":\"test\",\"kind\":\"qcloud-exporter\",\"spec\":{\"scrapeSpec\":{\"interval\":\"1m\",\"timeout\":\"1m\",\"relabelConfigs\":\"#metricRelabelings:\\n#- action: labeldrop\\n#  regex: tmp_test_label\\n\"},\"instanceSpec\":{\"region\":\"Guangzhou\",\"role\":\"CM_QCSLinkedRoleInTMP\",\"useRole\":true,\"authProvider\":{\"method\":1,\"presetRole\":\"CM_QCSLinkedRoleInTMP\"},\"rateLimit\":1000,\"delaySeconds\":0,\"rangeSeconds\":0,\"reload_interval_minutes\":10,\"uin\":\"100023201586\",\"tag_key_operation\":\"ToUnderLineAndLower\"},\"exporterSpec\":{\"cvm\":false,\"cbs\":true,\"imageRegistry\":\"ccr.ccs.tencentyun.com\",\"cpu\":\"0.25\",\"memory\":\"0.5Gi\"}},\"status\":{}}"
}
```

### Use es-exporter

```hcl
resource "tencentcloud_monitor_tmp_exporter_integration_v2" "example" {
  instance_id = "prom-gzg3f1em"
  kind        = "es-exporter"
  content = jsonencode({
    "name" : "ex-exporter-example",
    "kind" : "es-exporter",
    "spec" : {
      "instanceSpec" : {
        "user" : "root",
        "password" : "Password@123"
        "url" : "http://127.0.0.1:8080",
        "labels" : {
          "labelKey" : "labelValue"
        }
      },
      "exporterSpec" : {
        "all" : true,
        "indices" : true,
        "indicesSettings" : true,
        "shards" : true,
        "snapshots" : true,
        "clusterSettings" : true
      }
    }
  })
}
```

### Integration Center: CVM Scrape Job

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = "ap-guangzhou-6"
}

resource "tencentcloud_monitor_tmp_instance" "example" {
  instance_name       = "tf-example"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 15
  zone                = "ap-guangzhou-6"
  tags = {
    createdBy = "Terraform"
  }
}

# Integration Center: CVM Scrape Job
resource "tencentcloud_monitor_tmp_exporter_integration" "example" {
  instance_id = tencentcloud_monitor_tmp_instance.example.id
  kind        = "cvm-http-sd-exporter"
  content = jsonencode({
    "kind" : "cvm-http-sd-exporter",
    "spec" : {
      "job" : <<-EOT
        job_name: example-cvm-job-name
        metrics_path: /metrics
        cvm_sd_configs:
        - region: ap-guangzhou
          ports:
            - 9100
          filters:         
            - name: tag:YOUR_TAG_KEY
              values: 
              - YOUR_TAG_VALUE
        relabel_configs: 
          - source_labels: [__meta_cvm_instance_state]
            regex: RUNNING
            action: keep
          - regex: __meta_cvm_tag_(.*)
            replacement: $1
            action: labelmap
          - source_labels: [__meta_cvm_region]
            target_label: region
            action: replace
      EOT
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Integration config.
* `instance_id` - (Required, String) Instance ID.
* `kind` - (Required, String) Type.
* `cluster_id` - (Optional, String) Cluster ID.
* `kube_type` - (Optional, Int) Integration config. 1 - TKE; 2 - EKS; 3 - MEKS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



