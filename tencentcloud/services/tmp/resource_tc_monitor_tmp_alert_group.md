Provides a resource to create a monitor tmp_alert_group

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
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

resource "tencentcloud_monitor_tmp_instance" "example" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 30
  zone                = var.availability_zone
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_monitor_tmp_alert_group" "example" {
  group_name      = "tf-example"
  instance_id     = tencentcloud_monitor_tmp_instance.example.id
  repeat_interval = "5m"

  custom_receiver {
    type = "amp"
  }

  rules {
    duration  = "1m"
    expr      = "up{job=\"prometheus-agent\"} != 1"
    rule_name = "Agent health check"
    state     = 2

    annotations = {
      "summary"     = "Agent health check"
      "description" = "Agent {{$labels.instance}} is deactivated, please pay attention!"
    }

    labels = {
      "severity" = "critical"
    }
  }
}
```

Import

monitor tmp_alert_group can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_alert_group.example prom-34qkzwvs#alert-rfkkr6cw
```
