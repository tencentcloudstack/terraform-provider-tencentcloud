Provides a resource to create a monitor tmp_grafana_config

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

resource "tencentcloud_monitor_grafana_instance" "foo" {
  instance_name         = "tf-grafana"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_ids            = [tencentcloud_subnet.subnet.id]
  grafana_init_password = "1234567890"
  enable_internet       = false
  is_destroy            = true

  tags = {
    "createdBy" = "test"
  }
}

resource "tencentcloud_monitor_tmp_grafana_config" "foo" {
  config = jsonencode(
    {
      server = {
        http_port           = 8080
        root_url            = "https://cloud-grafana.woa.com/grafana-ffrdnrfa/"
        serve_from_sub_path = true
      }
    }
  )
  instance_id = tencentcloud_monitor_grafana_instance.foo.id
}
```

Import

monitor tmp_grafana_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_grafana_config.tmp_grafana_config tmp_grafana_config_id
```