package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverInstanceTDEResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverInstanceTDE,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_t_d_e.instance_t_d_e", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_instance_t_d_e.instance_t_d_e",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverInstanceTDE = `

resource "tencentcloud_sqlserver_instance_t_d_e" "instance_t_d_e" {
  instance_id = "mssql-i1z41iwd"
  certificate_attribution = ""
  quote_uin = ""
}

`
