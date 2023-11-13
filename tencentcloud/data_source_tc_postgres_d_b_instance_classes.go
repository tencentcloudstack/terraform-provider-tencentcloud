/*
Use this data source to query detailed information of postgres d_b_instance_classes

Example Usage

```hcl
data "tencentcloud_postgres_d_b_instance_classes" "d_b_instance_classes" {
  zone = ""
  d_b_engine = ""
  d_b_major_version = ""
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

func dataSourceTencentCloudPostgresDBInstanceClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresDBInstanceClassesRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "AZ ID, which can be obtained through the `DescribeZones` API.",
			},

			"d_b_engine": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database engines. Valid values:1. `postgresql` (TencentDB for PostgreSQL)2. `mssql_compatible` (MSSQL compatible-TencentDB for PostgreSQL).",
			},

			"d_b_major_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Major version of a database, such as 12 or 13, which can be obtained through the `DescribeDBVersions` API.",
			},

			"class_info_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of database specifications.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specification ID.",
						},
						"c_p_u": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of CPU cores.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size in MB.",
						},
						"max_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum storage capacity in GB supported by this specification.",
						},
						"min_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum storage capacity in GB supported by this specification.",
						},
						"q_p_s": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Estimated QPS for this specification.",
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

func dataSourceTencentCloudPostgresDBInstanceClassesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_postgres_d_b_instance_classes.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone"); ok {
		paramMap["Zone"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_b_engine"); ok {
		paramMap["DBEngine"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_b_major_version"); ok {
		paramMap["DBMajorVersion"] = helper.String(v.(string))
	}

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	var classInfoSet []*postgres.ClassInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresDBInstanceClassesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		classInfoSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(classInfoSet))
	tmpList := make([]map[string]interface{}, 0, len(classInfoSet))

	if classInfoSet != nil {
		for _, classInfo := range classInfoSet {
			classInfoMap := map[string]interface{}{}

			if classInfo.SpecCode != nil {
				classInfoMap["spec_code"] = classInfo.SpecCode
			}

			if classInfo.CPU != nil {
				classInfoMap["c_p_u"] = classInfo.CPU
			}

			if classInfo.Memory != nil {
				classInfoMap["memory"] = classInfo.Memory
			}

			if classInfo.MaxStorage != nil {
				classInfoMap["max_storage"] = classInfo.MaxStorage
			}

			if classInfo.MinStorage != nil {
				classInfoMap["min_storage"] = classInfo.MinStorage
			}

			if classInfo.QPS != nil {
				classInfoMap["q_p_s"] = classInfo.QPS
			}

			ids = append(ids, *classInfo.SpecCode)
			tmpList = append(tmpList, classInfoMap)
		}

		_ = d.Set("class_info_set", tmpList)
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
