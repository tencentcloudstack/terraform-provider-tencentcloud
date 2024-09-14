Provides a resource to create a CCN Route table associate instance config.

Example Usage

```hcl
variable "region" {
  default = "ap-guangzhou"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
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

# create ccn
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  tags = {
    createBy = "terraform"
  }
}

# create ccn route table
resource "tencentcloud_ccn_route_table" "example" {
  ccn_id      = tencentcloud_ccn.example.id
  name        = "tf-example"
  description = "desc."
}

# attachment instance
resource "tencentcloud_ccn_attachment" "attachment" {
  ccn_id          = tencentcloud_ccn.example.id
  instance_id     = tencentcloud_vpc.vpc.id
  instance_type   = "VPC"
  instance_region = var.region
  route_table_id  = tencentcloud_ccn_route_table.example.id
}

# route table associate instance
resource "tencentcloud_ccn_route_table_associate_instance_config" "example" {
  ccn_id         = tencentcloud_ccn.example.id
  route_table_id = tencentcloud_ccn_route_table.example.id
  instances {
    instance_id   = tencentcloud_vpc.vpc.id
    instance_type = "VPC"
  }

  depends_on = [tencentcloud_ccn_attachment.attachment]
}
```

Import

Ccn instance can be imported, e.g.

```
$ terraform import tencentcloud_ccn_route_table_associate_instance_config.example ccn-gr7nynbd#ccnrtb-jpf7bzn3
```
