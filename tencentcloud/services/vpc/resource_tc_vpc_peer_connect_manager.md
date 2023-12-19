Provides a resource to create a vpc peer_connect_manager

Example Usage

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  owner_uin = data.tencentcloud_user_info.info.owner_uin
}

resource "tencentcloud_vpc" "vpc" {
  name       = "tf-example-pcx"
  cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_vpc" "des_vpc" {
  name       = "tf-example-pcx-des"
  cidr_block = "172.16.0.0/16"
}
resource "tencentcloud_vpc_peer_connect_manager" "peer_connect_manager" {
  source_vpc_id = tencentcloud_vpc.vpc.id
  peering_connection_name = "example-iac"
  destination_vpc_id = tencentcloud_vpc.des_vpc.id
  destination_uin = local.owner_uin
  destination_region = "ap-guangzhou"
}
```

Import

vpc peer_connect_manager can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_peer_connect_manager.peer_connect_manager peer_connect_manager_id
```
