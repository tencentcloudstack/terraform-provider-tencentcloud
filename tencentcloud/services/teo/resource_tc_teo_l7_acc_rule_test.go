package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -test.run TestAccTencentCloudTeoL7AccRuleResource_basic -v
func TestAccTencentCloudTeoL7AccRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7AccRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "zone_id", "zone-39quuimqg8r6"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.#", "6"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.rule_name", "Web Acceleration"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.condition", "${http.request.host} in ['aaa.makn.cn']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.name", "Cache"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.cache_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.cache_parameters.0.custom_time.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.cache_parameters.0.custom_time.0.cache_time", "2592000"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.cache_parameters.0.custom_time.0.ignore_cache_control", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.cache_parameters.0.custom_time.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.name", "CacheKey"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.scheme", ""),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.query_string.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.query_string.0.action", ""),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.query_string.0.values.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.condition", "lower(${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.actions.0.name", "Cache"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.0.no_cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.0.no_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.condition", "${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.actions.0.name", "MaxAge"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.0.cache_time", "3600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.0.follow_origin", "off"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoL7AccRuleUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "zone_id", "zone-39quuimqg8r6"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.#", "5"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.rule_name", "Web Acceleration"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.condition", "${http.request.host} in ['aaa.makn.cn']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.name", "Cache"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.cache_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.cache_parameters.0.custom_time.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.cache_parameters.0.custom_time.0.cache_time", "2592000"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.cache_parameters.0.custom_time.0.ignore_cache_control", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.0.cache_parameters.0.custom_time.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.name", "CacheKey"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.scheme", ""),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.query_string.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.query_string.0.action", ""),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.actions.1.cache_key_parameters.0.query_string.0.values.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.condition", "lower(${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.actions.0.name", "Cache"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.0.no_cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.0.no_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.condition", "${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.actions.0.name", "MaxAge"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.0.cache_time", "3600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "rules.0.branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.0.follow_origin", "off"),
				),
			},
		},
	})
}

const testAccTeoL7AccRule = `

resource "tencentcloud_teo_l7_acc_rule" "teo_l7_acc_rule" {
  zone_id = "zone-39quuimqg8r6"
  rules {
    description = ["1"]
    rule_name   = "Web Acceleration"
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
          scheme         = null
          query_string {
            action = null
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
  rules {
    description = ["2"]
    rule_name   = "Live Video Streaming"
    branches {
      condition = "$${http.request.host} in ['aaa.makn.cn']"
      sub_rules {
        description = ["2-1"]
        branches {
          condition = "$${http.request.file_extension} in ['m3u8', 'mpd']"
          actions {
            name = "Cache"
            cache_parameters {
              custom_time {
                cache_time           = 1
                ignore_cache_control = "off"
                switch               = "on"
              }
            }
          }
        }
        branches {
          condition = "$${http.request.file_extension} in ['ts', 'mp4', 'm4a', 'm4s']"
          actions {
            name = "Cache"
            cache_parameters {
              custom_time {
                cache_time           = 86400
                ignore_cache_control = "off"
                switch               = "on"
              }
            }
          }
        }
        branches {
          condition = "*"
          actions {
            name = "Cache"
            cache_parameters {
              follow_origin {
                default_cache          = "on"
                default_cache_strategy = "on"
                default_cache_time     = 0
                switch                 = "on"
              }
            }
          }
        }
      }
    }
  }
  rules {
    description = ["3"]
    rule_name   = "Large File Download"
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
          full_url_cache = "off"
          ignore_case    = null
          scheme         = null
          query_string {
            action = null
            switch = "off"
            values = []
          }
        }
      }
      actions {
        name = "RangeOriginPull"
        range_origin_pull_parameters {
          switch = "on"
        }
      }
      sub_rules {
        description = ["3-1"]
        branches {
          condition = "$${http.request.file_extension} in ['php', 'jsp', 'asp', 'aspx']"
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
    }
  }
  rules {
    description = ["4"]
    rule_name   = "Video On Demand"
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
          full_url_cache = "off"
          ignore_case    = "off"
          scheme         = null
          query_string {
            action = null
            switch = "off"
            values = []
          }
        }
      }
      actions {
        name = "RangeOriginPull"
        range_origin_pull_parameters {
          switch = "on"
        }
      }
      sub_rules {
        description = ["4-1"]
        branches {
          condition = "$${http.request.file_extension} in ['php', 'jsp', 'asp', 'aspx']"
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
    }
  }
  rules {
    description = ["5"]
    rule_name   = "API Acceleration"
    branches {
      condition = "$${http.request.host} in ['aaa.makn.cn']"
      actions {
        name = "Cache"
        cache_parameters {
          no_cache {
            switch = "on"
          }
        }
      }
      actions {
        name = "SmartRouting"
        smart_routing_parameters {
          switch = "off"
        }
      }
    }
  }
  rules {
    description = ["6"]
    rule_name   = "WordPress Site"
    branches {
      condition = "$${http.request.host} in ['aaa.makn.cn']"
      sub_rules {
        description = ["6-1"]
        branches {
          condition = "$${http.request.file_extension} in ['gif', 'png', 'bmp', 'jpeg', 'tif', 'tiff', 'zip', 'exe', 'wmv', 'swf', 'mp3', 'wma', 'rar', 'css', 'flv', 'mp4', 'txt', 'ico', 'js']"
          actions {
            name = "Cache"
            cache_parameters {
              custom_time {
                cache_time           = 604800
                ignore_cache_control = "off"
                switch               = "on"
              }
            }
          }
        }
        branches {
          condition = "$${http.request.uri.path} in ['/']"
          actions {
            name = "Cache"
            cache_parameters {
              no_cache {
                switch = "on"
              }
            }
          }
        }
        branches {
          condition = "$${http.request.file_extension} in ['aspx', 'jsp', 'php', 'asp', 'do', 'dwr', 'cgi', 'fcgi', 'action', 'ashx', 'axd']"
          actions {
            name = "Cache"
            cache_parameters {
              no_cache {
                switch = "on"
              }
            }
          }
        }
        branches {
          condition = "$${http.request.uri.path} in ['/wp-admin/']"
          actions {
            name = "Cache"
            cache_parameters {
              no_cache {
                switch = "on"
              }
            }
          }
        }
        branches {
          condition = "*"
          actions {
            name = "Cache"
            cache_parameters {
              follow_origin {
                default_cache          = "on"
                default_cache_strategy = "on"
                default_cache_time     = 0
                switch                 = "on"
              }
            }
          }
        }
      }
    }
  }
}
`

