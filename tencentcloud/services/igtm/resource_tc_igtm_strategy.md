Provides a resource to create a IGTM strategy

Example Usage

```hcl
resource "tencentcloud_igtm_instance" "example" {
  domain            = "domain.com"
  access_type       = "CUSTOM"
  global_ttl        = 60
  package_type      = "STANDARD"
  instance_name     = "tf-example"
  access_domain     = "domain.com"
  access_sub_domain = "subDomain.com"
  remark            = "remark."
  resource_id       = "ins-lnpnnwvwqmr"
}

resource "tencentcloud_igtm_address_pool" "example1" {
  pool_name        = "tf-example1"
  traffic_strategy = "WEIGHT"
  address_set {
    addr      = "1.1.1.1"
    is_enable = "ENABLED"
    weight    = 90
  }

  address_set {
    addr      = "2.2.2.2"
    is_enable = "DISABLED"
    weight    = 50
  }
}

resource "tencentcloud_igtm_address_pool" "example2" {
  pool_name        = "tf-example2"
  traffic_strategy = "WEIGHT"
  address_set {
    addr      = "3.3.3.3"
    is_enable = "ENABLED"
    weight    = 90
  }

  address_set {
    addr      = "4.4.4.4"
    is_enable = "DISABLED"
    weight    = 50
  }
}

resource "tencentcloud_igtm_address_pool" "example3" {
  pool_name        = "tf-example3"
  traffic_strategy = "WEIGHT"
  address_set {
    addr      = "5.5.5.5"
    is_enable = "ENABLED"
    weight    = 90
  }

  address_set {
    addr      = "6.6.6.6"
    is_enable = "DISABLED"
    weight    = 50
  }
}

resource "tencentcloud_igtm_strategy" "example" {
  instance_id   = tencentcloud_igtm_instance.example.id
  strategy_name = "tf-example"
  source {
    dns_line_id = 1
    name        = "默认"
  }

  source {
    dns_line_id = 858
    name        = "电信"
  }

  source {
    dns_line_id = 859
    name        = "联通"
  }

  source {
    dns_line_id = 860
    name        = "移动"
  }

  main_address_pool_set {
    address_pools {
      pool_id = tencentcloud_igtm_address_pool.example1.pool_id
      weight  = 90
    }

    address_pools {
      pool_id = tencentcloud_igtm_address_pool.example2.pool_id
      weight  = 80
    }

    min_survive_num  = 1
    traffic_strategy = "WEIGHT"
  }

  fallback_address_pool_set {
    address_pools {
      pool_id = tencentcloud_igtm_address_pool.example3.pool_id
    }
  }

  keep_domain_records = "DISABLED"
  switch_pool_type    = "AUTO"
}
```

Import

IGTM strategy can be imported using the instanceId#strategyId, e.g.

```
terraform import tencentcloud_igtm_strategy.example gtm-uukztqtoaru#7556
```
