/*
Provides a mysql account privilege resource to grant different access privilege to different database. A database can be granted by multiple account.

Example Usage

```hcl
resource "tencentcloud_mysql_instance" "default" {
  mem_size          = 1000
  volume_size       = 25
  instance_name     = "guagua"
  engine_version    = "5.7"
  root_password     = "0153Y474"
  availability_zone = "ap-guangzhou-3"
  internet_service  = 1

}

resource "tencentcloud_mysql_account" "mysql_account2" {
  mysql_id    = tencentcloud_mysql_instance.default.id
  name        = "test11"
  password    = "test1234"
  description = "test from terraform"
}

resource "tencentcloud_mysql_privilege" "tttt" {
  mysql_id     = tencentcloud_mysql_instance.default.id
  account_name = tencentcloud_mysql_account.mysql_account2.name
  global       = ["TRIGGER"]
  database {
    privileges    = ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE"]
    database_name = "sys"
  }
  database {
    privileges    = ["SELECT"]
    database_name = "performance_schema"
  }

  table {
    privileges    = ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE"]
    database_name = "mysql"
    table_name    = "slow_log"
  }

  table {
    privileges    = ["SELECT", "INSERT", "UPDATE"]
    database_name = "mysql"
    table_name    = "user"
  }

  column {
    privileges    = ["SELECT", "INSERT", "UPDATE", "REFERENCES"]
    database_name = "mysql"
    table_name    = "user"
    column_name   = "host"
  }
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type resourceTencentCloudMysqlPrivilegeId struct {
	MysqlId     string
	AccountName string
	AccountHost string `json:"AccountHost,omitempty"`
}

func resourceTencentCloudMysqlPrivilegeHash(v interface{}) int {
	vmap := v.(map[string]interface{})
	hashMap := map[string]interface{}{}
	hashMap["database_name"] = vmap["database_name"]

	if vmap["table_name"] != nil {
		hashMap["table_name"] = vmap["table_name"]
	}
	if hashMap["column_name"] != nil {
		hashMap["column_name"] = vmap["column_name"]
	}
	slice := []string{}
	for _, v := range vmap["privileges"].(*schema.Set).List() {
		slice = append(slice, v.(string))
	}
	hashMap["privileges"] = slice
	b, _ := json.Marshal(hashMap)
	return hashcode.String(string(b))
}

func resourceTencentCloudMysqlPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlPrivilegeCreate,
		Read:   resourceTencentCloudMysqlPrivilegeRead,
		Update: resourceTencentCloudMysqlPrivilegeUpdate,
		Delete: resourceTencentCloudMysqlPrivilegeDelete,
		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Instance ID.",
			},
			"account_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) {
					if map[string]bool{
						"root":        true,
						"mysql.sys":   true,
						"tencentroot": true,
					}[v.(string)] {
						errs = append(errs, errors.New("account_name is forbidden to be "+v.(string)))
					}
					return
				},
				Description: "Account name.the forbidden value is:root,mysql.sys,tencentroot.",
			},
			"account_host": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     MYSQL_DEFAULT_ACCOUNT_HOST,
				Description: "Account host, default is `%`.",
			},
			"global": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
				Description: `Global privileges. available values for Privileges:` + strings.Join(MYSQL_GlOBAL_PRIVILEGE, ",") + ".",
			},
			"database": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Database privileges list.",
				Set:         resourceTencentCloudMysqlPrivilegeHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"privileges": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set: func(v interface{}) int {
								return hashcode.String(v.(string))
							},
							Description: `Database privilege.available values for Privileges:` + strings.Join(MYSQL_DATABASE_PRIVILEGE, ",") + ".",
						},
					},
				},
			},
			"table": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Table privileges list.",
				Set:         resourceTencentCloudMysqlPrivilegeHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Table name.",
						},
						"privileges": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set: func(v interface{}) int {
								return hashcode.String(v.(string))
							},
							Description: `Table privilege.available values for Privileges:` + strings.Join(MYSQL_TABLE_PRIVILEGE, ",") + ".",
						},
					},
				},
			},
			"column": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Column privileges list.",
				Set:         resourceTencentCloudMysqlPrivilegeHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Table name.",
						},
						"column_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Column name.",
						},
						"privileges": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set: func(v interface{}) int {
								return hashcode.String(v.(string))
							},
							Description: `Column privilege.available values for Privileges:` + strings.Join(MYSQL_COLUMN_PRIVILEGE, ",") + ".",
						},
					},
				},
			},
		},
	}
}

func (me *resourceTencentCloudMysqlPrivilegeId) update(ctx context.Context, d *schema.ResourceData, meta interface{}) error {

	if me.AccountHost == "" {
		me.AccountHost = MYSQL_DEFAULT_ACCOUNT_HOST
	}

	sp := func(str interface{}) *string { v := str.(string); return &v }

	request := cdb.NewModifyAccountPrivilegesRequest()
	request.InstanceId = &me.MysqlId

	var accountInfo = cdb.Account{User: &me.AccountName, Host: helper.String(me.AccountHost)}

	request.Accounts = []*cdb.Account{&accountInfo}

	if d != nil {
		sliceInterface := d.Get("global").(*schema.Set).List()
		if len(sliceInterface) > 0 {
			request.GlobalPrivileges = make([]*string, 0, len(sliceInterface))
			for _, v := range sliceInterface {
				ptr := sp(v)
				if !IsContains(MYSQL_GlOBAL_PRIVILEGE, *ptr) {
					return errors.New("global privileges not support " + *ptr)
				}
				request.GlobalPrivileges = append(request.GlobalPrivileges, ptr)
			}
		}

		same := map[string]bool{}

		sliceInterface = d.Get("database").(*schema.Set).List()
		if len(sliceInterface) > 0 {
			request.DatabasePrivileges = make([]*cdb.DatabasePrivilege, 0, len(sliceInterface))
			for _, v := range sliceInterface {
				vmap := v.(map[string]interface{})

				trace := *sp(vmap["database_name"])
				if same[trace] {
					return errors.New("can not assign two permissions to a database and an account," + trace)
				} else {
					same[trace] = true
				}

				p := &cdb.DatabasePrivilege{
					Database:   sp(vmap["database_name"]),
					Privileges: []*string{},
				}

				for _, privilege := range vmap["privileges"].(*schema.Set).List() {
					ptr := sp(privilege)
					if !IsContains(MYSQL_DATABASE_PRIVILEGE, *ptr) {
						return errors.New("database privileges not support:" + *ptr)
					}
					p.Privileges = append(p.Privileges, ptr)
				}
				request.DatabasePrivileges = append(request.DatabasePrivileges, p)
			}
		}

		sliceInterface = d.Get("table").(*schema.Set).List()
		if len(sliceInterface) > 0 {
			request.TablePrivileges = make([]*cdb.TablePrivilege, 0, len(sliceInterface))
			for _, v := range sliceInterface {
				vmap := v.(map[string]interface{})

				trace := *sp(vmap["database_name"]) + "." + *sp(vmap["table_name"])
				if same[trace] {
					return errors.New("can not assign two permissions to a table and an account," + trace)
				} else {
					same[trace] = true
				}

				p := &cdb.TablePrivilege{
					Database:   sp(vmap["database_name"]),
					Table:      sp(vmap["table_name"]),
					Privileges: []*string{},
				}
				for _, privilege := range vmap["privileges"].(*schema.Set).List() {
					ptr := sp(privilege)
					if !IsContains(MYSQL_TABLE_PRIVILEGE, *ptr) {
						return errors.New("table privileges not support:" + *ptr)
					}
					p.Privileges = append(p.Privileges, ptr)
				}
				request.TablePrivileges = append(request.TablePrivileges, p)
			}
		}

		sliceInterface = d.Get("column").(*schema.Set).List()
		if len(sliceInterface) > 0 {
			request.ColumnPrivileges = make([]*cdb.ColumnPrivilege, 0, len(sliceInterface))
			for _, v := range sliceInterface {
				vmap := v.(map[string]interface{})

				trace := *sp(vmap["database_name"]) + "." + *sp(vmap["table_name"]) + "." + *sp(vmap["column_name"])
				if same[trace] {
					return errors.New("can not assign two permissions to a column and an account," + trace)
				} else {
					same[trace] = true
				}

				p := &cdb.ColumnPrivilege{
					Database:   sp(vmap["database_name"]),
					Table:      sp(vmap["table_name"]),
					Column:     sp(vmap["column_name"]),
					Privileges: []*string{},
				}
				for _, privilege := range vmap["privileges"].(*schema.Set).List() {
					ptr := sp(privilege)
					if !IsContains(MYSQL_COLUMN_PRIVILEGE, *ptr) {
						return errors.New("column privileges not support:" + *ptr)
					}
					p.Privileges = append(p.Privileges, ptr)
				}
				request.ColumnPrivileges = append(request.ColumnPrivileges, p)
			}
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().ModifyAccountPrivileges(request)
	if err != nil {
		return err
	}
	if response.Response == nil || response.Response.AsyncRequestId == nil {
		return errors.New("sdk action ModifyAccountPrivileges return error,miss AsyncRequestId")
	}
	asyncRequestId := *response.Response.AsyncRequestId
	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return retryError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("modify account privilege  task  status is %s", taskStatus))
		}

		if taskStatus == MYSQL_TASK_STATUS_FAILED || taskStatus == MYSQL_TASK_STATUS_REMOVED {
			return resource.NonRetryableError(errors.New("sdk ModifyAccountPrivileges task running fail," + message))
		}
		err = fmt.Errorf("modify account privilege task status is %s,we won't wait for it finish ,it show message:%s\n", taskStatus, message)
		return resource.NonRetryableError(err)
	})
	return err
}

func resourceTencentCloudMysqlPrivilegeCreate(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_mysql_privilege.update")()

	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		mysqlId     = d.Get("mysql_id").(string)
		accountName = d.Get("account_name").(string)
		accountHost = d.Get("account_host").(string)

		privilegeId = resourceTencentCloudMysqlPrivilegeId{MysqlId: mysqlId,
			AccountName: accountName,
			AccountHost: accountHost}
	)
	privilegeIdStr, err := json.Marshal(privilegeId)
	if err != nil {
		return errors.New("json encode to id fail," + err.Error())
	}
	err = privilegeId.update(ctx, d, meta)
	if err != nil {
		return err
	}
	d.SetId(string(privilegeIdStr))
	return nil
}

func resourceTencentCloudMysqlPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_privilege.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var privilegeId resourceTencentCloudMysqlPrivilegeId

	if err := json.Unmarshal([]byte(d.Id()), &privilegeId); err != nil {
		err = fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
		log.Printf("[CRITAL]%s %s\n ", logId, err.Error())
		return err
	}

	if privilegeId.AccountHost == "" {
		privilegeId.AccountHost = MYSQL_DEFAULT_ACCOUNT_HOST
	}

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var accountInfo *cdb.AccountInfo = nil
	var onlineHas = true
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		accountInfos, err := mysqlService.DescribeAccounts(ctx, privilegeId.MysqlId)
		if err != nil {
			if sdkerr, ok := err.(*sdkError.TencentCloudSDKError); ok && sdkerr.GetCode() == "InvalidParameter" &&
				strings.Contains(sdkerr.GetMessage(), "instance not found") {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return retryError(err)
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

	request := cdb.NewDescribeAccountPrivilegesRequest()
	request.InstanceId = &privilegeId.MysqlId
	request.User = &privilegeId.AccountName
	request.Host = helper.String(privilegeId.AccountHost)

	var response *cdb.DescribeAccountPrivilegesResponse
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().DescribeAccountPrivileges(request)
		if err != nil {
			if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
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
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if !onlineHas {
		return nil
	}

	if response == nil || response.Response == nil {
		return errors.New("sdk DescribeAccountPrivileges return error,miss Response")
	}
	globals := make([]string, 0, len(response.Response.GlobalPrivileges))
	for _, v := range response.Response.GlobalPrivileges {
		globals = append(globals, *v)
	}

	databases := make([]map[string]interface{}, 0, len(response.Response.DatabasePrivileges))
	for _, v := range response.Response.DatabasePrivileges {
		privileges := make([]string, 0, len(v.Privileges))
		for _, p := range v.Privileges {
			privileges = append(privileges, *p)
		}
		m := map[string]interface{}{}
		m["database_name"] = *v.Database
		m["privileges"] = privileges
		databases = append(databases, m)
	}

	tables := make([]map[string]interface{}, 0, len(response.Response.TablePrivileges))
	for _, v := range response.Response.TablePrivileges {
		privileges := make([]string, 0, len(v.Privileges))
		for _, p := range v.Privileges {
			privileges = append(privileges, *p)
		}
		m := map[string]interface{}{}
		m["database_name"] = *v.Database
		m["table_name"] = *v.Table
		m["privileges"] = privileges
		tables = append(tables, m)
	}

	columns := make([]map[string]interface{}, 0, len(response.Response.ColumnPrivileges))

	for _, v := range response.Response.ColumnPrivileges {
		privileges := make([]string, 0, len(v.Privileges))
		for _, p := range v.Privileges {
			privileges = append(privileges, *p)
		}
		m := map[string]interface{}{}
		m["database_name"] = *v.Database
		m["table_name"] = *v.Table
		m["column_name"] = *v.Column
		m["privileges"] = privileges
		columns = append(columns, m)
	}
	_ = d.Set("global", globals)
	_ = d.Set("database", databases)
	_ = d.Set("table", tables)
	_ = d.Set("column", columns)
	_ = d.Set("mysql_id", privilegeId.MysqlId)
	_ = d.Set("account_name", privilegeId.AccountName)
	_ = d.Set("account_host", privilegeId.AccountHost)

	return nil
}

func resourceTencentCloudMysqlPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_privilege.update")()

	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		privilegeId = resourceTencentCloudMysqlPrivilegeId{}
	)
	if err := json.Unmarshal([]byte(d.Id()), &privilegeId); err != nil {
		err = fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
		log.Printf("[CRITAL]%s %s\n ", logId, err.Error())
		return err
	}
	if privilegeId.AccountHost == "" {
		privilegeId.AccountHost = MYSQL_DEFAULT_ACCOUNT_HOST
	}

	if d.HasChange("global") || d.HasChange("database") || d.HasChange("table") || d.HasChange("column") {
		err := privilegeId.update(ctx, d, meta)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTencentCloudMysqlPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_privilege.delete")()

	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		privilegeId = resourceTencentCloudMysqlPrivilegeId{}
	)
	if err := json.Unmarshal([]byte(d.Id()), &privilegeId); err != nil {
		err = fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
		log.Printf("[CRITAL]%s %s\n ", logId, err.Error())
		return err
	}
	if privilegeId.AccountHost == "" {
		privilegeId.AccountHost = MYSQL_DEFAULT_ACCOUNT_HOST
	}
	err := privilegeId.update(ctx, nil, meta)
	if err != nil {
		return err
	}
	return nil
}
