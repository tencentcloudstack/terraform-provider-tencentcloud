package oceanus

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOceanusSavepointList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusSavepointListRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job SerialId.",
			},
			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},
			//"record_types": {
			//	Optional:    true,
			//	Type:        schema.TypeList,
			//	Elem:        &schema.Schema{Type: schema.TypeInt},
			//	Description: "RecordTypes. 1 is triggering the savepoint, 2 is the checkpoint, and 3 is stopping the triggered savepoint",
			//},
			"savepoint": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Snapshot listNote: This field may return null, indicating that no valid value was found.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Primary keyNote: This field may return null, indicating that no valid value was found.",
						},
						"version_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Version numberNote: This field may return null, indicating that no valid value was found.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status: 1=Active; 2=Expired; 3=InProgress; 4=Failed; 5=TimeoutNote: This field may return null, indicating that no valid value was found.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creation timeNote: This field may return null, indicating that no valid value was found.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Update timeNote: This field may return null, indicating that no valid value was found.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PathNote: This field may return null, indicating that no valid value was found.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "SizeNote: This field may return null, indicating that no valid value was found.",
						},
						"record_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Snapshot type: 1=savepoint; 2=checkpoint; 3=cancelWithSavepointNote: This field may return null, indicating that no valid value was found.",
						},
						"job_runtime_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sequential ID of the running job instanceNote: This field may return null, indicating that no valid value was found.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DescriptionNote: This field may return null, indicating that no valid value was found.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Fixed timeoutNote: This field may return null, indicating that no valid value was found.",
						},
						"serial_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Snapshot SerialIdNote: This field may return null, indicating that no valid value was found.",
						},
						"time_consuming": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time consumptionNote: This field may return null, indicating that no valid value was found.",
						},
						"path_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Snapshot path status: 1=available; 2=unavailable;Note: This field may return null, indicating that no valid value was found.",
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

func dataSourceTencentCloudOceanusSavepointListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_oceanus_savepoint_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		savepointList []*oceanus.Savepoint
		jobId         string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["JobId"] = helper.String(v.(string))
		jobId = v.(string)
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record_types"); ok {
		recordTypesList := v.([]interface{})
		for _, item := range recordTypesList {
			recordTypesList = append(recordTypesList, item.(int))
		}

		paramMap["RecordTypes"] = recordTypesList
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusSavepointListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		savepointList = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(savepointList))

	if savepointList != nil {
		for _, savepoint := range savepointList {
			savepointMap := map[string]interface{}{}

			if savepoint.Id != nil {
				savepointMap["id"] = savepoint.Id
			}

			if savepoint.VersionId != nil {
				savepointMap["version_id"] = savepoint.VersionId
			}

			if savepoint.Status != nil {
				savepointMap["status"] = savepoint.Status
			}

			if savepoint.CreateTime != nil {
				savepointMap["create_time"] = savepoint.CreateTime
			}

			if savepoint.UpdateTime != nil {
				savepointMap["update_time"] = savepoint.UpdateTime
			}

			if savepoint.Path != nil {
				savepointMap["path"] = savepoint.Path
			}

			if savepoint.Size != nil {
				savepointMap["size"] = savepoint.Size
			}

			if savepoint.RecordType != nil {
				savepointMap["record_type"] = savepoint.RecordType
			}

			if savepoint.JobRuntimeId != nil {
				savepointMap["job_runtime_id"] = savepoint.JobRuntimeId
			}

			if savepoint.Description != nil {
				savepointMap["description"] = savepoint.Description
			}

			if savepoint.Timeout != nil {
				savepointMap["timeout"] = savepoint.Timeout
			}

			if savepoint.SerialId != nil {
				savepointMap["serial_id"] = savepoint.SerialId
			}

			if savepoint.TimeConsuming != nil {
				savepointMap["time_consuming"] = savepoint.TimeConsuming
			}

			if savepoint.PathStatus != nil {
				savepointMap["path_status"] = savepoint.PathStatus
			}

			tmpList = append(tmpList, savepointMap)
		}

		_ = d.Set("savepoint", tmpList)
	}

	d.SetId(jobId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
