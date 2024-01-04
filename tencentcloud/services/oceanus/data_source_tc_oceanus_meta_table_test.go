package oceanus_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusMetaTableDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusMetaTableDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusMetaTableDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_meta_table.example"),
					resource.TestCheckResourceAttr("data.tencentcloud_oceanus_meta_table.example", "work_space_id", "space-6w8eab6f"),
					resource.TestCheckResourceAttr("data.tencentcloud_oceanus_meta_table.example", "catalog", "_dc"),
					resource.TestCheckResourceAttr("data.tencentcloud_oceanus_meta_table.example", "database", "_db"),
					resource.TestCheckResourceAttr("data.tencentcloud_oceanus_meta_table.example", "table", "tf_table"),
				),
			},
		},
	})
}

const testAccOceanusMetaTableDataSource = `
data "tencentcloud_oceanus_meta_table" "example" {
  work_space_id = "space-6w8eab6f"
  catalog       = "_dc"
  database      = "_db"
  table         = "tf_table"
}
`
