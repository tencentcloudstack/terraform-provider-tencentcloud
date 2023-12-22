package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccDataSourceTencentCloudSqlserverZoneConfig_basic -v
func TestAccDataSourceTencentCloudSqlserverZoneConfig_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSqlserverZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_zone_config.foo"),
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
data "tencentcloud_sqlserver_zone_config" "foo" {}
`
