/*
Provides a resource to create a dts migrate_check_job

Example Usage

```hcl
resource "tencentcloud_dts_migrate_check_job" "migrate_check_job" {
  job_id = &lt;nil&gt;
        }
```

Import

dts migrate_check_job can be imported using the id, e.g.

```
terraform import tencentcloud_dts_migrate_check_job.migrate_check_job migrate_check_job_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceTencentCloudDtsMigrateCheckJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateCheckJobCreate,
		Read:   resourceTencentCloudDtsMigrateCheckJobRead,
		Update: resourceTencentCloudDtsMigrateCheckJobUpdate,
		Delete: resourceTencentCloudDtsMigrateCheckJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job id.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Status.",
			},

			"brief_msg": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Brief message of check results.",
			},

			"step_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Step info of check.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"step_no": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Step number.",
						},
						"step_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Step id.",
						},
						"step_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Step name.",
						},
						"step_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Step status.",
						},
						"step_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Step message.",
						},
						"detail_check_items": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detail of check items.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"check_item_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of check item.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of check item.",
									},
									"check_result": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The check result of item.",
									},
									"failure_reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The failed reason of check item.",
									},
									"solution": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The solution of check item.",
									},
									"error_log": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "The error log of check item.",
									},
									"help_doc": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "The help doc of check item.",
									},
									"skip_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The skip info of check item.",
									},
								},
							},
						},
						"has_skipped": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Has skipped.",
						},
					},
				},
			},

			"check_flag": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Check results, optional value is checkPass or checkNotPass.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateCheckJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_check_job.create")()
	defer inconsistentCheck(d, meta)()

	var jobId string
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}

	d.SetId(jobId)

	return resourceTencentCloudDtsMigrateCheckJobUpdate(d, meta)
}

func resourceTencentCloudDtsMigrateCheckJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_check_job.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	migrateCheckJobId := d.Id()

	migrateCheckJob, err := service.DescribeDtsMigrateCheckJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if migrateCheckJob == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsMigrateCheckJob` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if migrateCheckJob.JobId != nil {
		_ = d.Set("job_id", migrateCheckJob.JobId)
	}

	if migrateCheckJob.Status != nil {
		_ = d.Set("status", migrateCheckJob.Status)
	}

	if migrateCheckJob.BriefMsg != nil {
		_ = d.Set("brief_msg", migrateCheckJob.BriefMsg)
	}

	if migrateCheckJob.StepInfo != nil {
		stepInfoList := []interface{}{}
		for _, stepInfo := range migrateCheckJob.StepInfo {
			stepInfoMap := map[string]interface{}{}

			if migrateCheckJob.StepInfo.StepNo != nil {
				stepInfoMap["step_no"] = migrateCheckJob.StepInfo.StepNo
			}

			if migrateCheckJob.StepInfo.StepId != nil {
				stepInfoMap["step_id"] = migrateCheckJob.StepInfo.StepId
			}

			if migrateCheckJob.StepInfo.StepName != nil {
				stepInfoMap["step_name"] = migrateCheckJob.StepInfo.StepName
			}

			if migrateCheckJob.StepInfo.StepStatus != nil {
				stepInfoMap["step_status"] = migrateCheckJob.StepInfo.StepStatus
			}

			if migrateCheckJob.StepInfo.StepMessage != nil {
				stepInfoMap["step_message"] = migrateCheckJob.StepInfo.StepMessage
			}

			if migrateCheckJob.StepInfo.DetailCheckItems != nil {
				detailCheckItemsList := []interface{}{}
				for _, detailCheckItems := range migrateCheckJob.StepInfo.DetailCheckItems {
					detailCheckItemsMap := map[string]interface{}{}

					if detailCheckItems.CheckItemName != nil {
						detailCheckItemsMap["check_item_name"] = detailCheckItems.CheckItemName
					}

					if detailCheckItems.Description != nil {
						detailCheckItemsMap["description"] = detailCheckItems.Description
					}

					if detailCheckItems.CheckResult != nil {
						detailCheckItemsMap["check_result"] = detailCheckItems.CheckResult
					}

					if detailCheckItems.FailureReason != nil {
						detailCheckItemsMap["failure_reason"] = detailCheckItems.FailureReason
					}

					if detailCheckItems.Solution != nil {
						detailCheckItemsMap["solution"] = detailCheckItems.Solution
					}

					if detailCheckItems.ErrorLog != nil {
						detailCheckItemsMap["error_log"] = detailCheckItems.ErrorLog
					}

					if detailCheckItems.HelpDoc != nil {
						detailCheckItemsMap["help_doc"] = detailCheckItems.HelpDoc
					}

					if detailCheckItems.SkipInfo != nil {
						detailCheckItemsMap["skip_info"] = detailCheckItems.SkipInfo
					}

					detailCheckItemsList = append(detailCheckItemsList, detailCheckItemsMap)
				}

				stepInfoMap["detail_check_items"] = []interface{}{detailCheckItemsList}
			}

			if migrateCheckJob.StepInfo.HasSkipped != nil {
				stepInfoMap["has_skipped"] = migrateCheckJob.StepInfo.HasSkipped
			}

			stepInfoList = append(stepInfoList, stepInfoMap)
		}

		_ = d.Set("step_info", stepInfoList)

	}

	if migrateCheckJob.CheckFlag != nil {
		_ = d.Set("check_flag", migrateCheckJob.CheckFlag)
	}

	return nil
}

func resourceTencentCloudDtsMigrateCheckJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_check_job.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"job_id", "status", "brief_msg", "step_info", "check_flag"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudDtsMigrateCheckJobRead(d, meta)
}

func resourceTencentCloudDtsMigrateCheckJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_check_job.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
