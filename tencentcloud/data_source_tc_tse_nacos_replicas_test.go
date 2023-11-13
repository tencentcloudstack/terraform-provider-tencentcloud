package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseNacosReplicasDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseNacosReplicasDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_nacos_replicas.nacos_replicas")),
			},
		},
	})
}

const testAccTseNacosReplicasDataSource = `

data "tencentcloud_tse_nacos_replicas" "nacos_replicas" {
  instance_id = "ins-xxxxxx"
    tags = {
    "createdBy" = "terraform"
  }
}

`
