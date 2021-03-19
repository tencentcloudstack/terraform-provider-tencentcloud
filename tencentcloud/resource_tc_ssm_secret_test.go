package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudSsmSecret_basic(t *testing.T) {
	resourceName := "tencentcloud_ssm_secret.secret"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSsmSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccTencentCloudSsmSecret_basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSsmSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "secret_name", "unit-test"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test secret"),
					resource.TestCheckResourceAttr(resourceName, "init_secret.0.version_id", "v1"),
					resource.TestCheckResourceAttr(resourceName, "init_secret.0.secret_string", "123456789"),
					resource.TestCheckResourceAttrSet(resourceName, "kms_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: TestAccTencentCloudSsmSecret_modifyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSsmSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description modify"),
					resource.TestCheckResourceAttr(resourceName, "init_secret.0.version_id", "v2"),
					resource.TestCheckResourceAttr(resourceName, "init_secret.0.secret_string", "12345"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"recovery_window_in_days"},
			},
		},
	})
}

func testAccCheckSsmSecretDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	ssmService := SsmService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ssm_secret" {
			continue
		}

		secret, err := ssmService.DescribeSecretByName(ctx, rs.Primary.ID)
		if err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound" {
					return nil
				}
			}
			return err
		}
		if secret != nil && secret.status != SSM_STATUS_PENDINGDELETE {
			return fmt.Errorf("[CHECK][SSM secret][Destroy] check: SSM secret still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSsmSecretExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[CHECK][SSM secret][Exists] check: SSM secret %s is not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][SSM secret][Exists] check:SSM secret id is not set")
		}
		ssmService := SsmService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		secret, err := ssmService.DescribeSecretByName(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if secret == nil {
			return fmt.Errorf("[CHECK][SSM secret][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const TestAccTencentCloudSsmSecret_basicConfig = `
resource "tencentcloud_ssm_secret" "secret" {
  secret_name = "unit-test"
  description = "test secret"
  is_enabled = false

  init_secret {
    version_id = "v1"
    secret_string = "123456789"
  }

  tags = {
    test-tag = "test"
  }
}
`

const TestAccTencentCloudSsmSecret_modifyConfig = `
resource "tencentcloud_ssm_secret" "secret" {
  secret_name = "unit-test"
  description = "test description modify"
  is_enabled = true

  init_secret {
    version_id = "v2"
    secret_string = "12345"
  }

  tags = {
    test-tag = "test"
  }
}
`
