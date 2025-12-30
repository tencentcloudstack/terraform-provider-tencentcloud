package vcube_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudVcubeApplicationAndVideoResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVcubeApplicationAndVideo,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vcube_application_and_video.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_vcube_application_and_video.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVcubeApplicationAndVideo = `
resource "tencentcloud_vcube_application_and_video" "example" {
  app_name  = "tf-example"
  bundle_id = "com.example.appName"
}
`
