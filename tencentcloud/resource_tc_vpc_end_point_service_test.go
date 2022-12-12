package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestTencentCloudVpcEndPointServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointService,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point_service.end_point_service", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_end_point_service.end_point_service",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcEndPointService = `

resource "tencentcloud_vpc_end_point_service" "end_point_service" {
  vpc_id = ""
  end_point_service_name = ""
  auto_accept_flag = 
  service_instance_id = ""
  is_pass_service = 
}

`
