Provides a resource to create a monitor tmpExporterIntegration

~> **NOTE:** If you only want to upgrade the exporter version with same config, you can set `version` under `instanceSpec` with any value to trigger the change.

Example Usage

Use blackbox-exporter

```hcl
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegration" {
  instance_id = "prom-dko9d0nu"
  kind = "blackbox-exporter"
  content = "{\"name\":\"test\",\"kind\":\"blackbox-exporter\",\"spec\":{\"instanceSpec\":{\"module\":\"http_get\",\"urls\":[\"xx\"]}}}"
  kube_type = 1
  cluster_id = "cls-bmuaukfu"
}
```

Use es-exporter

```
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationEs" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id
  kind        = "es-exporter"
  content = jsonencode({
    "name": "ex-exporter-example",
    "kind": "es-exporter",
    "spec": {
      "instanceSpec": {
        "url": "http://127.0.0.1:9123",
        "labels": {
          "instance": "es-abcd"
        },
        "version": "1.70.1",
        "user": "fugiat Duis minim",
        "password": "exercitation cillum velit"
      },
      "exporterSpec": {
        "all": true,
        "indicesSettings": false,
        "snapshots": false,
        "indices": true,
        "shards": false
      }
    }
  })
  cluster_id = ""
  kube_type  = 3
}
```

Integration Center: CVM Scrape Job

```
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-eks-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "sub" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_monitor_tmp_instance" "tmpInstance" {
  instance_name       = "tf-test-tmp"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.sub.id
  data_retention_time = 15
  zone                = "ap-guangzhou-3"
  tags = {
    "createdBy" = "terraform"
  }
}

# Integration Center: CVM Scrape Job
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegration" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "cvm-http-sd-exporter"
  content     = jsonencode({
    "kind": "cvm-http-sd-exporter",
    "spec": {
      "job": <<-EOT
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
  kube_type   = 3
  cluster_id  = ""
}
```