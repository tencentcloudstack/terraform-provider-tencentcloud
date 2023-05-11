package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfDeliveryConfigByGroupIdDataSource_basic -v
func TestAccTencentCloudTsfDeliveryConfigByGroupIdDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDeliveryConfigByGroupIdDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_delivery_config_by_group_id.delivery_config_by_group_id"),
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
