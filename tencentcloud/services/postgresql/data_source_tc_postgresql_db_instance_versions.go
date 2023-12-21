package postgresql

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlDbInstanceVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlDbInstanceVersionsRead,
		Schema: map[string]*schema.Schema{
			"version_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of database versions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database engines. Valid values:1. `postgresql` (TencentDB for PostgreSQL)2. `mssql_compatible` (MSSQL compatible-TencentDB for PostgreSQL).",
						},
						"db_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database version, such as 12.4.",
						},
						"db_major_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database major version, such as 12.",
						},
						"db_kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database kernel version, such as v12.4_r1.3.",
						},
						"supported_feature_names": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of features supported by the database kernel, such as:TDE: Supports data encryption.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database version status. Valid values:`AVAILABLE`.`DEPRECATED`.",
						},
						"available_upgrade_target": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of versions to which this database version (`DBKernelVersion`) can be upgraded.",
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

func dataSourceTencentCloudPostgresqlDbInstanceVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_db_instance_versions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var versionSet []*postgresql.Version

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlDbInstanceVersionsByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		versionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(versionSet))
	tmpList := make([]map[string]interface{}, 0, len(versionSet))

	if versionSet != nil {
		for _, version := range versionSet {
			versionMap := map[string]interface{}{}

			if version.DBEngine != nil {
				versionMap["db_engine"] = version.DBEngine
			}

			if version.DBVersion != nil {
				versionMap["db_version"] = version.DBVersion
			}

			if version.DBMajorVersion != nil {
				versionMap["db_major_version"] = version.DBMajorVersion
			}

			if version.DBKernelVersion != nil {
				versionMap["db_kernel_version"] = version.DBKernelVersion
			}

			if version.SupportedFeatureNames != nil {
				versionMap["supported_feature_names"] = version.SupportedFeatureNames
			}

			if version.Status != nil {
				versionMap["status"] = version.Status
			}

			if version.AvailableUpgradeTarget != nil {
				versionMap["available_upgrade_target"] = version.AvailableUpgradeTarget
			}

			ids = append(ids, *version.DBEngine)
			tmpList = append(tmpList, versionMap)
		}

		_ = d.Set("version_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
