package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTkeTmpAlertPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeTmpAlertPolicy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tke_tmp_alert_policy.tmp_alert_policy", "id")),
			},
			{
				ResourceName:      "tencentcloud_tke_tmp_alert_policy.tmp_alert_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTkeTmpAlertPolicy = `

resource "tencentcloud_tke_tmp_alert_policy" "tmp_alert_policy" {
  instance_id = &lt;nil&gt;
  alert_rule {
		name = &lt;nil&gt;
		rules {
			name = &lt;nil&gt;
			rule = &lt;nil&gt;
			labels {
				name = &lt;nil&gt;
				value = &lt;nil&gt;
			}
			template = &lt;nil&gt;
			for = &lt;nil&gt;
			describe = &lt;nil&gt;
			annotations {
				name = &lt;nil&gt;
				value = &lt;nil&gt;
			}
			rule_state = &lt;nil&gt;
		}
		id = &lt;nil&gt;
		template_id = &lt;nil&gt;
		notification {
			enabled = &lt;nil&gt;
			type = &lt;nil&gt;
			web_hook = &lt;nil&gt;
			alert_manager {
				url = &lt;nil&gt;
				cluster_type = &lt;nil&gt;
				cluster_id = &lt;nil&gt;
			}
			repeat_interval = &lt;nil&gt;
			time_range_start = &lt;nil&gt;
			time_range_end = &lt;nil&gt;
			notify_way = &lt;nil&gt;
			receiver_groups = &lt;nil&gt;
			phone_notify_order = &lt;nil&gt;
			phone_circle_times = &lt;nil&gt;
			phone_inner_interval = &lt;nil&gt;
			phone_circle_interval = &lt;nil&gt;
			phone_arrive_notice = &lt;nil&gt;
		}
		updated_at = &lt;nil&gt;
		cluster_id = &lt;nil&gt;

  }
}

`
