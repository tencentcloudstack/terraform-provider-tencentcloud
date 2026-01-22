package postgresql

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlDbVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlDbVersionsRead,
		Schema: map[string]*schema.Schema{
			"db_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Version of the postgresql database engine.",
			},
			"db_major_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "PostgreSQL major version number.",
			},
			"db_kernel_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "PostgreSQL kernel version number.",
			},
			// computed
			"version_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of database versions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database engines. Valid values:\n1. `postgresql` (TencentDB for PostgreSQL)\n2. `mssql_compatible` (MSSQL compatible-TencentDB for PostgreSQL).",
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
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of features supported by the database kernel, such as:\nTDE: Supports data encryption.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database version status. Valid values:\n`AVAILABLE`.\n`DEPRECATED`.",
						},
						"available_upgrade_target": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of versions to which this database version (`DBKernelVersion`) can be upgraded, including minor and major version numbers available for upgrade (complete kernel version format example: v15.1_v1.6).",
							Elem:        &schema.Schema{Type: schema.TypeString},
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

func dataSourceTencentCloudPostgresqlDbVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_db_versions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("db_version"); ok {
		paramMap["DBVersion"] = v.(string)
	}

	if v, ok := d.GetOk("db_major_version"); ok {
		paramMap["DBMajorVersion"] = v.(string)
	}

	if v, ok := d.GetOk("db_kernel_version"); ok {
		paramMap["DBKernelVersion"] = v.(string)
	}

	var respData []*postgresql.Version
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlDbVersionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	ids := make([]string, 0, len(respData))
	versionSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, versionSet := range respData {
			versionSetMap := map[string]interface{}{}
			if versionSet.DBEngine != nil {
				versionSetMap["db_engine"] = versionSet.DBEngine
			}

			if versionSet.DBVersion != nil {
				versionSetMap["db_version"] = versionSet.DBVersion
			}

			if versionSet.DBMajorVersion != nil {
				versionSetMap["db_major_version"] = versionSet.DBMajorVersion
			}

			if versionSet.DBKernelVersion != nil {
				versionSetMap["db_kernel_version"] = versionSet.DBKernelVersion
			}

			if versionSet.SupportedFeatureNames != nil {
				versionSetMap["supported_feature_names"] = versionSet.SupportedFeatureNames
			}

			if versionSet.Status != nil {
				versionSetMap["status"] = versionSet.Status
			}

			if versionSet.AvailableUpgradeTarget != nil {
				versionSetMap["available_upgrade_target"] = versionSet.AvailableUpgradeTarget
			}

			ids = append(ids, *versionSet.DBEngine)
			versionSetList = append(versionSetList, versionSetMap)
		}

		_ = d.Set("version_set", versionSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), versionSetList); e != nil {
			return e
		}
	}

	return nil
}
