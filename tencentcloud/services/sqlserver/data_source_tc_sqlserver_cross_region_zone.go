package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverCrossRegionZone() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_cross_region_zone.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		crossRegion *sqlserver.DescribeCrossRegionZoneResponseParams
		instanceId  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverCrossRegionZoneByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
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
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
