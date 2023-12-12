Provides a resource to create a sqlserver business_intelligence_file

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_business_intelligence_instance" "example" {
  zone                = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory              = 4
  storage             = 100
  cpu                 = 2
  machine_type        = "CLOUD_PREMIUM"
  project_id          = 0
  subnet_id           = tencentcloud_subnet.subnet.id
  vpc_id              = tencentcloud_vpc.vpc.id
  db_version          = "201603"
  security_group_list  = [tencentcloud_security_group.security_group.id]
  weekly              = [1, 2, 3, 4, 5, 6, 7]
  start_time          = "00:00"
  span                = 6
  instance_name       = "tf_example"
}

resource "tencentcloud_sqlserver_business_intelligence_file" "example" {
  instance_id = tencentcloud_sqlserver_business_intelligence_instance.example.id
  file_url    = "https://tf-example-1208515315.cos.ap-guangzhou.myqcloud.com/sqlserver_business_intelligence_file.txt"
  file_type   = "FLAT"
  remark      = "desc."
}
```

Import

sqlserver business_intelligence_file can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_business_intelligence_file.example mssqlbi-fo2dwujt#test.xlsx
```