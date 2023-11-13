package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresCreateBaseBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresCreateBaseBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_create_base_backup.create_base_backup", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_create_base_backup.create_base_backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresCreateBaseBackup = `

resource "tencentcloud_postgres_create_base_backup" "create_base_backup" {
  d_b_instance_id = ""
  switch_tag = 
  tags = {
    "createdBy" = "terraform"
  }
}

`
