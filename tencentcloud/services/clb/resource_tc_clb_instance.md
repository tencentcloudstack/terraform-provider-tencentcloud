Provides a resource to create a CLB instance.

Example Usage

Create INTERNAL CLB

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
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
  is_multicast      = false
}

// create clb
resource "tencentcloud_clb_instance" "example" {
  network_type = "INTERNAL"
  clb_name     = "tf-example"
  project_id   = 0
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id

  tags = {
    tagKey = "tagValue"
  }
}
```

Create CLB with eip_address_id, Only support INTERNAL CLB

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
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
  is_multicast      = false
}

// create clb
resource "tencentcloud_clb_instance" "example" {
  network_type   = "INTERNAL"
  clb_name       = "tf-example"
  project_id     = 0
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
  eip_address_id = "eip-lt0w6jhq"

  tags = {
    tagKey = "tagValue"
  }
}
```

Create dedicated cluster clb

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
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
  cdc_id            = "cluster-lchwgxhs"
  is_multicast      = false
}

// create clb
resource "tencentcloud_clb_instance" "example" {
  network_type = "INTERNAL"
  clb_name     = "tf-example"
  project_id   = 0
  cluster_id   = "cluster-lchwgxhs"
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id

  tags = {
    tagKey = "tagValue"
  }
}
```

Create LCU-supported CLB

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
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
  is_multicast      = false
}

// create clb
resource "tencentcloud_clb_instance" "example" {
  network_type = "INTERNAL"
  clb_name     = "tf-example"
  project_id   = 0
  sla_type     = "clb.c3.medium"
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id

  tags = {
    tagKey = "tagValue"
  }
}
```

Create OPEN CLB

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0
  
  tags = {
    "example" = "test"
  }
}

// create clb
resource "tencentcloud_clb_instance" "example" {
  network_type    = "OPEN"
  clb_name        = "tf-example"
  project_id      = 0
  vpc_id          = tencentcloud_vpc.vpc.id
  security_groups = [tencentcloud_security_group.example.id]

  tags = {
    tagKey = "tagValue"
  }
}
```

Support CORS

```hcl
variable "zone" {
  default = "ap-guangzhou"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

// create clb
resource "tencentcloud_clb_instance" "example" {
  network_type              = "OPEN"
  clb_name                  = "tf-example"
  project_id                = 0
  vpc_id                    = tencentcloud_vpc.vpc.id
  security_groups           = [tencentcloud_security_group.example.id]
  target_region_info_region = var.zone
  target_region_info_vpc_id = tencentcloud_vpc.vpc.id

  tags = {
    tagKey = "tagValue"
  }
}
```

Open CLB with VipIsp

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create vpc bandwidth package
resource "tencentcloud_vpc_bandwidth_package" "example" {
  network_type           = "SINGLEISP_CMCC"
  charge_type            = "ENHANCED95_POSTPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  internet_max_bandwidth = 300
  egress                 = "center_egress1"

  tags = {
    createdBy = "terraform"
  }
}

// create clb
resource "tencentcloud_clb_instance" "example" {
  network_type         = "OPEN"
  clb_name             = "tf-example"
  project_id           = 0
  vip_isp              = "CMCC"
  internet_charge_type = "BANDWIDTH_PACKAGE"
  bandwidth_package_id = tencentcloud_vpc_bandwidth_package.example.id
  vpc_id               = tencentcloud_vpc.vpc.id

  tags = {
    tagKey = "tagValue"
  }
}
```

Dynamic Vip Instance

```hcl
variable "zone" {
  default = "ap-guangzhou"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
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
  is_multicast      = false
}

// create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

// create clb
resource "tencentcloud_clb_instance" "example" {
  network_type              = "OPEN"
  clb_name                  = "tf-example"
  project_id                = 0
  vpc_id                    = tencentcloud_vpc.vpc.id
  target_region_info_region = var.zone
  target_region_info_vpc_id = tencentcloud_vpc.vpc.id
  security_groups           = [tencentcloud_security_group.example.id]
  dynamic_vip               = true

  tags = {
    tagKey = "tagValue"
  }
}

output "domain" {
  value = tencentcloud_clb_instance.example.domain
}
```

Specified Vip Instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

// create clb
resource "tencentcloud_clb_instance" "example" {
  network_type    = "OPEN"
  clb_name        = "tf-example"
  project_id      = 0
  vpc_id          = tencentcloud_vpc.vpc.id
  security_groups = [tencentcloud_security_group.example.id]
  vip             = "111.230.4.204"

  tags = {
    tagKey = "tagValue"
  }
}

output "domain" {
  value = tencentcloud_clb_instance.example.domain
}
```

Default enable

```hcl
variable "zone" {
  default = "ap-guangzhou"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
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
  is_multicast      = false
}

// create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

// create clb
resource "tencentcloud_clb_instance" "example" {
  network_type                 = "OPEN"
  clb_name                     = "tf-example"
  project_id                   = 0
  load_balancer_pass_to_target = true
  vpc_id                       = tencentcloud_vpc.vpc.id
  security_groups              = [tencentcloud_security_group.example.id]
  target_region_info_vpc_id    = tencentcloud_vpc.vpc.id
  target_region_info_region    = var.zone

  tags = {
    tagKey = "tagValue"
  }
}
```

Create multiple instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_clb_instance" "example" {
  network_type   = "OPEN"
  clb_name       = "tf-example"
  master_zone_id = var.availability_zone
}
```

Create instance with log

```hcl
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
  is_multicast      = false
}

// create route table
resource "tencentcloud_route_table" "route" {
  name   = "route_table"
  vpc_id = tencentcloud_vpc.vpc.id
}

// create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_clb_log_set" "log" {
  period = 7
}

// create topic
resource "tencentcloud_clb_log_topic" "topic" {
  log_set_id = tencentcloud_clb_log_set.log.id
  topic_name = "clb-topic"
}

// create clb
resource "tencentcloud_clb_instance" "example" {
  network_type                 = "INTERNAL"
  clb_name                     = "tf-example"
  project_id                   = 0
  load_balancer_pass_to_target = true
  vpc_id                       = tencentcloud_vpc.vpc.id
  subnet_id                    = tencentcloud_subnet.subnet.id
  security_groups              = [tencentcloud_security_group.example.id]
  log_set_id                   = tencentcloud_clb_log_set.log.id
  log_topic_id                 = tencentcloud_clb_log_topic.topic.id

  tags = {
    tagKey = "tagValue"
  }
}
```

Create instance with associate endpoint

```hcl
resource "tencentcloud_clb_instance" "example" {
  network_type       = "OPEN"
  clb_name           = "tf-example"
  project_id         = 0
  vpc_id             = "vpc-e51ilko8"
  associate_endpoint = "vpce-du9ssd3z"
  tags = {
    createBy = "Terraform"
  }
}
```

Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_instance.example lb-7a0t6zqb
```