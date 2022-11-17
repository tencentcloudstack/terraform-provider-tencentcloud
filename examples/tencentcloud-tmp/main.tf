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

#Integration Center: CVM Scrape Job
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegration" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "cvm-http-sd-exporter"
  content     = jsonencode({
    "kind": "cvm-http-sd-exporter",
    "spec": {
      "job": <<-EOT
        job_name: example-job-name
        metrics_path: /metrics
        cvm_sd_configs:
        - region: ap-guangzhou
          ports:
            - 9100
          filters:         
            - name: tag:示例标签键
              values: 
              - 示例标签值
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

#Integration Center: Cloud Monitor
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationCloudMonitor" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "qcloud-exporter"
  content = jsonencode({
    "name": "qcloud-exporter-example",
    "kind": "qcloud-exporter",
    "spec": {
      "instanceSpec": {
        "region": "ap-singapore",
        "delaySeconds": 0,
        "useRole": true
      },
      "exporterSpec": {
        "mariadb": false,
        "redis": false,
        "vpngw": false,
        "lb": false,
        "nacos": false,
        "cos": false,
        "lb_public": false,
        "ces": false,
        "lighthouse": false,
        "cvm": false,
        "cynosdb_mysql": false,
        "vpnx": false,
        "dc": false,
        "dcg": false,
        "ckafka": false,
        "cmongo": false,
        "nat_gateway": false,
        "redis_mem": false,
        "zookeeper": false,
        "cdb": false,
        "self": false,
        "cdn": false,
        "tdmysql": false,
        "postgres": false,
        "dcx": false,
        "tdmq": false
      }
    }
  })
  cluster_id = ""
  kube_type  = 3
}

#Integration Center: Health Check
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationBlackbox" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "blackbox-exporter"
  content = jsonencode({
    "name": "blackbox-exporter-example",
    "kind": "blackbox-exporter",
    "spec": {
      "instanceSpec": {
        "module": "http_post",
        "urls": [
          "http://127.0.0.1:9123"
        ],
        "labels": {
          "instance": "instance-abcd"
        }
      },
      "scrapeSpec": {
        "interval": "15s"
      }
    }
  })
  cluster_id = ""
  kube_type  = 3
}

#Integration Center: Consul
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationConsul" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "consul-exporter"
  content = jsonencode({
    "name": "consul-exporter-example",
    "kind": "consul-exporter",
    "spec": {
      "instanceSpec": {
        "server": "127.0.0.1:9123",
        "labels": {
          "instance": "consul-abcd"
        }
      }
    }
  })
  cluster_id = ""
  kube_type  = 3
}

#Integration Center: ElasticSearch
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

#Integration Center: Kafka
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationKafka" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "kafka-exporter"
  content = jsonencode({
    "name": "kafka-exporter-example",
    "kind": "kafka-exporter",
    "spec": {
      "instanceSpec": {
        "servers": [
          "127.0.0.1:9123"
        ],
        "labels": {
          "instance": "ckafka-abcd"
        }
      },
      "exporterSpec": {}
    }
  })
  cluster_id = ""
  kube_type  = 3
}

#Integration Center: Memcached
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationMemcached" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "memcached-exporter"
  content = jsonencode({
    "name": "memcached-exporter-example",
    "kind": "memcached-exporter",
    "spec": {
      "instanceSpec": {
        "address": "127.0.0.1:9123",
        "labels": {
          "instance": "crs-abcd"
        }
      }
    }
  })
  cluster_id = ""
  kube_type  = 3
}

#Integration Center: MongoDB
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationMongodb" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "mongodb-exporter"
  content = jsonencode({
    "name": "mongodb-exporter-example",
    "kind": "mongodb-exporter",
    "spec": {
      "instanceSpec": {
        "user": "nisi ullamco eiusmod et ea",
        "password": "Duis",
        "servers": [
          "127.0.0.1:9123",
          "127.0.0.2:9123"
        ],
        "labels": {
          "instance": "cmgo-abcd"
        }
      },
      "exporterSpec": {
        "collectors": [
          "collection",
          "indexusage",
          "topmetrics"
        ]
      }
    }
  })
  cluster_id = ""
  kube_type  = 3
}

#Integration Center: Mysql
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationMysql" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "mysql-exporter"
  content = jsonencode({
    "name": "mysql-exporter",
    "kind": "mysql-exporter",
    "spec": {
      "instanceSpec": {
        "user": "est",
        "password": "id proident deserunt sint",
        "address": "127.0.0.1:9123",
        "labels": {
          "instance": "cdb-abcd"
        }
      }
    }
  })
  cluster_id = ""
  kube_type  = 3
}

#Integration Center: Postgres
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationPostgres" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "postgres-exporter"
  content = jsonencode({
    "name": "postgres-exporter-example",
    "kind": "postgres-exporter",
    "spec": {
      "instanceSpec": {
        "user": "laborum reprehenderit",
        "password": "pariatur",
        "address": "127.0.0.1:9123",
        "labels": {
          "instance": "postgres-abcd"
        }
      }
    }
  })
  cluster_id = ""
  kube_type  = 3
}

#Integration Center: Redis
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationRedis" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "redis-exporter"
  content = jsonencode({
    "name": "redis-exporter-example",
    "kind": "redis-exporter",
    "spec": {
      "instanceSpec": {
        "address": "127.0.0.1:9123",
        "password": "ea sed quis id",
        "labels": {
          "instance": "crs-abcd"
        }
      }
    }
  })
  cluster_id = ""
  kube_type  = 3
}

#Integration Center: Scrape Job
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegrationRaw" {
  instance_id = tencentcloud_monitor_tmp_instance.tmpInstance.id 
  kind        = "raw-job"
  content = jsonencode({
    "kind": "raw-job",
    "spec": {
      "job": <<-EOT
        job_name: example-raw-job-name
        metrics_path: /metrics
        static_configs:
          - targets:
            - 127.0.0.1:9090
      EOT
    }
  })
  cluster_id = ""
  kube_type  = 3
}