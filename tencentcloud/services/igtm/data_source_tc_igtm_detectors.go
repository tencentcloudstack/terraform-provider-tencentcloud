package igtm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIgtmDetectors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIgtmDetectorsRead,
		Schema: map[string]*schema.Schema{
			"detector_group_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detector group list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Line group ID GroupLineId.",
						},
						"group_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "bgp, international, isp.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group name.",
						},
						"internet_family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ipv4, ipv6.",
						},
						"package_set": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Supported package types.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func dataSourceTencentCloudIgtmDetectorsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_igtm_detectors.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	var respData []*igtmv20231024.DetectorGroup
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIgtmDetectorsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	detectorGroupSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, detectorGroupSet := range respData {
			detectorGroupSetMap := map[string]interface{}{}
			if detectorGroupSet.Gid != nil {
				detectorGroupSetMap["gid"] = detectorGroupSet.Gid
			}

			if detectorGroupSet.GroupType != nil {
				detectorGroupSetMap["group_type"] = detectorGroupSet.GroupType
			}

			if detectorGroupSet.GroupName != nil {
				detectorGroupSetMap["group_name"] = detectorGroupSet.GroupName
			}

			if detectorGroupSet.InternetFamily != nil {
				detectorGroupSetMap["internet_family"] = detectorGroupSet.InternetFamily
			}

			if detectorGroupSet.PackageSet != nil {
				detectorGroupSetMap["package_set"] = detectorGroupSet.PackageSet
			}

			detectorGroupSetList = append(detectorGroupSetList, detectorGroupSetMap)
		}

		_ = d.Set("detector_group_set", detectorGroupSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), detectorGroupSetList); e != nil {
			return e
		}
	}

	return nil
}
