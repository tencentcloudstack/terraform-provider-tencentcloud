package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -test.run TestAccTencentCloudTeoPlan_basic -v
func TestAccTencentCloudTeoPlan_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoPlan,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_plan.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_plan.example", "plan_type", "basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_plan.example", "prepaid_plan_param.0.period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_plan.example", "prepaid_plan_param.0.renew_flag", "on"),
				),
			},
			{
				Config: testAccTeoPlanUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_plan.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_plan.example", "plan_type", "standard"),
					resource.TestCheckResourceAttr("tencentcloud_teo_plan.example", "prepaid_plan_param.0.period", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_plan.example", "prepaid_plan_param.0.renew_flag", "off"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_plan.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoPlan = `
resource "tencentcloud_teo_plan" "example" {
  plan_type = "basic"
  prepaid_plan_param {
    period     = 1
    renew_flag = "on"
  }
}
`

const testAccTeoPlanUpdate = `
resource "tencentcloud_teo_plan" "example" {
  plan_type = "standard"
  prepaid_plan_param {
    period     = 2
    renew_flag = "off"
  }
}
`
