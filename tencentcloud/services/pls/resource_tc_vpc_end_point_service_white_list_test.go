package pls_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcEndPointServiceWhiteListResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointServiceWhiteList,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point_service_white_list.end_point_service_white_list", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_end_point_service_white_list.end_point_service_white_list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcEndPointServiceWhiteList = `

resource "tencentcloud_vpc_end_point_service_white_list" "end_point_service_white_list" {
  user_uin = "100020512675"
  end_point_service_id = "vpcsvc-98jddhcz"
  description = "terraform for test"
}

`
