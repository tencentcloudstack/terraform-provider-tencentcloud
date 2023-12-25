package tpulsar_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqPublisherSummaryDataSource_basic -v
func TestAccTencentCloudTdmqPublisherSummaryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqPublisherSummaryDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_publisher_summary.publisher_summary"),
				),
			},
		},
	})
}

const testAccTdmqPublisherSummaryDataSource = `
data "tencentcloud_tdmq_publisher_summary" "publisher_summary" {
  cluster_id = "pulsar-9n95ax58b9vn"
  namespace  = "keep-ns"
  topic      = "keep-topic"
}
`
