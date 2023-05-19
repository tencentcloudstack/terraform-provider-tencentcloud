package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverDBCharsetsDataSource_basic -v
func TestAccTencentCloudSqlserverDBCharsetsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDBCharsetsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_db_charsets.db_charsets")),
			},
		},
	})
}

const testAccSqlserverDBCharsetsDataSource = `
data "tencentcloud_sqlserver_db_charsets" "db_charsets" {
  instance_id = "mssql-qelbzgwf"
}
`
