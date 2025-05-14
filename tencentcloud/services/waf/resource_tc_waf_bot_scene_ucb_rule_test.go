package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentClouddWafBotSceneUCBRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccdWafBotSceneUCBRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_ucb_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_ucb_rule.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_ucb_rule.example", "scene_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_ucb_rule.example", "rule.#"),
				),
			},
			{
				Config: testAccdWafBotSceneUCBRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_ucb_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_ucb_rule.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_ucb_rule.example", "scene_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_ucb_rule.example", "rule.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_bot_scene_ucb_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccdWafBotSceneUCBRule = `
resource "tencentcloud_waf_bot_scene_ucb_rule" "example" {
  domain   = "news.bots.icu"
  scene_id = "3000000791"
  rule {
    domain = "news.bots.icu"
    name   = "tf-example"
    rule {
      key  = "post_value"
      op   = "prefix"
      lang = "cn"
      value {
        multi_value = [
          "terraform",
          "provider"
        ]
      }
    }

    action        = "intercept"
    on_off        = "on"
    rule_type     = 0
    prior         = 100
    label         = "恶意BOT"
    appid         = 1256704386
    addition_arg  = "none"
    desc          = "rule desc."
    pre_define    = true
    block_page_id = 71
    job_type      = "cron_week"
    job_date_time {
      cron {
        w_days     = [1, 2, 3, 4, 5]
        start_time = "00:00:00"
        end_time   = "23:59:59"
      }

      time_t_zone = "UTC+8"
    }
  }
}
`

const testAccdWafBotSceneUCBRuleUpdate = `
resource "tencentcloud_waf_bot_scene_ucb_rule" "example" {
  domain   = "news.bots.icu"
  scene_id = "3000000791"
  rule {
    domain = "news.bots.icu"
    name   = "tf-example-update"
    rule {
      key  = "post_value"
      op   = "prefix"
      lang = "cn"
      value {
        multi_value = [
          "terraform",
          "provider"
        ]
      }
    }

    action        = "intercept"
    on_off        = "on"
    rule_type     = 0
    prior         = 100
    label         = "恶意BOT"
    appid         = 1256704386
    addition_arg  = "none"
    desc          = "rule desc."
    pre_define    = true
    block_page_id = 71
    job_type      = "cron_week"
    job_date_time {
      cron {
        w_days     = [1, 2, 3, 4, 5]
        start_time = "01:00:00"
        end_time   = "20:00:00"
      }

      time_t_zone = "UTC+8"
    }
  }
}
`
