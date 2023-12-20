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
				Config: testAccCvmLaunchTemplateDefaultVersion,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_launch_template_default_version.launch_template_default_version", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_launch_template_default_version.launch_template_default_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmLaunchTemplateDefaultVersion = `

resource "tencentcloud_cvm_launch_template_default_version" "launch_template_default_version" {
  launch_template_id = "lt-9e1znnsa"
  default_version = 4
}

`
