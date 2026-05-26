---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_l7_acc_rule_v2"
sidebar_current: "docs-tencentcloud-resource-teo_l7_acc_rule_v2"
description: |-
  Provides a resource to create a TEO l7 acc rule v2
---

# tencentcloud_teo_l7_acc_rule_v2

Provides a resource to create a TEO l7 acc rule v2

~> **NOTE:** Compared to tencentcloud_teo_l7_acc_rule, tencentcloud_teo_l7_acc_rule_v2 is simpler to use but is limited to managing a single rule and lacks the ability to maintain rule ordering. It is best suited for scenarios where you need to manage multiple rules independently and priority/sequencing is not a concern.

## Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rule_v2" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  description = ["description"]
  rule_name   = "Web Acceleration"
  status      = "enable"
  branches {
    condition = "$${http.request.host} in ['www.example.com']"
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

    actions {
      name = "ModifyRequestHeader"
      modify_request_header_parameters {
        header_actions {
          action = "set"
          name   = "EO-Client-OS"
          value  = "*"
        }

        header_actions {
          action = "add"
          name   = "O-Client-Browser"
          value  = "*"
        }

        header_actions {
          action = "del"
          name   = "Eo-Client-Device"
        }
      }
    }

    actions {
      name = "ContentCompression"
      content_compression_parameters {
        switch = "on"
      }
    }

    actions {
      name = "Vary"
      vary_parameters {
        switch = "on"
      }
    }

    actions {
      name = "OriginAuthentication"
      origin_authentication_parameters {
        request_properties {
          type  = "Header"
          name  = "Authorization"
          value = "Bearer token123"
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
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Zone id.
* `branches` - (Optional, List) Sub-Rule branch. this list currently supports filling in only one rule; multiple entries are invalid.
* `description` - (Optional, List: [`String`]) Rule annotation. multiple annotations can be added.
* `rule_name` - (Optional, String) Rule name. The name length limit is 255 characters.
* `status` - (Optional, String) Rule status. The possible values are: `enable`: enabled; `disable`: disabled.

The `branches` object supports the following:

* `actions` - (Optional, List) 
* `condition` - (Optional, String) 
* `sub_rules` - (Optional, List) 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID. Unique identifier of the rule.
* `rule_priority` - Rule priority. only used as an output parameter.


## Import

TEO l7 acc rule v2 can be imported using the {zone_id}#{rule_id}, e.g.

````
terraform import tencentcloud_teo_l7_acc_rule_v2.example zone-3fkff38fyw8s#rule-3ft1xeuhlj1b
````

