package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverRollbackTimeDataSource_basic -v
func TestAccTencentCloudSqlserverRollbackTimeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRollbackTimeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_rollback_time.example")),
			},
		},
	})
}

const testAccSqlserverRollbackTimeDataSource = `
data "tencentcloud_sqlserver_rollback_time" "example" {
  instance_id = "mssql-qelbzgwf"
  dbs         = ["keep_pubsub_db"]
}
`
