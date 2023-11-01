package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchDiagnoseDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchDiagnoseDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_elasticsearch_diagnose.diagnose")),
			},
		},
	})
}

const testAccElasticsearchDiagnoseDataSource = `
data "tencentcloud_elasticsearch_diagnose" "diagnose" {
	instance_id = "es-5wn36he6"
	date = "20231030"
	limit = 1
}
`
