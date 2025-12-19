---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instance"
sidebar_current: "docs-tencentcloud-resource-clb_instance"
description: |-
  Provides a resource to create a CLB instance.
---

# tencentcloud_clb_instance

Provides a resource to create a CLB instance.

## Example Usage

### Create INTERNAL CLB

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

### Create CLB with eip_address_id, Only support INTERNAL CLB

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

### Create dedicated cluster clb

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

### Create LCU-supported CLB

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

### Create OPEN CLB

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

### Support CORS

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

### Open CLB with VipIsp

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

### Dynamic Vip Instance

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

### Specified Vip Instance

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

### Default enable

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

### Create multiple instance

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

### Create instance with log

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

### Create instance with associate endpoint

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

## Argument Reference

The following arguments are supported:

* `clb_name` - (Required, String) Name of the CLB. The name can only contain Chinese characters, English letters, numbers, underscore and hyphen '-'.
* `network_type` - (Required, String, ForceNew) Type of CLB instance. Valid values: `OPEN` and `INTERNAL`.
* `address_ip_version` - (Optional, String) It's only applicable to public network CLB instances. IP version. Values: `IPV4`, `IPV6` and `IPv6FullChain` (case-insensitive). Default: `IPV4`. Note: IPV6 indicates IPv6 NAT64, while IPv6FullChain indicates IPv6.
* `associate_endpoint` - (Optional, String) The associated terminal node ID; passing an empty string indicates unassociating the node.
* `bandwidth_package_id` - (Optional, String) Bandwidth package id. If set, the `internet_charge_type` must be `BANDWIDTH_PACKAGE`.
* `cluster_id` - (Optional, String, ForceNew) Cluster ID.
* `delete_protect` - (Optional, Bool) Whether to enable delete protection.
* `dynamic_vip` - (Optional, Bool) If create dynamic vip CLB instance, `true` or `false`.
* `eip_address_id` - (Optional, String) The unique ID of the EIP, such as eip-1v2rmbwk, is only applicable to the intranet load balancing binding EIP. During the EIP change, there may be a brief network interruption.
* `internet_bandwidth_max_out` - (Optional, Int) Max bandwidth out, only applicable to open CLB. Valid value ranges is [1, 2048]. Unit is Mbps.
* `internet_charge_type` - (Optional, String) Internet charge type, only applicable to open CLB. Valid values are `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
* `load_balancer_pass_to_target` - (Optional, Bool) Whether the target allow flow come from clb. If value is true, only check security group of clb, or check both clb and backend instance security group.
* `log_set_id` - (Optional, String) The id of log set.
* `log_topic_id` - (Optional, String) The id of log topic.
* `master_zone_id` - (Optional, String) Setting master zone id of cross available zone disaster recovery, only applicable to open CLB.
* `project_id` - (Optional, Int) ID of the project within the CLB instance, `0` - Default Project.
* `security_groups` - (Optional, List: [`String`]) Security groups of the CLB instance. Supports both `OPEN` and `INTERNAL` CLBs.
* `sla_type` - (Optional, String) This parameter is required to create LCU-supported instances. Values:`SLA`: Super Large 4. When you have activated Super Large models, `SLA` refers to Super Large 4; `clb.c2.medium`: Standard; `clb.c3.small`: Advanced 1; `clb.c3.medium`: Advanced 1; `clb.c4.small`: Super Large 1; `clb.c4.medium`: Super Large 2; `clb.c4.large`: Super Large 3; `clb.c4.xlarge`: Super Large 4. For more details, see [Instance Specifications](https://intl.cloud.tencent.com/document/product/214/84689?from_cn_redirect=1).
* `slave_zone_id` - (Optional, String) Setting slave zone id of cross available zone disaster recovery, only applicable to open CLB. this zone will undertake traffic when the master is down.
* `snat_ips` - (Optional, List) Snat Ip List, required with `snat_pro=true`. NOTE: This argument cannot be read and modified here because dynamic ip is untraceable, please import resource `tencentcloud_clb_snat_ip` to handle fixed ips.
* `snat_pro` - (Optional, Bool) Indicates whether Binding IPs of other VPCs feature switch.
* `subnet_id` - (Optional, String, ForceNew) In the case of purchasing a `INTERNAL` clb instance, the subnet id must be specified. The VIP of the `INTERNAL` clb instance will be generated from this subnet.
* `tags` - (Optional, Map) The available tags within this CLB.
* `target_region_info_region` - (Optional, String) Region information of backend services are attached the CLB instance. Only supports `OPEN` CLBs.
* `target_region_info_vpc_id` - (Optional, String) Vpc information of backend services are attached the CLB instance. Only supports `OPEN` CLBs.
* `vip_isp` - (Optional, String, ForceNew) Network operator, only applicable to open CLB. Valid values are `CMCC`(China Mobile), `CTCC`(Telecom), `CUCC`(China Unicom) and `BGP`. If this ISP is specified, network billing method can only use the bandwidth package billing (BANDWIDTH_PACKAGE).
* `vip` - (Optional, String, ForceNew) Specifies the VIP for the application of a CLB instance. This parameter is optional. If you do not specify this parameter, the system automatically assigns a value for the parameter. IPv4 and IPv6 CLB instances support this parameter, but IPv6 NAT64 CLB instances do not.
* `vpc_id` - (Optional, String, ForceNew) VPC ID of the CLB.
* `zone_id` - (Optional, String) Available zone id, only applicable to open CLB.

The `snat_ips` object supports the following:

* `subnet_id` - (Required, String) Snat subnet ID.
* `ip` - (Optional, String) Snat IP address, If set to empty will auto allocated.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `address_ipv6` - The IPv6 address of the load balancing instance.
* `clb_vips` - The virtual service address table of the CLB.
* `domain` - Domain name of the CLB instance.
* `ipv6_mode` - This field is meaningful when the IP address version is ipv6, `IPv6Nat64` | `IPv6FullChain`.


## Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_instance.example lb-7a0t6zqb
```

