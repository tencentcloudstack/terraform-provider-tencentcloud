package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchInstancePluginListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstancePluginListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_elasticsearch_instance_plugin_list.instance_plugin_list")),
			},
		},
	})
}

const testAccElasticsearchInstancePluginListDataSource = `
data "tencentcloud_elasticsearch_instance_plugin_list" "instance_plugin_list" {
	instance_id = "es-5wn36he6"
}
`
