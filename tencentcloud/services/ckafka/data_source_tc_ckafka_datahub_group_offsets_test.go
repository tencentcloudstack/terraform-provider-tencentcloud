package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaDatahubGroupOffsetsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaDatahubGroupOffsetsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_datahub_group_offsets.datahub_group_offsets")),
			},
		},
	})
}

const testAccCkafkaDatahubGroupOffsetsDataSource = `

data "tencentcloud_ckafka_datahub_group_offsets" "datahub_group_offsets" {
	name = "1308726196-keep-topic"
	group = "topic-lejrlafu-test"
}
`
