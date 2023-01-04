/*
Provides a resource to create a dcdb account_privileges

Example Usage

```hcl
resource "tencentcloud_dcdb_account_privileges" "account_privileges" {
  instance_id = "%s"
  account {
		user = "tf_test"
		host = "%s"
  }
  global_privileges = ["SHOW DATABASES","SHOW VIEW"]
  database_privileges {
		privileges = ["SELECT","INSERT","UPDATE","DELETE","CREATE"]
		database = "tf_test_db"
  }

  table_privileges {
		database = "tf_test_db"
		table = "tf_test_table"
		privileges = ["SELECT","INSERT","UPDATE","DELETE","CREATE"]

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
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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

			"account": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The account of the database, including username and host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "account name.",
						},
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "account host.",
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
				Description: "&amp;quot;Global permissions. Among them, the optional value of the permission in GlobalPrivileges is:&amp;quot;&amp;quot;SELECT, INSERT, UPDATE, DELETE, CREATE, PROCESS, DROP, REFERENCES, INDEX, ALTER, SHOW DATABASES,&amp;quot;&amp;quot;CREATE TEMPORARY TABLES, LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER.&amp;quot;&amp;quot;Note that if this parameter is not passed, it means that the existing permissions are reserved. If it needs to be cleared, pass an empty array in this field.&amp;quot;.",
			},

			"database_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "&amp;quot;Database permissions. Optional values for the Privileges permission are:&amp;quot;&amp;quot;SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE TEMPORARY TABLES,&amp;quot;&amp;quot;LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER.&amp;quot;&amp;quot;Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.&amp;quot;.",
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
				Description: "&amp;quot;Permissions for tables in the database. Optional values for the Privileges permission are:&amp;quot;&amp;quot;SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE VIEW, SHOW VIEW, TRIGGER.&amp;quot;&amp;quot;Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.&amp;quot;.",
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
				Description: "&amp;quot;Permissions for columns in database tables. Optional values for the Privileges permission are:&amp;quot;&amp;quot;SELECT, INSERT, UPDATE, REFERENCES.&amp;quot;&amp;quot;Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.&amp;quot;.",
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
				Description: "&amp;quot;Permissions for database views. Optional values for the Privileges permission are:&amp;quot;&amp;quot;SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE VIEW, SHOW VIEW, TRIGGER.&amp;quot;&amp;quot;Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.&amp;quot;.",
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
	var (
		instanceId string
		userName   string
		host       string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "account"); ok {
		if v, ok := dMap["user"]; ok {
			userName = v.(string)
		}
		if v, ok := dMap["host"]; ok {
			host = v.(string)
		}
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

	// query global privileges
	dbName := helper.String("*")
	globalPrivileges, err := service.DescribeDcdbAccountPrivilegesById(ctx, d.Id(), dbName, nil, nil, nil)
	if err != nil {
		return err
	}
	if globalPrivileges == nil {
		d.SetId("")
		return fmt.Errorf("resource `DcdbAccountPrivileges` %s does not exist", d.Id())
	}

	if globalPrivileges.Privileges != nil {
		_ = d.Set("global_privileges", helper.StringsInterfaces(globalPrivileges.Privileges))
	}

	// set account and ins id
	if globalPrivileges.InstanceId != nil {
		_ = d.Set("instance_id", globalPrivileges.InstanceId)
	}

	accountMap := make(map[string]interface{})
	if globalPrivileges.UserName != nil {
		accountMap["user"] = globalPrivileges.UserName
	}
	if globalPrivileges.Host != nil {
		accountMap["host"] = globalPrivileges.Host
	}
	_ = d.Set("account", accountMap)

	// query database_privileges
	// dbPrivileges, err = service.DescribeDcdbAccountPrivilegesById(ctx, d.Id(), nil, queryType, nil, nil)
	// if err != nil {
	// 	return err
	// }
	// if globalPrivileges == nil {
	// 	d.SetId("")
	// 	return fmt.Errorf("resource `DcdbAccountPrivileges` %s does not exist", d.Id())
	// }

	// if accountPrivileges.DatabasePrivileges != nil {
	// 	databasePrivilegesList := []interface{}{}
	// 	for _, databasePriviprivilegesleges := range accountPrivileges.DatabasePrivileges {
	// 		databasePrivilegesMap := map[string]interface{}{}

	// 		if accountPrivileges.DatabasePrivileges.Privileges != nil {
	// 			databasePrivilegesMap["privileges"] = accountPrivileges.DatabasePrivileges.Privileges
	// 		}

	// 		if accountPrivileges.DatabasePrivileges.Database != nil {
	// 			databasePrivilegesMap["database"] = accountPrivileges.DatabasePrivileges.Database
	// 		}

	// 		databasePrivilegesList = append(databasePrivilegesList, databasePrivilegesMap)
	// 	}

	// 	_ = d.Set("database_privileges", databasePrivilegesList)

	// }

	// // query table privileges
	// if accountPrivileges.TablePrivileges != nil {
	// 	tablePrivilegesList := []interface{}{}
	// 	for _, tablePrivileges := range accountPrivileges.TablePrivileges {
	// 		tablePrivilegesMap := map[string]interface{}{}

	// 		if accountPrivileges.TablePrivileges.Database != nil {
	// 			tablePrivilegesMap["database"] = accountPrivileges.TablePrivileges.Database
	// 		}

	// 		if accountPrivileges.TablePrivileges.Table != nil {
	// 			tablePrivilegesMap["table"] = accountPrivileges.TablePrivileges.Table
	// 		}

	// 		if accountPrivileges.TablePrivileges.Privileges != nil {
	// 			tablePrivilegesMap["privileges"] = accountPrivileges.TablePrivileges.Privileges
	// 		}

	// 		tablePrivilegesList = append(tablePrivilegesList, tablePrivilegesMap)
	// 	}

	// 	_ = d.Set("table_privileges", tablePrivilegesList)

	// }

	// // query column privileges
	// if accountPrivileges.ColumnPrivileges != nil {
	// 	columnPrivilegesList := []interface{}{}
	// 	for _, columnPrivileges := range accountPrivileges.ColumnPrivileges {
	// 		columnPrivilegesMap := map[string]interface{}{}

	// 		if accountPrivileges.ColumnPrivileges.Database != nil {
	// 			columnPrivilegesMap["database"] = accountPrivileges.ColumnPrivileges.Database
	// 		}

	// 		if accountPrivileges.ColumnPrivileges.Table != nil {
	// 			columnPrivilegesMap["table"] = accountPrivileges.ColumnPrivileges.Table
	// 		}

	// 		if accountPrivileges.ColumnPrivileges.Column != nil {
	// 			columnPrivilegesMap["column"] = accountPrivileges.ColumnPrivileges.Column
	// 		}

	// 		if accountPrivileges.ColumnPrivileges.Privileges != nil {
	// 			columnPrivilegesMap["privileges"] = accountPrivileges.ColumnPrivileges.Privileges
	// 		}

	// 		columnPrivilegesList = append(columnPrivilegesList, columnPrivilegesMap)
	// 	}

	// 	_ = d.Set("column_privileges", columnPrivilegesList)

	// }

	// // query view privileges
	// if accountPrivileges.ViewPrivileges != nil {
	// 	viewPrivilegesList := []interface{}{}
	// 	for _, viewPrivileges := range accountPrivileges.ViewPrivileges {
	// 		viewPrivilegesMap := map[string]interface{}{}

	// 		if accountPrivileges.ViewPrivileges.Database != nil {
	// 			viewPrivilegesMap["database"] = accountPrivileges.ViewPrivileges.Database
	// 		}

	// 		if accountPrivileges.ViewPrivileges.View != nil {
	// 			viewPrivilegesMap["view"] = accountPrivileges.ViewPrivileges.View
	// 		}

	// 		if accountPrivileges.ViewPrivileges.Privileges != nil {
	// 			viewPrivilegesMap["privileges"] = accountPrivileges.ViewPrivileges.Privileges
	// 		}

	// 		viewPrivilegesList = append(viewPrivilegesList, viewPrivilegesMap)
	// 	}

	// 	_ = d.Set("view_privileges", viewPrivilegesList)

	// }

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
	request.Accounts = []*dcdb.Account{
		{
			User: helper.String(userName),
			Host: helper.String(host),
		},
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
				dMap := item.(map[string]interface{})
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
				dMap := item.(map[string]interface{})
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
				dMap := item.(map[string]interface{})
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
				dMap := item.(map[string]interface{})
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
