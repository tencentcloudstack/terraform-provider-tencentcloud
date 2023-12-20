package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudReservedInstanceConfigsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_INTERNATIONAL) },
		Providers: tcacctest.AccProviders,
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
