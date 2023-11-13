package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfGroupConfigReleaseDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfGroupConfigReleaseDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_group_config_release.group_config_release")),
			},
		},
	})
}

const testAccTsfGroupConfigReleaseDataSource = `

data "tencentcloud_tsf_group_config_release" "group_config_release" {
  group_id = ""
  }

`
