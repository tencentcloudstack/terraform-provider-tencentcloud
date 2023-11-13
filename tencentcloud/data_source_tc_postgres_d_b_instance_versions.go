/*
Use this data source to query detailed information of postgres d_b_instance_versions

Example Usage

```hcl
data "tencentcloud_postgres_d_b_instance_versions" "d_b_instance_versions" {
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudPostgresDBInstanceVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresDBInstanceVersionsRead,
		Schema: map[string]*schema.Schema{
			"version_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of database versions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"d_b_engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database engines. Valid values:1. `postgresql` (TencentDB for PostgreSQL)2. `mssql_compatible` (MSSQL compatible-TencentDB for PostgreSQL).",
						},
						"d_b_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database version, such as 12.4.",
						},
						"d_b_major_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database major version, such as 12.",
						},
						"d_b_kernel_version": {
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

func dataSourceTencentCloudPostgresDBInstanceVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_postgres_d_b_instance_versions.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	var versionSet []*postgres.Version

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresDBInstanceVersionsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
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
				versionMap["d_b_engine"] = version.DBEngine
			}

			if version.DBVersion != nil {
				versionMap["d_b_version"] = version.DBVersion
			}

			if version.DBMajorVersion != nil {
				versionMap["d_b_major_version"] = version.DBMajorVersion
			}

			if version.DBKernelVersion != nil {
				versionMap["d_b_kernel_version"] = version.DBKernelVersion
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
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
