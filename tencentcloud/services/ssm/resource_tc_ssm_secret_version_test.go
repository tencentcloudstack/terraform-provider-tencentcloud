package ssm_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcssm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudSsmSecretVersion_basic -v
func TestAccTencentCloudSsmSecretVersion_basic(t *testing.T) {
	t.Parallel()
	resourceV1Name := "tencentcloud_ssm_secret_version.v1"
	resourceV2Name := "tencentcloud_ssm_secret_version.v2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSsmSecretVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccTencentCloudSsmSecretVersionBinaryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSsmSecretVersionExists(resourceV1Name),
					resource.TestCheckResourceAttr(resourceV1Name, "secret_name", "tf-example-secret"),
					resource.TestCheckResourceAttr(resourceV1Name, "version_id", "v1"),
					resource.TestCheckResourceAttr(resourceV1Name, "secret_binary", "MTIzMTIzMTIzMTIzMTIzQQ=="),
				),
			},
			{
				ResourceName:      resourceV1Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: TestAccTencentCloudSsmSecretVersionStringConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSsmSecretVersionExists(resourceV2Name),
					resource.TestCheckResourceAttr(resourceV2Name, "secret_name", "tf-example-secret"),
					resource.TestCheckResourceAttr(resourceV2Name, "version_id", "v2"),
					resource.TestCheckResourceAttr(resourceV2Name, "secret_string", "this is secret string"),
				),
			},
			{
				ResourceName:      resourceV2Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSsmSecretVersionDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	ssmService := svcssm.NewSsmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ssm_secret_version" {
			continue
		}

		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[CHECK][SSM secret version][Exists] check: SSM secret version %s is not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][SSM secret version][Exists] check:SSM secret version id is not set")
		}
		ssmService := svcssm.NewSsmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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

const TestAccTencentCloudSsmSecretVersionBinaryConfig = `
resource "tencentcloud_ssm_secret" "example" {
  secret_name             = "tf-example-secret"
  description             = "desc."
  recovery_window_in_days = 0
  is_enabled              = true

  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ssm_secret_version" "v1" {
  secret_name   = tencentcloud_ssm_secret.example.secret_name
  version_id    = "v1"
  secret_binary = "MTIzMTIzMTIzMTIzMTIzQQ=="
}
`

const TestAccTencentCloudSsmSecretVersionStringConfig = `
resource "tencentcloud_ssm_secret" "example" {
  secret_name             = "tf-example-secret"
  description             = "desc."
  recovery_window_in_days = 0
  is_enabled              = true

  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ssm_secret_version" "v2" {
  secret_name   = tencentcloud_ssm_secret.example.secret_name
  version_id    = "v2"
  secret_string = "this is secret string"
}
`
