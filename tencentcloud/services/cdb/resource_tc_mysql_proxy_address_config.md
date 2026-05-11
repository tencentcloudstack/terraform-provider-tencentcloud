Provides a resource to manage MySQL database proxy address configuration.

Example Usage

```hcl
resource "tencentcloud_mysql_proxy_address_config" "example" {
  instance_id       = "cdb-o2t7gmjl"
  proxy_group_id    = "proxy-ov7dqp8n"
  proxy_address_id  = "proxyaddr-y8dnlfs0"
  weight_mode       = "system"
  is_kick_out       = true
  min_count         = 0
  max_delay         = 10
  fail_over         = true
  auto_add_ro       = true
  read_only         = false
  trans_split       = false
  connection_pool   = true
  auto_load_balance = true
  access_mode       = "nearby"
  proxy_allocation {
    region = "ap-guangzhou"
    zone   = "ap-guangzhou-6"

    proxy_instance {
      instance_id = "cdb-o2t7gmjl"
      weight      = 0
    }
  }

  proxy_allocation {
    region = "ap-guangzhou"
    zone   = "ap-guangzhou-7"

    proxy_instance {
      instance_id = "cdb-o2t7gmjl"
      weight      = 0
    }
  }
}
```

Import

MySQL proxy address config can be imported using the instanceId#proxyGroupId#proxyAddressId, e.g.

```
terraform import tencentcloud_mysql_proxy_address_config.example cdb-o2t7gmjl#proxy-ov7dqp8n#proxyaddr-y8dnlfs0
```
