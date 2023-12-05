package tencentcloud

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOceanusRunJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusRunJobCreate,
		Read:   resourceTencentCloudOceanusRunJobRead,
		Delete: resourceTencentCloudOceanusRunJobDelete,

		Schema: map[string]*schema.Schema{
			"run_job_descriptions": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "The description information for batch job startup.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Job ID.",
						},
						"run_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The type of the run. 1 indicates start, and 2 indicates resume.",
						},
						"start_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Compatible with the startup parameters of the old SQL type job: specify the start time point of data source consumption (recommended to pass the value)Ensure that the parameter is LATEST, EARLIEST, T+Timestamp (example: T1557394288000).",
						},
						"job_config_version": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "A certain version of the current job(Not passed by default as a non-draft job version).",
						},
						"savepoint_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Savepoint path.",
						},
						"savepoint_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Savepoint ID.",
						},
						"use_old_system_connector": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Use the historical version of the system dependency.",
						},
						"custom_timestamp": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Custom timestamp.",
						},
					},
				},
			},
			"work_space_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},
		},
	}
}

func resourceTencentCloudOceanusRunJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_run_job.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = oceanus.NewRunJobsRequest()
		jobIds  []string
	)

	if v, ok := d.GetOk("run_job_descriptions"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			runJobDescription := oceanus.RunJobDescription{}
			if v, ok := dMap["job_id"]; ok {
				runJobDescription.JobId = helper.String(v.(string))
				jobIds = append(jobIds, v.(string))
			}

			if v, ok := dMap["run_type"]; ok {
				runJobDescription.RunType = helper.IntInt64(v.(int))
			}

			if v, ok := dMap["start_mode"]; ok {
				runJobDescription.StartMode = helper.String(v.(string))
			}

			if v, ok := dMap["job_config_version"]; ok {
				runJobDescription.JobConfigVersion = helper.IntUint64(v.(int))
			}

			if v, ok := dMap["savepoint_path"]; ok {
				runJobDescription.SavepointPath = helper.String(v.(string))
			}

			if v, ok := dMap["savepoint_id"]; ok {
				runJobDescription.SavepointId = helper.String(v.(string))
			}

			if v, ok := dMap["use_old_system_connector"]; ok {
				runJobDescription.UseOldSystemConnector = helper.Bool(v.(bool))
			}

			if v, ok := dMap["custom_timestamp"]; ok {
				runJobDescription.CustomTimestamp = helper.IntInt64(v.(int))
			}

			request.RunJobDescriptions = append(request.RunJobDescriptions, &runJobDescription)
		}
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		request.WorkSpaceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOceanusClient().RunJobs(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate oceanus runJob failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join(jobIds, FILED_SP))

	return resourceTencentCloudOceanusRunJobRead(d, meta)
}

func resourceTencentCloudOceanusRunJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_run_job.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudOceanusRunJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_run_job.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
