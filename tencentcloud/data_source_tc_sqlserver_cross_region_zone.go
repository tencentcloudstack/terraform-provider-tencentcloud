/*
Use this data source to query detailed information of sqlserver datasource_cross_region_zone

Example Usage

```hcl
data "tencentcloud_sqlserver_cross_region_zone" "cross_region_zone" {
  instance_id = "mssql-qelbzgwf"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSqlserverCrossRegionZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverCrossRegionZoneRead,
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

func dataSourceTencentCloudSqlserverCrossRegionZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_cross_region_zone.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		service     = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		crossRegion *sqlserver.DescribeCrossRegionZoneResponseParams
		instanceId  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverCrossRegionZoneByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		crossRegion = result
		return nil
	})

	if err != nil {
		return err
	}

	if crossRegion.Region != nil {
		_ = d.Set("region", crossRegion.Region)
	}

	if crossRegion.Zone != nil {
		_ = d.Set("zone", crossRegion.Zone)
	}

	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
