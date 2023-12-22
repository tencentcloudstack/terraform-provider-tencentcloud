package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcModifyUserTypOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcModifyUserTypOperation,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_modify_user_typ_operation.modify_user_typ_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_modify_user_typ_operation.modify_user_typ_operation", "user_id", "100032676511"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_modify_user_typ_operation.modify_user_typ_operation", "user_type", "ADMIN")),
			},
			{
				Config: testAccDlcModifyUserTypOperationRecover,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_modify_user_typ_operation.modify_user_typ_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_modify_user_typ_operation.modify_user_typ_operation", "user_id", "100032676511"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_modify_user_typ_operation.modify_user_typ_operation", "user_type", "COMMON")),
			},
		},
	})
}

const testAccDlcModifyUserTypOperation = `

resource "tencentcloud_dlc_modify_user_typ_operation" "modify_user_typ_operation" {
  user_id = "100032676511"
  user_type = "ADMIN"
}

`
const testAccDlcModifyUserTypOperationRecover = `

resource "tencentcloud_dlc_modify_user_typ_operation" "modify_user_typ_operation" {
  user_id = "100032676511"
  user_type = "COMMON"
}

`
