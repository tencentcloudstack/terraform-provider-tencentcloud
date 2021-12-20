package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudSqlserverZoneConfig_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSqlserverZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_zone_config.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.availability_zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.specinfo_list.0.spec_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.specinfo_list.0.machine_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.specinfo_list.0.db_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.specinfo_list.0.db_version_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.specinfo_list.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.specinfo_list.0.cpu"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.specinfo_list.0.min_storage_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.specinfo_list.0.max_storage_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.specinfo_list.0.qps"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_zone_config.foo", "zone_list.0.specinfo_list.0.charge_type"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudSqlserverZoneConfig = `
data "tencentcloud_sqlserver_zone_config" "foo" {
}
`
