package tpulsar_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqInstanceResource_basic -v
func TestAccTencentCloudTdmqInstanceResource_basic(t *testing.T) {
	terraformId := "tencentcloud_tdmq_instance.example"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(terraformId, "cluster_name", "tf_example"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark."),
					resource.TestCheckResourceAttr(terraformId, "tags.createdBy", "terraform"),
				),
			},
			{
				Config: testAccTdmqInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(terraformId, "cluster_name", "tf_example_update"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark update."),
					resource.TestCheckResourceAttr(terraformId, "tags.createdByUpdate", "terraformUpdate"),
				),
			},
			{
				ResourceName:      terraformId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqInstance = `
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags = {
    "createdBy" = "terraform"
  }
}
`

const testAccTdmqInstanceUpdate = `
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example_update"
  remark       = "remark update."
  tags = {
    "createdByUpdate" = "terraformUpdate"
  }
}
`
