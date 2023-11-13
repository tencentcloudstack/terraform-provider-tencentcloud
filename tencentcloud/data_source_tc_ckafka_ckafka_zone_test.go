package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaCkafkaZoneDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaCkafkaZoneDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_ckafka_zone.ckafka_zone")),
			},
		},
	})
}

const testAccCkafkaCkafkaZoneDataSource = `

data "tencentcloud_ckafka_ckafka_zone" "ckafka_zone" {
  cdc_id = "id"
  }

`
