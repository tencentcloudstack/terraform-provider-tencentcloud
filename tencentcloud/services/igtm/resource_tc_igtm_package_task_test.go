package igtm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIgtmPackageTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIgtmPackageTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_package_task.example", "id"),
				),
			},
			{
				Config: testAccIgtmPackageTaskUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_package_task.example", "id"),
				),
			},
		},
	})
}

const testAccIgtmPackageTask = `
resource "tencentcloud_igtm_package_task" "example" {
  task_detection_quantity = 100
  auto_renew              = 2
  time_span               = 1
  auto_voucher            = 1
}
`

const testAccIgtmPackageTaskUpdate = `
resource "tencentcloud_igtm_package_task" "example" {
  task_detection_quantity = 200
  auto_renew              = 1
  time_span               = 2
  auto_voucher            = 0
}
`
