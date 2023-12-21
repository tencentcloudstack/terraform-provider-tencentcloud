package es_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchDescribeIndexListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchDescribeIndexListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_elasticsearch_describe_index_list.describe_index_list"),
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
