Provides a resource to create a tke kubernetes cluster master attachment

Example Usage

```hcl
resource "tencentcloud_kubernetes_cluster_master_attachment" "example" {
  cluster_id                  = "cls-fp5o961e"
  instance_id                 = "ins-7d6tpbyg"
  node_role                   = "MASTER_ETCD"
  enhanced_security_service   = true
  enhanced_monitor_service    = true
  enhanced_automation_service = true
  password                    = "Password@123"
  security_group_ids          = ["sg-hjs685q9"]

  master_config {
    mount_target      = "/var/data"
    docker_graph_path = "/var/lib/containerd"
    unschedulable     = 0
    labels {
      name  = "key"
      value = "value"
    }

    data_disk {
      file_system           = "ext4"
      auto_format_and_mount = true
      mount_target          = "/var/data"
      disk_partition        = "/dev/vdb"
    }

    extra_args {
      kubelet = ["root-dir=/root"]
    }

    taints {
      key    = "key"
      value  = "value"
      effect = "NoSchedule"
    }
  }
}
```
