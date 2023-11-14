/*
Provides a resource to create a cynosdb account_privileges

Example Usage

```hcl
resource "tencentcloud_cynosdb_account_privileges" "account_privileges" {
  cluster_id = "xxx"
  account {
		account_name = ""
		host = ""

  }
  global_privileges =
  database_privileges {
		db = ""
		privileges =

  }
  table_privileges {
		db = ""
		table_name = ""
		privileges =

  }
}
```

Import

cynosdb account_privileges can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_account_privileges.account_privileges account_privileges_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCynosdbAccountPrivileges() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbAccountPrivilegesCreate,
		Read:   resourceTencentCloudCynosdbAccountPrivilegesRead,
		Update: resourceTencentCloudCynosdbAccountPrivilegesUpdate,
		Delete: resourceTencentCloudCynosdbAccountPrivilegesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"account": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Account information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account.",
						},
						"host": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Host, default &amp;#39;%&amp;#39;.",
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
				Description: "Global permission array.",
			},

			"database_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Database permission array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database.",
						},
						"privileges": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Permission List.",
						},
					},
				},
			},

			"table_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Table permission array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Table Name.",
						},
						"privileges": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Permission List.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCynosdbAccountPrivilegesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_account_privileges.create")()
	defer inconsistentCheck(d, meta)()

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	var account string
	if v, ok := d.GetOk("account"); ok {
		account = v.(string)
	}

	d.SetId(strings.Join([]string{clusterId, account}, FILED_SP))

	return resourceTencentCloudCynosdbAccountPrivilegesUpdate(d, meta)
}

func resourceTencentCloudCynosdbAccountPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_account_privileges.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	account := idSplit[1]

	accountPrivileges, err := service.DescribeCynosdbAccountPrivilegesById(ctx, clusterId, account)
	if err != nil {
		return err
	}

	if accountPrivileges == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbAccountPrivileges` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if accountPrivileges.ClusterId != nil {
		_ = d.Set("cluster_id", accountPrivileges.ClusterId)
	}

	if accountPrivileges.Account != nil {
		accountMap := map[string]interface{}{}

		if accountPrivileges.Account.AccountName != nil {
			accountMap["account_name"] = accountPrivileges.Account.AccountName
		}

		if accountPrivileges.Account.Host != nil {
			accountMap["host"] = accountPrivileges.Account.Host
		}

		_ = d.Set("account", []interface{}{accountMap})
	}

	if accountPrivileges.GlobalPrivileges != nil {
		_ = d.Set("global_privileges", accountPrivileges.GlobalPrivileges)
	}

	if accountPrivileges.DatabasePrivileges != nil {
		databasePrivilegesList := []interface{}{}
		for _, databasePrivileges := range accountPrivileges.DatabasePrivileges {
			databasePrivilegesMap := map[string]interface{}{}

			if accountPrivileges.DatabasePrivileges.Db != nil {
				databasePrivilegesMap["db"] = accountPrivileges.DatabasePrivileges.Db
			}

			if accountPrivileges.DatabasePrivileges.Privileges != nil {
				databasePrivilegesMap["privileges"] = accountPrivileges.DatabasePrivileges.Privileges
			}

			databasePrivilegesList = append(databasePrivilegesList, databasePrivilegesMap)
		}

		_ = d.Set("database_privileges", databasePrivilegesList)

	}

	if accountPrivileges.TablePrivileges != nil {
		tablePrivilegesList := []interface{}{}
		for _, tablePrivileges := range accountPrivileges.TablePrivileges {
			tablePrivilegesMap := map[string]interface{}{}

			if accountPrivileges.TablePrivileges.Db != nil {
				tablePrivilegesMap["db"] = accountPrivileges.TablePrivileges.Db
			}

			if accountPrivileges.TablePrivileges.TableName != nil {
				tablePrivilegesMap["table_name"] = accountPrivileges.TablePrivileges.TableName
			}

			if accountPrivileges.TablePrivileges.Privileges != nil {
				tablePrivilegesMap["privileges"] = accountPrivileges.TablePrivileges.Privileges
			}

			tablePrivilegesList = append(tablePrivilegesList, tablePrivilegesMap)
		}

		_ = d.Set("table_privileges", tablePrivilegesList)

	}

	return nil
}

func resourceTencentCloudCynosdbAccountPrivilegesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_account_privileges.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cynosdb.NewModifyAccountPrivilegesRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	account := idSplit[1]

	request.ClusterId = &clusterId
	request.Account = &account

	immutableArgs := []string{"cluster_id", "account", "global_privileges", "database_privileges", "table_privileges"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_id") {
		if v, ok := d.GetOk("cluster_id"); ok {
			request.ClusterId = helper.String(v.(string))
		}
	}

	if d.HasChange("account") {
		if dMap, ok := helper.InterfacesHeadMap(d, "account"); ok {
			inputAccount := cynosdb.InputAccount{}
			if v, ok := dMap["account_name"]; ok {
				inputAccount.AccountName = helper.String(v.(string))
			}
			if v, ok := dMap["host"]; ok {
				inputAccount.Host = helper.String(v.(string))
			}
			request.Account = &inputAccount
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
				databasePrivileges := cynosdb.DatabasePrivileges{}
				if v, ok := dMap["db"]; ok {
					databasePrivileges.Db = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					privilegesSet := v.(*schema.Set).List()
					for i := range privilegesSet {
						privileges := privilegesSet[i].(string)
						databasePrivileges.Privileges = append(databasePrivileges.Privileges, &privileges)
					}
				}
				request.DatabasePrivileges = append(request.DatabasePrivileges, &databasePrivileges)
			}
		}
	}

	if d.HasChange("table_privileges") {
		if v, ok := d.GetOk("table_privileges"); ok {
			for _, item := range v.([]interface{}) {
				tablePrivileges := cynosdb.TablePrivileges{}
				if v, ok := dMap["db"]; ok {
					tablePrivileges.Db = helper.String(v.(string))
				}
				if v, ok := dMap["table_name"]; ok {
					tablePrivileges.TableName = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					privilegesSet := v.(*schema.Set).List()
					for i := range privilegesSet {
						privileges := privilegesSet[i].(string)
						tablePrivileges.Privileges = append(tablePrivileges.Privileges, &privileges)
					}
				}
				request.TablePrivileges = append(request.TablePrivileges, &tablePrivileges)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyAccountPrivileges(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb accountPrivileges failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbAccountPrivilegesRead(d, meta)
}

func resourceTencentCloudCynosdbAccountPrivilegesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_account_privileges.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
