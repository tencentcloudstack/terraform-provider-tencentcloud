Provides a resource to create a tse instance

Example Usage

Create zookeeper standard version
```hcl
resource "tencentcloud_tse_instance" "zookeeper_standard" {
  engine_type = "zookeeper"
  engine_version = "3.5.9.4"
  engine_product_version = "STANDARD"
  engine_region = "ap-guangzhou"
  engine_name = "zookeeper-test"
  trade_type = 0
  engine_resource_spec = "spec-qvj6k7t4q"
  engine_node_num = 3
  vpc_id = "vpc-4owdpnwr"
  subnet_id = "subnet-dwj7ipnc"

  tags = {
    "createdBy" = "terraform"
  }
}
```

Create zookeeper professional version
```hcl
resource "tencentcloud_tse_instance" "zookeeper_professional" {
  engine_type = "zookeeper"
  engine_version = "3.5.9.4"
  engine_product_version = "PROFESSIONAL"
  engine_region = "ap-guangzhou"
  engine_name = "zookeeper-test"
  trade_type = 0
  engine_resource_spec = "spec-qvj6k7t4q"
  engine_node_num = 3
  vpc_id = "vpc-4owdpnwr"
  subnet_id = "subnet-dwj7ipnc"

  engine_region_infos {
    engine_region = "ap-guangzhou"
    replica       = 3

    vpc_infos {
        subnet_id = "subnet-dwj7ipnc"
        vpc_id    = "vpc-4owdpnwr"
    }
    vpc_infos {
        subnet_id = "subnet-403mgks4"
        vpc_id    = "vpc-b1puef4z"
    }
  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Create nacos standard version
```hcl
resource "tencentcloud_tse_instance" "nacos" {
    enable_client_internet_access = false
    engine_name                   = "test"
    engine_node_num               = 3
    engine_product_version        = "STANDARD"
    engine_region                 = "ap-guangzhou"
    engine_resource_spec          = "spec-1160a35a"
    engine_type                   = "nacos"
    engine_version                = "2.0.3.4"
    subnet_id                     = "subnet-5vpegquy"
    trade_type                    = 0
    vpc_id                        = "vpc-99xmasf9"

	tags = {
    	"createdBy" = "terraform"
    }
}
```

Create polaris base version
```hcl
resource "tencentcloud_tse_instance" "polaris" {
    enable_client_internet_access = false
    engine_name                   = "test"
    engine_node_num               = 2
    engine_product_version        = "BASE"
    engine_region                 = "ap-guangzhou"
    engine_resource_spec          = "spec-c160bas1"
    engine_type                   = "polaris"
    engine_version                = "1.16.0.1"
    subnet_id                     = "subnet-5vpegquy"
    trade_type                    = 0
    vpc_id                        = "vpc-99xmasf9"
	tags = {
    	"createdBy" = "terraform"
    }
}
```

Import

tse instance can be imported using the id, e.g.

```
terraform import tencentcloud_tse_instance.instance instance_id
```