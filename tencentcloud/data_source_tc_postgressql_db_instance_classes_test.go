package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgressqlDbInstanceClassesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgressqlDbInstanceClassesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgressql_db_instance_classes.db_instance_classes")),
			},
		},
	})
}

const testAccPostgressqlDbInstanceClassesDataSource = `

data "tencentcloud_postgressql_db_instance_classes" "db_instance_classes" {
  zone = ""
  d_b_engine = ""
  d_b_major_version = ""
  }

`
