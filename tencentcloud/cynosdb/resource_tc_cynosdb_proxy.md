Provides a resource to create a cynosdb proxy

Example Usage

```hcl
resource "tencentcloud_cynosdb_proxy" "proxy" {
  cluster_id               = "cynosdbmysql-bws8h88b"
  cpu                      = 2
  mem                      = 4000
  unique_vpc_id            = "vpc-k1t8ickr"
  unique_subnet_id         = "subnet-jdi5xn22"
  connection_pool_type     = "SessionConnectionPool"
  open_connection_pool     = "yes"
  connection_pool_time_out = 30
  security_group_ids       = ["sg-baxfiao5"]
  description              = "desc sample"
  proxy_zones {
    proxy_node_zone  = "ap-guangzhou-7"
    proxy_node_count = 2
  }
}
```