package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaVersionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaVersionDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_version.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_version.example", "kafka_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_version.example", "cur_broker_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_version.example", "latest_broker_versions.#"),
				),
			},
		},
	})
}

const testAccCkafkaVersionDataSource = `
data "tencentcloud_ckafka_version" "example" {
  instance_id = "ckafka-8j4raxv8"
}
`
