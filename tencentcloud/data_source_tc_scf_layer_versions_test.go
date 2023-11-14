package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfLayerVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfLayerVersionsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_layer_versions.layer_versions")),
			},
		},
	})
}

const testAccScfLayerVersionsDataSource = `

data "tencentcloud_scf_layer_versions" "layer_versions" {
  layer_name = ""
  compatible_runtime = 
  }

`
