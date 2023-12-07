package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchDescribeIndexListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchDescribeIndexListDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_elasticsearch_describe_index_list.describe_index_list"),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_describe_index_list.describe_index_list", "index_meta_fields.#", "1"),
				),
			},
		},
	})
}

const testAccElasticsearchDescribeIndexListDataSource = `
data "tencentcloud_elasticsearch_describe_index_list" "describe_index_list" {
  index_type  = "normal"
  instance_id = "es-nni6pm4s"
}
`
