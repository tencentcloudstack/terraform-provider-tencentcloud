Provides a resource to create a THPC workspaces

-> **Note:** If space_charge_type is UNDERWRITE, Not currently supported for deletion.

Example Usage

Create a PREPAID THPC workspaces

```hcl
variable "availability_zone" {
  default = "ap-nanjing-1"
}

data "tencentcloud_images" "images" {
  image_type = ["PUBLIC_IMAGE"]
  os_name = "TencentOS Server 3.1 (TK4) UEFI"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "172.16.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  is_multicast      = false
}

# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "security group desc."

  tags = {
    "createBy" = "Terraform"
  }
}

# create thpc workspaces
resource "tencentcloud_thpc_workspaces" "example" {
  space_name        = "tf-example"
  space_charge_type = "PREPAID"
  space_type        = "96A.96XLARGE2304"
  hpc_cluster_id    = "hpc-l9anqcbl"
  image_id          = data.tencentcloud_images.images.images.0.image_id
  security_group_ids = [tencentcloud_security_group.example.id]
  placement {
    zone       = var.availability_zone
    project_id = 0
  }

  space_charge_prepaid {
    period     = 1
    renew_flag = "NOTIFY_AND_AUTO_RENEW"
  }

  system_disk {
    disk_size = 100
    disk_type = "CLOUD_HSSD"
  }

  data_disks {
    disk_size = 200
    disk_type = "CLOUD_HSSD"
    encrypt   = false
  }

  virtual_private_cloud {
    vpc_id             = tencentcloud_vpc.vpc.id
    subnet_id          = tencentcloud_subnet.subnet.id
    as_vpc_gateway     = false
    ipv6_address_count = 0
  }

  internet_accessible {
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 200
    public_ip_assigned         = true
  }

  login_settings {
    password = "Password@123"
  }

  enhanced_service {
    security_service {
      enabled = true
    }

    monitor_service {
      enabled = true
    }

    automation_service {
      enabled = true
    }
  }
}
```

Create a UNDERWRITE THPC workspaces

```hcl
variable "availability_zone" {
  default = "ap-nanjing-1"
}

data "tencentcloud_images" "images" {
  image_type = ["PUBLIC_IMAGE"]
  os_name = "TencentOS Server 3.1 (TK4) UEFI"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "172.16.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  is_multicast      = false
}

# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "security group desc."

  tags = {
    "createBy" = "Terraform"
  }
}

# create thpc workspaces
resource "tencentcloud_thpc_workspaces" "example" {
  space_name        = "tf-example"
  space_charge_type = "UNDERWRITE"
  space_type        = "96A.96XLARGE2304"
  hpc_cluster_id    = "hpc-l9anqcbl"
  image_id          = data.tencentcloud_images.images.images.0.image_id
  security_group_ids = [tencentcloud_security_group.example.id]
  placement {
    zone       = var.availability_zone
    project_id = 0
  }

  space_charge_prepaid {
    period     = 12
    renew_flag = "NOTIFY_AND_AUTO_RENEW"
  }

  system_disk {
    disk_size = 100
    disk_type = "CLOUD_HSSD"
  }

  data_disks {
    disk_size = 200
    disk_type = "CLOUD_HSSD"
    encrypt   = false
  }

  virtual_private_cloud {
    vpc_id             = tencentcloud_vpc.vpc.id
    subnet_id          = tencentcloud_subnet.subnet.id
    as_vpc_gateway     = false
    ipv6_address_count = 0
    private_ip_addresses = ["172.16.0.2"]
  }

  internet_accessible {
    internet_charge_type       = "BANDWIDTH_PREPAID"
    internet_max_bandwidth_out = 200
    public_ip_assigned         = true
  }

  login_settings {
    key_ids = ["skey-qxjpz7uj"]
  }

  enhanced_service {
    security_service {
      enabled = true
    }

    monitor_service {
      enabled = true
    }

    automation_service {
      enabled = true
    }
  }
}
```

Import

THPC workspaces can be imported using the id, e.g.

```
terraform import tencentcloud_thpc_workspaces.example wks-gwg3ygz1
```
