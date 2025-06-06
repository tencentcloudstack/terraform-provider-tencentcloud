Provides a resource to create a MQTT instance

Example Usage

Create a POSTPAID instance

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
  instance_type = "PRO"
  name          = "tf-example"
  sku_code      = "pro_6k_1"
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
```

Create a PREPAID instance

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
  instance_type = "PRO"
  name          = "tf-example"
  sku_code      = "pro_10k_2"
  remark        = "remarks."
  vpc_list {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
  }
  pay_mode             = 1
  time_span            = 1
  renew_flag           = 1
  force_delete         = false
  automatic_activation = true
  authorization_policy = true
  tags = {
    createBy = "Terraform"
  }
}
```
