package es_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchInstanceLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstanceLogsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_elasticsearch_instance_logs.elasticsearch_instance_logs")),
			},
		},
	})
}

const testAccElasticsearchInstanceLogsDataSource = `
data "tencentcloud_elasticsearch_instance_logs" "elasticsearch_instance_logs" {
	instance_id = "es-5wn36he6"
}
`
