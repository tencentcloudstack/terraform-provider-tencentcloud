Provides a resource to create a CLB instance.

Example Usage

INTERNAL CLB

```hcl
resource "tencentcloud_clb_instance" "internal_clb" {
  network_type = "INTERNAL"
  clb_name     = "myclb"
  project_id   = 0
  vpc_id       = "vpc-7007ll7q"
  subnet_id    = "subnet-12rastkr"

  tags = {
    test = "tf"
  }
}
```

LCU-supported CLB

```hcl
resource "tencentcloud_clb_instance" "internal_clb" {
  network_type = "INTERNAL"
  clb_name     = "myclb"
  project_id   = 0
  sla_type     = "clb.c3.medium"
  vpc_id       = "vpc-2hfyray3"
  subnet_id    = "subnet-o3a5nt20"

  tags = {
    test = "tf"
  }
}
```

OPEN CLB

```hcl
resource "tencentcloud_clb_instance" "open_clb" {
  network_type              = "OPEN"
  clb_name                  = "myclb"
  project_id                = 0
  vpc_id                    = "vpc-da7ffa61"
  security_groups           = ["sg-o0ek7r93"]
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = "vpc-da7ffa61"

  tags = {
    test = "tf"
  }
}
```

OPNE CLB with VipIsp

```hcl
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

resource "tencentcloud_clb_instance" "open_clb" {
  network_type         = "OPEN"
  clb_name             = "my-open-clb"
  project_id           = 0
  vpc_id               = "vpc-4owdpnwr"
  vip_isp              = "CMCC"
  internet_charge_type = "BANDWIDTH_PACKAGE"
  bandwidth_package_id = tencentcloud_vpc_bandwidth_package.example.id

  tags = {
    test = "open"
  }
}
```

Dynamic Vip Instance

```hcl
resource "tencentcloud_security_group" "foo" {
  name = "clb-instance-open-sg"
}

resource "tencentcloud_vpc" "foo" {
  name       = "clb-instance-open-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_clb_instance" "clb_open" {
  network_type              = "OPEN"
  clb_name                  = "clb-instance-open"
  project_id                = 0
  vpc_id                    = tencentcloud_vpc.foo.id
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = tencentcloud_vpc.foo.id
  security_groups           = [tencentcloud_security_group.foo.id]

  dynamic_vip = true

  tags = {
    test = "tf"
  }
}

output "domain" {
  value = tencentcloud_clb_instance.clb_open.domain
}
```

Default enable

```hcl
resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-1"
  name              = "sdk-feature-test"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_security_group" "sglab" {
  name        = "sg_o0ek7r93"
  description = "favourite sg"
  project_id  = 0
}

resource "tencentcloud_vpc" "foo" {
  name         = "for-my-open-clb"
  cidr_block   = "10.0.0.0/16"

  tags = {
    "test" = "mytest"
  }
}

resource "tencentcloud_clb_instance" "open_clb" {
  network_type                 = "OPEN"
  clb_name                     = "my-open-clb"
  project_id                   = 0
  vpc_id                       = tencentcloud_vpc.foo.id
  load_balancer_pass_to_target = true

  security_groups              = [tencentcloud_security_group.sglab.id]
  target_region_info_region    = "ap-guangzhou"
  target_region_info_vpc_id    = tencentcloud_vpc.foo.id

  tags = {
    test = "open"
  }
}
```

CREATE multiple instance

```hcl
resource "tencentcloud_clb_instance" "open_clb1" {
  network_type              = "OPEN"
  clb_name = "hello"
  master_zone_id = "ap-guangzhou-3"
}
```

CREATE instance with log
```hcl
resource "tencentcloud_vpc" "vpc_test" {
  name = "clb-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "rtb_test" {
  name = "clb-test"
  vpc_id = "${tencentcloud_vpc.vpc_test.id}"
}

resource "tencentcloud_subnet" "subnet_test" {
  name = "clb-test"
  cidr_block = "10.0.1.0/24"
  availability_zone = "ap-guangzhou-3"
  vpc_id = "${tencentcloud_vpc.vpc_test.id}"
  route_table_id = "${tencentcloud_route_table.rtb_test.id}"
}

resource "tencentcloud_clb_log_set" "set" {
  period = 7
}

resource "tencentcloud_clb_log_topic" "topic" {
  log_set_id = "${tencentcloud_clb_log_set.set.id}"
  topic_name = "clb-topic"
}

resource "tencentcloud_clb_instance" "internal_clb" {
  network_type = "INTERNAL"
  clb_name = "myclb"
  project_id = 0
  vpc_id = "${tencentcloud_vpc.vpc_test.id}"
  subnet_id = "${tencentcloud_subnet.subnet_test.id}"
  load_balancer_pass_to_target = true
  log_set_id = "${tencentcloud_clb_log_set.set.id}"
  log_topic_id = "${tencentcloud_clb_log_topic.topic.id}"

  tags = {
    test = "tf"
  }
}

```

Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_instance.foo lb-7a0t6zqb
```