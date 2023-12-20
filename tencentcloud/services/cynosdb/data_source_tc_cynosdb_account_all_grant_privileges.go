package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbAccountAllGrantPrivileges() *schema.Resource {
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
				Description: "account information.",
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
							Description: "Host, default `%`.",
						},
					},
				},
			},
			"privilege_statements": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Permission statement note: This field may return null, indicating that a valid value cannot be obtained.",
			},
			"global_privileges": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
							Description: "database.",
						},
						"privileges": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
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
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_account_all_grant_privileges.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                     = tccommon.GetLogId(tccommon.ContextNil)
		ctx                       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service                   = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		accountAllGrantPrivileges *cynosdb.DescribeAccountAllGrantPrivilegesResponseParams
		clusterId                 string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "account"); ok {
		inputAccount := cynosdb.InputAccount{}
		if v, ok := dMap["account_name"]; ok {
			inputAccount.AccountName = helper.String(v.(string))
		}
		if v, ok := dMap["host"]; ok {
			inputAccount.Host = helper.String(v.(string))
		}
		paramMap["Account"] = &inputAccount
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbAccountAllGrantPrivilegesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		accountAllGrantPrivileges = result
		return nil
	})

	if err != nil {
		return err
	}

	if accountAllGrantPrivileges.PrivilegeStatements != nil {
		_ = d.Set("privilege_statements", accountAllGrantPrivileges.PrivilegeStatements)
	}

	if accountAllGrantPrivileges.GlobalPrivileges != nil {
		_ = d.Set("global_privileges", accountAllGrantPrivileges.GlobalPrivileges)
	}

	if accountAllGrantPrivileges.DatabasePrivileges != nil {
		tmpList := []interface{}{}
		for _, databasePrivileges := range accountAllGrantPrivileges.DatabasePrivileges {
			databasePrivilegesMap := map[string]interface{}{}

			if databasePrivileges.Db != nil {
				databasePrivilegesMap["db"] = databasePrivileges.Db
			}

			if databasePrivileges.Privileges != nil {
				databasePrivilegesMap["privileges"] = databasePrivileges.Privileges
			}

			tmpList = append(tmpList, databasePrivilegesMap)
		}

		_ = d.Set("database_privileges", tmpList)
	}

	if accountAllGrantPrivileges.TablePrivileges != nil {
		tmpList := []interface{}{}
		for _, tablePrivileges := range accountAllGrantPrivileges.TablePrivileges {
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

			tmpList = append(tmpList, tablePrivilegesMap)
		}

		_ = d.Set("table_privileges", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
