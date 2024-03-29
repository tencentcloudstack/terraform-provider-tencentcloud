package common

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAvailabilityRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAvailabilityRegionsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When specified, only the region with the exactly name match will be returned. `default` value means it consistent with the provider region.",
			},
			"include_unavailable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A bool variable indicates that the query will include `UNAVAILABLE` regions.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values.
			"regions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of regions will be exported and its every element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the region, like `ap-guangzhou`.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the region, like `Guangzhou Region`.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the region, indicate availability using `AVAILABLE` and `UNAVAILABLE` values.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAvailabilityRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_availability_regions.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var name string
	var includeUnavailable = false
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if name == "default" {
		name = meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
	}
	if v, ok := d.GetOkExists("include_unavailable"); ok {
		includeUnavailable = v.(bool)
	}

	var regions []*cvm.RegionInfo
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		regions, errRet = cvmService.DescribeRegions(ctx)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	regionList := make([]map[string]interface{}, 0, len(regions))
	ids := make([]string, 0, len(regions))
	for _, region := range regions {
		if name != "" && name != *region.Region {
			continue
		}
		if !includeUnavailable && *region.RegionState == svccvm.ZONE_STATE_UNAVAILABLE {
			continue
		}
		mapping := map[string]interface{}{
			"name":        region.Region,
			"description": region.RegionName,
			"state":       region.RegionState,
		}
		regionList = append(regionList, mapping)
		ids = append(ids, *region.Region)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("regions", regionList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set regions list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), regionList); err != nil {
			return err
		}
	}

	return nil
}
