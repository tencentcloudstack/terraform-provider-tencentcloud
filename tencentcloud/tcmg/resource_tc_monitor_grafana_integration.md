Provides a resource to create a monitor grafanaIntegration

Example Usage

Create a grafan instance and integrate the configuration

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

resource "tencentcloud_monitor_grafana_integration" "grafanaIntegration" {
  instance_id = tencentcloud_monitor_grafana_instance.foo.id
  kind        = "tencentcloud-monitor-app"
  content     = "{\"kind\":\"tencentcloud-monitor-app\",\"spec\":{\"dataSourceSpec\":{\"authProvider\":{\"__anyOf\":\"使用密钥\",\"useRole\":true,\"secretId\":\"arunma@tencent.com\",\"secretKey\":\"12345678\"},\"name\":\"uint-test\"},\"grafanaSpec\":{\"organizationIds\":[]}}}"
}
```