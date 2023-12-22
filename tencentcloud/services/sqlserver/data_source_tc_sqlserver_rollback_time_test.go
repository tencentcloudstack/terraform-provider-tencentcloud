package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverRollbackTimeDataSource_basic -v
func TestAccTencentCloudSqlserverRollbackTimeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRollbackTimeDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_rollback_time.example")),
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
