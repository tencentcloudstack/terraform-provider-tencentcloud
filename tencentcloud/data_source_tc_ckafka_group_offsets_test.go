package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaGroupOffsetsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaGroupOffsetsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_group_offsets.group_offsets")),
			},
		},
	})
}

const testAccCkafkaGroupOffsetsDataSource = `

data "tencentcloud_ckafka_group_offsets" "group_offsets" {
  instance_id = "InstanceId"
  group = "groupName"
  topics = 
  search_word = "topicName"
  }

`
