package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbDbImportResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbDbImport,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_db_import.db_import", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_db_import.db_import",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbDbImport = `

resource "tencentcloud_cdb_db_import" "db_import" {
  instance_id = ""
  user = ""
  file_name = ""
  password = ""
  db_name = ""
  cos_url = ""
}

`
