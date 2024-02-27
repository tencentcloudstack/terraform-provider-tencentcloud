package vod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVodSubApplicationResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSubApplication,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vod_sub_application.foo", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.foo", "name", "foo"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.foo", "status", "On"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.foo", "description", "this is sub application"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_sub_application.foo", "create_time"),
				),
			},
			{
				ResourceName:            "tencentcloud_vod_sub_application.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

const testAccVodSubApplication = `
resource  "tencentcloud_vod_sub_application" "foo" {
	name = "foo"
	status = "On"
	description = "this is sub application"
  }
`
