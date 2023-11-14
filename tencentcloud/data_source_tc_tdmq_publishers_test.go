package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqPublishersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqPublishersDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_publishers.publishers")),
			},
		},
	})
}

const testAccTdmqPublishersDataSource = `

data "tencentcloud_tdmq_publishers" "publishers" {
  cluster_id = ""
  namespace = ""
  topic = ""
  filters {
		name = ""
		values = 

  }
  sort {
		name = ""
		order = ""

  }
  }

`
