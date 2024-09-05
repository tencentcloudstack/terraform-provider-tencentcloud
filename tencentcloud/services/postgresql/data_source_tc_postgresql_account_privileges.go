package postgresql

import (
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
)

func DataSourceTencentCloudPostgresqlAccountPrivileges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlAccountPrivilegesRead,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"user_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance username.",
			},
			"database_object_set": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Instance database object info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Supported database object types: account, database, schema, sequence, procedure, type, function, table, view, matview, column. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"object_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database object Name.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"database_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database name to which the database object belongs. This parameter is mandatory when ObjectType is not database.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"schema_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Schema name to which the database object belongs. This parameter is mandatory when ObjectType is not database or schema.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Table name to which the database object belongs. This parameter is mandatory when ObjectType is column.Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},
			// computed
			"privilege_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Privilege list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Database object.If ObjectType is database, DatabaseName/SchemaName/TableName can be null.If ObjectType is schema, SchemaName/TableName can be null.If ObjectType is table, TableName can be null.If ObjectType is column, DatabaseName/SchemaName/TableName can&amp;#39;t be null.In all other cases, DatabaseName/SchemaName/TableName can be null. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Supported database object types: account, database, schema, sequence, procedure, type, function, table, view, matview, column. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"object_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database object Name. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"database_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database name to which the database object belongs. This parameter is mandatory when ObjectType is not database. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"schema_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Schema name to which the database object belongs. This parameter is mandatory when ObjectType is not database or schema. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"table_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Table name to which the database object belongs. This parameter is mandatory when ObjectType is column. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"privilege_set": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Privileges the specific account has on database object. Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudPostgresqlAccountPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_account_privileges.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		request      = postgresql.NewDescribeAccountPrivilegesRequest()
		privilegeSet []*postgresql.DatabasePrivilege
		dBInstanceId string
		userName     string
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = &dBInstanceId
		dBInstanceId = v.(string)
	}

	if v, ok := d.GetOk("user_name"); ok {
		request.UserName = &userName
		userName = v.(string)
	}

	if v, ok := d.GetOk("database_object_set"); ok {
		for _, item := range v.([]interface{}) {
			objectMap := item.(map[string]interface{})
			databaseObject := postgresql.DatabaseObject{}
			if v, ok := objectMap["object_type"]; ok {
				databaseObject.ObjectType = helper.String(v.(string))
			}

			if v, ok := objectMap["object_name"]; ok {
				databaseObject.ObjectName = helper.String(v.(string))
			}

			if v, ok := objectMap["database_name"]; ok {
				databaseObject.DatabaseName = helper.String(v.(string))
			}

			if v, ok := objectMap["schema_name"]; ok {
				databaseObject.SchemaName = helper.String(v.(string))
			}

			if v, ok := objectMap["table_name"]; ok {
				databaseObject.TableName = helper.String(v.(string))
			}

			request.DatabaseObjectSet = append(request.DatabaseObjectSet, &databaseObject)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DescribeAccountPrivileges(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		privilegeSet = result.Response.PrivilegeSet
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(privilegeSet))
	for _, item := range privilegeSet {
		databasePrivilegeMap := map[string]interface{}{}
		if item.Object != nil {
			objectMap := map[string]interface{}{}
			if item.Object.ObjectType != nil {
				objectMap["object_type"] = item.Object.ObjectType
			}

			if item.Object.ObjectName != nil {
				objectMap["object_name"] = item.Object.ObjectName
			}

			if item.Object.DatabaseName != nil {
				objectMap["database_name"] = item.Object.DatabaseName
			}

			if item.Object.SchemaName != nil {
				objectMap["schema_name"] = item.Object.SchemaName
			}

			if item.Object.TableName != nil {
				objectMap["table_name"] = item.Object.TableName
			}

			databasePrivilegeMap["object"] = []interface{}{objectMap}
		}

		if item.PrivilegeSet != nil {
			privilegeList := make([]string, 0, len(item.PrivilegeSet))
			for _, privilege := range item.PrivilegeSet {
				privilegeList = append(privilegeList, *privilege)
			}

			databasePrivilegeMap["privilege_set"] = privilegeList
		}

		tmpList = append(tmpList, databasePrivilegeMap)
	}

	_ = d.Set("privilege_set", tmpList)

	d.SetId(strings.Join([]string{dBInstanceId, userName}, tccommon.FILED_SP))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
