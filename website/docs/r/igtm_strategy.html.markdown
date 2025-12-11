---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_strategy"
sidebar_current: "docs-tencentcloud-resource-igtm_strategy"
description: |-
  Provides a resource to create a IGTM strategy
---

# tencentcloud_igtm_strategy

Provides a resource to create a IGTM strategy

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `fallback_address_pool_set` - (Required, List) Fallback address pool set, only one level allowed and address pool count must be 1.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `main_address_pool_set` - (Required, List) Main address pool set, up to four levels allowed.
* `source` - (Required, List) Resolution lines.
* `strategy_name` - (Required, String) Strategy name, cannot be duplicated.
* `keep_domain_records` - (Optional, String) Whether to enable policy forced retention of default lines disabled, enabled, default is disabled and only one policy can be enabled.
* `switch_pool_type` - (Optional, String) Policy scheduling mode: AUTO default switching; STOP only pause without switching.

The `address_pools` object of `fallback_address_pool_set` supports the following:

* `pool_id` - (Required, Int) Address pool ID.
* `weight` - (Optional, Int) Weight.

The `address_pools` object of `main_address_pool_set` supports the following:

* `pool_id` - (Required, Int) Address pool ID.
* `weight` - (Optional, Int) Weight.

The `fallback_address_pool_set` object supports the following:

* `address_pools` - (Required, List) Address pool IDs and weights in the set, array.
* `main_address_pool_id` - (Optional, Int) Address pool set ID.
* `min_survive_num` - (Optional, Int) Switch threshold, cannot exceed the total number of addresses in the main set.
* `traffic_strategy` - (Optional, String) Switch strategy: ALL resolves all addresses; WEIGHT: load balancing. When ALL, the weight value of resolved addresses is 1; when WEIGHT, weight is address pool weight * address weight.

The `main_address_pool_set` object supports the following:

* `address_pools` - (Required, List) Address pool IDs and weights in the set, array.
* `main_address_pool_id` - (Optional, Int) Address pool set ID.
* `min_survive_num` - (Optional, Int) Switch threshold, cannot exceed the total number of addresses in the main set.
* `traffic_strategy` - (Optional, String) Switch strategy: ALL resolves all addresses; WEIGHT: load balancing. When ALL, the weight value of resolved addresses is 1; when WEIGHT, weight is address pool weight * address weight.

The `source` object supports the following:

* `dns_line_id` - (Required, Int) Resolution request source line ID.
* `name` - (Optional, String) Resolution request source line name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `strategy_id` - Strategy ID.


## Import

IGTM strategy can be imported using the instanceId#strategyId, e.g.

```
terraform import tencentcloud_igtm_strategy.example gtm-uukztqtoaru#7556
```

