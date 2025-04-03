Provides a resource to create a gwlb gwlb_target_group_register_instances

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_gwlb_instance" "gwlb_instance" {
  vpc_id = tencentcloud_vpc.vpc.id
  subnet_id = tencentcloud_subnet.subnet.id
  load_balancer_name = "tf-test"
  lb_charge_type = "POSTPAID_BY_HOUR"
  tags {
    tag_key = "test_key"
    tag_value = "tag_value"
  }
}

resource "tencentcloud_gwlb_target_group" "gwlb_target_group" {
  target_group_name = "tf-test"
  vpc_id = tencentcloud_vpc.vpc.id
  port = 6081
  health_check {
    health_switch = true
    protocol = "tcp"
    port = 6081
    timeout = 2
    interval_time = 5
    health_num = 3
    un_health_num = 3
  }
}

resource "tencentcloud_instance" "foo" {
  system_disk_type = "CLOUD_PREMIUM"
  instance_name = "tf-test"
  image_id = data.tencentcloud_images.default.images.0.image_id
  instance_type = "S5.MEDIUM2"
  system_disk_size = 100
  subnet_id = tencentcloud_subnet.subnet.id
  vpc_id = tencentcloud_vpc.vpc.id
  hostname = "tf-test"
  disable_security_service = true
  allocate_public_ip = true
  internet_max_bandwidth_out = 5
  availability_zone = var.availability_zone
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_gwlb_target_group_register_instances" "gwlb_target_group_register_instances" {
  target_group_id = tencentcloud_gwlb_target_group.gwlb_target_group.id
  target_group_instances {
   bind_ip = tencentcloud_instance.foo.private_ip
   port = 6081
   weight = 0
  }
}
```

Import

gwlb gwlb_target_group_register_instances can be imported using the id, e.g.

```
terraform import tencentcloud_gwlb_target_group_register_instances.gwlb_target_group_register_instances gwlb_target_group_register_instances_id
```
