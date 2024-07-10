Provides a resource to create a CLB instance.

Example Usage

INTERNAL CLB

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-4"
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
}

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

LCU-supported CLB

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-4"
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
}

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

OPEN CLB

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

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
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_clb_instance" "example" {
  network_type              = "OPEN"
  clb_name                  = "tf-example"
  project_id                = 0
  target_region_info_region = "ap-guangzhou"
  security_groups           = [tencentcloud_security_group.example.id]
  vpc_id                    = tencentcloud_vpc.vpc.id

  tags = {
    tagKey = "tagValue"
  }
}
```

OPEN CLB with VipIsp

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

resource "tencentcloud_vpc_bandwidth_package" "example" {
  network_type           = "SINGLEISP_CMCC"
  charge_type            = "ENHANCED95_POSTPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  internet_max_bandwidth = 300
  egress                 = "center_egress1"

  tags = {
    "createdBy" = "terraform"
  }
}

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
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_clb_instance" "example" {
  network_type              = "OPEN"
  clb_name                  = "tf-example"
  project_id                = 0
  vpc_id                    = tencentcloud_vpc.vpc.id
  target_region_info_vpc_id = tencentcloud_vpc.vpc.id
  security_groups           = [tencentcloud_security_group.example.id]
  target_region_info_region = "ap-guangzhou"
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
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

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
  value = tencentcloud_clb_instance.example
}
```

Default enable

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_clb_instance" "example" {
  network_type                 = "OPEN"
  clb_name                     = "tf-example"
  project_id                   = 0
  load_balancer_pass_to_target = true
  target_region_info_region    = "ap-guangzhou"
  vpc_id                       = tencentcloud_vpc.vpc.id
  target_region_info_vpc_id    = tencentcloud_vpc.vpc.id
  security_groups              = [tencentcloud_security_group.example.id]

  tags = {
    tagKey = "tagValue"
  }
}
```

CREATE multiple instance

```hcl
resource "tencentcloud_clb_instance" "example" {
  network_type   = "OPEN"
  clb_name       = "tf-example"
  master_zone_id = "ap-guangzhou-3"
}
```

Create instance with log

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

resource "tencentcloud_route_table" "example" {
  name   = "tf-example"
  vpc_id = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  route_table_id    = tencentcloud_route_table.example.id
  availability_zone = "ap-guangzhou-4"
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_clb_log_set" "example" {
  period = 7
}

resource "tencentcloud_clb_log_topic" "example" {
  log_set_id = tencentcloud_clb_log_set.example.id
  topic_name = "clb-topic"
}

resource "tencentcloud_clb_instance" "example" {
  network_type                 = "INTERNAL"
  clb_name                     = "tf-example"
  project_id                   = 0
  vpc_id                       = tencentcloud_vpc.vpc.id
  subnet_id                    = tencentcloud_subnet.subnet.id
  log_set_id                   = tencentcloud_clb_log_set.example.id
  log_topic_id                 = tencentcloud_clb_log_topic.example.id
  load_balancer_pass_to_target = true

  tags = {
    tagKey = "tagValue"
  }
}
```

Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_instance.example lb-7a0t6zqb
```
