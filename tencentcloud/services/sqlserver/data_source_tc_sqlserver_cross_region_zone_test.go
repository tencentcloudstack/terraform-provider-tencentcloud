package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverCrossRegionZoneDataSource_basic -v
func TestAccTencentCloudSqlserverCrossRegionZoneDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverCrossRegionZoneDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_cross_region_zone.example")),
			},
		},
	})
}

const testAccSqlserverCrossRegionZoneDataSource = `
data "tencentcloud_sqlserver_cross_region_zone" "example" {
  instance_id = "mssql-qelbzgwf"
}
`
