package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCbsStorage_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsStorageConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.my_storage"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.my_storage", "storage_name", "testAccCbsStorageTest"),
				),
			},
			{
				Config: testAccCbsStorageConfigChanged,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.my_storage"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.my_storage", "storage_name", "testAccCbsStorageTest-2"),
				),
			},
		},
	})
}

func testAccCheckStorageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		provider := testAccProvider
		// Ignore if Meta is empty, this can happen for validation providers
		if provider.Meta() == nil {
			return fmt.Errorf("Provider Meta is nil")
		}

		client := provider.Meta().(*TencentCloudClient).commonConn
		_, _, err := describeCbsStorage(rs.Primary.ID, client)

		if err == nil {
			return err
		}
		return fmt.Errorf("Error finding CBS Storage %s", rs.Primary.ID)
	}
}

const testAccCbsStorageConfig = `
resource "tencentcloud_cbs_storage" "my_storage" {
  availability_zone = "ap-guangzhou-3"
  storage_size      = 50
  storage_type      = "cloudPremium"
  period            = 1
  storage_name      = "testAccCbsStorageTest"
}
`

const testAccCbsStorageConfigChanged = `
resource "tencentcloud_cbs_storage" "my_storage" {
  availability_zone = "ap-guangzhou-3"
  storage_size      = 50
  storage_type      = "cloudPremium"
  period            = 1
  storage_name      = "testAccCbsStorageTest-2"
}
`
