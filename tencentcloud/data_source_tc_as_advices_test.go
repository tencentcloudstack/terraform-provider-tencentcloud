package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAsAdvicesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsAdvicesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_as_advices.advices")),
			},
		},
	})
}

const testAccAsAdvicesDataSource = `

data "tencentcloud_as_advices" "advices" {
  auto_scaling_group_ids = ["asc-lo0b94oy"]
}

`
