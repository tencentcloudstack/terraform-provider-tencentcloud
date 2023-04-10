---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_param"
sidebar_current: "docs-tencentcloud-resource-redis_param"
description: |-
  Provides a resource to create a redis param
---

# tencentcloud_redis_param

Provides a resource to create a redis param

## Example Usage

```hcl
resource "tencentcloud_redis_param" "param" {
  instance_id = "crs-c1nl9rpv"
  instance_params = {
    "cluster-node-timeout"          = "15000"
    "disable-command-list"          = "\"\""
    "hash-max-ziplist-entries"      = "512"
    "hash-max-ziplist-value"        = "64"
    "hz"                            = "10"
    "lazyfree-lazy-eviction"        = "yes"
    "lazyfree-lazy-expire"          = "yes"
    "lazyfree-lazy-server-del"      = "yes"
    "maxmemory-policy"              = "noeviction"
    "notify-keyspace-events"        = "\"\""
    "proxy-slowlog-log-slower-than" = "500"
    "replica-lazy-flush"            = "yes"
    "sentineauth"                   = "no"
    "set-max-intset-entries"        = "512"
    "slowlog-log-slower-than"       = "10"
    "timeout"                       = "31536000"
    "zset-max-ziplist-entries"      = "128"
    "zset-max-ziplist-value"        = "64"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of instance.
* `instance_params` - (Required, Map) A list of parameters modified by the instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis param can be imported using the instanceId, e.g.

```
terraform import tencentcloud_redis_param.param crs-c1nl9rpv
```

