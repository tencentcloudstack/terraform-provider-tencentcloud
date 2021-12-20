package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_kms_key", &resource.Sweeper{
		Name: "tencentcloud_kms_key",
		F:    testSweepKmsKeys,
	})
}

func testSweepKmsKeys(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sharedClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(*TencentCloudClient)

	kmsService := KmsService{
		client: client.apiV3Conn,
	}

	param := make(map[string]interface{})
	param["search_key_alias"] = "tf-testacc-kms-key-"

	keys, err := kmsService.DescribeKeysByFilter(ctx, param)
	if err != nil {
		return fmt.Errorf("list KMS keys error: %s", err.Error())
	}
	for _, v := range keys {
		keyId := *v.KeyId
		if *v.KeyState == KMS_KEY_STATE_PENDINGDELETE {
			// Skip keys which are already scheduled for deletion
			continue
		}
		if *v.KeyState == KMS_KEY_STATE_ENABLED {
			if err := kmsService.DisableKey(ctx, keyId); err != nil {
				log.Printf("[ERROR] modify KMS key %s state error: %s", keyId, err.Error())
			}
		}
		if err := kmsService.DeleteKey(ctx, keyId, 7); err != nil {
			log.Printf("[ERROR] delete KMS key %s error: %s", keyId, err.Error())
		}
	}
	return nil
}

func TestAccKmsKey_basic(t *testing.T) {
	t.Parallel()
	rName := fmt.Sprintf("tf-testacc-kms-key-%s", acctest.RandString(13))
	resourceName := "tencentcloud_kms_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "key_rotation_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "key_usage", "ENCRYPT_DECRYPT"),
					resource.TestCheckResourceAttr(resourceName, "tags.test-tag", "unit-test"),
				),
			},
			{
				Config: testAccKmsKey_disabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "key_rotation_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "key_usage", "ENCRYPT_DECRYPT"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pending_delete_window_in_days", "is_enabled", "is_archived"},
			},
		},
	})
}

func TestAccKmsKey_asymmetricKey(t *testing.T) {
	t.Parallel()
	rName := fmt.Sprintf("tf-testacc-kms-key-%s", acctest.RandString(13))
	resourceName := "tencentcloud_kms_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_asymmetric(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "key_usage", "ASYMMETRIC_DECRYPT_RSA_2048"),
				),
			},
		},
	})
}

func testAccCheckKmsKeyDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	kmsService := KmsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_kms_key" {
			continue
		}

		key, err := kmsService.DescribeKeyById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if key != nil && *key.KeyState != KMS_KEY_STATE_PENDINGDELETE {
			return fmt.Errorf("[CHECK][KMS key][Destroy] check: Kms key still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckKmsKeyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[CHECK][KMS key][Exists] check: KMS key %s is not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][KMS key][Exists] check:KMS key id is not set")
		}
		kmsService := KmsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		keyId := rs.Primary.ID
		key, err := kmsService.DescribeKeyById(ctx, keyId)
		if err != nil {
			return err
		}
		if key == nil {
			return fmt.Errorf("[CHECK][KMS key][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

func testAccKmsKey_basic(rName string) string {
	return fmt.Sprintf(`
resource "tencentcloud_kms_key" "test" {
	alias = %[1]q
	description = %[1]q
  	key_rotation_enabled = true
	is_enabled = true

	tags = {
    "test-tag" = "unit-test"
  }
}
`, rName)
}

func testAccKmsKey_asymmetric(rName string) string {
	return fmt.Sprintf(`
resource "tencentcloud_kms_key" "test" {
	alias = %[1]q
	description = %[1]q
	key_usage = "ASYMMETRIC_DECRYPT_RSA_2048"
  	is_enabled = false
}
`, rName)
}

func testAccKmsKey_disabled(rName string) string {
	return fmt.Sprintf(`
resource "tencentcloud_kms_key" "test" {
 	alias = %[1]q
	description = %[1]q
	key_rotation_enabled = false
  	is_enabled = false

	tags = {
    "test-tag" = "unit-test"
  	}
}
`, rName)
}
