package mongodb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMongodbZoneConfigDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config:    testAccMongodbZoneConfigDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_zone_config.zone_config", "list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_zone_config.zone_config", "list.0.available_zone", tcacctest.DefaultAZone),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_zone_config.zone_config", "list.0.cluster_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_zone_config.zone_config", "list.0.machine_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_zone_config.zone_config", "list.0.cpu"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_zone_config.zone_config", "list.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_zone_config.zone_config", "list.0.default_storage"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_zone_config.zone_config", "list.0.min_storage"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_zone_config.zone_config", "list.0.max_storage"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_zone_config.zone_config", "list.0.engine_version"),
				),
			},
		},
	})
}

const testAccMongodbZoneConfigDataSource = tcacctest.DefaultAzVariable + `
data "tencentcloud_mongodb_zone_config" "zone_config" {
	available_zone = var.default_az
}
`
