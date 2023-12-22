package es_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchInstanceOperationsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstanceOperationsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_elasticsearch_instance_operations.instance_operations")),
			},
		},
	})
}

const testAccElasticsearchInstanceOperationsDataSource = `
data "tencentcloud_elasticsearch_instance_operations" "instance_operations" {
	instance_id = "es-5wn36he6"
	start_time = "2018-01-01 00:00:00"
	end_time = "2023-10-31 10:12:45"
}
`
