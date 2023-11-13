package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaDatahubGroupOffsetsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaDatahubGroupOffsetsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_datahub_group_offsets.datahub_group_offsets")),
			},
		},
	})
}

const testAccCkafkaDatahubGroupOffsetsDataSource = `

data "tencentcloud_ckafka_datahub_group_offsets" "datahub_group_offsets" {
  name = "1300xxxx-topicName"
  group = "datahub-task-lzp7qb7e"
  search_word = ""
  }

`
