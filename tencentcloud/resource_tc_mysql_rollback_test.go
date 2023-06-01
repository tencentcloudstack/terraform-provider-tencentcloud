package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlRollbackResource_basic -v
func TestAccTencentCloudMysqlRollbackResource_basic(t *testing.T) {
	t.Parallel()

	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -1).In(loc).Format("2006-01-02 15:04:05")
	timeUnix := time.Now().AddDate(0, 0, -1).Unix()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMysqlRollback, startTime, timeUnix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_rollback.rollback", "id"),
				),
			},
		},
	})
}

const testAccMysqlRollbackVar = `
variable "instance_id" {
  default = "` + defaultDbBrainInstanceId + `"
}
`

const testAccMysqlRollback = testAccMysqlRollbackVar + `

resource "tencentcloud_mysql_rollback" "rollback" {
	instance_id = var.instance_id
	strategy = "full"
	rollback_time = "%v"
	databases {
	  database_name = "tf_ci_test"
	  new_database_name = "tf_ci_test_%v"
	}

}

`
