/*
Provides a mysql account privilege resource to grant different access privilege to different database. A database can be granted by multiple account.

Example Usage

```hcl
resource "tencentcloud_mysql_account_privilege" "default" {
  mysql_id = "my-test-database"
  account_name= "tf_account"
  privileges = ["SELECT"]
  database_names = ["instance.name"]
}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

type resourceTencentCloudMysqlAccountPrivilegeId struct {
	MysqlId     string
	AccountName string
}

func resourceTencentCloudMysqlAccountPrivilege() *schema.Resource {

	return &schema.Resource{
		Create: resourceTencentCloudMysqlAccountPrivilegeCreate,
		Read:   resourceTencentCloudMysqlAccountPrivilegeRead,
		Update: resourceTencentCloudMysqlAccountPrivilegeUpdate,
		Delete: resourceTencentCloudMysqlAccountPrivilegeDelete,

		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Instance ID.",
			},
			"account_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Account name.",
			},
			"privileges": {
				Optional:    true,
				Type:        schema.TypeSet,
				Description: `Database permissions. Available values for Privileges: "SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "REFERENCES", "INDEX", "ALTER", "CREATE TEMPORARY TABLES", "LOCK TABLES","EXECUTE", "CREATE VIEW", "SHOW VIEW", "CREATE ROUTINE", "ALTER ROUTINE", "EVENT", and "TRIGGER".`,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAllowedStringValueIgnoreCase(MYSQL_DATABASE_PRIVILEGE),
				},
				Set: func(v interface{}) int {
					return hashcode.String(strings.ToLower(v.(string)))
				},
			},
			"database_names": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of specified database name.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
		},
	}
}

func resourceTencentCloudMysqlAccountPrivilegeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_account_privilege.create")()

	var (
		mysqlId     = d.Get("mysql_id").(string)
		accountName = d.Get("account_name").(string)
		privilegeId = resourceTencentCloudMysqlAccountPrivilegeId{MysqlId: mysqlId,
			AccountName: accountName}
	)

	privilegeIdStr, _ := json.Marshal(privilegeId)

	d.SetId(string(privilegeIdStr))

	return resourceTencentCloudMysqlAccountPrivilegeUpdate(d, meta)
}

func resourceTencentCloudMysqlAccountPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_account_privilege.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	var privilegeId resourceTencentCloudMysqlAccountPrivilegeId

	if err := json.Unmarshal([]byte(d.Id()), &privilegeId); err != nil {
		err = fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
		log.Printf("[CRITAL]%s %s\n ", logId, err.Error())
		return err
	}
	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	//check if the account is delete
	var accountInfo *cdb.AccountInfo = nil
	var onlineHas bool = true
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		accountInfos, e := mysqlService.DescribeAccounts(ctx, privilegeId.MysqlId)
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return retryError(e)
		}
		for _, account := range accountInfos {
			if *account.User == privilegeId.AccountName {
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

	dbNames := make([]string, 0, d.Get("database_names").(*schema.Set).Len())
	for _, v := range d.Get("database_names").(*schema.Set).List() {
		dbNames = append(dbNames, v.(string))
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		privileges, e := mysqlService.DescribeAccountPrivileges(ctx, privilegeId.MysqlId,
			privilegeId.AccountName, dbNames)

		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return retryError(e)
		}

		var finalPrivileges = make([]string, 0, len(privileges))

		var allowPrivileges = make(map[string]struct{})
		for _, allowPrivilege := range MYSQL_DATABASE_PRIVILEGE {
			allowPrivileges[allowPrivilege] = struct{}{}
		}

		for _, getPrivilege := range privileges {
			if getPrivilege == MYSQL_DATABASE_MUST_PRIVILEGE {
				continue
			}
			if _, ok := allowPrivileges[getPrivilege]; ok {
				finalPrivileges = append(finalPrivileges, getPrivilege)
			}
		}
		d.Set("privileges", finalPrivileges)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Describe mysql acounts privilege fails, reaseon %s", err.Error())
	}
	return nil
}

func resourceTencentCloudMysqlAccountPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_account_privilege.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var privilegeId resourceTencentCloudMysqlAccountPrivilegeId

	if err := json.Unmarshal([]byte(d.Id()), &privilegeId); err != nil {
		err = fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
		log.Printf("[CRITAL]%s %s\n ", logId, err.Error())
		return err
	}

	if d.HasChange("privileges") || d.HasChange("database_names") {
		d.Partial(true)
		d.Get("privileges").(*schema.Set).Add(MYSQL_DATABASE_MUST_PRIVILEGE)
		privileges := d.Get("privileges").(*schema.Set).List()

		log.Println(privileges)
		upperPrivileges := make([]string, len(privileges))

		for i := range privileges {
			upperPrivileges[i] = strings.ToUpper(privileges[i].(string))
		}

		dbNames := make([]string, 0, d.Get("database_names").(*schema.Set).Len())
		for _, v := range d.Get("database_names").(*schema.Set).List() {
			dbNames = append(dbNames, v.(string))
		}

		asyncRequestId, err := mysqlService.ModifyAccountPrivileges(ctx, privilegeId.MysqlId, privilegeId.AccountName, dbNames, upperPrivileges)
		if err != nil {
			return err
		}

		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("modify account privilege   task  status is %s", taskStatus))
			}
			log.Printf("modify account privilege task status is %s,we won't wait for it finish ,it show message:%s\n", taskStatus, message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify account privilege fail, reason:%s\n ", logId, err.Error())
			return err
		}
		d.SetPartial("privileges")
		d.SetPartial("db_names")
		d.Partial(false)

	}

	return resourceTencentCloudMysqlAccountPrivilegeRead(d, meta)
}

func resourceTencentCloudMysqlAccountPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_account_privilege.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var privilegeId resourceTencentCloudMysqlAccountPrivilegeId

	if err := json.Unmarshal([]byte(d.Id()), &privilegeId); err != nil {
		err = fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
		log.Printf("[CRITAL]%s %s\n ", logId, err.Error())
		return err
	}

	upperPrivileges := []string{MYSQL_DATABASE_MUST_PRIVILEGE}

	dbNames := make([]string, 0, d.Get("database_names").(*schema.Set).Len())
	for _, v := range d.Get("database_names").(*schema.Set).List() {
		dbNames = append(dbNames, v.(string))
	}

	asyncRequestId, err := mysqlService.ModifyAccountPrivileges(ctx, privilegeId.MysqlId, privilegeId.AccountName, dbNames, upperPrivileges)
	if err != nil {
		return err
	}
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("delete account privilege   task  status is %s", taskStatus))
		}
		err = fmt.Errorf("delete account privilege  task status is %s,we won't wait for it finish ,it show message:%s", taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete account privilege fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
