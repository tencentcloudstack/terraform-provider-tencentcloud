package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfGroupConfigReleaseDataSource_basic -v
func TestAccTencentCloudTsfGroupConfigReleaseDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfGroupConfigReleaseDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_group_config_release.group_config_release"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.package_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.package_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.package_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.application_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.cluster_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.cluster_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.config_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.config_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.config_release_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.config_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.namespace_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.namespace_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_config_release.group_config_release", "result.0.config_release_list.0.release_time"),
				),
			},
		},
	})
}

const testAccTsfGroupConfigReleaseDataSourceVar = `
variable "group_id" {
	default = "` + defaultTsfGroupId + `"
}
`

const testAccTsfGroupConfigReleaseDataSource = testAccTsfGroupConfigReleaseDataSourceVar + `

data "tencentcloud_tsf_group_config_release" "group_config_release" {
	group_id = var.group_id
}

`
