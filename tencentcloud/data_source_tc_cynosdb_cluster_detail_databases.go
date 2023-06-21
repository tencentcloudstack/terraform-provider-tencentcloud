/*
Use this data source to query detailed information of cynosdb cluster_detail_databases

Example Usage

```hcl
data "tencentcloud_cynosdb_cluster_detail_databases" "cluster_detail_databases" {
  cluster_id = "cynosdbmysql-bws8h88b"
  db_name    = "users"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbClusterDetailDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClusterDetailDatabasesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"db_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database Name.",
			},
			"db_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Database information note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database Name.",
						},
						"character_set": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Character Set Type.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database Status.",
						},
						"collate_rule": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Capture Rules.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"user_host_privileges": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User permission note: This field may return null, indicating that a valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_user_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DbUserName.",
									},
									"db_host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database host.",
									},
									"db_privilege": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User permission note: This field may return null, indicating that a valid value cannot be obtained.",
									},
								},
							},
						},
						"db_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Database ID note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User appid note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User Uin note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster Id note: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudCynosdbClusterDetailDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_cluster_detail_databases.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		dbInfos   []*cynosdb.DbInfo
		clusterId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("db_name"); ok {
		paramMap["DbName"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbClusterDetailDatabasesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		dbInfos = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(dbInfos))

	if dbInfos != nil {
		for _, dbInfo := range dbInfos {
			dbInfoMap := map[string]interface{}{}

			if dbInfo.DbName != nil {
				dbInfoMap["db_name"] = dbInfo.DbName
			}

			if dbInfo.CharacterSet != nil {
				dbInfoMap["character_set"] = dbInfo.CharacterSet
			}

			if dbInfo.Status != nil {
				dbInfoMap["status"] = dbInfo.Status
			}

			if dbInfo.CollateRule != nil {
				dbInfoMap["collate_rule"] = dbInfo.CollateRule
			}

			if dbInfo.Description != nil {
				dbInfoMap["description"] = dbInfo.Description
			}

			if dbInfo.UserHostPrivileges != nil {
				userHostPrivilegesList := []interface{}{}
				for _, userHostPrivileges := range dbInfo.UserHostPrivileges {
					userHostPrivilegesMap := map[string]interface{}{}

					if userHostPrivileges.DbUserName != nil {
						userHostPrivilegesMap["db_user_name"] = userHostPrivileges.DbUserName
					}

					if userHostPrivileges.DbHost != nil {
						userHostPrivilegesMap["db_host"] = userHostPrivileges.DbHost
					}

					if userHostPrivileges.DbPrivilege != nil {
						userHostPrivilegesMap["db_privilege"] = userHostPrivileges.DbPrivilege
					}

					userHostPrivilegesList = append(userHostPrivilegesList, userHostPrivilegesMap)
				}

				dbInfoMap["user_host_privileges"] = userHostPrivilegesList
			}

			if dbInfo.DbId != nil {
				dbInfoMap["db_id"] = dbInfo.DbId
			}

			if dbInfo.CreateTime != nil {
				dbInfoMap["create_time"] = dbInfo.CreateTime
			}

			if dbInfo.UpdateTime != nil {
				dbInfoMap["update_time"] = dbInfo.UpdateTime
			}

			if dbInfo.AppId != nil {
				dbInfoMap["app_id"] = dbInfo.AppId
			}

			if dbInfo.Uin != nil {
				dbInfoMap["uin"] = dbInfo.Uin
			}

			if dbInfo.ClusterId != nil {
				dbInfoMap["cluster_id"] = dbInfo.ClusterId
			}

			tmpList = append(tmpList, dbInfoMap)
		}

		_ = d.Set("db_infos", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
