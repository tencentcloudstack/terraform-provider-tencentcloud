Provides a resource to create a cynosdb proxy_end_point

Example Usage

```hcl
resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id       = "cynosdbmysql-bws8h88b"
  unique_vpc_id    = "vpc-4owdpnwr"
  unique_subnet_id = "subnet-dwj7ipnc"
  instance_weights {
    instance_id = "cynosdbmysql-ins-afqx1hy0"
    weight      = 1
  }
}
```

Set `vip` and `vport`

```hcl
resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id       = "cynosdbmysql-bws8h88b"
  unique_vpc_id    = "vpc-4owdpnwr"
  unique_subnet_id = "subnet-dwj7ipnc"
  vip              = "172.16.112.108"
  vport            = "3306"
  instance_weights {
    instance_id = "cynosdbmysql-ins-afqx1hy0"
    weight      = 1
  }
}
```

Open connection pool

```hcl
resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id               = "cynosdbmysql-bws8h88b"
  unique_vpc_id            = "vpc-4owdpnwr"
  unique_subnet_id         = "subnet-dwj7ipnc"
  vip                      = "172.16.112.108"
  vport                    = "3306"
  connection_pool_type     = "SessionConnectionPool"
  open_connection_pool     = "yes"
  connection_pool_time_out = 30
  instance_weights {
    instance_id = "cynosdbmysql-ins-afqx1hy0"
    weight      = 1
  }
}
```

Close connection pool

```hcl
resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id               = "cynosdbmysql-bws8h88b"
  unique_vpc_id            = "vpc-4owdpnwr"
  unique_subnet_id         = "subnet-dwj7ipnc"
  vip                      = "172.16.112.108"
  vport                    = "3306"
  open_connection_pool     = "no"
  instance_weights {
    instance_id = "cynosdbmysql-ins-afqx1hy0"
    weight      = 1
  }
}
```

if `rw_type` is `READWRITE`

```hcl
resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id               = "cynosdbmysql-bws8h88b"
  unique_vpc_id            = "vpc-4owdpnwr"
  unique_subnet_id         = "subnet-dwj7ipnc"
  vip                      = "172.16.112.108"
  vport                    = "3306"
  open_connection_pool     = "no"
  fail_over                = "yes"
  consistency_type         = "global"
  rw_type                  = "READWRITE"
  consistency_time_out     = 30
  instance_weights {
    instance_id = "cynosdbmysql-ins-afqx1hy0"
    weight      = 1
  }
}
```

if `rw_type` is `READONLY`

```hcl
resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id               = "cynosdbmysql-bws8h88b"
  unique_vpc_id            = "vpc-4owdpnwr"
  unique_subnet_id         = "subnet-dwj7ipnc"
  vip                      = "172.16.112.108"
  vport                    = "3306"
  open_connection_pool     = "no"
  rw_type                  = "READONLY"
  instance_weights {
    instance_id = "cynosdbmysql-ins-rikr6z4o"
    weight      = 1
  }
}
```

Comprehensive parameter examples

```hcl
resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id               = "cynosdbmysql-bws8h88b"
  unique_vpc_id            = "vpc-4owdpnwr"
  unique_subnet_id         = "subnet-dwj7ipnc"
  vip                      = "172.16.112.118"
  vport                    = "3306"
  connection_pool_type     = "SessionConnectionPool"
  open_connection_pool     = "yes"
  connection_pool_time_out = 30
  security_group_ids       = ["sg-7kpsbxdb"]
  description              = "desc value"
  weight_mode              = "system"
  auto_add_ro              = "yes"
  fail_over                = "yes"
  consistency_type         = "global"
  rw_type                  = "READWRITE"
  consistency_time_out     = 30
  trans_split              = true
  access_mode              = "nearby"
  instance_weights {
    instance_id = "cynosdbmysql-ins-afqx1hy0"
    weight      = 1
  }
}
```