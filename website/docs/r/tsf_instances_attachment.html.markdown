---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_instances_attachment"
sidebar_current: "docs-tencentcloud-resource-tsf_instances_attachment"
description: |-
  Provides a resource to create a tsf instances_attachment
---

# tencentcloud_tsf_instances_attachment

Provides a resource to create a tsf instances_attachment

## Example Usage

```hcl
resource "tencentcloud_tsf_instances_attachment" "instances_attachment" {
  cluster_id           = "cluster-123456"
  instance_id_list     = [""]
  os_name              = "Ubuntu 20.04"
  image_id             = "img-123456"
  password             = "MyP@ssw0rd"
  key_id               = "key-123456"
  sg_id                = "sg-123456"
  instance_import_mode = "R"
  os_customize_type    = "my_customize"
  feature_id_list      =
  instance_advanced_settings {
    mount_target      = "/mnt/data"
    docker_graph_path = "/var/lib/docker"
  }
  security_group_ids = [""]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `instance_id` - (Required, String, ForceNew) Cloud server ID.
* `feature_id_list` - (Optional, Set: [`String`], ForceNew) Image feature ID list.
* `image_id` - (Optional, String, ForceNew) Operating system image ID.
* `instance_advanced_settings` - (Optional, List, ForceNew) Additional instance parameter information.
* `instance_import_mode` - (Optional, String, ForceNew) Cloud server import mode, required for virtual machine clusters, not required for container clusters. R: Reinstall TSF system image, M: Manual installation of agent.
* `key_id` - (Optional, String, ForceNew) Associated key for system reinstallation.
* `os_customize_type` - (Optional, String, ForceNew) Image customization type.
* `os_name` - (Optional, String, ForceNew) Operating system name.
* `password` - (Optional, String, ForceNew) Reset system password.
* `security_group_ids` - (Optional, Set: [`String`], ForceNew) Security group.
* `sg_id` - (Optional, String, ForceNew) Security group setting.

The `instance_advanced_settings` object supports the following:

* `docker_graph_path` - (Required, String) Dockerd --graph specifies the value, default is /var/lib/docker Note: This field may return null, indicating that no valid values can be obtained.
* `mount_target` - (Required, String) Data disk mount point, data disks are not mounted by default. Data disks with formatted ext3, ext4, xfs file systems will be mounted directly, other file systems or unformatted data disks will be automatically formatted as ext4 and mounted. Please back up your data! This setting does not take effect for cloud servers with no data disks or multiple data disks. Note: This field may return null, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



