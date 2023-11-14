package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoSecurityPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoSecurityPolicy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_security_policy.security_policy", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_security_policy.security_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoSecurityPolicy = `

resource "tencentcloud_teo_security_policy" "security_policy" {
  zone_id = &lt;nil&gt;
  entity = &lt;nil&gt;
  config {
		waf_config {
			switch = &lt;nil&gt;
			level = &lt;nil&gt;
			mode = &lt;nil&gt;
			waf_rules {
				switch = &lt;nil&gt;
				block_rule_i_ds = &lt;nil&gt;
				observe_rule_i_ds = &lt;nil&gt;
			}
			ai_rule {
				mode = &lt;nil&gt;
			}
		}
		rate_limit_config {
			switch = &lt;nil&gt;
			user_rules {
				rule_name = &lt;nil&gt;
				threshold = &lt;nil&gt;
				period = &lt;nil&gt;
				action = &lt;nil&gt;
				punish_time = &lt;nil&gt;
				punish_time_unit = &lt;nil&gt;
				rule_status = &lt;nil&gt;
				freq_fields = &lt;nil&gt;
				conditions {
					match_from = &lt;nil&gt;
					match_param = &lt;nil&gt;
					operator = &lt;nil&gt;
					match_content = &lt;nil&gt;
				}
				rule_priority = &lt;nil&gt;
			}
			template {
				mode = &lt;nil&gt;
				detail {
					mode = &lt;nil&gt;
					i_d = &lt;nil&gt;
					action = &lt;nil&gt;
					punish_time = &lt;nil&gt;
					threshold = &lt;nil&gt;
					period = &lt;nil&gt;
				}
			}
			intelligence {
				switch = &lt;nil&gt;
				action = &lt;nil&gt;
			}
		}
		acl_config {
			switch = &lt;nil&gt;
			user_rules {
				rule_name = &lt;nil&gt;
				action = &lt;nil&gt;
				rule_status = &lt;nil&gt;
				conditions {
					match_from = &lt;nil&gt;
					match_param = &lt;nil&gt;
					operator = &lt;nil&gt;
					match_content = &lt;nil&gt;
				}
				rule_priority = &lt;nil&gt;
				punish_time = &lt;nil&gt;
				punish_time_unit = &lt;nil&gt;
				name = &lt;nil&gt;
				page_id = &lt;nil&gt;
				redirect_url = &lt;nil&gt;
				response_code = &lt;nil&gt;
			}
		}
		bot_config {
			switch = &lt;nil&gt;
			managed_rule {
				action = &lt;nil&gt;
				punish_time = &lt;nil&gt;
				punish_time_unit = &lt;nil&gt;
				name = &lt;nil&gt;
				page_id = &lt;nil&gt;
				redirect_url = &lt;nil&gt;
				response_code = &lt;nil&gt;
				trans_managed_ids = &lt;nil&gt;
				alg_managed_ids = &lt;nil&gt;
				cap_managed_ids = &lt;nil&gt;
				mon_managed_ids = &lt;nil&gt;
				drop_managed_ids = &lt;nil&gt;
			}
			portrait_rule {
				alg_managed_ids = &lt;nil&gt;
				cap_managed_ids = &lt;nil&gt;
				mon_managed_ids = &lt;nil&gt;
				drop_managed_ids = &lt;nil&gt;
				switch = &lt;nil&gt;
			}
			intelligence_rule {
				switch = &lt;nil&gt;
				items {
					label = &lt;nil&gt;
					action = &lt;nil&gt;
				}
			}
		}
		switch_config {
			web_switch = &lt;nil&gt;
		}
		ip_table_config {
			switch = &lt;nil&gt;
			rules {
				action = &lt;nil&gt;
				match_from = &lt;nil&gt;
				match_content = &lt;nil&gt;
			}
		}
		except_config {
			switch = &lt;nil&gt;
			except_user_rules {
				action = &lt;nil&gt;
				rule_status = &lt;nil&gt;
				rule_priority = &lt;nil&gt;
				except_user_rule_conditions {
					match_from = &lt;nil&gt;
					match_param = &lt;nil&gt;
					operator = &lt;nil&gt;
					match_content = &lt;nil&gt;
				}
				except_user_rule_scope {
					modules = &lt;nil&gt;
				}
			}
		}
		drop_page_config {
			switch = &lt;nil&gt;
			waf_drop_page_detail {
				page_id = &lt;nil&gt;
				status_code = &lt;nil&gt;
				name = &lt;nil&gt;
				type = &lt;nil&gt;
			}
			acl_drop_page_detail {
				page_id = &lt;nil&gt;
				status_code = &lt;nil&gt;
				name = &lt;nil&gt;
				type = &lt;nil&gt;
			}
		}

  }
}

`
