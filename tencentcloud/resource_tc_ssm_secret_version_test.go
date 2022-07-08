package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudSsmSecretVersion_basic(t *testing.T) {
	t.Parallel()
	resourceName := "tencentcloud_ssm_secret_version.v1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSsmSecretVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccTencentCloudSsmSecretVersion_basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSsmSecretVersionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "secret_name", "unit-test-for-version"),
					resource.TestCheckResourceAttr(resourceName, "version_id", "v1"),
					resource.TestCheckResourceAttr(resourceName, "secret_binary", "MTIzMTIzMTIzMTIzMTIzQQ=="),
				),
			},
			{
				Config: TestAccTencentCloudSsmSecretVersion_secretStringConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSsmSecretVersionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "secret_name", "unit-test-for-version"),
					resource.TestCheckResourceAttr(resourceName, "version_id", "v1"),
					resource.TestCheckResourceAttr(resourceName, "secret_string", "123456"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSsmSecretVersionDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	ssmService := SsmService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ssm_secret_version" {
			continue
		}

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		secretName := items[0]
		versionId := items[1]
		secretVersion, err := ssmService.DescribeSecretVersion(ctx, secretName, versionId)
		if err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound" {
					return nil
				}
			}
			return err
		}
		if secretVersion != nil {
			return fmt.Errorf("[CHECK][SSM secret version][Destroy] check: SSM secret version still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSsmSecretVersionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[CHECK][SSM secret version][Exists] check: SSM secret version %s is not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][SSM secret version][Exists] check:SSM secret version id is not set")
		}
		ssmService := SsmService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		secretName := items[0]
		versionId := items[1]
		secretVersion, err := ssmService.DescribeSecretVersion(ctx, secretName, versionId)
		if err != nil {
			return err
		}
		if secretVersion == nil {
			return fmt.Errorf("[CHECK][SSM secret version][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const TestAccTencentCloudSsmSecretVersion_basicConfig = `
resource "tencentcloud_ssm_secret" "secret" {
  secret_name = "unit-test-for-version"
  description = "test secret"

  tags = {
    test-tag = "test"
  }
}

resource "tencentcloud_ssm_secret_version" "v1" {
  secret_name = tencentcloud_ssm_secret.secret.secret_name
  version_id = "v1"
  secret_binary = "MTIzMTIzMTIzMTIzMTIzQQ=="
}
`

const TestAccTencentCloudSsmSecretVersion_secretStringConfig = `
resource "tencentcloud_ssm_secret" "secret" {
  secret_name = "unit-test-for-version"
  description = "test secret"

  tags = {
    test-tag = "test"
  }
}

resource "tencentcloud_ssm_secret_version" "v1" {
  secret_name = tencentcloud_ssm_secret.secret.secret_name
  version_id = "v1"
  secret_string = "123456"
}
`
