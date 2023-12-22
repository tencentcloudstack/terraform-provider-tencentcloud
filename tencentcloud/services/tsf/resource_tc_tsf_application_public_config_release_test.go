package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationPublicConfigReleaseResource_basic -v
func TestAccTencentCloudTsfApplicationPublicConfigReleaseResource_basic(t *testing.T) {
	// t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationPublicConfigRelease,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_public_config_release.application_public_config_release", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_public_config_release.application_public_config_release", "release_desc", "v1"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_application_public_config_release.application_public_config_release",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApplicationPublicConfigRelease = testAccTsfApplicationPublicConfig + testAccTsfNamespace + `

resource "tencentcloud_tsf_application_public_config_release" "application_public_config_release" {
  config_id = tencentcloud_tsf_application_public_config.application_public_config.id
  namespace_id = tencentcloud_tsf_namespace.namespace.id
  release_desc = "v1"
}

`
