package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationFileConfigReleaseResource_basic -v
func TestAccTencentCloudTsfApplicationFileConfigReleaseResource_basic(t *testing.T) {
	// t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationFileConfigRelease,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_file_config_release.application_file_config_release", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config_release.application_file_config_release", "release_desc", "product release"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_application_file_config_release.application_file_config_release",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApplicationFileConfigReleaseVar = `
variable "config_id" {
	default = "` + tcacctest.DefaultTsfFileConfigId + `"
}
`
const testAccTsfApplicationFileConfigRelease = testAccTsfGroup + testAccTsfApplicationFileConfigReleaseVar + `

resource "tencentcloud_tsf_application_file_config_release" "application_file_config_release" {
  config_id = var.config_id
  group_id = tencentcloud_tsf_group.group.id
  release_desc = "product release"
}

`
