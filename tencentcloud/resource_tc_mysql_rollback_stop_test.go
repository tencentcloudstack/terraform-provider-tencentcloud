package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMysqlRollbackStopResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRollbackStop,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_rollback_stop.rollback_stop", "id")),
			},
		},
	})
}

const testAccMysqlRollbackStopVar = `
variable "instance_id" {
  default = "` + defaultDbBrainInstanceId + `"
}
`

const testAccMysqlRollbackStop = testAccMysqlRollbackStopVar + `

resource "tencentcloud_mysql_rollback_stop" "rollback_stop" {
	instance_id = var.instance_id
}

`
