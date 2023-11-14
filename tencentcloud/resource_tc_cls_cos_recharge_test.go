package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClsCosRechargeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
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
  topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  logset_id = "5cd3a17e-fb0b-418c-afd7-77b365397427"
  name = "test"
  bucket = "test-12345677"
  bucket_region = "ap-guangzhou"
  prefix = "/path"
  log_type = "json_log"
  compress = "gzip"
  extract_rule_info {
		time_key = "time"
		time_format = "YYYY-MM-DD HH:MM:SS"
		delimiter = ","
		log_regex = "*"
		begin_regex = "^*"
		keys = 
		filter_key_regex {
			key = "testKey"
			regex = "testValue"
		}
		un_match_up_load_switch = false
		un_match_log_key = "test"
		backtracking = -1
		is_g_b_k = 0
		json_standard = 1
		protocol = "tcp"
		address = "127.0.0.1:9000"
		parse_protocol = "rfc3164"
		metadata_type = 0
		path_regex = "null"
		meta_tags {
			key = "testKey"
			value = "testValue"
		}

  }
}

`
