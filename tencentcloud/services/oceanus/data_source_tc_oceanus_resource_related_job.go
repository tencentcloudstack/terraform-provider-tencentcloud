package oceanus

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOceanusResourceRelatedJob() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusResourceRelatedJobRead,
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Resource ID.",
			},
			"desc_by_job_config_create_time": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      DESC_BY_JOB_CONFIG_CREATETIME_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(DESC_BY_JOB_CONFIG_CREATETIME),
				Description:  "Default:0; 1:sort by job version creation time in descending order.",
			},
			"resource_config_version": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Resource version number.",
			},
			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},
			"ref_job_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Associated job information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Job ID.",
						},
						"job_config_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Job configuration version.",
						},
						"resource_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource version.",
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

func dataSourceTencentCloudOceanusResourceRelatedJobRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_oceanus_resource_related_job.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		refJobInfos []*oceanus.ResourceRefJobInfo
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("resource_id"); ok {
		paramMap["ResourceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("desc_by_job_config_create_time"); ok {
		paramMap["DESCByJobConfigCreateTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("resource_config_version"); ok {
		paramMap["ResourceConfigVersion"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusResourceRelatedJobByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		refJobInfos = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(refJobInfos))
	tmpList := make([]map[string]interface{}, 0, len(refJobInfos))

	if refJobInfos != nil {
		for _, resourceRefJobInfo := range refJobInfos {
			resourceRefJobInfoMap := map[string]interface{}{}

			if resourceRefJobInfo.JobId != nil {
				resourceRefJobInfoMap["job_id"] = resourceRefJobInfo.JobId
			}

			if resourceRefJobInfo.JobConfigVersion != nil {
				resourceRefJobInfoMap["job_config_version"] = resourceRefJobInfo.JobConfigVersion
			}

			if resourceRefJobInfo.ResourceVersion != nil {
				resourceRefJobInfoMap["resource_version"] = resourceRefJobInfo.ResourceVersion
			}

			ids = append(ids, *resourceRefJobInfo.JobId)
			tmpList = append(tmpList, resourceRefJobInfoMap)
		}

		_ = d.Set("ref_job_infos", tmpList)
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
