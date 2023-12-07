Provides a resource to create a cynosdb read_only_instance_exclusive_access

Example Usage

```hcl
variable "cynosdb_cluster_id" {
  default = "default_cynosdb_cluster"
}
variable "cynosdb_cluster_instance_id" {
  default = "default_cluster_instance"
}

variable "cynosdb_cluster_security_group_id" {
  default = "default_security_group_id"
}

data "tencentcloud_vpc_subnets" "gz3" {
  availability_zone = var.default_az
  is_default        = true
}

locals {
  vpc_id    = data.tencentcloud_vpc_subnets.gz3.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.subnet_id
}

resource "tencentcloud_cynosdb_read_only_instance_exclusive_access" "read_only_instance_exclusive_access" {
  cluster_id         = var.cynosdb_cluster_id
  instance_id        = var.cynosdb_cluster_instance_id
  vpc_id             = local.vpc_id
  subnet_id          = local.subnet_id
  port               = 1234
  security_group_ids = [var.cynosdb_cluster_security_group_id]
}
```