package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoDeployConfigGroupVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDeployConfigGroupVersion,
				Check: resource.ComposeTestCheckFunc(
					// Basic resource attributes
					resource.TestCheckResourceAttrSet("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "id"),

					// Input parameters validation
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "zone_id", "zone-2xkazzl8yf6k"),
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "env_id", "env-3lchxiq1h855"),
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "description", "Deploy config group version for production"),

					// Config group version infos validation
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "config_group_version_infos.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "config_group_version_infos.0.version_id", "ver-3lchxizh2mqn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "config_group_version_infos.1.version_id", "ver-3lchxjdciuzx"),

					// Computed attributes validation
					resource.TestCheckResourceAttrSet("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "record_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "deploy_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "status"),

					// Optional computed attributes (may or may not be present)
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "status", "success"),
				),
			},
		},
	})
}

const testAccTeoDeployConfigGroupVersion = `

resource "tencentcloud_teo_deploy_config_group_version" "teo_deploy_config_group_version" {
  zone_id = "zone-2xkazzl8yf6k"
  env_id = "env-3lchxiq1h855"
  description = "Deploy config group version for production"
  # l7_acceleration
  config_group_version_infos {
    version_id = "ver-3lchxizh2mqn"
  }
  # edge_functions
  config_group_version_infos {
    version_id = "ver-3lchxjdciuzx"
  }
}
`
