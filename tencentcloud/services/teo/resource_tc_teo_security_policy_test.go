package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoSecurityPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoSecurityPolicy,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_security_policy.teo_security_policy", "id")),
		}, {
			ResourceName:      "tencentcloud_teo_security_policy.teo_security_policy",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccTeoSecurityPolicy = `

resource "tencentcloud_teo_security_policy" "teo_security_policy" {
  security_config = {
    waf_config = {
      waf_rule = {
      }
      ai_rule = {
      }
    }
    rate_limit_config = {
      rate_limit_user_rules = {
        acl_conditions = {
        }
      }
      rate_limit_template = {
        rate_limit_template_detail = {
        }
      }
      rate_limit_intelligence = {
      }
      rate_limit_customizes = {
        acl_conditions = {
        }
      }
    }
    acl_config = {
      acl_user_rules = {
        acl_conditions = {
        }
      }
      customizes = {
        acl_conditions = {
        }
      }
    }
    bot_config = {
      bot_managed_rule = {
      }
      bot_portrait_rule = {
      }
      intelligence_rule = {
        intelligence_rule_items = {
        }
      }
      bot_user_rules = {
        extend_actions = {
        }
        acl_conditions = {
        }
      }
      alg_detect_rule = {
        alg_conditions = {
        }
        alg_detect_session = {
          alg_detect_results = {
          }
          session_behaviors = {
          }
        }
        alg_detect_js = {
          alg_detect_results = {
          }
        }
      }
      customizes = {
        extend_actions = {
        }
        acl_conditions = {
        }
      }
    }
    switch_config = {
    }
    ip_table_config = {
      ip_table_rules = {
      }
    }
    except_config = {
      except_user_rules = {
        except_user_rule_conditions = {
        }
        except_user_rule_scope = {
          partial_modules = {
          }
          skip_conditions = {
          }
        }
      }
    }
    drop_page_config = {
      waf_drop_page_detail = {
      }
      acl_drop_page_detail = {
      }
    }
    template_config = {
    }
    slow_post_config = {
      first_part_config = {
      }
      slow_rate_config = {
      }
    }
    detect_length_limit_config = {
      detect_length_limit_rules = {
        conditions = {
        }
      }
    }
  }
  security_policy = {
    custom_rules = {
      rules = {
        action = {
          block_ip_action_parameters = {
          }
          return_custom_page_action_parameters = {
          }
          redirect_action_parameters = {
          }
        }
      }
    }
    managed_rules = {
      auto_update = {
      }
      managed_rule_groups = {
        action = {
          block_ip_action_parameters = {
          }
          return_custom_page_action_parameters = {
          }
          redirect_action_parameters = {
          }
        }
        rule_actions = {
          action = {
            block_ip_action_parameters = {
            }
            return_custom_page_action_parameters = {
            }
            redirect_action_parameters = {
            }
          }
        }
        meta_data = {
          rule_details = {
          }
        }
      }
    }
  }
}
`
