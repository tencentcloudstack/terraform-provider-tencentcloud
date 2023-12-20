package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func ResourceTencentCloudMysqlAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlAccountCreate,
		Read:   resourceTencentCloudMysqlAccountRead,
		Update: resourceTencentCloudMysqlAccountUpdate,
		Delete: resourceTencentCloudMysqlAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Instance ID to which the account belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Account name.",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     MYSQL_DEFAULT_ACCOUNT_HOST,
				Description: "Account host, default is `%`.",
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: tccommon.ValidateMysqlPassword,
				Description:  "Operation password.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "--",
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 200),
				Description:  "Database description.",
			},
			"max_user_connections": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of available connections for a new account, the default value is 10240, and the maximum value that can be set is 10240.",
			},
		},
	}
}

func resourceTencentCloudMysqlAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_account.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		mysqlId            = d.Get("mysql_id").(string)
		accountName        = d.Get("name").(string)
		accountHost        = d.Get("host").(string)
		accountPassword    = d.Get("password").(string)
		accountDescription = d.Get("description").(string)
		maxUserConnections = int64(d.Get("max_user_connections").(int))
	)

	asyncRequestId, err := mysqlService.CreateAccount(ctx, mysqlId, accountName, accountHost, accountPassword, accountDescription, maxUserConnections)
	if err != nil {
		return err
	}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s create account %s.%s task  status is %s", mysqlId, accountName, accountHost, taskStatus))
		}
		err = fmt.Errorf("%s create account task status is %s,we won't wait for it finish ,it show message:%s", mysqlId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql account fail, reason:%s\n ", logId, err.Error())
		return err
	}

	resourceId := fmt.Sprintf("%s%s%s", mysqlId, tccommon.FILED_SP, accountName)

	if accountHost != MYSQL_DEFAULT_ACCOUNT_HOST {
		resourceId += tccommon.FILED_SP + accountHost
	}

	d.SetId(resourceId)

	return resourceTencentCloudMysqlAccountRead(d, meta)
}

func resourceTencentCloudMysqlAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_account.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	items := strings.Split(d.Id(), tccommon.FILED_SP)

	var (
		mysqlId                      = items[0]
		accountName                  = items[1]
		accountHost                  = MYSQL_DEFAULT_ACCOUNT_HOST
		accountInfo *cdb.AccountInfo = nil
	)

	if len(items) == 3 {
		accountHost = items[2]
	}

	var onlineHas = true
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		allAccounts, e := mysqlService.DescribeAccounts(ctx, mysqlId)
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return tccommon.RetryError(e)
		}
		for _, account := range allAccounts {
			if *account.User == accountName && *account.Host == accountHost {
				accountInfo = account
				break
			}
		}
		if accountInfo == nil {
			d.SetId("")
			onlineHas = false
			return nil
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Describe mysql acounts fails, reason %s", err.Error())
	}
	if !onlineHas {
		return nil
	}
	if *accountInfo.Notes == "" {
		_ = d.Set("description", "--")
	} else {
		_ = d.Set("description", *accountInfo.Notes)
	}
	_ = d.Set("mysql_id", mysqlId)
	_ = d.Set("host", *accountInfo.Host)
	_ = d.Set("name", *accountInfo.User)
	_ = d.Set("max_user_connections", *accountInfo.MaxUserConnections)

	return nil
}
func resourceTencentCloudMysqlAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_account.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	items := strings.Split(d.Id(), tccommon.FILED_SP)

	var (
		mysqlId     = items[0]
		accountName = items[1]
		accountHost = MYSQL_DEFAULT_ACCOUNT_HOST
	)

	if len(items) == 3 {
		accountHost = items[2]
	}

	d.Partial(true)

	if d.HasChange("description") {

		asyncRequestId, err := mysqlService.ModifyAccountDescription(ctx, mysqlId, accountName, accountHost, d.Get("description").(string))
		if err != nil {
			return err
		}

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("%s modify account  description %s.%s task  status is %s", mysqlId, accountName, accountHost, taskStatus))
			}
			err = fmt.Errorf("modify mysql account description task status is %s,we won't wait for it finish ,it show message:%s", taskStatus, message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify mysql account description fail, reason:%s\n ", logId, err.Error())
			return err
		}

	}

	if d.HasChange("password") {

		asyncRequestId, err := mysqlService.ModifyAccountPassword(ctx, mysqlId, accountName, accountHost, d.Get("password").(string))
		if err != nil {
			return err
		}

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("%s modify mysql account password %s.%s task  status is %s", mysqlId, accountName, accountHost, taskStatus))
			}
			err = fmt.Errorf("modify mysql account password task status is %s,we won't wait for it finish ,it show message:%s", taskStatus,
				message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify mysql account password fail, reason:%s\n ", logId, err.Error())
			return err
		}

	}

	if d.HasChange("max_user_connections") {
		var maxUserConnections int64
		if v, ok := d.GetOkExists("max_user_connections"); ok {
			maxUserConnections = int64(v.(int))
		}
		asyncRequestId, err := mysqlService.ModifyAccountMaxUserConnections(ctx, mysqlId, accountName, accountHost, maxUserConnections)
		if err != nil {
			return err
		}

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("%s modify mysql account maxUserConnections %s task  status is %s", mysqlId, accountName, taskStatus))
			}
			err = fmt.Errorf("modify mysql account maxUserConnections task status is %s,we won't wait for it finish ,it show message:%s", taskStatus, message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify mysql account maxUserConnections fail, reason:%s\n ", logId, err.Error())
			return err
		}

	}

	if d.HasChange("host") {
		oldHost, newHost := d.GetChange("host")
		asyncRequestId, err := mysqlService.ModifyAccountHost(ctx, mysqlId, accountName, oldHost.(string), newHost.(string))
		if err != nil {
			return err
		}

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("%s modify account  host %s.%s task  status is %s", mysqlId, accountName, accountHost, taskStatus))
			}
			err = fmt.Errorf("modify mysql account host task status is %s,we won't wait for it finish ,it show message:%s", taskStatus, message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify mysql account host fail, reason:%s\n ", logId, err.Error())
			return err
		}

		resourceId := fmt.Sprintf("%s%s%s", mysqlId, tccommon.FILED_SP, accountName)

		if newHost.(string) != MYSQL_DEFAULT_ACCOUNT_HOST {
			resourceId += tccommon.FILED_SP + newHost.(string)
		}

		d.SetId(resourceId)
	}

	d.Partial(false)

	return nil
}

func resourceTencentCloudMysqlAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_account.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	items := strings.Split(d.Id(), tccommon.FILED_SP)

	var (
		mysqlId     = items[0]
		accountName = items[1]
		accountHost = MYSQL_DEFAULT_ACCOUNT_HOST
	)
	if len(items) == 3 {
		accountHost = items[2]
	}

	asyncRequestId, err := mysqlService.DeleteAccount(ctx, mysqlId, accountName, accountHost)
	if err != nil {
		return err
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s delete mysql account %s.%s task  status is %s", mysqlId, accountName, accountHost, taskStatus))
		}
		err = fmt.Errorf("delete mysql account  task status is %s,we won't wait for it finish ,it show message:%s", taskStatus,
			message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
