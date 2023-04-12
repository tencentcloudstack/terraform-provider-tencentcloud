package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixCkafkaConnectResourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaConnectResource,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ckafka_connect_resource.connect_resource", "id")),
			},
			{
				ResourceName:            "tencentcloud_ckafka_connect_resource.connect_resource",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"mysql_connect_param.0.password"},
			},
		},
	})
}

const testAccCkafkaConnectResource = `

resource "tencentcloud_ckafka_connect_resource" "connect_resource" {
  description   = "for terraform test"
  resource_name = "terraform-test"
  type          = "MYSQL"

  mysql_connect_param {
    is_update   = false
    port        = 3306
    resource    = "cdb-fitq5t9h"
    password    = "QtaTest@123"
    self_built  = false
    service_vip = "172.16.80.59"
    uniq_vpc_id = "vpc-4owdpnwr"
    user_name   = "root"
  }
}

`
