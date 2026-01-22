Provides a resource to create a CLB target group attachment is bound to the load balancing listener or forwarding rule.

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

// create clb instance
resource "tencentcloud_clb_instance" "example" {
  clb_name     = "tf-example"
  network_type = "INTERNAL"
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id
}

// create clb listener
resource "tencentcloud_clb_listener" "example" {
  clb_id        = tencentcloud_clb_instance.example.id
  listener_name = "tf-example"
  port          = 8080
  protocol      = "HTTP"
}

// create clb listener rule
resource "tencentcloud_clb_listener_rule" "example" {
  clb_id              = tencentcloud_clb_instance.example.id
  listener_id         = tencentcloud_clb_listener.example.listener_id
  domain              = "example.com"
  url                 = "/"
  session_expire_time = 60
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}

// create clb target group
resource "tencentcloud_clb_target_group" "example" {
  target_group_name = "tf-example"
  vpc_id            = tencentcloud_vpc.vpc.id
}

// create clb target group attachment
resource "tencentcloud_clb_target_group_attachment" "example" {
  clb_id          = tencentcloud_clb_instance.example.id
  target_group_id = tencentcloud_clb_target_group.example.id
  listener_id     = tencentcloud_clb_listener.example.listener_id
  rule_id         = tencentcloud_clb_listener_rule.example.rule_id
}
```

Import

CLB target group attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group_attachment.example lbtg-odareyb2#lbl-bicjmx3i#lb-cv0iz74c#loc-ac6uk7b6
```