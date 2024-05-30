package cvm_test

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudReservedInstancesDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_INTERNATIONAL)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccReservedInstancesDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_reserved_instances.instances"), resource.TestCheckResourceAttrSet("data.tencentcloud_reserved_instances.instances", "reserved_instance_list.#")),
			},
		},
	})
}

const testAccReservedInstancesDataSource_BasicCreate = `

data "tencentcloud_reserved_instances" "instances" {
    availability_zone = "ap-guangzhou-3"
}

`
