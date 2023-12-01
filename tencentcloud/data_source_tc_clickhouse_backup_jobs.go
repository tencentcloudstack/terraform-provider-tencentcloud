/*
Use this data source to query detailed information of clickhouse backup jobs

Example Usage

```hcl
data "tencentcloud_clickhouse_backup_jobs" "backup_jobs" {
  instance_id = "cdwch-xxxxxx"
}
```
*/
package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clickhouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClickhouseBackupJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClickhouseBackupJobsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"begin_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Begin time.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},

			"back_up_jobs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Back up jobs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Back up job id.",
						},
						"snapshot": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Back up job name.",
						},
						"back_up_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Back up type.",
						},
						"back_up_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Back up size.",
						},
						"back_up_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Back up create time.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Back up expire time.",
						},
						"job_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Job status.",
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

func dataSourceTencentCloudClickhouseBackupJobsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clickhouse_backup_jobs.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("begin_time"); ok {
		paramMap["begin_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["end_time"] = helper.String(v.(string))
	}

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}

	var backUpJobs []*clickhouse.BackUpJobDisplay

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClickhouseBackupJobsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		backUpJobs = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(backUpJobs))
	tmpList := make([]map[string]interface{}, 0, len(backUpJobs))

	if backUpJobs != nil {
		for _, backUpJobDisplay := range backUpJobs {
			backUpJobDisplayMap := map[string]interface{}{}

			if backUpJobDisplay.JobId != nil {
				backUpJobDisplayMap["job_id"] = backUpJobDisplay.JobId
			}

			if backUpJobDisplay.Snapshot != nil {
				backUpJobDisplayMap["snapshot"] = backUpJobDisplay.Snapshot
			}

			if backUpJobDisplay.BackUpType != nil {
				backUpJobDisplayMap["back_up_type"] = backUpJobDisplay.BackUpType
			}

			if backUpJobDisplay.BackUpSize != nil {
				backUpJobDisplayMap["back_up_size"] = backUpJobDisplay.BackUpSize
			}

			if backUpJobDisplay.BackUpTime != nil {
				backUpJobDisplayMap["back_up_time"] = backUpJobDisplay.BackUpTime
			}

			if backUpJobDisplay.ExpireTime != nil {
				backUpJobDisplayMap["expire_time"] = backUpJobDisplay.ExpireTime
			}

			if backUpJobDisplay.JobStatus != nil {
				backUpJobDisplayMap["job_status"] = backUpJobDisplay.JobStatus
			}

			ids = append(ids, strconv.FormatInt(*backUpJobDisplay.JobId, 10))
			tmpList = append(tmpList, backUpJobDisplayMap)
		}

		_ = d.Set("back_up_jobs", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
