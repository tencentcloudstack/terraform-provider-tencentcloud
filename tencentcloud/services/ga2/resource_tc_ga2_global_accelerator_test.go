package ga2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudGa2GlobalAcceleratorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGa2GlobalAccelerator,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_global_accelerator.example", "id"),
				),
			},
			{
				Config: testAccGa2GlobalAcceleratorUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_global_accelerator.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_ga2_global_accelerator.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGa2GlobalAccelerator = `
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                 = "tf-example"
  instance_charge_type = "POSTPAID"
  description          = "tf example global accelerator"

  tags = {
    createdBy = "Terraform"
  }
}
`

const testAccGa2GlobalAcceleratorUpdate = `
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                 = "tf-example-update"
  instance_charge_type = "POSTPAID"
  description          = "tf example global accelerator"

  tags = {
    createdBy = "Terraform1"
  }
}
`
