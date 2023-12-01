package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverCrossRegionZoneDataSource_basic -v
func TestAccTencentCloudSqlserverCrossRegionZoneDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverCrossRegionZoneDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_cross_region_zone.example")),
			},
		},
	})
}

const testAccSqlserverCrossRegionZoneDataSource = `
data "tencentcloud_sqlserver_cross_region_zone" "example" {
  instance_id = "mssql-qelbzgwf"
}
`
