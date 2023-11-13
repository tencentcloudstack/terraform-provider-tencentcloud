/*
Use this data source to query detailed information of sqlserver datasource_cross_region_zone

Example Usage

```hcl
data "tencentcloud_sqlserver_datasource_cross_region_zone" "datasource_cross_region_zone" {
  instance_id = "mssql-j8kv137v"
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

func dataSourceTencentCloudSqlserverDatasourceCrossRegionZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverDatasourceCrossRegionZoneRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of mssql-j8kv137v.",
			},

			"region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The string ID of the region where the standby machine is located, such as: ap-guangzhou.",
			},

			"zone": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The string ID of the availability zone where the standby machine is located, such as: ap-guangzhou-1.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSqlserverDatasourceCrossRegionZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_datasource_cross_region_zone.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverDatasourceCrossRegionZoneByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		region = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(region))
	if region != nil {
		_ = d.Set("region", region)
	}

	if zone != nil {
		_ = d.Set("zone", zone)
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
