package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudClbFunctionTargetsAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbFunctionTargetsAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_function_targets_attachment.function_targets", "id")),
			},
			{
				ResourceName:      "tencentcloud_clb_function_targets_attachment.function_targets",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClbFunctionTargetsAttachment = `

resource "tencentcloud_clb_function_targets_attachment" "function_targets" {
  domain           = "xxx.com"
  listener_id      = "lbl-nonkgvc2"
  load_balancer_id = "lb-5dnrkgry"
  url              = "/"

  function_targets {
    weight = 10

    function {
      function_name           = "keep-tf-test-1675954233"
      function_namespace      = "default"
      function_qualifier      = "$LATEST"
      function_qualifier_type = "VERSION"
    }
  }
}

`
