package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresDBInstanceClassesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresDBInstanceClassesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgres_d_b_instance_classes.d_b_instance_classes")),
			},
		},
	})
}

const testAccPostgresDBInstanceClassesDataSource = `

data "tencentcloud_postgres_d_b_instance_classes" "d_b_instance_classes" {
  zone = ""
  d_b_engine = ""
  d_b_major_version = ""
  }

`
