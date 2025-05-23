package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoL7AccRulePriorityOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7AccRulePriorityOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule_priority_operation.teo_l7_acc_rule_priority_operation", "id"),
				),
			},
		},
	})
}

const testAccTeoL7AccRulePriorityOperation = `
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
`
