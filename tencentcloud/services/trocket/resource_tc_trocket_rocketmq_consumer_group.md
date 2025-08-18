Provides a resource to create a TROCKET rocketmq consumer group

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
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

// create rocketmq instance
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-example"
  instance_type = "BASIC"
  sku_code      = "basic_2k"
  remark        = "remark."
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  tags = {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}

// create consumer group
resource "tencentcloud_trocket_rocketmq_consumer_group" "example" {
  instance_id             = tencentcloud_trocket_rocketmq_instance.example.id
  consumer_group          = "tf-example"
  max_retry_times         = 20
  consume_enable          = false
  consume_message_orderly = true
  remark                  = "remark."
  tags = {
    createBy = "Terraform"
  }
}
```

Import

TROCKET rocketmq consumer group can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_consumer_group.example rmq-1n58qbwg3#tf-example
```
