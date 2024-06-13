package cvm_test

import (
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
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
			client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

			cvmService := svccvm.NewCvmService(client)
			keyPairs, err := cvmService.DescribeKeyPairByFilter(ctx, "", "", nil)
			if err != nil {
				return fmt.Errorf("get instance list error: %s", err.Error())
			}

			// add scanning resources
			var resources, nonKeepResources []*tccommon.ResourceInstance
			for _, v := range keyPairs {
				if !tccommon.CheckResourcePersist(*v.KeyName, *v.CreatedTime) {
					nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
						Id:   *v.KeyId,
						Name: *v.KeyName,
					})
				}
				resources = append(resources, &tccommon.ResourceInstance{
					Id:         *v.KeyId,
					Name:       *v.KeyName,
					CreateTime: *v.CreatedTime,
				})
			}
			tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateKeyPair")

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

func TestAccTencentCloudCvmKeyPairResource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmKeyPairResource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmKeyPairExists("tencentcloud_key_pair.foo"), resource.TestCheckResourceAttr("tencentcloud_key_pair.foo", "key_name", "test_terraform")),
			},
			{
				ResourceName:      "tencentcloud_key_pair.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmKeyPairResource_BasicCreate = `

resource "tencentcloud_key_pair" "foo" {
    key_name = "test_terraform"
}

`

func TestAccTencentCloudCvmKeyPairResource_PublicKey(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmKeyPairResource_PublicKeyCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmKeyPairExists("tencentcloud_key_pair.foo1"), resource.TestCheckResourceAttr("tencentcloud_key_pair.foo1", "key_name", "from_terraform")),
			},
			{
				ResourceName:      "tencentcloud_key_pair.foo1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmKeyPairResource_PublicKeyCreate = `

resource "tencentcloud_key_pair" "foo1" {
    key_name = "from_terraform"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}

`

func testAccCheckCvmKeyPairExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := common.GetLogId(common.ContextNil)
		ctx := context.WithValue(context.TODO(), common.LogIdKey, logId)
		service := cvm.NewCvmService(acctest.AccProvider.Meta().(common.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource `%s` is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource `%s` is not found", n)
		}
		id := rs.Primary.ID

		result, err := service.DescribeKeyPairById(ctx, id)
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("resource `%s` create failed", id)
		}
		return nil
	}
}
func testAccCheckCvmKeyPairDestroy(s *terraform.State) error {
	logId := common.GetLogId(common.ContextNil)
	ctx := context.WithValue(context.TODO(), common.LogIdKey, logId)
	service := cvm.NewCvmService(acctest.AccProvider.Meta().(common.ProviderMeta).GetAPIV3Conn())

	for _, rs := range s.RootModule().Resources {
		id := rs.Primary.ID
		if rs.Type != "tencentcloud_cvm_key_pair" {
			continue
		}
		result, err := service.DescribeKeyPairById(ctx, id)
		if err != nil {
			return err
		}
		if result != nil {
			return fmt.Errorf("resource `%s` still exist", id)
		}
	}
	return nil
}

// used in data_source_tc_key_pairs_test.go
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

// used in data_source_tc_key_pairs_test.go
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
