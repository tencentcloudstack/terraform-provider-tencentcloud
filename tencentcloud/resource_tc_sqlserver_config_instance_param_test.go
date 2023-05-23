package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverConfigInstanceParamResource_basic -v
func TestAccTencentCloudSqlserverConfigInstanceParamResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigInstanceParam,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_instance_param.config_instance_param", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_instance_param.config_instance_param",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigInstanceParam = `
resource "tencentcloud_sqlserver_config_instance_param" "config_instance_param" {
  instance_id = "mssql-qelbzgwf"
  param_list {
    name = "fill factor(%)"
    current_value = "90"
  }
}
`
