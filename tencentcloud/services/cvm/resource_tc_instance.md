Provides a CVM instance resource.

~> **NOTE:** You can launch an CVM instance for a VPC network via specifying parameter `vpc_id`. One instance can only belong to one VPC.

~> **NOTE:** At present, `PREPAID` instance cannot be deleted directly and must wait it to be outdated and released automatically.

~> **NOTE:** Currently, the `placement_group_id` field only supports setting and modification, but not deletion.

~> **NOTE:** When creating a CVM instance using a `launch_template_id`, if you set other parameter values ​​at the same time, the template definition values ​​will be overwritten.

~> **NOTE:** It is recommended to use resource `tencentcloud_eip` to create a AntiDDos Eip, and then call resource `tencentcloud_eip_association` to bind it to resource `tencentcloud_instance`.

~> **NOTE:** When creating a prepaid CVM instance and binding a data disk, you need to explicitly set `delete_with_instance` to `false`.

Example Usage

Create a general POSTPAID_BY_HOUR CVM instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

data "tencentcloud_images" "images" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "OpenCloudOS Server"
}

data "tencentcloud_instance_types" "types" {
  filter {
    name   = "instance-family"
    values = ["S1", "S2", "S3", "S4", "S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
}

// create CVM instance
resource "tencentcloud_instance" "example" {
  instance_name        = "tf-example"
  availability_zone    = var.availability_zone
  image_id             = data.tencentcloud_images.images.images.0.image_id
  instance_type        = data.tencentcloud_instance_types.types.instance_types.0.instance_type
  system_disk_type     = "CLOUD_PREMIUM"
  system_disk_size     = 50
  hostname             = "user"
  project_id           = 0
  vpc_id               = tencentcloud_vpc.vpc.id
  subnet_id            = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

Create a general PREPAID CVM instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

data "tencentcloud_images" "images" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "OpenCloudOS Server"
}

data "tencentcloud_instance_types" "types" {
  filter {
    name   = "instance-family"
    values = ["S1", "S2", "S3", "S4", "S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
}

// create CVM instance
resource "tencentcloud_instance" "example" {
  instance_name                           = "tf-example"
  availability_zone                       = var.availability_zone
  image_id                                = data.tencentcloud_images.images.images.0.image_id
  instance_type                           = data.tencentcloud_instance_types.types.instance_types.0.instance_type
  system_disk_type                        = "CLOUD_PREMIUM"
  system_disk_size                        = 50
  hostname                                = "user"
  project_id                              = 0
  vpc_id                                  = tencentcloud_vpc.vpc.id
  subnet_id                               = tencentcloud_subnet.subnet.id
  instance_charge_type                    = "PREPAID"
  instance_charge_type_prepaid_period     = 1
  instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  force_delete                            = true
  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }

  timeouts {
    create = "30m"
  }
}
```

Create a dedicated cluster CVM instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

data "tencentcloud_images" "images" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "OpenCloudOS Server"
}

data "tencentcloud_instance_types" "types" {
  filter {
    name   = "instance-family"
    values = ["S1", "S2", "S3", "S4", "S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
  cdc_id            = "cluster-262n63e8"
  is_multicast      = false
}

// create CVM instance
resource "tencentcloud_instance" "example" {
  instance_name        = "tf-example"
  availability_zone    = var.availability_zone
  image_id             = data.tencentcloud_images.images.images.0.image_id
  instance_type        = data.tencentcloud_instance_types.types.instance_types.0.instance_type
  dedicated_cluster_id = "cluster-262n63e8"
  instance_charge_type = "CDCPAID"
  system_disk_type     = "CLOUD_SSD"
  system_disk_size     = 50
  hostname             = "user"
  project_id           = 0
  vpc_id               = tencentcloud_vpc.vpc.id
  subnet_id            = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_SSD"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

Create CVM instance with placement_group_id

```hcl
resource "tencentcloud_instance" "example" {
  instance_name                    = "tf-example"
  availability_zone                = "ap-guangzhou-6"
  image_id                         = "img-eb30mz89"
  instance_type                    = "S5.MEDIUM4"
  system_disk_size                 = 50
  system_disk_name                 = "sys_disk_1"
  hostname                         = "user"
  project_id                       = 0
  vpc_id                           = "vpc-i5yyodl9"
  subnet_id                        = "subnet-hhi88a58"
  placement_group_id               = "ps-ejt4brtz"
  force_replace_placement_group_id = false

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 100
    encrypt        = false
    data_disk_name = "data_disk_1"
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

Create CVM instance with template

```hcl
resource "tencentcloud_instance" "example" {
  launch_template_id      = "lt-b20scl2a"
  launch_template_version = 1
}
```

Create CVM instance with AntiDDos Eip

```hcl
resource "tencentcloud_instance" "example" {
  instance_name              = "tf-example"
  availability_zone          = "ap-guangzhou-6"
  image_id                   = "img-eb30mz89"
  instance_type              = "S5.MEDIUM4"
  system_disk_type           = "CLOUD_HSSD"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = "vpc-i5yyodl9"
  subnet_id                  = "subnet-hhi88a58"
  orderly_security_groups    = ["sg-l222vn6w"]
  internet_charge_type       = "BANDWIDTH_PACKAGE"
  bandwidth_package_id       = "bwp-rp2nx3ab"
  ipv4_address_type          = "AntiDDoSEIP"
  anti_ddos_package_id       = "bgp-31400fvq"
  allocate_public_ip         = true
  internet_max_bandwidth_out = 100
  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 100
    encrypt        = false
  }
  tags = {
    tagKey = "tagValue"
  }
}
```

Create CVM instance with setting running flag

```hcl
resource "tencentcloud_instance" "example" {
  instance_name                           = "tf-example"
  availability_zone                       = "ap-guangzhou-6"
  image_id                                = "img-eb30mz89"
  instance_type                           = "S5.MEDIUM4"
  system_disk_type                        = "CLOUD_HSSD"
  system_disk_size                        = 50
  hostname                                = "user"
  project_id                              = 0
  vpc_id                                  = "vpc-i5yyodl9"
  subnet_id                               = "subnet-hhi88a58"
  orderly_security_groups                 = ["sg-ma82yjwp"]
  running_flag                            = false
  stop_type                               = "SOFT_FIRST"
  stopped_mode                            = "KEEP_CHARGING"
  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 100
    encrypt        = false
  }
  tags = {
    tagKey = "tagValue"
  }
}
```

Import

CVM instance can be imported using the id, e.g.

```
terraform import tencentcloud_instance.example ins-2qol3a80
```
