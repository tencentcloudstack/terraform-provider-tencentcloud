package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchIndexResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchIndex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_index.index", "id"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_index.index", "instance_id", "es-nni6pm4s"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_index.index", "index_type", "normal"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_index.index", "index_name", "test-es-index"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_index.index", "index_meta_json", "{\"mappings\":{},\"settings\":{\"index.number_of_replicas\":1,\"index.number_of_shards\":1,\"index.refresh_interval\":\"30s\"}}"),
				),
			},
			{
				Config: testAccElasticsearchIndexUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_index.index", "id"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_index.index", "instance_id", "es-nni6pm4s"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_index.index", "index_type", "normal"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_index.index", "index_name", "test-es-index"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_index.index", "index_meta_json", "{\"mappings\":{},\"settings\":{\"index.number_of_replicas\":1,\"index.number_of_shards\":1,\"index.refresh_interval\":\"36s\"}}"),
				),
			},
			{
				ResourceName:      "tencentcloud_elasticsearch_index.index",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccElasticsearchIndex = `
resource "tencentcloud_elasticsearch_index" "index" {
	instance_id = "es-nni6pm4s"
	index_type = "normal"
	index_name = "test-es-index"
	index_meta_json = "{\"mappings\":{},\"settings\":{\"index.number_of_replicas\":1,\"index.number_of_shards\":1,\"index.refresh_interval\":\"30s\"}}"
}
`

const testAccElasticsearchIndexUpdate = `
resource "tencentcloud_elasticsearch_index" "index" {
	instance_id = "es-nni6pm4s"
	index_type = "normal"
	index_name = "test-es-index"
	index_meta_json = "{\"mappings\":{},\"settings\":{\"index.number_of_replicas\":1,\"index.number_of_shards\":1,\"index.refresh_interval\":\"36s\"}}"
}
`
