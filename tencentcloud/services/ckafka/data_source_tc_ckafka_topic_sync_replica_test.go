package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaTopicSyncReplicaDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaTopicSyncReplicaDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_topic_sync_replica.topic_sync_replica")),
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
