/*
Provides a resource to create a mariadb account_privileges

Example Usage

```hcl
resource "tencentcloud_mariadb_account_privileges" "account_privileges" {
  instance_id = "tdsql-9vqvls95"
  accounts {
		user = "keep-modify-privileges"
		host = "127.0.0.1"
  }
  global_privileges = ["ALTER", "CREATE", "DELETE", "SELECT", "UPDATE", "DROP"]
}
```

Import

mariadb account_privileges can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_account_privileges.account_privileges account_privileges_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbAccountPrivileges() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbAccountPrivilegesCreate,
		Read:   resourceTencentCloudMariadbAccountPrivilegesRead,
		Update: resourceTencentCloudMariadbAccountPrivilegesUpdate,
		Delete: resourceTencentCloudMariadbAccountPrivilegesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},
			"accounts": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "account information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "user name.",
						},
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "user host.",
						},
					},
				},
			},
			"global_privileges": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Global permission. Valid values of `GlobalPrivileges`: `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `PROCESS`, `DROP`, `REFERENCES`, `INDEX`, `ALTER`, `SHOW DATABASES`, `CREATE TEMPORARY TABLES`, `LOCK TABLES`, `EXECUTE`, `CREATE VIEW`, `SHOW VIEW`, `CREATE ROUTINE`, `ALTER ROUTINE`, `EVENT`, `TRIGGER`.Note: if the parameter is left empty, no change will be made to the granted global permissions. To clear the granted global permissions, set the parameter to an empty array.",
			},
			"database_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Database permission. Valid values of `Privileges`: `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP`, `REFERENCES`, `INDEX`, `ALTER`, `CREATE TEMPORARY TABLES`, `LOCK TABLES`, `EXECUTE`, `CREATE VIEW`, `SHOW VIEW`, `CREATE ROUTINE`, `ALTER ROUTINE`, `EVENT`, `TRIGGER`.Note: if the parameter is left empty, no change will be made to the granted database permissions. To clear the granted database permissions, set `Privileges` to an empty array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"privileges": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Permission information.",
						},
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
					},
				},
			},
			"table_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "`SELECT`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP`, `REFERENCES`, `INDEX`, `ALTER`, `CREATE VIEW`, `SHOW VIEW`, `TRIGGER`.Note: if the parameter is not passed in, no change will be made to the granted table permissions. To clear the granted table permissions, set `Privileges` to an empty array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"table": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Table name.",
						},
						"privileges": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Permission information.",
						},
					},
				},
			},
			"column_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Column permission. Valid values of `Privileges`: `SELECT`, `INSERT`, `UPDATE`, `REFERENCES`.Note: if the parameter is left empty, no change will be made to the granted column permissions. To clear the granted column permissions, set `Privileges` to an empty array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"table": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Table name.",
						},
						"column": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Column name.",
						},
						"privileges": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Permission information.",
						},
					},
				},
			},
			"view_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Database view permission. Valid values of `Privileges`: `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP`, `REFERENCES`, `INDEX`, `ALTER`, `CREATE VIEW`, `SHOW VIEW`, `TRIGGER`.Note: if the parameter is not passed in, no change will be made to the granted view permissions. To clear the granted view permissions, set `Privileges` to an empty array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"view": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "View name.",
						},
						"privileges": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Permission information.",
						},
					},
				},
			},
			"function_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Database function permissions. Valid values of `Privileges`: `ALTER ROUTINE`, `EXECUTE`.Note: if the parameter is not passed in, no change will be made to the granted function permissions. To clear the granted function permissions, set `Privileges` to an empty array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"function_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Function name.",
						},
						"privileges": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Permission information.",
						},
					},
				},
			},
			"procedure_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Database stored procedure permission. Valid values of `Privileges`: `ALTER ROUTINE`, `EXECUTE`.Note: if the parameter is not passed in, no change will be made to the granted stored procedure permissions. To clear the granted stored procedure permissions, set `Privileges` to an empty array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"procedure": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Procedure name.",
						},
						"privileges": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Permission information.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMariadbAccountPrivilegesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_account_privileges.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
		User       string
		Host       string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("accounts"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			if v, ok := dMap["user"]; ok {
				User = v.(string)
			}
			if v, ok := dMap["host"]; ok {
				Host = v.(string)
			}
		}
	}

	d.SetId(strings.Join([]string{instanceId, User, Host}, FILED_SP))

	return resourceTencentCloudMariadbAccountPrivilegesUpdate(d, meta)
}

func resourceTencentCloudMariadbAccountPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_account_privileges.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	User := idSplit[1]
	Host := idSplit[2]

	accountPrivileges, err := service.DescribeMariadbAccountPrivilegesById(ctx, instanceId, User, Host)
	if err != nil {
		return err
	}

	if accountPrivileges == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MariadbAccountPrivileges` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if accountPrivileges.InstanceId != nil {
		_ = d.Set("instance_id", accountPrivileges.InstanceId)
	}

	if accountPrivileges.UserName != nil && accountPrivileges.Host != nil {
		accountsList := []interface{}{}
		accountsMap := map[string]interface{}{}

		accountsMap["user"] = accountPrivileges.UserName
		accountsMap["host"] = accountPrivileges.Host

		accountsList = append(accountsList, accountsMap)

		_ = d.Set("accounts", accountsList)
	}

	if accountPrivileges.Privileges != nil {
		_ = d.Set("global_privileges", accountPrivileges.Privileges)
	}

	return nil
}

func resourceTencentCloudMariadbAccountPrivilegesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_account_privileges.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request = mariadb.NewModifyAccountPrivilegesRequest()
		flowId  int64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	User := idSplit[1]
	Host := idSplit[2]

	needChange := false

	mutableArgs := []string{"global_privileges", "database_privileges", "table_privileges", "column_privileges", "view_privileges", "function_privileges", "procedure_privileges"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request.InstanceId = &instanceId
		request.Accounts = []*mariadb.Account{
			{
				User: common.StringPtr(User),
				Host: common.StringPtr(Host),
			},
		}

		if v, ok := d.GetOk("global_privileges"); ok {
			globalPrivileges := v.(*schema.Set).List()
			for i := range globalPrivileges {
				backupCycle := globalPrivileges[i].(string)
				request.GlobalPrivileges = append(request.GlobalPrivileges, helper.String(backupCycle))
			}
		}

		if v, ok := d.GetOk("database_privileges"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				parameter := mariadb.DatabasePrivilege{}
				if v, ok := dMap["privileges"]; ok {
					parameter.Privileges = helper.Strings(v.([]string))
				}
				if v, ok := dMap["database"]; ok {
					parameter.Database = helper.String(v.(string))
				}
				request.DatabasePrivileges = append(request.DatabasePrivileges, &parameter)
			}
		}

		if v, ok := d.GetOk("table_privileges"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				parameter := mariadb.TablePrivilege{}
				if v, ok := dMap["database"]; ok {
					parameter.Database = helper.String(v.(string))
				}
				if v, ok := dMap["table"]; ok {
					parameter.Table = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					parameter.Privileges = helper.Strings(v.([]string))
				}
				request.TablePrivileges = append(request.TablePrivileges, &parameter)
			}
		}

		if v, ok := d.GetOk("column_privileges"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				parameter := mariadb.ColumnPrivilege{}
				if v, ok := dMap["database"]; ok {
					parameter.Database = helper.String(v.(string))
				}
				if v, ok := dMap["table"]; ok {
					parameter.Table = helper.String(v.(string))
				}
				if v, ok := dMap["column"]; ok {
					parameter.Column = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					parameter.Privileges = helper.Strings(v.([]string))
				}
				request.ColumnPrivileges = append(request.ColumnPrivileges, &parameter)
			}
		}

		if v, ok := d.GetOk("view_privileges"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				parameter := mariadb.ViewPrivileges{}
				if v, ok := dMap["database"]; ok {
					parameter.Database = helper.String(v.(string))
				}
				if v, ok := dMap["view"]; ok {
					parameter.View = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					parameter.Privileges = helper.Strings(v.([]string))
				}
				request.ViewPrivileges = append(request.ViewPrivileges, &parameter)
			}
		}

		if v, ok := d.GetOk("function_privileges"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				parameter := mariadb.FunctionPrivilege{}
				if v, ok := dMap["database"]; ok {
					parameter.Database = helper.String(v.(string))
				}
				if v, ok := dMap["function_name"]; ok {
					parameter.FunctionName = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					parameter.Privileges = helper.Strings(v.([]string))
				}
				request.FunctionPrivileges = append(request.FunctionPrivileges, &parameter)
			}
		}

		if v, ok := d.GetOk("procedure_privileges"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				parameter := mariadb.ProcedurePrivilege{}
				if v, ok := dMap["database"]; ok {
					parameter.Database = helper.String(v.(string))
				}
				if v, ok := dMap["procedure"]; ok {
					parameter.Procedure = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					parameter.Privileges = helper.Strings(v.([]string))
				}
				request.ProcedurePrivileges = append(request.ProcedurePrivileges, &parameter)
			}
		}

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyAccountPrivileges(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update mariadb accountPrivileges failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeFlowById(ctx, flowId)
		if e != nil {
			return retryError(e)
		}

		if *result.Status == MARIADB_TASK_SUCCESS {
			return nil
		} else if *result.Status == MARIADB_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("mariadb accountPrivileges status is running"))
		} else if *result.Status == MARIADB_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("mariadb accountPrivileges status is fail"))
		} else {
			e = fmt.Errorf("mariadb accountPrivileges status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s update mariadb accountPrivileges task failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMariadbAccountPrivilegesRead(d, meta)
}

func resourceTencentCloudMariadbAccountPrivilegesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_account_privileges.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
