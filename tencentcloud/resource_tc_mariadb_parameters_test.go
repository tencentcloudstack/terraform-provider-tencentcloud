package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbParametersResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbParameters,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_parameters.parameters", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_parameters.parameters",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbParameters = `

resource "tencentcloud_mariadb_parameters" "parameters" {
  instance_id = &lt;nil&gt;
  params {
		param = &lt;nil&gt;
		value = &lt;nil&gt;

  }
}

`
