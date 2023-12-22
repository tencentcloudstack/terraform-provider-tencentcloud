package pls_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcEndPointServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
  vpc_id = "vpc-4owdpnwr"
  end_point_service_name = "terraform-endpoint-service"
  auto_accept_flag = false
  service_instance_id = "lb-6ksdvyqw"
  service_type = "CLB"
}

`
