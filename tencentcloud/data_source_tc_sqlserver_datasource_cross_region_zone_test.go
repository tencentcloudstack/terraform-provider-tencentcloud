package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverDatasourceCrossRegionZoneDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDatasourceCrossRegionZoneDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_datasource_cross_region_zone.datasource_cross_region_zone")),
			},
		},
	})
}

const testAccSqlserverDatasourceCrossRegionZoneDataSource = `

data "tencentcloud_sqlserver_datasource_cross_region_zone" "datasource_cross_region_zone" {
  instance_id = "mssql-j8kv137v"
    }

`
