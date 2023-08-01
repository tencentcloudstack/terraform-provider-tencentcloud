/*
Provides a mysql account privilege resource to grant different access privilege to different database. A database can be granted by multiple account.

~> **NOTE:** It has been deprecated and replaced by  tencentcloud_mysql_privilege.

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-mysql"
  description = "mysql test"
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  slave_sync_mode   = 1
  instance_name     = "tf-example-mysql"
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [tencentcloud_security_group.security_group.id]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_mysql_account" "example" {
  mysql_id             = tencentcloud_mysql_instance.example.id
  name                 = "tf_example"
  password             = "Qwer@234"
  description          = "desc."
  max_user_connections = 10
}

resource "tencentcloud_mysql_account_privilege" "default" {
  mysql_id       = tencentcloud_mysql_instance.example.id
  account_name   = tencentcloud_mysql_account.example.name
  privileges     = ["SELECT", "INSERT", "UPDATE", "DELETE"]
  database_names = ["dbname1", "dbname2"]
}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

type resourceTencentCloudMysqlAccountPrivilegeId struct {
	MysqlId     string
	AccountName string
	AccountHost string `json:"AccountHost,omitempty"`
}

func resourceTencentCloudMysqlAccountPrivilege() *schema.Resource {

	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.26.0. Please use 'tencentcloud_mysql_privilege' instead.",
		Create:             resourceTencentCloudMysqlAccountPrivilegeCreate,
		Read:               resourceTencentCloudMysqlAccountPrivilegeRead,
		Update:             resourceTencentCloudMysqlAccountPrivilegeUpdate,
		Delete:             resourceTencentCloudMysqlAccountPrivilegeDelete,

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
			"account_host": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     MYSQL_DEFAULT_ACCOUNT_HOST,
				Description: "Account host, default is `%`.",
			},
			"privileges": {
				Optional:    true,
				Type:        schema.TypeSet,
				Description: "Database permissions. Valid values: `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP`, `REFERENCES`, `INDEX`, `ALTER`, `CREATE TEMPORARY TABLES`, `LOCK TABLES`, `EXECUTE`, `CREATE VIEW`, `SHOW VIEW`, `CREATE ROUTINE`, `ALTER ROUTINE`, `EVENT` and `TRIGGER``.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAllowedStringValueIgnoreCase(MYSQL_DATABASE_PRIVILEGE),
				},
				Set: func(v interface{}) int {
					return helper.HashString(strings.ToLower(v.(string)))
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
		accountHost = d.Get("account_host").(string)
		privilegeId = resourceTencentCloudMysqlAccountPrivilegeId{MysqlId: mysqlId,
			AccountName: accountName}
	)

	if accountHost != MYSQL_DEFAULT_ACCOUNT_HOST {
		privilegeId.AccountHost = accountHost
	}

	privilegeIdStr, _ := json.Marshal(privilegeId)

	d.SetId(string(privilegeIdStr))

	return resourceTencentCloudMysqlAccountPrivilegeUpdate(d, meta)
}

func resourceTencentCloudMysqlAccountPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_account_privilege.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var privilegeId resourceTencentCloudMysqlAccountPrivilegeId

	if err := json.Unmarshal([]byte(d.Id()), &privilegeId); err != nil {
		err = fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
		log.Printf("[CRITAL]%s %s\n ", logId, err.Error())
		return err
	}

	if privilegeId.AccountHost == "" {
		privilegeId.AccountHost = MYSQL_DEFAULT_ACCOUNT_HOST
	}

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	//check if the account is delete
	var accountInfo *cdb.AccountInfo = nil
	var onlineHas = true
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
			if *account.User == privilegeId.AccountName && *account.Host == privilegeId.AccountHost {
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
		privileges, e := mysqlService.DescribeAccountPrivileges(ctx,
			privilegeId.MysqlId,
			privilegeId.AccountName,
			privilegeId.AccountHost,
			dbNames)

		if e != nil {
			if sdkErr, ok := e.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == MysqlInstanceIdNotFound {
					onlineHas = false
				}
				if sdkErr.Code == "InvalidParameter" && strings.Contains(sdkErr.GetMessage(), "instance not found") {
					onlineHas = false
				}
				if sdkErr.Code == "InternalError.TaskError" && strings.Contains(sdkErr.Message, "User does not exist") {
					onlineHas = false
				}
				if !onlineHas {
					return nil
				}
			}
			return retryError(e)
		}

		if !onlineHas {
			d.SetId("")
			return nil
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
		_ = d.Set("privileges", finalPrivileges)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Describe mysql acounts privilege fails, reaseon %s", err.Error())
	}
	_ = d.Set("account_name", privilegeId.AccountName)
	_ = d.Set("account_host", privilegeId.AccountHost)
	_ = d.Set("mysql_id", privilegeId.MysqlId)
	return nil
}

func resourceTencentCloudMysqlAccountPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_account_privilege.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var privilegeId resourceTencentCloudMysqlAccountPrivilegeId

	if err := json.Unmarshal([]byte(d.Id()), &privilegeId); err != nil {
		err = fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
		log.Printf("[CRITAL]%s %s\n ", logId, err.Error())
		return err
	}

	if privilegeId.AccountHost == "" {
		privilegeId.AccountHost = MYSQL_DEFAULT_ACCOUNT_HOST
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

		asyncRequestId, err := mysqlService.ModifyAccountPrivileges(ctx,
			privilegeId.MysqlId,
			privilegeId.AccountName,
			privilegeId.AccountHost,
			dbNames,
			upperPrivileges)

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
				return resource.RetryableError(fmt.Errorf("modify account privilege   task  status is %s", taskStatus))
			}
			if taskStatus == MYSQL_TASK_STATUS_FAILED || taskStatus == MYSQL_TASK_STATUS_REMOVED {
				return resource.NonRetryableError(errors.New("sdk ModifyAccountPrivileges task running fail," + message))
			}
			err = fmt.Errorf("modify account privilege task status is %s,we won't wait for it finish ,it show message:%s\n", taskStatus, message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify account privilege fail, reason:%s\n ", logId, err.Error())
			return err
		}

		d.Partial(false)

	}

	return resourceTencentCloudMysqlAccountPrivilegeRead(d, meta)
}

func resourceTencentCloudMysqlAccountPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_account_privilege.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var privilegeId resourceTencentCloudMysqlAccountPrivilegeId

	if err := json.Unmarshal([]byte(d.Id()), &privilegeId); err != nil {
		err = fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
		log.Printf("[CRITAL]%s %s\n ", logId, err.Error())
		return err
	}

	if privilegeId.AccountHost == "" {
		privilegeId.AccountHost = MYSQL_DEFAULT_ACCOUNT_HOST
	}

	upperPrivileges := []string{MYSQL_DATABASE_MUST_PRIVILEGE}

	dbNames := make([]string, 0, d.Get("database_names").(*schema.Set).Len())
	for _, v := range d.Get("database_names").(*schema.Set).List() {
		dbNames = append(dbNames, v.(string))
	}

	asyncRequestId, err := mysqlService.ModifyAccountPrivileges(ctx,
		privilegeId.MysqlId,
		privilegeId.AccountName,
		privilegeId.AccountHost,
		dbNames,
		upperPrivileges)

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
