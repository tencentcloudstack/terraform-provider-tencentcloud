---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_l7_acc_rule_priority_operation"
sidebar_current: "docs-tencentcloud-resource-teo_l7_acc_rule_priority_operation"
description: |-
  Provides a resource to set TEO l7 acc rules priority
---

# tencentcloud_teo_l7_acc_rule_priority_operation

Provides a resource to set TEO l7 acc rules priority

## Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rule_v2" "rule1" {
  zone_id     = "zone-39quuimqg8r6"
  description = ["1"]
  rule_name   = "网站加速1"
  status      = "enable"
  branches {
    condition = "$${http.request.host} in ['aaa.makn.cn']"
    actions {
      name = "Cache"
      cache_parameters {
        custom_time {
          cache_time           = 2592000
          ignore_cache_control = "off"
          switch               = "on"
        }
      }
    }

    actions {
      name = "CacheKey"
      cache_key_parameters {
        full_url_cache = "on"
        ignore_case    = "off"
        query_string {
          switch = "off"
          values = []
        }
      }
    }

    sub_rules {
      description = ["1-1"]
      branches {
        condition = "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"
        actions {
          name = "Cache"
          cache_parameters {
            no_cache {
              switch = "on"
            }
          }
        }
      }
    }

    sub_rules {
      description = ["1-2"]
      branches {
        condition = "$${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']"
        actions {
          name = "MaxAge"
          max_age_parameters {
            cache_time    = 3600
            follow_origin = "off"
          }
        }
      }
    }
  }
}

resource "tencentcloud_teo_l7_acc_rule_v2" "rule2" {
  zone_id     = "zone-39quuimqg8r6"
  description = ["2"]
  rule_name   = "网站加速2"
  status      = "enable"
  branches {
    condition = "$${http.request.host} in ['aaa.makn.cn']"
    actions {
      name = "Cache"
      cache_parameters {
        custom_time {
          cache_time           = 2592000
          ignore_cache_control = "off"
          switch               = "on"
        }
      }
    }

    actions {
      name = "CacheKey"
      cache_key_parameters {
        full_url_cache = "on"
        ignore_case    = "off"
        query_string {
          switch = "off"
          values = []
        }
      }
    }

    sub_rules {
      description = ["1-1"]
      branches {
        condition = "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"
        actions {
          name = "Cache"
          cache_parameters {
            no_cache {
              switch = "on"
            }
          }
        }
      }
    }

    sub_rules {
      description = ["1-2"]
      branches {
        condition = "$${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']"
        actions {
          name = "MaxAge"
          max_age_parameters {
            cache_time    = 3600
            follow_origin = "off"
          }
        }
      }
    }
  }
}

resource "tencentcloud_teo_l7_acc_rule_priority_operation" "teo_l7_acc_rule_priority_operation" {
  zone_id  = "zone-39quuimqg8r6"
  rule_ids = [resource.tencentcloud_teo_l7_acc_rule_v2.rule2.rule_id, resource.tencentcloud_teo_l7_acc_rule_v2.rule1.rule_id]
}
```

## Argument Reference

The following arguments are supported:

* `rule_ids` - (Required, List: [`String`], ForceNew) Complete list of rule IDs under site ID.
* `zone_id` - (Required, String, ForceNew) Zone id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



