package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqPublisherSummaryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqPublisherSummaryDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_publisher_summary.publisher_summary")),
			},
		},
	})
}

const testAccTdmqPublisherSummaryDataSource = `

data "tencentcloud_tdmq_publisher_summary" "publisher_summary" {
  cluster_id = ""
  namespace = ""
  topic = ""
        }

`
