/*
Use this data source to query detailed information of cynosdb account_all_grant_privileges

Example Usage

```hcl
data "tencentcloud_cynosdb_account_all_grant_privileges" "account_all_grant_privileges" {
  cluster_id = "xxx"
  account {
		account_name = ""
		host = ""

  }
        }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbAccountAllGrantPrivileges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbAccountAllGrantPrivilegesRead,
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

			"privilege_statements": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Permission statement note: This field may return null, indicating that a valid value cannot be obtained.",
			},

			"global_privileges": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Global permission note: This field may return null, indicating that a valid value cannot be obtained.",
			},

			"database_privileges": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Database permissions note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database.",
						},
						"privileges": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Permission List.",
						},
					},
				},
			},

			"table_privileges": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Database table permissions note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Table Name.",
						},
						"privileges": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Permission List.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCynosdbAccountAllGrantPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_account_all_grant_privileges.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "account"); ok {
		inputAccount := cynosdb.InputAccount{}
		if v, ok := dMap["account_name"]; ok {
			inputAccount.AccountName = helper.String(v.(string))
		}
		if v, ok := dMap["host"]; ok {
			inputAccount.Host = helper.String(v.(string))
		}
		paramMap["account"] = &inputAccount
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var privilegeStatements []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbAccountAllGrantPrivilegesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		privilegeStatements = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(privilegeStatements))
	if privilegeStatements != nil {
		_ = d.Set("privilege_statements", privilegeStatements)
	}

	if globalPrivileges != nil {
		_ = d.Set("global_privileges", globalPrivileges)
	}

	if databasePrivileges != nil {
		for _, databasePrivileges := range databasePrivileges {
			databasePrivilegesMap := map[string]interface{}{}

			if databasePrivileges.Db != nil {
				databasePrivilegesMap["db"] = databasePrivileges.Db
			}

			if databasePrivileges.Privileges != nil {
				databasePrivilegesMap["privileges"] = databasePrivileges.Privileges
			}

			ids = append(ids, *databasePrivileges.TableName)
			tmpList = append(tmpList, databasePrivilegesMap)
		}

		_ = d.Set("database_privileges", tmpList)
	}

	if tablePrivileges != nil {
		for _, tablePrivileges := range tablePrivileges {
			tablePrivilegesMap := map[string]interface{}{}

			if tablePrivileges.Db != nil {
				tablePrivilegesMap["db"] = tablePrivileges.Db
			}

			if tablePrivileges.TableName != nil {
				tablePrivilegesMap["table_name"] = tablePrivileges.TableName
			}

			if tablePrivileges.Privileges != nil {
				tablePrivilegesMap["privileges"] = tablePrivileges.Privileges
			}

			ids = append(ids, *tablePrivileges.TableName)
			tmpList = append(tmpList, tablePrivilegesMap)
		}

		_ = d.Set("table_privileges", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
