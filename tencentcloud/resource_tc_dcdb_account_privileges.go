/*
Provides a resource to create a dcdb account_privileges

Example Usage

```hcl
resource "tencentcloud_dcdb_account_privileges" "account_privileges" {
  instance_id = "tdsql-c1nl9rpv"
  accounts {
		user = &lt;nil&gt;
		host = &lt;nil&gt;

  }
  global_privileges = &lt;nil&gt;
  database_privileges {
		privileges = &lt;nil&gt;
		database = &lt;nil&gt;

  }
  table_privileges {
		database = &lt;nil&gt;
		table = &lt;nil&gt;
		privileges = &lt;nil&gt;

  }
  column_privileges {
		database = &lt;nil&gt;
		table = &lt;nil&gt;
		column = &lt;nil&gt;
		privileges = &lt;nil&gt;

  }
  view_privileges {
		database = &lt;nil&gt;
		view = &lt;nil&gt;
		privileges = &lt;nil&gt;

  }
}
```

Import

dcdb account_privileges can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_account_privileges.account_privileges account_privileges_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudDcdbAccountPrivileges() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbAccountPrivilegesCreate,
		Read:   resourceTencentCloudDcdbAccountPrivilegesRead,
		Update: resourceTencentCloudDcdbAccountPrivilegesUpdate,
		Delete: resourceTencentCloudDcdbAccountPrivilegesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"accounts": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "The account of the database, including username and host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account name.",
						},
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account host.",
						},
					},
				},
			},

			"global_privileges": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Global permissions. Among them, the optional value of the permission in GlobalPrivileges is:SELECT, INSERT, UPDATE, DELETE, CREATE, PROCESS, DROP, REFERENCES, INDEX, ALTER, SHOW DATABASES,CREATE TEMPORARY TABLES, LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER.Note that if this parameter is not passed, it means that the existing permissions are reserved. If it needs to be cleared, pass an empty array in this field.",
			},

			"database_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Database permissions. Optional values for the Privileges permission are:SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE TEMPORARY TABLES,LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER.Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"privileges": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Permission information.",
						},
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of database.",
						},
					},
				},
			},

			"table_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Permissions for tables in the database. Optional values for the Privileges permission are:SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE VIEW, SHOW VIEW, TRIGGER.Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of database.",
						},
						"table": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database table name.",
						},
						"privileges": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Permission information.",
						},
					},
				},
			},

			"column_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Permissions for columns in database tables. Optional values for the Privileges permission are:SELECT, INSERT, UPDATE, REFERENCES.Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of database.",
						},
						"table": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database table name.",
						},
						"column": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database column name.",
						},
						"privileges": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Permission information.",
						},
					},
				},
			},

			"view_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Permissions for database views. Optional values for the Privileges permission are:SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE VIEW, SHOW VIEW, TRIGGER.Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of database.",
						},
						"view": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database view name.",
						},
						"privileges": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Permission information.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDcdbAccountPrivilegesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account_privileges.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	var userName string
	if v, ok := d.GetOk("user_name"); ok {
		userName = v.(string)
	}

	var host string
	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, userName, host}, FILED_SP))

	return resourceTencentCloudDcdbAccountPrivilegesUpdate(d, meta)
}

func resourceTencentCloudDcdbAccountPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account_privileges.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]
	host := idSplit[2]

	accountPrivileges, err := service.DescribeDcdbAccountPrivilegesById(ctx, instanceId, userName, host)
	if err != nil {
		return err
	}

	if accountPrivileges == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcdbAccountPrivileges` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if accountPrivileges.InstanceId != nil {
		_ = d.Set("instance_id", accountPrivileges.InstanceId)
	}

	if accountPrivileges.Accounts != nil {
		accountsList := []interface{}{}
		for _, accounts := range accountPrivileges.Accounts {
			accountsMap := map[string]interface{}{}

			if accountPrivileges.Accounts.User != nil {
				accountsMap["user"] = accountPrivileges.Accounts.User
			}

			if accountPrivileges.Accounts.Host != nil {
				accountsMap["host"] = accountPrivileges.Accounts.Host
			}

			accountsList = append(accountsList, accountsMap)
		}

		_ = d.Set("accounts", accountsList)

	}

	if accountPrivileges.GlobalPrivileges != nil {
		_ = d.Set("global_privileges", accountPrivileges.GlobalPrivileges)
	}

	if accountPrivileges.DatabasePrivileges != nil {
		databasePrivilegesList := []interface{}{}
		for _, databasePrivileges := range accountPrivileges.DatabasePrivileges {
			databasePrivilegesMap := map[string]interface{}{}

			if accountPrivileges.DatabasePrivileges.Privileges != nil {
				databasePrivilegesMap["privileges"] = accountPrivileges.DatabasePrivileges.Privileges
			}

			if accountPrivileges.DatabasePrivileges.Database != nil {
				databasePrivilegesMap["database"] = accountPrivileges.DatabasePrivileges.Database
			}

			databasePrivilegesList = append(databasePrivilegesList, databasePrivilegesMap)
		}

		_ = d.Set("database_privileges", databasePrivilegesList)

	}

	if accountPrivileges.TablePrivileges != nil {
		tablePrivilegesList := []interface{}{}
		for _, tablePrivileges := range accountPrivileges.TablePrivileges {
			tablePrivilegesMap := map[string]interface{}{}

			if accountPrivileges.TablePrivileges.Database != nil {
				tablePrivilegesMap["database"] = accountPrivileges.TablePrivileges.Database
			}

			if accountPrivileges.TablePrivileges.Table != nil {
				tablePrivilegesMap["table"] = accountPrivileges.TablePrivileges.Table
			}

			if accountPrivileges.TablePrivileges.Privileges != nil {
				tablePrivilegesMap["privileges"] = accountPrivileges.TablePrivileges.Privileges
			}

			tablePrivilegesList = append(tablePrivilegesList, tablePrivilegesMap)
		}

		_ = d.Set("table_privileges", tablePrivilegesList)

	}

	if accountPrivileges.ColumnPrivileges != nil {
		columnPrivilegesList := []interface{}{}
		for _, columnPrivileges := range accountPrivileges.ColumnPrivileges {
			columnPrivilegesMap := map[string]interface{}{}

			if accountPrivileges.ColumnPrivileges.Database != nil {
				columnPrivilegesMap["database"] = accountPrivileges.ColumnPrivileges.Database
			}

			if accountPrivileges.ColumnPrivileges.Table != nil {
				columnPrivilegesMap["table"] = accountPrivileges.ColumnPrivileges.Table
			}

			if accountPrivileges.ColumnPrivileges.Column != nil {
				columnPrivilegesMap["column"] = accountPrivileges.ColumnPrivileges.Column
			}

			if accountPrivileges.ColumnPrivileges.Privileges != nil {
				columnPrivilegesMap["privileges"] = accountPrivileges.ColumnPrivileges.Privileges
			}

			columnPrivilegesList = append(columnPrivilegesList, columnPrivilegesMap)
		}

		_ = d.Set("column_privileges", columnPrivilegesList)

	}

	if accountPrivileges.ViewPrivileges != nil {
		viewPrivilegesList := []interface{}{}
		for _, viewPrivileges := range accountPrivileges.ViewPrivileges {
			viewPrivilegesMap := map[string]interface{}{}

			if accountPrivileges.ViewPrivileges.Database != nil {
				viewPrivilegesMap["database"] = accountPrivileges.ViewPrivileges.Database
			}

			if accountPrivileges.ViewPrivileges.View != nil {
				viewPrivilegesMap["view"] = accountPrivileges.ViewPrivileges.View
			}

			if accountPrivileges.ViewPrivileges.Privileges != nil {
				viewPrivilegesMap["privileges"] = accountPrivileges.ViewPrivileges.Privileges
			}

			viewPrivilegesList = append(viewPrivilegesList, viewPrivilegesMap)
		}

		_ = d.Set("view_privileges", viewPrivilegesList)

	}

	return nil
}

func resourceTencentCloudDcdbAccountPrivilegesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account_privileges.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dcdb.NewModifyAccountPrivilegesRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]
	host := idSplit[2]

	request.InstanceId = &instanceId
	request.UserName = &userName
	request.Host = &host

	immutableArgs := []string{"instance_id", "accounts", "global_privileges", "database_privileges", "table_privileges", "column_privileges", "view_privileges"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_id") {
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}
	}

	if d.HasChange("accounts") {
		if v, ok := d.GetOk("accounts"); ok {
			for _, item := range v.([]interface{}) {
				account := dcdb.Account{}
				if v, ok := dMap["user"]; ok {
					account.User = helper.String(v.(string))
				}
				if v, ok := dMap["host"]; ok {
					account.Host = helper.String(v.(string))
				}
				request.Accounts = append(request.Accounts, &account)
			}
		}
	}

	if d.HasChange("global_privileges") {
		if v, ok := d.GetOk("global_privileges"); ok {
			globalPrivilegesSet := v.(*schema.Set).List()
			for i := range globalPrivilegesSet {
				globalPrivileges := globalPrivilegesSet[i].(string)
				request.GlobalPrivileges = append(request.GlobalPrivileges, &globalPrivileges)
			}
		}
	}

	if d.HasChange("database_privileges") {
		if v, ok := d.GetOk("database_privileges"); ok {
			for _, item := range v.([]interface{}) {
				databasePrivilege := dcdb.DatabasePrivilege{}
				if v, ok := dMap["privileges"]; ok {
					privilegesSet := v.(*schema.Set).List()
					for i := range privilegesSet {
						privileges := privilegesSet[i].(string)
						databasePrivilege.Privileges = append(databasePrivilege.Privileges, &privileges)
					}
				}
				if v, ok := dMap["database"]; ok {
					databasePrivilege.Database = helper.String(v.(string))
				}
				request.DatabasePrivileges = append(request.DatabasePrivileges, &databasePrivilege)
			}
		}
	}

	if d.HasChange("table_privileges") {
		if v, ok := d.GetOk("table_privileges"); ok {
			for _, item := range v.([]interface{}) {
				tablePrivilege := dcdb.TablePrivilege{}
				if v, ok := dMap["database"]; ok {
					tablePrivilege.Database = helper.String(v.(string))
				}
				if v, ok := dMap["table"]; ok {
					tablePrivilege.Table = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					privilegesSet := v.(*schema.Set).List()
					for i := range privilegesSet {
						privileges := privilegesSet[i].(string)
						tablePrivilege.Privileges = append(tablePrivilege.Privileges, &privileges)
					}
				}
				request.TablePrivileges = append(request.TablePrivileges, &tablePrivilege)
			}
		}
	}

	if d.HasChange("column_privileges") {
		if v, ok := d.GetOk("column_privileges"); ok {
			for _, item := range v.([]interface{}) {
				columnPrivilege := dcdb.ColumnPrivilege{}
				if v, ok := dMap["database"]; ok {
					columnPrivilege.Database = helper.String(v.(string))
				}
				if v, ok := dMap["table"]; ok {
					columnPrivilege.Table = helper.String(v.(string))
				}
				if v, ok := dMap["column"]; ok {
					columnPrivilege.Column = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					privilegesSet := v.(*schema.Set).List()
					for i := range privilegesSet {
						privileges := privilegesSet[i].(string)
						columnPrivilege.Privileges = append(columnPrivilege.Privileges, &privileges)
					}
				}
				request.ColumnPrivileges = append(request.ColumnPrivileges, &columnPrivilege)
			}
		}
	}

	if d.HasChange("view_privileges") {
		if v, ok := d.GetOk("view_privileges"); ok {
			for _, item := range v.([]interface{}) {
				viewPrivileges := dcdb.ViewPrivileges{}
				if v, ok := dMap["database"]; ok {
					viewPrivileges.Database = helper.String(v.(string))
				}
				if v, ok := dMap["view"]; ok {
					viewPrivileges.View = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					privilegesSet := v.(*schema.Set).List()
					for i := range privilegesSet {
						privileges := privilegesSet[i].(string)
						viewPrivileges.Privileges = append(viewPrivileges.Privileges, &privileges)
					}
				}
				request.ViewPrivileges = append(request.ViewPrivileges, &viewPrivileges)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyAccountPrivileges(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dcdb accountPrivileges failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbAccountPrivilegesRead(d, meta)
}

func resourceTencentCloudDcdbAccountPrivilegesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account_privileges.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
