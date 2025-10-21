---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_scaling_group_status"
sidebar_current: "docs-tencentcloud-resource-as_scaling_group_status"
description: |-
  Provides a resource to set as scaling_group status
---

# tencentcloud_as_scaling_group_status

Provides a resource to set as scaling_group status

## Example Usage

### Deactivate Scaling Group

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "as"
}

data "tencentcloud_images" "image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = data.tencentcloud_images.image.images.0.image_id
  instance_types     = ["SA1.SMALL1", "SA2.SMALL1", "SA2.SMALL2", "SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name"
  }
}

resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name = "tf-example"
  configuration_id   = tencentcloud_as_scaling_config.example.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_as_scaling_group_status" "scaling_group_status" {
  auto_scaling_group_id = tencentcloud_as_scaling_group.example.id
  enable                = false
}
```

### Enable Scaling Group

```hcl
resource "tencentcloud_as_scaling_group_status" "scaling_group_status" {
  auto_scaling_group_id = tencentcloud_as_scaling_group.example.id
  enable                = true
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_group_id` - (Required, String, ForceNew) Scaling group ID.
* `enable` - (Required, Bool) If enable auto scaling group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

as scaling_group_status can be imported using the id, e.g.

```
terraform import tencentcloud_as_scaling_group_status.scaling_group_status scaling_group_id
```

