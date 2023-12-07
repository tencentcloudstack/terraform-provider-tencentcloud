Provides a resource to create a sqlserver general_cloud_ro_instance

Example Usage

If read_only_group_type value is 1 - Ship according to one instance and one read-only group:

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

resource "tencentcloud_sqlserver_general_cloud_instance" "example" {
  name                 = "tf_example"
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_HSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  db_version           = "2008R2"
  security_group_list  = [tencentcloud_security_group.security_group.id]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
}

resource "tencentcloud_sqlserver_general_cloud_ro_instance" "example" {
  instance_id          = tencentcloud_sqlserver_general_cloud_instance.example.id
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  read_only_group_type = 1
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_BSSD"
  instance_charge_type = "POSTPAID"
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  security_group_list  = [tencentcloud_security_group.security_group.id]
  collation            = "Chinese_PRC_CI_AS"
  time_zone            = "China Standard Time"
  resource_tags        = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
}
```

If read_only_group_type value is 2 - Ship after creating a read-only group, all instances are under this read-only group:

```hcl
resource "tencentcloud_sqlserver_general_cloud_ro_instance" "example" {
  instance_id                      = tencentcloud_sqlserver_general_cloud_instance.example.id
  zone                             = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  read_only_group_type             = 2
  read_only_group_name             = "test-ro-group"
  read_only_group_is_offline_delay = 1
  read_only_group_max_delay_time   = 10
  read_only_group_min_in_group     = 1
  memory                           = 4
  storage                          = 100
  cpu                              = 2
  machine_type                     = "CLOUD_BSSD"
  instance_charge_type             = "POSTPAID"
  subnet_id                        = tencentcloud_subnet.subnet.id
  vpc_id                           = tencentcloud_vpc.vpc.id
  security_group_list              = [tencentcloud_security_group.security_group.id]
  collation                        = "Chinese_PRC_CI_AS"
  time_zone                        = "China Standard Time"
  resource_tags                    = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
}
```

If read_only_group_type value is 3 - All instances shipped are in the existing Some read-only groups below:

```hcl
resource "tencentcloud_sqlserver_general_cloud_ro_instance" "example" {
  instance_id          = tencentcloud_sqlserver_general_cloud_instance.example.id
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  read_only_group_type = 3
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_BSSD"
  read_only_group_id   = "mssqlrg-clboghrj"
  instance_charge_type = "POSTPAID"
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  security_group_list  = [tencentcloud_security_group.security_group.id]
  collation            = "Chinese_PRC_CI_AS"
  time_zone            = "China Standard Time"
  resource_tags        = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
}
```