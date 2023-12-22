package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfDeliveryConfigsDataSource_basic -v
func TestAccTencentCloudTsfDeliveryConfigsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDeliveryConfigsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_delivery_configs.delivery_configs"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_configs.delivery_configs", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_configs.delivery_configs", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_configs.delivery_configs", "result.0.content.0.config_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_configs.delivery_configs", "result.0.content.0.config_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_configs.delivery_configs", "result.0.content.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_configs.delivery_configs", "result.0.content.0.enable_auth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_configs.delivery_configs", "result.0.content.0.kafka_infos.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_configs.delivery_configs", "result.0.content.0.kafka_infos.0.path.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_configs.delivery_configs", "result.0.content.0.kafka_infos.0.topic"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_configs.delivery_configs", "result.0.content.0.enable_global_line_rule"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_delivery_configs.delivery_configs", "result.0.content.0.line_rule"),
				),
			},
		},
	})
}

const testAccTsfDeliveryConfigsDataSource = `

data "tencentcloud_tsf_delivery_configs" "delivery_configs" {
  search_word = "test"
}

`
