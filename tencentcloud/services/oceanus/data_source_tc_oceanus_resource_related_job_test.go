package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusResourceRelatedJobDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusResourceRelatedJobDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusResourceRelatedJobDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_resource_related_job.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_resource_related_job.example", "resource_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_resource_related_job.example", "desc_by_job_config_create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_resource_related_job.example", "resource_config_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_resource_related_job.example", "work_space_id"),
				),
			},
		},
	})
}

const testAccOceanusResourceRelatedJobDataSource = `
data "tencentcloud_oceanus_resource_related_job" "example" {
  resource_id                    = "resource-8y9lzcuz"
  desc_by_job_config_create_time = 0
  resource_config_version        = 1
  work_space_id                  = "space-2idq8wbr"
}
`
