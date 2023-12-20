package cdb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"

	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_mysql_account
	resource.AddTestSweepers("tencentcloud_mysql_account", &resource.Sweeper{
		Name: "tencentcloud_mysql_account",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()

			service := localcdb.NewMysqlService(client)

			request := cdb.NewDescribeDBInstancesRequest()
			request.InstanceNames = []*string{
				helper.String(tcacctest.DefaultMySQLName),
			}
			response, err := client.UseMysqlClient().DescribeDBInstances(request)
			if err != nil {
				log.Printf("[CRITICAL] [%s] fail, request: %s, reason: %s", request.GetAction(), request.ToJsonString(), err.Error())
				return err
			}

			if len(response.Response.Items) == 0 {
				return nil
			}

			instance := response.Response.Items[0]
			id := instance.InstanceId

			accounts, err := service.DescribeAccounts(ctx, *id)

			if err != nil {
				return err
			}

			for i := range accounts {
				item := accounts[i]
				name := *item.User
				host := *item.Host
				created, err := time.Parse(time.RFC3339, *item.CreateTime)
				if err != nil {
					created = time.Time{}
				}
				if tcacctest.IsResourcePersist(name, &created) {
					continue
				}
				if !strings.Contains(name, "test") {
					continue
				}
				log.Printf("Will delete %s %s@%s", *id, name, host)
				_, err = service.DeleteAccount(ctx, *id, name, host)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudMysqlAccountResource_basic -v
func TestAccTencentCloudMysqlAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMysqlAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlAccount(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlAccountExists("tencentcloud_mysql_account.mysql_account"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_account.mysql_account", "mysql_id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "description", "test from terraform"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "max_user_connections", "10"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_account.mysql_account",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
			{
				Config: testAccMysqlAccountUp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlAccountExists("tencentcloud_mysql_account.mysql_account"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_account.mysql_account", "mysql_id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "description", "test from terraform"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "max_user_connections", "10"),
				),
			},
		},
	})
}

func testAccCheckMysqlAccountExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		mysqlService := localcdb.NewMysqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		split := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(split) < 2 {
			return fmt.Errorf("mysql account is not set")
		}

		var accountInfos []*cdb.AccountInfo
		var inErr, outErr error
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			accountInfos, inErr = mysqlService.DescribeAccounts(ctx, split[0])
			if inErr != nil {
				sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError)
				if ok && sdkErr.Code == localcdb.MysqlInstanceIdNotFound {
					return resource.NonRetryableError(fmt.Errorf("mysql account %s is not found", rs.Primary.ID))
				}
				return tccommon.RetryError(inErr, tccommon.InternalError)

			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
		for _, account := range accountInfos {
			if *account.User == split[1] && *account.Host == split[2] {
				return nil
			}
		}
		return fmt.Errorf("mysql account %s is not found", rs.Primary.ID)
	}
}

func testAccCheckMysqlAccountDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := localcdb.NewMysqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, r := range s.RootModule().Resources {
		if r.Type != "tencentcloud_mysql_account" {
			continue
		}

		split := strings.Split(r.Primary.ID, tccommon.FILED_SP)
		if len(split) < 2 {
			continue
		}
		instance, err := mysqlService.DescribeDBInstanceById(ctx, split[0])
		if err == nil && instance == nil {
			return nil
		}

		var accountInfos []*cdb.AccountInfo
		var inErr, outErr error
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			accountInfos, inErr = mysqlService.DescribeAccounts(ctx, split[0])
			if inErr != nil {
				sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError)
				if ok && sdkErr.Code == localcdb.MysqlInstanceIdNotFound {
					return nil
				}
				return tccommon.RetryError(inErr, tccommon.InternalError)

			}
			return nil
		})

		if outErr != nil {
			return outErr
		}
		if accountInfos == nil {
			return nil
		}
		for _, account := range accountInfos {
			if *account.User == split[1] && *account.Host == split[2] {
				return fmt.Errorf("mysql account %s is still existing", split[1])
			}
		}
	}
	return nil
}

func testAccMysqlAccount() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_mysql_account" "mysql_account" {
	mysql_id = local.mysql_id
	name    = "terraform_test"
    host = "192.168.0.%%"
	password = "Test@123456#"
	description = "test from terraform"
	max_user_connections = 10
}
	`, tcacctest.CommonPresetMysql)
}

func testAccMysqlAccountUp() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_mysql_account" "mysql_account" {
	mysql_id = local.mysql_id
	name    = "terraform_test"
    host = "192.168.1.%%"
	password = "Test@123456#"
	description = "test from terraform"
	max_user_connections = 10
}
	`, tcacctest.CommonPresetMysql)
}
