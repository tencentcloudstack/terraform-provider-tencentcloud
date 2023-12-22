package ssm_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcssm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssm"

	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_ssm_secret
	resource.AddTestSweepers("tencentcloud_ssm_secret", &resource.Sweeper{
		Name: "tencentcloud_ssm_secret",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svcssm.NewSsmService(client)

			secrets, err := service.DescribeSecretsByFilter(ctx, nil)

			if err != nil {
				return err
			}

			for i := range secrets {
				ss := secrets[i]
				name := *ss.SecretName
				createTime := ss.CreateTime
				created := time.Time{}
				if createTime != nil {
					created = time.Unix(int64(*createTime), 0)
				}
				if tcacctest.IsResourcePersist(name, &created) {
					continue
				}
				err = service.DisableSecret(ctx, name)
				if err != nil {
					continue
				}
				err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
					err := service.DeleteSecret(ctx, name, 0)
					if err != nil {
						return tccommon.RetryError(err, ssm.FAILEDOPERATION)
					}
					return nil
				})
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudSsmSecret_basic -v
func TestAccTencentCloudSsmSecret_basic(t *testing.T) {
	t.Parallel()
	resourceName := "tencentcloud_ssm_secret.example"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSsmSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccTencentCloudSsmSecretBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSsmSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "secret_name", "tf-example"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "desc."),
					resource.TestCheckResourceAttrSet(resourceName, "kms_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: TestAccTencentCloudSsmSecretBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSsmSecretExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "update desc."),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	ssmService := svcssm.NewSsmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
		if secret != nil && secret.Status() != svcssm.SSM_STATUS_PENDINGDELETE {
			return fmt.Errorf("[CHECK][SSM secret][Destroy] check: SSM secret still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSsmSecretExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[CHECK][SSM secret][Exists] check: SSM secret %s is not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][SSM secret][Exists] check:SSM secret id is not set")
		}
		ssmService := svcssm.NewSsmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const TestAccTencentCloudSsmSecretBasic = `
resource "tencentcloud_ssm_secret" "example" {
  secret_name = "tf-example"
  description = "desc."
  is_enabled  = false

  tags = {
    createBy = "terraform"
  }
}
`

const TestAccTencentCloudSsmSecretBasicUpdate = `
resource "tencentcloud_ssm_secret" "example" {
  secret_name = "tf-example"
  description = "update desc."
  is_enabled  = true

  tags = {
    createBy = "terraform"
  }
}
`
