package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudEbEventTargetResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbEventTarget,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_eb_event_target.event_target", "id")),
			},
			{
				ResourceName:      "tencentcloud_eb_event_target.event_target",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccEbEventTarget = `

resource "tencentcloud_eb_event_target" "event_target" {
  event_bus_id = "eb-xxx"
  type = "scf"
  target_description {
		resource_description = "qcs::scf:ap-guangzhou:uin/xxxxxxxx:namespace/xxxxxx/function/xxxxx/x"
		s_c_f_params {
			batch_timeout = 1
			batch_event_count = 1
			enable_batch_delivery = true
		}
		ckafka_target_params {
			topic_name = ""
			retry_policy {
				retry_interval = 1
				max_retry_attempts = 1
			}
		}
		e_s_target_params {
			net_mode = ""
			index_prefix = ""
			rotation_interval = ""
			output_mode = ""
			index_suffix_mode = ""
			index_template_type = ""
		}

  }
  rule_id = ""
}

`
