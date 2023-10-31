package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudEsElasticsearchInstanceLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEsElasticsearchInstanceLogsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_elasticsearch_instance_logs.elasticsearch_instance_logs")),
			},
		},
	})
}

const testAccEsElasticsearchInstanceLogsDataSource = `
data "tencentcloud_elasticsearch_instance_logs" "elasticsearch_instance_logs" {
	instance_id = "es-5wn36he6"
}
`
