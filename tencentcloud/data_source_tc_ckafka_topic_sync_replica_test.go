package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaTopicSyncReplicaDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaTopicSyncReplicaDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_topic_sync_replica.topic_sync_replica")),
			},
		},
	})
}

const testAccCkafkaTopicSyncReplicaDataSource = `

data "tencentcloud_ckafka_topic_sync_replica" "topic_sync_replica" {
  instance_id = "InstanceId"
  topic_name = "TopicName"
  out_of_sync_replica_only = true
  }

`
