---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_proxy_end_point"
sidebar_current: "docs-tencentcloud-resource-cynosdb_proxy_end_point"
description: |-
  Provides a resource to create a cynosdb proxy_end_point
---

# tencentcloud_cynosdb_proxy_end_point

Provides a resource to create a cynosdb proxy_end_point

## Example Usage

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

### Open connection pool

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

### Close connection pool

```hcl
resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id           = "cynosdbmysql-bws8h88b"
  unique_vpc_id        = "vpc-4owdpnwr"
  unique_subnet_id     = "subnet-dwj7ipnc"
  vip                  = "172.16.112.108"
  vport                = "3306"
  open_connection_pool = "no"
  instance_weights {
    instance_id = "cynosdbmysql-ins-afqx1hy0"
    weight      = 1
  }
}
```



```hcl
resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id           = "cynosdbmysql-bws8h88b"
  unique_vpc_id        = "vpc-4owdpnwr"
  unique_subnet_id     = "subnet-dwj7ipnc"
  vip                  = "172.16.112.108"
  vport                = "3306"
  open_connection_pool = "no"
  fail_over            = "yes"
  consistency_type     = "global"
  rw_type              = "READWRITE"
  consistency_time_out = 30
  instance_weights {
    instance_id = "cynosdbmysql-ins-afqx1hy0"
    weight      = 1
  }
}
```



```hcl
resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id           = "cynosdbmysql-bws8h88b"
  unique_vpc_id        = "vpc-4owdpnwr"
  unique_subnet_id     = "subnet-dwj7ipnc"
  vip                  = "172.16.112.108"
  vport                = "3306"
  open_connection_pool = "no"
  rw_type              = "READONLY"
  instance_weights {
    instance_id = "cynosdbmysql-ins-rikr6z4o"
    weight      = 1
  }
}
```

### Comprehensive parameter examples

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

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `unique_subnet_id` - (Required, String) The private network subnet ID is consistent with the cluster subnet ID by default.
* `unique_vpc_id` - (Required, String) Private network ID, which is consistent with the cluster private network ID by default.
* `access_mode` - (Optional, String) Connection mode: nearby, balance.
* `auto_add_ro` - (Optional, String) Do you want to automatically add read-only instances? Yes - Yes, no - Do not automatically add.
* `connection_pool_time_out` - (Optional, Int) Connection pool threshold: unit (second).
* `connection_pool_type` - (Optional, String) Connection pool type: SessionConnectionPool (session level Connection pool).
* `consistency_time_out` - (Optional, Int) Consistency timeout.
* `consistency_type` - (Optional, String) Consistency type: event, global, session.
* `description` - (Optional, String) Description.
* `fail_over` - (Optional, String) Enable Failover. yes or no.
* `instance_weights` - (Optional, List) Instance Weight.
* `open_connection_pool` - (Optional, String) Whether to enable Connection pool, yes - enable, no - do not enable.
* `rw_type` - (Optional, String) Read and write attributes: READWRITE, READONLY.
* `security_group_ids` - (Optional, Set: [`String`]) Security Group ID Array.
* `trans_split` - (Optional, Bool) Transaction splitting.
* `vip` - (Optional, String) VIP Information.
* `vport` - (Optional, Int) Port Information.
* `weight_mode` - (Optional, String) Weight mode: system system allocation, custom customization.

The `instance_weights` object supports the following:

* `instance_id` - (Required, String) Instance Id.
* `weight` - (Required, Int) Instance Weight.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_group_id` - Instance Group ID.
* `proxy_group_id` - Proxy Group ID.


