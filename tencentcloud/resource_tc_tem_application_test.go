package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTemApplicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemApplication,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tem_application.application", "id")),
			},
			{
				ResourceName:      "tencentcloud_tem_application.application",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemApplication = `

resource "tencentcloud_tem_application" "application" {
  application_name = "xxx"
  description = "xxx"
  coding_language = "JAVA"
  use_default_image_service = 1
  repo_type = 0
  repo_server = &lt;nil&gt;
  repo_name = &lt;nil&gt;
  instance_id = &lt;nil&gt;
  tags {
		tag_key = "key"
		tag_value = "tag value"

  }
}

`
