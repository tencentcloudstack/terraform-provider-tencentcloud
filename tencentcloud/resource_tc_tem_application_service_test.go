package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudTemApplicationServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemApplicationService,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tem_application_service.application_service", "id")),
			},
			{
				ResourceName:      "tencentcloud_tem_application_service.application_service",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemApplicationService = `

resource "tencentcloud_tem_application_service" "application_service" {
  environment_id = "en-xxx"
  application_id = "xxx"
  service {
		type = "CLUSTER"
		service_name = "consumer"
		port_mapping_item_list {
			port = 80
			target_port = 80
			protocol = "tcp"
		}

  }
}

`
