package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmLaunchTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmLaunchTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_launch_template.launch_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_launch_template.launch_template", "image_id", "img-9qrfy1xt"),
				),
			},
		},
	})
}

const testAccCvmLaunchTemplate = `
resource "tencentcloud_cvm_launch_template" "launch_template" {
	launch_template_name = "test_launch_template"
	placement {
	  zone = "ap-guangzhou-6"
	  project_id = 0
	  host_ids = []
	  host_ips = []
	}
	image_id = "img-9qrfy1xt"
  }
`
