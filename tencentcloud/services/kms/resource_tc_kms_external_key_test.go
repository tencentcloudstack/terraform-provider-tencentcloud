package kms_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKmsExternalKey_basic(t *testing.T) {
	t.Parallel()
	rName := fmt.Sprintf("tf-testacc-kms-key-%s", acctest.RandString(13))
	resourceName := "tencentcloud_kms_external_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckKmsKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsExternalKey_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "alias"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "key_state", "PendingImport"),
					resource.TestCheckResourceAttr(resourceName, "tags.test-tag", "unit-test"),
				),
			},
			{
				Config: testAccKmsExternalKey_import(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_state", "Enabled"),
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

	tags = {
    "test-tag" = "unit-test"
  }
}
`, rName)
}

func testAccKmsExternalKey_import(rName string) string {
	return fmt.Sprintf(`
resource "tencentcloud_kms_external_key" "test" {
 	alias = %[1]q
	description = %[1]q
	wrapping_algorithm = "RSAES_PKCS1_V1_5"
	key_material_base64 = "MTIzMTIzMTIzMTIzMTIzQQ=="
  	is_enabled = true

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
