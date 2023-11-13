package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverConfigInstanceParamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigInstanceParam,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_instance_param.config_instance_param", "id")),
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
  instance_ids = 
  param_list {
		name = "fill factor(%)"
		current_value = "90"

  }
  wait_switch = 0
}

`
