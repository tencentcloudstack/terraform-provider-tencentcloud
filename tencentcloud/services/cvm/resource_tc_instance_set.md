Provides a CVM instance set resource.

~> **NOTE:** You can launch an CVM instance for a VPC network via specifying parameter `vpc_id`. One instance can only belong to one VPC.

~> **NOTE:** This resource is designed to cater for the scenario of creating CVM in large batches.

~> **NOTE:** After run command `terraform apply`, must wait all cvms is ready, then run command `terraform plan`, either it will cause state change.

Example Usage

```hcl
data "tencentcloud_images" "my_favorite_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "Tencent Linux release 3.2 (Final)"
}

data "tencentcloud_instance_types" "my_favorite_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }

  cpu_core_count = 1
  memory_size    = 1
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

// Create 10 CVM instances to host awesome_app
resource "tencentcloud_instance_set" "my_awesome_app" {
  timeouts {
			create = "5m"
			read   = "20s"
			delete = "1h"
  }

  instance_count             = 10
  instance_name              = "awesome_app"
  availability_zone          = data.tencentcloud_availability_zones.my_favorite_zones.zones.0.name
  image_id                   = data.tencentcloud_images.my_favorite_image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.app.id
  subnet_id                  = tencentcloud_subnet.app.id
}
```