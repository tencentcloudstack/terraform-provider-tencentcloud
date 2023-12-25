package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverDBCharsetsDataSource_basic -v
func TestAccTencentCloudSqlserverDBCharsetsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDBCharsetsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_db_charsets.example")),
			},
		},
	})
}

const testAccSqlserverDBCharsetsDataSource = `
data "tencentcloud_sqlserver_db_charsets" "example" {
  instance_id = "mssql-qelbzgwf"
}
`
