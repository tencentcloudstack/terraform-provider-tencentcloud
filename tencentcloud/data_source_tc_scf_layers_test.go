package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudScfLayersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfLayersDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_layers.layers")),
			},
		},
	})
}

const testAccScfLayersDataSource = `

data "tencentcloud_scf_layers" "layers" {}
`
