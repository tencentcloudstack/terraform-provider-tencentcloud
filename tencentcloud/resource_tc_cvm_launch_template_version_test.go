package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmLaunchTemplateVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmLaunchTemplateVersion,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_launch_template_version.launch_template_version", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_launch_template_version.launch_template_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmLaunchTemplateVersion = `
resource "tencentcloud_cvm_launch_template_version" "launch_template_version" {
	placement {
		  zone = "ap-guangzhou-6"
		  project_id = 0
  
	}
	launch_template_id = "lt-9e1znnsa"
	launch_template_version_description = "version description"
	disable_api_termination = false
	instance_type = "S5.MEDIUM4"
	image_id = "img-9qrfy1xt"
}
`
