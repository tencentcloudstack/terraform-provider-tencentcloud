package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixSqlserverInstanceTDEResource_basic -v
func TestAccTencentCloudNeedFixSqlserverInstanceTDEResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverInstanceTDE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_tde.instance_tde", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_tde.instance_tde", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_tde.instance_tde", "certificate_attribution"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_tde.instance_tde", "quote_uin"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_instance_tde.instance_tde",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverInstanceTDE = `
resource "tencentcloud_sqlserver_instance_tde" "instance_tde" {
  instance_id             = "mssql-qelbzgwf"
  certificate_attribution = "self"
}
`
