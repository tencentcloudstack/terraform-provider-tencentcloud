/*
Provides a MySQL account resource for database management. A MySQL instance supports multiple database account.

Example Usage

```hcl
resource "tencentcloud_mysql_account" "default" {
  mysql_id = "my-test-database"
  name = "tf_account"
  password = "********"
  description = "My test account"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func resourceTencentCloudMysqlAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlAccountCreate,
		Read:   resourceTencentCloudMysqlAccountRead,
		Update: resourceTencentCloudMysqlAccountUpdate,
		Delete: resourceTencentCloudMysqlAccountDelete,

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
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validateMysqlPassword,
				Description:  "Operation password.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "--",
				ValidateFunc: validateStringLengthInRange(1, 200),
				Description:  "Database description.",
			},
		},
	}
}

func resourceTencentCloudMysqlAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_account.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		mysqlId            = d.Get("mysql_id").(string)
		accountName        = d.Get("name").(string)
		accountPassword    = d.Get("password").(string)
		accountDescription = d.Get("description").(string)
	)

	asyncRequestId, err := mysqlService.CreateAccount(ctx, mysqlId, accountName, accountPassword, accountDescription)
	if err != nil {
		return err
	}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("create account task  status is %s", taskStatus))
		}
		err = fmt.Errorf("create account task status is %s,we won't wait for it finish ,it show message:%s", ",", message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql account fail, reason:%s\n ", logId, err.Error())
		return err
	}
	d.SetId(fmt.Sprintf("%s%s%s", mysqlId, FILED_SP, accountName))
	return resourceTencentCloudMysqlAccountRead(d, meta)
}

func resourceTencentCloudMysqlAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_account.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	items := strings.Split(d.Id(), FILED_SP)

	var (
		mysqlId                      = items[0]
		accountName                  = items[1]
		accountInfo *cdb.AccountInfo = nil
	)
	var onlineHas = true
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		allAccounts, e := mysqlService.DescribeAccounts(ctx, mysqlId)
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return retryError(e)
		}
		for _, account := range allAccounts {
			if *account.User == accountName {
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
		return fmt.Errorf("Describe mysql acounts fails, reaseon %s", err.Error())
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
	_ = d.Set("name", *accountInfo.User)
	return nil
}
func resourceTencentCloudMysqlAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_account.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	items := strings.Split(d.Id(), FILED_SP)

	var (
		mysqlId     = items[0]
		accountName = items[1]
	)

	d.Partial(true)

	if d.HasChange("description") {

		asyncRequestId, err := mysqlService.ModifyAccountDescription(ctx, mysqlId, accountName, d.Get("description").(string))
		if err != nil {
			return err
		}

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("create account task  status is %s", taskStatus))
			}
			err = fmt.Errorf("modify mysql account description task status is %s,we won't wait for it finish ,it show message:%s", taskStatus, message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify mysql account description fail, reason:%s\n ", logId, err.Error())
			return err
		}

		d.SetPartial("description")
	}

	if d.HasChange("password") {

		asyncRequestId, err := mysqlService.ModifyAccountPassword(ctx, mysqlId, accountName, d.Get("password").(string))
		if err != nil {
			return err
		}

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("create account task  status is %s", taskStatus))
			}
			err = fmt.Errorf("modify mysql account password task status is %s,we won't wait for it finish ,it show message:%s", taskStatus,
				message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify mysql account password fail, reason:%s\n ", logId, err.Error())
			return err
		}
		d.SetPartial("password")
	}

	d.Partial(false)

	return nil
}

func resourceTencentCloudMysqlAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_account.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	items := strings.Split(d.Id(), FILED_SP)

	var (
		mysqlId     = items[0]
		accountName = items[1]
	)

	asyncRequestId, err := mysqlService.DeleteAccount(ctx, mysqlId, accountName)
	if err != nil {
		return err
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("create account task  status is %s", taskStatus))
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
