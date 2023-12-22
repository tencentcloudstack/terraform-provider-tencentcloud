package rum_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixRumReleaseFileResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumReleaseFile,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_rum_release_file.release_file", "id")),
			},
			{
				ResourceName:      "tencentcloud_rum_release_file.release_file",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRumReleaseFile = `

resource "tencentcloud_rum_release_file" "release_file" {
  project_id      = 120000
  version         = "1.0"
  file_key        = "120000-last-1632921299138-index.js.map"
  file_name       = "index.js.map"
  file_hash 	  = "b148c43fd81d845ba7cc6907928ce430"
  release_file_id = 1
}

`
