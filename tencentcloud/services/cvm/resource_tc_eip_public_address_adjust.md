Provides a resource to create a eip public_address_adjust

~> **NOTE:** This interface is used to change the IP address. It supports changing the common public IP of the CVM instance and the EIP of the monthly bandwidth. `address_id` and `instance_id` cannot exist at the same time. When `address_id` is passed, only the EIP of the monthly bandwidth is supported.

Example Usage

```hcl
# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create cvm
resource "tencentcloud_instance" "example" {
  instance_name              = "tf_example"
  availability_zone          = "ap-guangzhou-6"
  image_id                   = "img-9qrfy1xt"
  instance_type              = "SA3.MEDIUM4"
  system_disk_type           = "CLOUD_HSSD"
  system_disk_size           = 100
  hostname                   = "example"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.vpc.id
  subnet_id                  = tencentcloud_subnet.subnet.id
  allocate_public_ip         = true
  internet_max_bandwidth_out = 10

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

# create eip
resource "tencentcloud_eip" "example" {
  name = "tf-example"
}

resource "tencentcloud_eip_public_address_adjust" "example" {
  instance_id = tencentcloud_instance.example.id
  address_id  = tencentcloud_eip.example.id
}
```