Provides a resource to create a monitor grafana ssoAccount

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_monitor_grafana_instance" "foo" {
  instance_name         = "test-grafana"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_ids            = [tencentcloud_subnet.subnet.id]
  grafana_init_password = "1234567890"
  enable_internet = false

  tags = {
    "createdBy" = "test"
  }
}

resource "tencentcloud_monitor_grafana_sso_account" "ssoAccount" {
  instance_id = tencentcloud_monitor_grafana_instance.foo.id
  user_id     = "111"
  notes       = "desc12222"
  role {
    organization  = "Main Org."
    role          = "Admin"
  }
}

```
Import

monitor grafana ssoAccount can be imported using the instance_id#user_id, e.g.
```
$ terraform import tencentcloud_monitor_grafana_sso_account.ssoAccount grafana-50nj6v00#111
```