package tpulsar_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRoleResource_basic -v
func TestAccTencentCloudTdmqRoleResource_basic(t *testing.T) {
	terraformId := "tencentcloud_tdmq_role.example"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(terraformId, "role_name", "tf_example"),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark."),
				),
			},
			{
				Config: testAccTdmqRoleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(terraformId, "role_name", "tf_example"),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark update."),
				),
			},
		},
	})
}

const testAccTdmqRole = `
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_role" "example" {
  role_name  = "tf_example"
  cluster_id = tencentcloud_tdmq_instance.example.id
  remark     = "remark."
}
`

const testAccTdmqRoleUpdate = `
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_role" "example" {
  role_name  = "tf_example"
  cluster_id = tencentcloud_tdmq_instance.example.id
  remark     = "remark update."
}
`
