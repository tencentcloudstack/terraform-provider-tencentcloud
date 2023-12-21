package mps_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// NeedFix: only abort task once
func TestAccTencentCloudNeedFixMpsManageTaskOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsManageTaskOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_manage_task_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_manage_task_operation.operation", "task_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_manage_task_operation.operation", "operation_type", "Abort"),
				),
			},
		},
	})
}

const testAccMpsManageTaskOperation = tcacctest.UserInfoData + `

resource "tencentcloud_mps_manage_task_operation" "operation" {
  operation_type = "Abort"
  task_id = "2600010949-LiveScheduleTask-322343d93884db7c1cc252d7f7145d7att7"
}

`
