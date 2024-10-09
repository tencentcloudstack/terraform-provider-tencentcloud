Provides a resource to create a CVM instance action timer

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

data "tencentcloud_images" "images" {
  image_type = ["PUBLIC_IMAGE"]
  image_name_regex = "TencentOS Server"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create cvm
resource "tencentcloud_instance" "example" {
  instance_name     = "tf_example"
  availability_zone = var.availability_zone
  image_id          = data.tencentcloud_images.images.images.0.image_id
  instance_type     = "SA3.MEDIUM4"
  system_disk_type  = "CLOUD_HSSD"
  system_disk_size  = 100
  hostname          = "example"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    createBy = "terraform"
  }
}

# create cvm action timer
resource "tencentcloud_cvm_action_timer" "example" {
  instance_id = tencentcloud_instance.example.id

  action_timer {
    timer_action = "TerminateInstances"
    action_time  = "2024-11-11T11:26:40Z"
  }
}
```