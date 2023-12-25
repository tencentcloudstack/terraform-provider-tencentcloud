package tpulsar_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqPublishersDataSource_basic -v
func TestAccTencentCloudTdmqPublishersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqPublishersDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_publishers.publishers"),
				),
			},
		},
	})
}

const testAccTdmqPublishersDataSource = `
data "tencentcloud_tdmq_publishers" "publishers" {
  cluster_id = "pulsar-9n95ax58b9vn"
  namespace  = "keep-ns"
  topic      = "keep-topic"
  filters {
    name   = "ProducerName"
    values = ["test"]
  }
  sort {
    name  = "ProducerName"
    order = "DESC"
  }
}
`
