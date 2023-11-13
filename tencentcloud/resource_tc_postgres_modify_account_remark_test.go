package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresModifyAccountRemarkResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresModifyAccountRemark,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_modify_account_remark.modify_account_remark", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_modify_account_remark.modify_account_remark",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresModifyAccountRemark = `

resource "tencentcloud_postgres_modify_account_remark" "modify_account_remark" {
  d_b_instance_id = ""
  user_name = ""
  remark = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
