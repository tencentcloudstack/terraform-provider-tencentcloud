package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlDefaultParametersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlDefaultParametersDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_default_parameters.default_parameters")),
			},
		},
	})
}

const testAccPostgresqlDefaultParametersDataSource = `

data "tencentcloud_postgresql_default_parameters" "default_parameters" {
  db_major_version = ""
  db_engine = ""
  }

`
