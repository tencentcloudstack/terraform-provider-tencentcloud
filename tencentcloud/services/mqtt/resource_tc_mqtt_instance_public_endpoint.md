Provides a resource to create a MQTT instance public endpoint

~> **NOTE:** This resource must exclusive in one MQTT instance, do not declare additional public endpoint resources of this instance elsewhere.

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

// create mqtt instance
resource "tencentcloud_mqtt_instance" "example" {
  instance_type = "BASIC"
  name          = "tf-example"
  sku_code      = "basic_2k"
  remark        = "remarks."
  vpc_list {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
  }
  pay_mode = 0
  tags = {
    createBy = "Terraform"
  }
}

// create public endpoint
resource "tencentcloud_mqtt_instance_public_endpoint" "example" {
  instance_id = tencentcloud_mqtt_instance.example.id
  bandwidth   = 100
  rules {
    ip_rule = "192.168.1.0/24"
    remark  = "Remark."
  }

  rules {
    ip_rule = "172.16.1.0/24"
    remark  = "Remark."
  }
}
```

Import

MQTT instance public endpoint can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_instance_public_endpoint.example mqtt-emb2v5wk
```
