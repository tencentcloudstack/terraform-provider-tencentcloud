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
terraform import tencentcloud_dcdb_account_privileges.account_privileges instanceId#userName#host#dbName#tabName#viewName#colName
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Description: "Global permissions. Among them, the optional value of the permission in GlobalPrivileges is: SELECT, INSERT, UPDATE, DELETE, CREATE, PROCESS, DROP, REFERENCES, INDEX, ALTER, SHOW DATABASES,  CREATE TEMPORARY TABLES, LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER.  Note that if this parameter is not passed, it means that the existing permissions are reserved. If it needs to be cleared, pass an empty array in this field.",
			},

			"database_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Database permissions. Optional values for the Privileges permission are: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE TEMPORARY TABLES,  LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER.  Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.",
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
				MaxItems:    1,
				Description: "Permissions for tables in the database. Optional values for the Privileges permission are: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE VIEW, SHOW VIEW, TRIGGER. Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.",
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
				MaxItems:    1,
				Description: "Permissions for columns in database tables. Optional values for the Privileges permission are:  SELECT, INSERT, UPDATE, REFERENCES.  Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.",
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
				MaxItems:    1,
				Description: "Permissions for database views. Optional values for the Privileges permission are:  SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE VIEW, SHOW VIEW, TRIGGER.  Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.",
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

	return resourceTencentCloudDcdbAccountPrivilegesUpdate(d, meta)
}

func resourceTencentCloudDcdbAccountPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account_privileges.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 7 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	// userName := idSplit[1]
	// host := idSplit[2]
	dbName := helper.String(idSplit[3])
	tabName := helper.String(idSplit[4])
	viewName := helper.String(idSplit[5])
	colName := helper.String(idSplit[6])

	// query global_privileges
	globalPrivileges, err := service.DescribeDcdbAccountPrivilegesById(ctx, d.Id(), helper.String("*"), nil, nil, nil)
	if err != nil {
		return err
	}
	if globalPrivileges == nil {
		d.SetId("")
		return fmt.Errorf("resource `globalPrivileges` %s does not exist", d.Id())
	}

	if globalPrivileges.Privileges != nil {
		log.Printf("[DEBUG]%s read globalPrivileges. Privileges:[%v]\n", logId, globalPrivileges.Privileges)
		_ = d.Set("global_privileges", helper.StringsInterfaces(globalPrivileges.Privileges))
	}

	// set common info
	if globalPrivileges.InstanceId != nil {
		_ = d.Set("instance_id", globalPrivileges.InstanceId)
	}

	accountList := []interface{}{}
	accountMap := make(map[string]interface{})
	if globalPrivileges.UserName != nil {
		accountMap["user"] = globalPrivileges.UserName
	}
	if globalPrivileges.Host != nil {
		accountMap["host"] = globalPrivileges.Host
	}
	accountList = append(accountList, accountMap)
	_ = d.Set("account", accountList)

	// query db list
	dbs, err := service.DescribeDcdbDatabases(ctx, instanceId)
	if err != nil {
		return err
	}
	if dbs == nil {
		return fmt.Errorf("%s the dbs get from DescribeDcdbDatabases is nil!", logId)
	}

	// query database_privileges
	dbPrivilegesList := make([]interface{}, 0, len(dbs.Databases))
	for _, db := range dbs.Databases {
		if *db.DbName == *dbName { // only support single name
			dbPrivileges, err := service.DescribeDcdbAccountPrivilegesById(ctx, d.Id(), dbName, helper.String("*"), nil, nil)
			if err != nil {
				return err
			}
			if dbPrivileges == nil {
				d.SetId("")
				return fmt.Errorf("resource `dbPrivileges` %s does not exist", d.Id())
			}

			dbPrivilegesMap := map[string]interface{}{}
			if dbPrivileges.Privileges != nil {
				dbPrivilegesMap["privileges"] = dbPrivileges.Privileges
			}
			dbPrivilegesMap["database"] = dbName
			dbPrivilegesList = append(dbPrivilegesList, dbPrivilegesMap)

			// query db objects
			objs, err := service.DescribeDcdbDBObjects(ctx, instanceId, *dbName)
			if err != nil {
				return err
			}
			if objs == nil {
				return fmt.Errorf("%s the objs get from DescribeDcdbDBObjects is nil!", logId)
			}

			// query table_privileges
			tabPrivilegesList := make([]interface{}, 0, len(objs.Tables))
			for _, tab := range objs.Tables {
				if *tab.Table == *tabName { // only support single name
					tabPrivileges, err := service.DescribeDcdbAccountPrivilegesById(ctx, d.Id(), dbName, helper.String("table"), tabName, helper.String("*"))
					if err != nil {
						return err
					}
					if tabPrivileges == nil {
						d.SetId("")
						return fmt.Errorf("resource `tabPrivileges` %s does not exist", d.Id())
					}

					tabPrivilegesMap := map[string]interface{}{}
					if tabPrivileges.Privileges != nil {
						tabPrivilegesMap["privileges"] = tabPrivileges.Privileges
					}
					tabPrivilegesMap["database"] = dbName
					tabPrivilegesMap["table"] = tabName
					tabPrivilegesList = append(tabPrivilegesList, tabPrivilegesMap)

					// query db tables
					tabs, err := service.DescribeDcdbDBTables(ctx, instanceId, *dbName, *tabName)
					if err != nil {
						return err
					}
					if tabs == nil {
						return fmt.Errorf("%s the tabs get from DescribeDcdbDBTables is nil!", logId)
					}

					// query column_privileges
					colPrivilegesList := make([]interface{}, 0, len(objs.Tables))
					for _, col := range tabs.Cols {
						if *col.Col == *colName { // only support single name
							colPrivileges, err := service.DescribeDcdbAccountPrivilegesById(ctx, d.Id(), dbName, helper.String("table"), tabName, colName)
							if err != nil {
								return err
							}
							if colPrivileges == nil {
								d.SetId("")
								return fmt.Errorf("resource `colPrivileges` %s does not exist", d.Id())
							}

							colPrivilegesMap := map[string]interface{}{}
							if colPrivileges.Privileges != nil {
								colPrivilegesMap["privileges"] = colPrivileges.Privileges
							}
							colPrivilegesMap["database"] = dbName
							colPrivilegesMap["table"] = tabName
							colPrivilegesMap["column"] = colName
							colPrivilegesList = append(colPrivilegesList, colPrivilegesMap)
						}
					} // cols
					_ = d.Set("column_privileges", colPrivilegesList)
				}
			} // tabs
			_ = d.Set("table_privileges", tabPrivilegesList)

			// query view_privileges
			viewPrivilegesList := make([]interface{}, 0, len(objs.Views))
			for _, view := range objs.Views {
				if *view.View == *viewName { // only support single name
					viewPrivileges, err := service.DescribeDcdbAccountPrivilegesById(ctx, d.Id(), dbName, helper.String("view"), viewName, nil)
					if err != nil {
						return err
					}
					if viewPrivileges == nil {
						d.SetId("")
						return fmt.Errorf("resource `viewPrivileges` %s does not exist", d.Id())
					}

					viewPrivilegesMap := map[string]interface{}{}
					if viewPrivileges.Privileges != nil {
						viewPrivilegesMap["privileges"] = viewPrivileges.Privileges
					}
					viewPrivilegesMap["database"] = dbName
					viewPrivilegesMap["view"] = viewName
					viewPrivilegesList = append(viewPrivilegesList, viewPrivilegesMap)
				}
			} // views
			_ = d.Set("view_privileges", viewPrivilegesList)
		}
	} // dbs
	_ = d.Set("database_privileges", dbPrivilegesList)

	return nil
}