const testAccTeoL7AccRuleUp = `

resource "tencentcloud_teo_l7_acc_rule" "teo_l7_acc_rule" {
  zone_id = "zone-39quuimqg8r6"
  rules {
    description = ["1"]
    rule_name   = "Web Acceleration"
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
          scheme         = null
          query_string {
            action = null
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
  rules {
    description = ["2"]
    rule_name   = "Live Video Streaming"
    branches {
      condition = "$${http.request.host} in ['aaa.makn.cn']"
      sub_rules {
        description = ["2-1"]
        branches {
          condition = "$${http.request.file_extension} in ['m3u8', 'mpd']"
          actions {
            name = "Cache"
            cache_parameters {
              custom_time {
                cache_time           = 1
                ignore_cache_control = "off"
                switch               = "on"
              }
            }
          }
        }
        branches {
          condition = "$${http.request.file_extension} in ['ts', 'mp4', 'm4a', 'm4s']"
          actions {
            name = "Cache"
            cache_parameters {
              custom_time {
                cache_time           = 86400
                ignore_cache_control = "off"
                switch               = "on"
              }
            }
          }
        }
        branches {
          condition = "*"
          actions {
            name = "Cache"
            cache_parameters {
              follow_origin {
                default_cache          = "on"
                default_cache_strategy = "on"
                default_cache_time     = 0
                switch                 = "on"
              }
            }
          }
        }
      }
    }
  }
  rules {
    description = ["3"]
    rule_name   = "Large File Download"
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
          full_url_cache = "off"
          ignore_case    = null
          scheme         = null
          query_string {
            action = null
            switch = "off"
            values = []
          }
        }
      }
      actions {
        name = "RangeOriginPull"
        range_origin_pull_parameters {
          switch = "on"
        }
      }
      sub_rules {
        description = ["3-1"]
        branches {
          condition = "$${http.request.file_extension} in ['php', 'jsp', 'asp', 'aspx']"
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
    }
  }
  rules {
    description = ["6"]
    rule_name   = "WordPress Site"
    branches {
      condition = "$${http.request.host} in ['aaa.makn.cn']"
      sub_rules {
        description = ["6-1"]
        branches {
          condition = "$${http.request.file_extension} in ['gif', 'png', 'bmp', 'jpeg', 'tif', 'tiff', 'zip', 'exe', 'wmv', 'swf', 'mp3', 'wma', 'rar', 'css', 'flv', 'mp4', 'txt', 'ico', 'js']"
          actions {
            name = "Cache"
            cache_parameters {
              custom_time {
                cache_time           = 604800
                ignore_cache_control = "off"
                switch               = "on"
              }
            }
          }
        }
        branches {
          condition = "$${http.request.uri.path} in ['/']"
          actions {
            name = "Cache"
            cache_parameters {
              no_cache {
                switch = "on"
              }
            }
          }
        }
        branches {
          condition = "$${http.request.file_extension} in ['aspx', 'jsp', 'php', 'asp', 'do', 'dwr', 'cgi', 'fcgi', 'action', 'ashx', 'axd']"
          actions {
            name = "Cache"
            cache_parameters {
              no_cache {
                switch = "on"
              }
            }
          }
        }
        branches {
          condition = "$${http.request.uri.path} in ['/wp-admin/']"
          actions {
            name = "Cache"
            cache_parameters {
              no_cache {
                switch = "on"
              }
            }
          }
        }
        branches {
          condition = "*"
          actions {
            name = "Cache"
            cache_parameters {
              follow_origin {
                default_cache          = "on"
                default_cache_strategy = "on"
                default_cache_time     = 0
                switch                 = "on"
              }
            }
          }
        }
      }
    }
  }
    rules {
    description = ["4"]
    rule_name   = "Video On Demand"
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
          full_url_cache = "off"
          ignore_case    = "off"
          scheme         = null
          query_string {
            action = null
            switch = "off"
            values = []
          }
        }
      }
      actions {
        name = "RangeOriginPull"
        range_origin_pull_parameters {
          switch = "on"
        }
      }
      sub_rules {
        description = ["4-1"]
        branches {
          condition = "$${http.request.file_extension} in ['php', 'jsp', 'asp', 'aspx']"
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
    }
  }
}
`
