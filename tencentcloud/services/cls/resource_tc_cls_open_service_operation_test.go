package cls_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClsOpenServiceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsOpenServiceOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_open_service_operation.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_open_service_operation.example", "status"),
				),
			},
		},
	})
}

const testAccClsOpenServiceOperation = `

resource "tencentcloud_cls_open_service_operation" "example" {}

`
