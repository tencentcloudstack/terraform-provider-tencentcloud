package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudReservedInstanceConfigsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_INTERNATION) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccReservedInstanceConfigsDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_reserved_instance_configs.configs", "config_list.#"),
				),
			},
		},
	})
}

const testAccReservedInstanceConfigsDataSource = `
data "tencentcloud_reserved_instance_configs" "configs" {
  availability_zone = "ap-guangzhou-2"
}
`
