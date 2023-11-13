package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_rollback_time.rollback_time")),
			},
		},
	})
}

const testAccSqlserverRollbackTimeDataSource = `

data "tencentcloud_sqlserver_rollback_time" "rollback_time" {
  instance_id = "mssql-j8kv137v"
  d_bs = 
  }

`
