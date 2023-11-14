package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresDescribeDefaultParametersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresDescribeDefaultParametersDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgres_describe_default_parameters.describe_default_parameters")),
			},
		},
	})
}

const testAccPostgresDescribeDefaultParametersDataSource = `

data "tencentcloud_postgres_describe_default_parameters" "describe_default_parameters" {
  d_b_major_version = ""
  d_b_engine = ""
  }

`
