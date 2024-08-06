Provides a CVM instance resource.

~> **NOTE:** You can launch an CVM instance for a VPC network via specifying parameter `vpc_id`. One instance can only belong to one VPC.

~> **NOTE:** At present, 'PREPAID' instance cannot be deleted directly and must wait it to be outdated and released automatically.

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

Import

CVM instance can be imported using the id, e.g.

```
terraform import tencentcloud_instance.example ins-2qol3a80
```
