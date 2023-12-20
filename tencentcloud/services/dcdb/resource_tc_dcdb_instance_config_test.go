package dcdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbInstanceConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_instance_config.instance_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_instance_config.instance_config", "rs_access_strategy", "0"),
				),
			},
			{
				Config: testAccDcdbInstanceConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_instance_config.instance_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_instance_config.instance_config", "rs_access_strategy", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_instance_config.instance_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbInstanceConfig = tcacctest.CommonPresetDcdb + `

resource "tencentcloud_dcdb_instance_config" "instance_config" {
  instance_id = local.dcdb_id
  rs_access_strategy = 0
}

`

const testAccDcdbInstanceConfig_update = tcacctest.CommonPresetDcdb + `

resource "tencentcloud_dcdb_instance_config" "instance_config" {
  instance_id = local.dcdb_id
  rs_access_strategy = 1
}

`
