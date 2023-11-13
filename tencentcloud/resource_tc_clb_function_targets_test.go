package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClbFunctionTargetsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbFunctionTargets,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_function_targets.function_targets", "id")),
			},
			{
				ResourceName:      "tencentcloud_clb_function_targets.function_targets",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClbFunctionTargets = `

resource "tencentcloud_clb_function_targets" "function_targets" {
  load_balancer_id = &lt;nil&gt;
  listener_id = &lt;nil&gt;
  function_targets {
		function {
			function_namespace = &lt;nil&gt;
			function_name = &lt;nil&gt;
			function_qualifier = &lt;nil&gt;
			function_qualifier_type = &lt;nil&gt;
		}
		weight = &lt;nil&gt;

  }
  location_id = &lt;nil&gt;
  domain = &lt;nil&gt;
  url = &lt;nil&gt;
}

`
