package cvm_test

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudReservedInstanceConfigsDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_INTERNATIONAL)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccReservedInstanceConfigsDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_reserved_instance_configs.configs"), resource.TestCheckResourceAttrSet("data.tencentcloud_reserved_instance_configs.configs", "config_list.#")),
			},
		},
	})
}

const testAccReservedInstanceConfigsDataSource_BasicCreate = `

data "tencentcloud_reserved_instance_configs" "configs" {
    availability_zone = "ap-guangzhou-2"
}

`
