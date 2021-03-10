package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKmsExternalKey_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-testacc-kms-key-%s", acctest.RandString(13))
	resourceName := "tencentcloud_kms_external_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsExternalKey_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.test-tag", "unit-test"),
				),
			},
			{
				Config: testAccKmsExternalKey_disabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pending_delete_window_in_days", "key_material_base64", "wrapping_algorithm", "is_enabled", "is_archived"},
			},
		},
	})
}

func testAccKmsExternalKey_basic(rName string) string {
	return fmt.Sprintf(`
resource "tencentcloud_kms_external_key" "test" {
	alias = %[1]q
	description = %[1]q
	wrapping_algorithm = "RSAES_PKCS1_V1_5"
	key_material_base64 = "MTIzMTIzMTIzMTIzMTIzQQ=="

	tags = {
    "test-tag" = "unit-test"
  }
}
`, rName)
}

func testAccKmsExternalKey_disabled(rName string) string {
	return fmt.Sprintf(`
resource "tencentcloud_kms_external_key" "test" {
 	alias = %[1]q
	description = %[1]q
	wrapping_algorithm = "RSAES_PKCS1_V1_5"
	key_material_base64 = "MTIzMTIzMTIzMTIzMTIzQQ=="
  	is_enabled = false

	tags = {
    "test-tag" = "unit-test"
  }
}
`, rName)
}
