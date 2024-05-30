Provides a resource to create a cvm renew_instance

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
  instance_name                           = "tf_example"
  availability_zone                       = "ap-guangzhou-6"
  image_id                                = "img-9qrfy1xt"
  instance_type                           = "SA3.MEDIUM4"
  system_disk_type                        = "CLOUD_HSSD"
  system_disk_size                        = 100
  hostname                                = "example"
  project_id                              = 0
  vpc_id                                  = tencentcloud_vpc.vpc.id
  subnet_id                               = tencentcloud_subnet.subnet.id
  force_delete                            = true
  instance_charge_type                    = "PREPAID"
  instance_charge_type_prepaid_period     = 1
  instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

# renew instance
resource "tencentcloud_cvm_renew_instance" "example" {
  instance_id              = tencentcloud_instance.example.id
  renew_portable_data_disk = true

  instance_charge_prepaid {
    period     = 1
    renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  }
}
```