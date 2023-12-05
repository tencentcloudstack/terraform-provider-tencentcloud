Provides an elasticsearch instance resource.

Example Usage

Create a basic version of elasticsearch instance paid by the hour

```hcl
data "tencentcloud_availability_zones_by_product" "availability_zone" {
  product = "es"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_es_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_availability_zones_by_product.availability_zone.zones.0.name
  name              = "tf_es_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_elasticsearch_instance" "example" {
  instance_name       = "tf_example_es"
  availability_zone   = data.tencentcloud_availability_zones_by_product.availability_zone.zones.0.name
  version             = "7.10.1"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  password            = "Test12345"
  license_type        = "basic"
  basic_security_type = 2

  web_node_type_info {
    node_num  = 1
    node_type = "ES.S1.MEDIUM4"
  }

  node_info_list {
    node_num  = 2
    node_type = "ES.S1.MEDIUM8"
    encrypt   = false
  }

  es_acl {
    # black_list = [
    #   "9.9.9.9",
    #   "8.8.8.8",
    # ]
    white_list = [
      "127.0.0.1",
    ]
  }

  tags = {
    test = "test"
  }
}
```

Create a basic version of elasticsearch instance for multi-availability zone deployment

```hcl
data "tencentcloud_availability_zones_by_product" "availability_zone" {
  product = "es"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_es_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_availability_zones_by_product.availability_zone.zones.0.name
  name              = "tf_es_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_subnet" "subnet_multi_zone" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_availability_zones_by_product.availability_zone.zones.1.name
  name              = "tf_es_subnet"
  cidr_block        = "10.0.2.0/24"
}

resource "tencentcloud_elasticsearch_instance" "example_multi_zone" {
  instance_name       = "tf_example_es"
  availability_zone   = "-"
  version             = "7.10.1"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = "-"
  password            = "Test12345"
  license_type        = "basic"
  basic_security_type = 2
  deploy_mode         = 1

  multi_zone_infos {
    availability_zone = data.tencentcloud_availability_zones_by_product.availability_zone.zones.0.name
    subnet_id = tencentcloud_subnet.subnet.id
  }

  multi_zone_infos {
    availability_zone = data.tencentcloud_availability_zones_by_product.availability_zone.zones.1.name
    subnet_id = tencentcloud_subnet.subnet_multi_zone.id
  }

  web_node_type_info {
    node_num  = 1
    node_type = "ES.S1.MEDIUM4"
  }

  node_info_list {
    type = "dedicatedMaster"
    node_num  = 3
    node_type = "ES.S1.MEDIUM8"
    encrypt   = false
  }

  node_info_list {
    type = "hotData"
    node_num  = 2
    node_type = "ES.S1.MEDIUM8"
    encrypt   = false
  }

  es_acl {
    # black_list = [
    #   "9.9.9.9",
    #   "8.8.8.8",
    # ]
    white_list = [
      "127.0.0.1",
    ]
  }

  tags = {
    test = "test"
  }
}
```

Import

Elasticsearch instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_elasticsearch_instance.foo es-17634f05
```