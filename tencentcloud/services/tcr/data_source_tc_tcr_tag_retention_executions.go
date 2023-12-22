package tcr

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTcrTagRetentionExecutions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrTagRetentionExecutionsRead,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"retention_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "retention id.",
			},

			"retention_execution_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "list of version retention execution records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"execution_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "execution id.",
						},
						"retention_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "retention id.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "execution start time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "execution end time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "execution status: Failed, Succeed, Stopped, InProgress.",
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

func dataSourceTencentCloudTcrTagRetentionExecutionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcr_tag_retention_executions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("registry_id"); ok {
		paramMap["RegistryId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("retention_id"); v != nil {
		paramMap["RetentionId"] = helper.IntInt64(v.(int))
	}

	service := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var retentionExecutionList []*tcr.RetentionExecution

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcrTagRetentionExecutionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		retentionExecutionList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(retentionExecutionList))
	tmpList := make([]map[string]interface{}, 0, len(retentionExecutionList))

	if retentionExecutionList != nil {
		for _, retentionExecution := range retentionExecutionList {
			retentionExecutionMap := map[string]interface{}{}

			if retentionExecution.ExecutionId != nil {
				retentionExecutionMap["execution_id"] = retentionExecution.ExecutionId
			}

			if retentionExecution.RetentionId != nil {
				retentionExecutionMap["retention_id"] = retentionExecution.RetentionId
			}

			if retentionExecution.StartTime != nil {
				retentionExecutionMap["start_time"] = retentionExecution.StartTime
			}

			if retentionExecution.EndTime != nil {
				retentionExecutionMap["end_time"] = retentionExecution.EndTime
			}

			if retentionExecution.Status != nil {
				retentionExecutionMap["status"] = retentionExecution.Status
			}

			ids = append(ids, helper.Int64ToStr(*retentionExecution.ExecutionId))
			tmpList = append(tmpList, retentionExecutionMap)
		}

		_ = d.Set("retention_execution_list", tmpList)
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
