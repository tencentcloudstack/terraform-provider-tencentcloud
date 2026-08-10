Use this resource to create TcaplusDB cluster.

~> **NOTE:** TcaplusDB now only supports the following regions: `ap-shanghai,ap-hongkong,na-siliconvalley,ap-singapore,ap-seoul,ap-tokyo,eu-frankfurt, and na-ashburn`.

Example Usage

Create a new tcaplus cluster instance

```hcl
locals {
  vpc_id    = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_example_tcaplus_cluster"
  vpc_id                   = local.vpc_id
  subnet_id                = local.subnet_id
  password                 = "your_pw_123111"
  old_password_expire_last = 3600
}
```

Create a dedicated tcaplus cluster instance with resource tags, server list and proxy list

```hcl
locals {
  vpc_id    = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_tcaplus_cluster" "dedicated_example" {
  idl_type     = "PROTO"
  cluster_name = "tf_example_dedicated_cluster"
  vpc_id       = local.vpc_id
  subnet_id    = local.subnet_id
  password     = "your_pw_123111"
  cluster_type = 2

  resource_tags {
    tag_key   = "env"
    tag_value = "prod"
  }

  resource_tags {
    tag_key   = "owner"
    tag_value = "terraform"
  }

  server_list {
    machine_type = "S5.LARGE8"
    machine_num = 2
  }

  proxy_list {
    machine_type = "S5.LARGE8"
    machine_num = 1
  }
}
```

Import

tcaplus cluster can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcaplus_cluster.example cluster_id
```