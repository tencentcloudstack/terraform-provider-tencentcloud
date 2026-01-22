package sqlserver_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

var testAccSqlserverCollationTimeZoneDataSourceName = "data.tencentcloud_sqlserver_collation_time_zone.sqlserver_collation_time_zone"

func TestAccTencentCloudSqlserverCollationTimeZoneDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverBasicInstancesBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAccSqlserverCollationTimeZoneDataSourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccSqlserverCollationTimeZoneDataSourceName, "collation.#"),
					resource.TestCheckResourceAttrSet(testAccSqlserverCollationTimeZoneDataSourceName, "time_zone.#"),
				),
			},
		},
	})
}

const testAccSqlserverCollationTimeZoneDataSource = `

data "tencentcloud_sqlserver_collation_time_zone" "sqlserver_collation_time_zone" {
}
`
