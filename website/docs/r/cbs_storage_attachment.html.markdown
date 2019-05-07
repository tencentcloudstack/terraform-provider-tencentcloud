---
layout: "tencentcloud"
page_title: "tencentcloud: tencentcloud_cbs_storage_attachment"
sidebar_current: "docs-tencentcloud-resource-cbs-storage-attachment"
description: |-
  Provides an tencentcloud CBS stoarge attachment as a resource, to attach and detach storage from CVM Instances.
---

# tencentcloud_cbs_storage_attachment

Provides CBS stoarge attachment resource.

## Example Usage

```hcl
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"

  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }

  cpu_core_count = 2
  memory_size    = 4
}

data "tencentcloud_availability_zones" "my_favorate_zones" {}

resource "tencentcloud_instance" "instance-without-specified-image-id-example" {
  instance_name     = "my-instance"
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type     = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
}

resource "tencentcloud_cbs_storage" "my-storage" {
  storage_type      = "cloudBasic"
  storage_size      = 10
  period            = 1
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  storage_name      = "my-storage"
}

resource "tencentcloud_cbs_storage_attachment" "my-attachment" {
  storage_id  = "${tencentcloud_cbs_storage.my-storage.id}"
  instance_id = "${tencentcloud_instance.instance-without-specified-image-id-example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `storage_id` - (Required, Forces new resource) ID of the storage to be attached.
* `instance_id` - (Required, Forces new resource) ID of the CVM instance to attache to.


## Attributes Reference

The following attributes are exported:

* `storage_id` - ID of the storage.
* `instance_id` - ID of the CVM instance.

