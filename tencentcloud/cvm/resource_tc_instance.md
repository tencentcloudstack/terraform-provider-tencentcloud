Provides a CVM instance resource.

~> **NOTE:** You can launch an CVM instance for a VPC network via specifying parameter `vpc_id`. One instance can only belong to one VPC.

~> **NOTE:** At present, 'PREPAID' instance cannot be deleted directly and must wait it to be outdated and released automatically.

Example Usage

```hcl
data "tencentcloud_images" "my_favorite_image" {
  image_type = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

data "tencentcloud_instance_types" "my_favorite_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S1", "S2", "S3", "S4", "S5"]
  }

  cpu_core_count = 2
  exclude_sold_out = true
}

data "tencentcloud_availability_zones" "my_favorite_zones" {
}

// Create VPC resource
resource "tencentcloud_vpc" "app" {
  cidr_block = "10.0.0.0/16"
  name       = "awesome_app_vpc"
}

resource "tencentcloud_subnet" "app" {
  vpc_id            = tencentcloud_vpc.app.id
  availability_zone = data.tencentcloud_availability_zones.my_favorite_zones.zones.0.name
  name              = "awesome_app_subnet"
  cidr_block        = "10.0.1.0/24"
}

// Create a POSTPAID_BY_HOUR CVM instance
resource "tencentcloud_instance" "cvm_postpaid" {
  instance_name              = "cvm_postpaid"
  availability_zone          = data.tencentcloud_availability_zones.my_favorite_zones.zones.0.name
  image_id                   = data.tencentcloud_images.my_favorite_image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.app.id
  subnet_id                  = tencentcloud_subnet.app.id

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
	encrypt = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

// Create a PREPAID CVM instance
resource "tencentcloud_instance" "cvm_prepaid" {
  timeouts {
    create = "30m"
  }
  instance_name                           = "cvm_prepaid"
  availability_zone                       = data.tencentcloud_availability_zones.my_favorite_zones.zones.0.name
  image_id                                = data.tencentcloud_images.my_favorite_image.images.0.image_id
  instance_type                           = data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type
  system_disk_type                        = "CLOUD_PREMIUM"
  system_disk_size                        = 50
  hostname                                = "user"
  project_id                              = 0
  vpc_id                                  = tencentcloud_vpc.app.id
  subnet_id                               = tencentcloud_subnet.app.id
  instance_charge_type                    = "PREPAID"
  instance_charge_type_prepaid_period     = 1
  instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt        = false
  }
  force_delete = true
  tags         = {
    tagKey = "tagValue"
  }
}
```

Import

CVM instance can be imported using the id, e.g.

```
terraform import tencentcloud_instance.foo ins-2qol3a80
```