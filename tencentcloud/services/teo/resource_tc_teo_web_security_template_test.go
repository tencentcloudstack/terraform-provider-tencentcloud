package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoWebSecurityTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoWebSecurityTemplate,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_web_security_template.teo_web_security_template", "id")),
		}, {
			ResourceName:      "tencentcloud_teo_web_security_template.teo_web_security_template",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccTeoWebSecurityTemplate = `

resource "tencentcloud_teo_web_security_template" "teo_web_security_template" {
  security_policy = {
    custom_rules = {
      rules = {
        action = {
          deny_action_parameters = {
          }
          redirect_action_parameters = {
          }
          allow_action_parameters = {
          }
          challenge_action_parameters = {
          }
          block_ip_action_parameters = {
          }
          return_custom_page_action_parameters = {
          }
        }
      }
    }
    managed_rules = {
      auto_update = {
      }
      managed_rule_groups = {
        action = {
          deny_action_parameters = {
          }
          redirect_action_parameters = {
          }
          allow_action_parameters = {
          }
          challenge_action_parameters = {
          }
          block_ip_action_parameters = {
          }
          return_custom_page_action_parameters = {
          }
        }
        rule_actions = {
          action = {
            deny_action_parameters = {
            }
            redirect_action_parameters = {
            }
            allow_action_parameters = {
            }
            challenge_action_parameters = {
            }
            block_ip_action_parameters = {
            }
            return_custom_page_action_parameters = {
            }
          }
        }
        meta_data = {
          rule_details = {
          }
        }
      }
    }
    http_d_do_s_protection = {
      adaptive_frequency_control = {
        action = {
          deny_action_parameters = {
          }
          redirect_action_parameters = {
          }
          allow_action_parameters = {
          }
          challenge_action_parameters = {
          }
          block_ip_action_parameters = {
          }
          return_custom_page_action_parameters = {
          }
        }
      }
      client_filtering = {
        action = {
          deny_action_parameters = {
          }
          redirect_action_parameters = {
          }
          allow_action_parameters = {
          }
          challenge_action_parameters = {
          }
          block_ip_action_parameters = {
          }
          return_custom_page_action_parameters = {
          }
        }
      }
      bandwidth_abuse_defense = {
        action = {
          deny_action_parameters = {
          }
          redirect_action_parameters = {
          }
          allow_action_parameters = {
          }
          challenge_action_parameters = {
          }
          block_ip_action_parameters = {
          }
          return_custom_page_action_parameters = {
          }
        }
      }
      slow_attack_defense = {
        action = {
          deny_action_parameters = {
          }
          redirect_action_parameters = {
          }
          allow_action_parameters = {
          }
          challenge_action_parameters = {
          }
          block_ip_action_parameters = {
          }
          return_custom_page_action_parameters = {
          }
        }
        minimal_request_body_transfer_rate = {
        }
        request_body_transfer_timeout = {
        }
      }
    }
    rate_limiting_rules = {
      rules = {
        action = {
          deny_action_parameters = {
          }
          redirect_action_parameters = {
          }
          allow_action_parameters = {
          }
          challenge_action_parameters = {
          }
          block_ip_action_parameters = {
          }
          return_custom_page_action_parameters = {
          }
        }
      }
    }
    exception_rules = {
      rules = {
        request_fields_for_exception = {
        }
      }
    }
    bot_management = {
      client_attestation_rules = {
        rules = {
          device_profiles = {
            high_risk_request_action = {
              deny_action_parameters = {
              }
              redirect_action_parameters = {
              }
              allow_action_parameters = {
              }
              challenge_action_parameters = {
              }
              block_ip_action_parameters = {
              }
              return_custom_page_action_parameters = {
              }
            }
            medium_risk_request_action = {
              deny_action_parameters = {
              }
              redirect_action_parameters = {
              }
              allow_action_parameters = {
              }
              challenge_action_parameters = {
              }
              block_ip_action_parameters = {
              }
              return_custom_page_action_parameters = {
              }
            }
          }
          invalid_attestation_action = {
            deny_action_parameters = {
            }
            redirect_action_parameters = {
            }
            allow_action_parameters = {
            }
            challenge_action_parameters = {
            }
            block_ip_action_parameters = {
            }
            return_custom_page_action_parameters = {
            }
          }
        }
      }
    }
  }
}
`
