Provides a resource to create a Trocket rocketmq instance

~> **NOTE:** It only supports create postpaid rocketmq 5.x instance.

Example Usage

Create Basic Instance

```hcl
# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create rocketmq instance
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-example"
  instance_type = "PRO"
  sku_code      = "pro_4k"
  remark        = "remark"
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  tags = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}
```

Create Enable Public Network Instance

```hcl
# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create rocketmq instance
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-example"
  instance_type = "PRO"
  sku_code      = "pro_4k"
  remark        = "remark"
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  enable_public = true
  bandwidth     = 10
  ip_rules {
    ip     = "1.1.1.1"
    allow  = true
    remark = "remark message."
  }

  ip_rules {
    ip     = "2.2.2.2"
    allow  = false
    remark = "remark message."
  }

  tags = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}
```

Import

Trocket rocketmq instance can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_instance.example rmq-n5qado7m
```
