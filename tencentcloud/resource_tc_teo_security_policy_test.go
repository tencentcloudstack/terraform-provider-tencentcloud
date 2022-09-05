package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTeoSecurityPolicy_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoSecurityPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_policy.securityPolicy", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_security_policy.securityPolicy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoSecurityPolicy = `

resource "tencentcloud_teo_security_policy" "securityPolicy" {
  zone_id = ""
  entity  = ""
  config {
    waf_config {
      switch = ""
      level  = ""
      mode   = ""
      waf_rules {
        switch            = ""
        block_rule_i_ds   = ""
        observe_rule_i_ds = ""
      }
      ai_rule {
        mode = ""
      }
    }
    rate_limit_config {
      switch = ""
      user_rules {
        rule_name        = ""
        threshold        = ""
        period           = ""
        action           = ""
        punish_time      = ""
        punish_time_unit = ""
        rule_status      = ""
        freq_fields      = ""
        conditions {
          match_from    = ""
          match_param   = ""
          operator      = ""
          match_content = ""
        }
        rule_priority    = ""
      }
      template {
        mode = ""
        detail {
          mode        = ""
          i_d         = ""
          action      = ""
          punish_time = ""
          threshold   = ""
          period      = ""
        }
      }
      intelligence {
        switch = ""
        action = ""
      }
    }
    acl_config {
      switch = ""
      user_rules {
        rule_name        = ""
        action           = ""
        rule_status      = ""
        conditions {
          match_from    = ""
          match_param   = ""
          operator      = ""
          match_content = ""
        }
        rule_priority    = ""
        punish_time      = ""
        punish_time_unit = ""
        name             = ""
        page_id          = ""
        redirect_url     = ""
        response_code    = ""
      }
    }
    bot_config {
      switch = ""
      managed_rule {
        rule_i_d          = ""
        action            = ""
        punish_time       = ""
        punish_time_unit  = ""
        name              = ""
        page_id           = ""
        redirect_url      = ""
        response_code     = ""
        trans_managed_ids = ""
        alg_managed_ids   = ""
        cap_managed_ids   = ""
        mon_managed_ids   = ""
        drop_managed_ids  = ""
      }
      portrait_rule {
        rule_i_d         = ""
        alg_managed_ids  = ""
        cap_managed_ids  = ""
        mon_managed_ids  = ""
        drop_managed_ids = ""
        switch           = ""
      }
      intelligence_rule {
        switch = ""
        items {
          label  = ""
          action = ""
        }
      }
    }
    switch_config {
      web_switch = ""
    }
    ip_table_config {
      switch = ""
      rules {
        action        = ""
        match_from    = ""
        match_content = ""
        rule_i_d      = ""
      }
    }

  }
}

`