func resourceTencentCloudDcdbAccountPrivilegesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account_privileges.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dcdb.NewModifyAccountPrivilegesRequest()

	var (
		instanceId string
		userName   string
		host       string
		dbName     string
		tabName    string
		viewName   string
		colName    string
		flowId     *int64
		service    = DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}
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
					dbName = v.(string)
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
					dbName = v.(string)
				}
				if v, ok := dMap["table"]; ok {
					tablePrivilege.Table = helper.String(v.(string))
					tabName = v.(string)
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
					dbName = v.(string)
				}
				if v, ok := dMap["table"]; ok {
					columnPrivilege.Table = helper.String(v.(string))
					tabName = v.(string)
				}
				if v, ok := dMap["column"]; ok {
					columnPrivilege.Column = helper.String(v.(string))
					colName = v.(string)
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
					dbName = v.(string)
				}
				if v, ok := dMap["view"]; ok {
					viewPrivileges.View = helper.String(v.(string))
					viewName = v.(string)
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

	request.InstanceId = &instanceId
	request.Accounts = []*dcdb.Account{
		{
			User: helper.String(userName),
			Host: helper.String(host),
		},
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyAccountPrivileges(request)
		flowId = response.Response.FlowId
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dcdb accountPrivileges failed, reason:%+v", logId, err)
		return err
	}

	if flowId != nil {
		// need to wait modify operation success
		// 0:success; 1:failed, 2:running
		conf := BuildStateChangeConf([]string{}, []string{"0"}, 3*readRetryTimeout, time.Second, service.DcdbDbInstanceStateRefreshFunc(flowId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	d.SetId(strings.Join([]string{instanceId, userName, host, dbName, tabName, viewName, colName}, FILED_SP))

	return resourceTencentCloudDcdbAccountPrivilegesRead(d, meta)
}

func resourceTencentCloudDcdbAccountPrivilegesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account_privileges.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
