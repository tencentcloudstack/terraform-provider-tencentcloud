/*
Use this data source to query detailed information of cdb supported_privileges

Example Usage

```hcl
data "tencentcloud_cdb_supported_privileges" "supported_privileges" {
  instance_id = ""
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

func dataSourceTencentCloudCdbSupportedPrivileges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbSupportedPrivilegesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.",
			},

			"global_supported_privileges": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Global permissions supported by the instance.",
			},

			"database_supported_privileges": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Database permissions supported by the instance.",
			},

			"table_supported_privileges": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Database table permissions supported by the instance.",
			},

			"column_supported_privileges": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The database column permissions supported by the instance.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCdbSupportedPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_supported_privileges.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var globalSupportedPrivileges []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbSupportedPrivilegesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		globalSupportedPrivileges = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(globalSupportedPrivileges))
	if globalSupportedPrivileges != nil {
		_ = d.Set("global_supported_privileges", globalSupportedPrivileges)
	}

	if databaseSupportedPrivileges != nil {
		_ = d.Set("database_supported_privileges", databaseSupportedPrivileges)
	}

	if tableSupportedPrivileges != nil {
		_ = d.Set("table_supported_privileges", tableSupportedPrivileges)
	}

	if columnSupportedPrivileges != nil {
		_ = d.Set("column_supported_privileges", columnSupportedPrivileges)
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
