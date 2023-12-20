package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaInstancesDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCkafkaUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceCkafkaInstances,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_instances.foo", "instance_list.0.instance_id", tcacctest.DefaultKafkaInstanceId),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceCkafkaInstances = tcacctest.DefaultKafkaVariable + `
data "tencentcloud_ckafka_instances" "foo" {
	instance_ids=[var.instance_id]
}
`
