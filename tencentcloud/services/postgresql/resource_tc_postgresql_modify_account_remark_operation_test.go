package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const TestObjectPgModifyAccountRemark = "tencentcloud_postgresql_modify_account_remark_operation.modify_account_remark_operation"

func TestAccTencentCloudPostgresqlModifyAccountRemarkOperationResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlModifyAccountRemarkOperation,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TestObjectPgModifyAccountRemark, "id"),
					resource.TestCheckResourceAttrSet(TestObjectPgModifyAccountRemark, "db_instance_id"),
					resource.TestCheckResourceAttr(TestObjectPgModifyAccountRemark, "user_name", "root"),
					resource.TestCheckResourceAttr(TestObjectPgModifyAccountRemark, "remark", "hello_world"),
				),
			},
		},
	})
}

const testAccPostgresqlModifyAccountRemarkOperation = tcacctest.OperationPresetPGSQL + `

resource "tencentcloud_postgresql_modify_account_remark_operation" "modify_account_remark_operation" {
  db_instance_id = local.pgsql_id
  user_name = "root"
  remark = "hello_world"
}

`
