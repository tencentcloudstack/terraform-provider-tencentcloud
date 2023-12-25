package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfDeliveryConfigByGroupIdDataSource_basic -v
func TestAccTencentCloudTsfDeliveryConfigByGroupIdDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDeliveryConfigByGroupIdDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_delivery_config_by_group_id.delivery_config_by_group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_config_by_group_id.delivery_config_by_group_id", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_config_by_group_id.delivery_config_by_group_id", "result.0.config_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_config_by_group_id.delivery_config_by_group_id", "result.0.config_name"),
				),
			},
		},
	})
}

const testAccTsfDeliveryConfigByGroupIdDataSource = `

data "tencentcloud_tsf_delivery_config_by_group_id" "delivery_config_by_group_id" {
	group_id = "group-yrjkln9v"
}

`
