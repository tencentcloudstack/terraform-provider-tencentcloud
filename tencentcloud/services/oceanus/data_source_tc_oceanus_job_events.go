package oceanus

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOceanusJobEvents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusJobEventsRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job ID.",
			},
			"start_timestamp": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Filter condition:Start Unix timestamp (seconds).",
			},
			"end_timestamp": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Filter condition:End Unix timestamp (seconds).",
			},
			"types": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Event types. If not passed, data of all types will be returned.",
			},
			"work_space_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},
			// Computed
			"running_order_ids": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Array of running instance IDs.",
			},
			"events": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of events within the specified range for this jobNote: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internally defined event type.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description text of the event type.",
						},
						"timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unix timestamp (seconds) when the event occurred.",
						},
						"running_order_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Running ID when the event occurredNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Some optional explanations of the eventNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"solution_link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Troubleshooting manual link for the abnormal eventNote: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudOceanusJobEventsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_oceanus_job_events.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		JobEvents *oceanus.DescribeJobEventsResponseParams
		jobId     string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["JobId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("start_timestamp"); ok {
		paramMap["StartTimestamp"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("end_timestamp"); ok {
		paramMap["EndTimestamp"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("types"); ok {
		typesSet := v.(*schema.Set).List()
		paramMap["Types"] = helper.InterfacesStringsPoint(typesSet)
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusJobEventsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			e = fmt.Errorf("oceanus Job events not exists")
			return resource.NonRetryableError(e)
		}

		JobEvents = result
		return nil
	})

	if err != nil {
		return err
	}

	if JobEvents.RunningOrderIds != nil {
		_ = d.Set("running_order_ids", JobEvents.RunningOrderIds)
	}

	if JobEvents.Events != nil {
		tmpList := make([]map[string]interface{}, 0, len(JobEvents.Events))

		for _, jobEvent := range JobEvents.Events {
			jobEventMap := map[string]interface{}{}

			if jobEvent.Type != nil {
				jobEventMap["type"] = jobEvent.Type
			}

			if jobEvent.Description != nil {
				jobEventMap["description"] = jobEvent.Description
			}

			if jobEvent.Timestamp != nil {
				jobEventMap["timestamp"] = jobEvent.Timestamp
			}

			if jobEvent.RunningOrderId != nil {
				jobEventMap["running_order_id"] = jobEvent.RunningOrderId
			}

			if jobEvent.Message != nil {
				jobEventMap["message"] = jobEvent.Message
			}

			if jobEvent.SolutionLink != nil {
				jobEventMap["solution_link"] = jobEvent.SolutionLink
			}

			tmpList = append(tmpList, jobEventMap)
		}

		_ = d.Set("events", tmpList)
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
