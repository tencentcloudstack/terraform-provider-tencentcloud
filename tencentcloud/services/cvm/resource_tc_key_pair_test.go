package cvm_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cvm_key_pair
	resource.AddTestSweepers("tencentcloud_cvm_key_pair", &resource.Sweeper{
		Name: "tencentcloud_cvm_key_pair",
		F: func(region string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			sharedClient, err := tcacctest.SharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
			}
			client := sharedClient.(tccommon.ProviderMeta)

			cvmService := svccvm.NewCvmService(client.GetAPIV3Conn())
			keyPairs, err := cvmService.DescribeKeyPairByFilter(ctx, "", "", nil)
			if err != nil {
				return fmt.Errorf("get instance list error: %s", err.Error())
			}
			for _, keyPair := range keyPairs {
				instanceId := *keyPair.KeyId
				instanceName := *keyPair.KeyName
				createTime := tccommon.StringToTime(*keyPair.CreatedTime)
				now := time.Now()
				interval := now.Sub(createTime).Minutes()

				if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
					continue
				}

				if tccommon.NeedProtect == 1 && int64(interval) < 30 {
					continue
				}

				if err = cvmService.DeleteKeyPair(ctx, instanceId); err != nil {
					log.Printf("[ERROR] sweep keyPair instance %s error: %s", instanceId, err.Error())
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudKeyPairResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: "tencentcloud_key_pair.foo",
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairPublicKeyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists("tencentcloud_key_pair.foo"),
					resource.TestCheckResourceAttr("tencentcloud_key_pair.foo", "key_name", "test_terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_key_pair.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccTencentCloudKeyPairResource_publicKey(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: "tencentcloud_key_pair.foo",
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairPublicKeyImport,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists("tencentcloud_key_pair.foo1"),
					resource.TestCheckResourceAttr("tencentcloud_key_pair.foo1", "key_name", "from_terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_key_pair.foo1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKeyPairExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("key pair %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("key pair id is not set")
		}
		cvmService := svccvm.NewCvmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		keyPair, err := cvmService.DescribeKeyPairById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				keyPair, err = cvmService.DescribeKeyPairById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if keyPair == nil {
			return fmt.Errorf("key pair is not found")
		}
		return nil
	}
}

func testAccCheckKeyPairDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cvmService := svccvm.NewCvmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_key_pair" {
			continue
		}

		keyPair, err := cvmService.DescribeKeyPairById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				keyPair, err = cvmService.DescribeKeyPairById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if keyPair != nil {
			return fmt.Errorf("key pair still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccKeyPairPublicKeyBasic = `
resource "tencentcloud_key_pair" "foo" {
  key_name   = "test_terraform"
}
`

const testAccKeyPairPublicKeyImport = `
resource "tencentcloud_key_pair" "foo1" {
  key_name   = "from_terraform"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}
`
