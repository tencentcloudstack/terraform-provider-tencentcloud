package cls_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixClsCosRechargeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsCosRecharge,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cls_cos_recharge.cos_recharge", "id")),
			},
			{
				ResourceName:      "tencentcloud_cls_cos_recharge.cos_recharge",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsCosRecharge = `

resource "tencentcloud_cls_cos_recharge" "cos_recharge" {
  bucket        = "cos-lock-1308919341"
  bucket_region = "ap-guangzhou"
  log_type      = "minimalist_log"
  logset_id     = "dd426d1a-95bc-4bca-b8c2-baa169261812"
  name          = "cos_recharge_for_test"
  prefix        = "test"
  topic_id      = "7e34a3a7-635e-4da8-9005-88106c1fde69"

  extract_rule_info {
    backtracking            = 0
    is_gbk                  = 0
    json_standard           = 0
    keys                    = []
    metadata_type           = 0
    un_match_up_load_switch = false

    filter_key_regex {
      key   = "__CONTENT__"
      regex = "dasd"
    }
  }
}


`
