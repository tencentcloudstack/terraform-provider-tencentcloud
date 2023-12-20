package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaGroupOffsetsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaGroupOffsetsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_group_offsets.group_offsets")),
			},
		},
	})
}

const testAccCkafkaGroupOffsetsDataSource = `
data "tencentcloud_ckafka_group_offsets" "group_offsets" {
	instance_id = "ckafka-vv7wpvae"
	group = "keep-group"
}
`
