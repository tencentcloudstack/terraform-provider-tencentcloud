package teo_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"testing"
)

func TestAccTencentCloudTeoL7AccRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoL7AccRule,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule", "id")),
		}, {
			ResourceName:      "tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccTeoL7AccRule = `

resource "tencentcloud_teo_l7_acc_rule" "teo_l7_acc_rule" {
  rules = {
    status = ""
    rule_name = ""
    rule_priority = ""
    branches = {
      actions = {
        cache_parameters = {
          follow_origin = {
          }
          no_cache = {
          }
          custom_time = {
          }
        }
        cache_key_parameters = {
          query_string = {
          }
          header = {
          }
          cookie = {
          }
        }
        cache_prefresh_parameters = {
        }
        access_url_redirect_parameters = {
          host_name = {
          }
          url_path = {
          }
          query_string = {
          }
        }
        upstream_url_rewrite_parameters = {
        }
        quic_parameters = {
        }
        web_socket_parameters = {
        }
        authentication_parameters = {
        }
        max_age_parameters = {
        }
        status_code_cache_parameters = {
          status_code_cache_params = {
          }
        }
        offline_cache_parameters = {
        }
        smart_routing_parameters = {
        }
        range_origin_pull_parameters = {
        }
        upstream_http2_parameters = {
        }
        host_header_parameters = {
        }
        force_redirect_https_parameters = {
        }
        compression_parameters = {
        }
        hsts_parameters = {
        }
        client_ip_header_parameters = {
        }
        ocsp_stapling_parameters = {
        }
        http2_parameters = {
        }
        post_max_size_parameters = {
        }
        client_ip_country_parameters = {
        }
        upstream_follow_redirect_parameters = {
        }
        upstream_request_parameters = {
          query_string = {
          }
          cookie = {
          }
        }
        tls_config_parameters = {
        }
        modify_origin_parameters = {
          private_parameters = {
          }
        }
        http_upstream_timeout_parameters = {
        }
        http_response_parameters = {
        }
        error_page_parameters = {
          error_page_params = {
          }
        }
        modify_response_header_parameters = {
          header_actions = {
          }
        }
        modify_request_header_parameters = {
          header_actions = {
          }
        }
        response_speed_limit_parameters = {
        }
        set_content_identifier_parameters = {
        }
      }
      sub_rules = {
      }
    }
  }
}
`
