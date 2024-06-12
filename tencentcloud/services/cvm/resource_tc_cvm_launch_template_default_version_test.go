package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmLaunchTemplateDefaultVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmLaunchTemplateDefaultVersionBase + testAccCvmLaunchTemplateDefaultVersion1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cvm_launch_template_default_version.test_launch_tpl_default", "default_version", "2"),
				),
			},
			{
				Config: testAccCvmLaunchTemplateDefaultVersionBase + testAccCvmLaunchTemplateDefaultVersion2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cvm_launch_template_default_version.test_launch_tpl_default", "default_version", "1"),
				),
			},
		},
	})
}

const testAccCvmLaunchTemplateDefaultVersionBase = `
data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}
resource "tencentcloud_cvm_launch_template" "test_launch_tpl" {
  launch_template_name = "test"
  image_id             = data.tencentcloud_images.default.images.0.image_id
  placement {
    zone = "ap-guangzhou-7"
  }
  instance_name = "v1"
}
resource "tencentcloud_cvm_launch_template_version" "test_launch_tpl_v2" {
  launch_template_id = tencentcloud_cvm_launch_template.test_launch_tpl.id
  placement {
    zone = "ap-guangzhou-7"
  }
  instance_name = "v2"
}`

const testAccCvmLaunchTemplateDefaultVersion1 = `
resource "tencentcloud_cvm_launch_template_default_version" "test_launch_tpl_default" {
  launch_template_id = tencentcloud_cvm_launch_template.test_launch_tpl.id
  default_version = tencentcloud_cvm_launch_template_version.test_launch_tpl_v2.launch_template_version
}
`

const testAccCvmLaunchTemplateDefaultVersion2 = `
resource "tencentcloud_cvm_launch_template_default_version" "test_launch_tpl_default" {
  launch_template_id = tencentcloud_cvm_launch_template.test_launch_tpl.id
  default_version = 1
}
`
