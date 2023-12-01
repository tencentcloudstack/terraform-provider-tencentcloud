package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaGroupInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaGroupInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_group_info.group_info")),
			},
		},
	})
}

const testAccCkafkaGroupInfoDataSource = `

data "tencentcloud_ckafka_group_info" "group_info" {
	instance_id = "ckafka-vv7wpvae"
	group_list = ["keep-group"]
}
`
