package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaTopicSyncReplicaDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
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
	instance_id = "ckafka-vv7wpvae"
	topic_name = "keep-topic"
}
`
