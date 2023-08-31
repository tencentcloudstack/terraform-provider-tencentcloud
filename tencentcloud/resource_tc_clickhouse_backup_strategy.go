/*
Provides a resource to create a clickhouse backup strategy

Example Usage

```hcl
resource "tencentcloud_clickhouse_backup" "backup" {
	instance_id = "cdwch-xxxxxx"
	cos_bucket_name = "xxxxxx"
}

resource "tencentcloud_clickhouse_backup_strategy" "backup_strategy" {
  instance_id = "cdwch-xxxxxx"
  data_backup_strategy {
    week_days = "3"
    retain_days = 2
    execute_hour = 1
    back_up_tables {
      database = "iac"
      table = "my_table"
      total_bytes = 0
      v_cluster = "default_cluster"
      ips = "10.0.0.35"
    }
  }
  meta_backup_strategy {
	week_days = "1"
	retain_days = 2
	execute_hour = 3
  }
}
```

Import

clickhouse backup_strategy can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_backup_strategy.backup_strategy instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudClickhouseBackupStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClickhouseBackupStrategyCreate,
		Read:   resourceTencentCloudClickhouseBackupStrategyRead,
		Update: resourceTencentCloudClickhouseBackupStrategyUpdate,
		Delete: resourceTencentCloudClickhouseBackupStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"data_backup_strategy": {
				Required:    true,
				MinItems:    1,
				Type:        schema.TypeList,
				Description: "Data backup strategy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"week_days": {
							Required:    true,
							Type:        schema.TypeString,
							Description: "The day of the week is separated by commas. For example: 2 represents Tuesday.",
						},

						"execute_hour": {
							Required:    true,
							Type:        schema.TypeInt,
							Description: "Execution hour.",
						},

						"retain_days": {
							Required:    true,
							Type:        schema.TypeInt,
							Description: "Retention days.",
						},
						"back_up_tables": {
							Required:    true,
							Type:        schema.TypeList,
							Description: "Back up the list of tables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Database.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Table.",
									},
									"total_bytes": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Back up the list of tables.",
									},
									"v_cluster": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Virtual clusters.",
									},
									"ips": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Table ip.",
									},
									"zoo_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "ZK path.",
									},
									"rip": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Ip address of cvm.",
									},
								},
							},
						},
					},
				},
			},

			"meta_backup_strategy": {
				Optional:    true,
				Computed:    true,
				MinItems:    1,
				Type:        schema.TypeList,
				Description: "Metadata backup strategy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"week_days": {
							Optional:    true,
							Type:        schema.TypeString,
							Description: "The day of the week is separated by commas. For example: 2 represents Tuesday.",
						},

						"execute_hour": {
							Optional:    true,
							Type:        schema.TypeInt,
							Description: "Execution hour.",
						},

						"retain_days": {
							Optional:    true,
							Type:        schema.TypeInt,
							Description: "Retention days.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClickhouseBackupStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_backup_strategy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	instanceId := d.Get("instance_id").(string)

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
	backUpSchedule, err := service.DescribeBackUpScheduleById(ctx, instanceId)
	if err != nil {
		return err
	}
	if backUpSchedule == nil || backUpSchedule.MetaStrategy == nil || backUpSchedule.MetaStrategy.ScheduleId == nil {
		return fmt.Errorf("Can't fetch scheduleId!")
	}
	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("data_backup_strategy"); ok && len(v.([]interface{})) > 0 {
		dataBackupStrategy := v.([]interface{})[0].(map[string]interface{})
		paramMap["operation_type"] = OPERATION_TYPE_CREATE
		paramMap["schedule_type"] = SCHEDULE_TYPE_DATA
		if v, ok := dataBackupStrategy["week_days"]; ok {
			paramMap["week_days"] = v.(string)
		}
		if v, ok := dataBackupStrategy["execute_hour"]; ok {
			paramMap["execute_hour"] = v.(int)
		}
		if v, ok := dataBackupStrategy["retain_days"]; ok {
			paramMap["retain_days"] = v.(int)
		}
		if v, ok := dataBackupStrategy["back_up_tables"]; ok {
			paramMap["back_up_tables"] = v.([]interface{})
		}
	}
	err = service.CreateBackUpSchedule(ctx, instanceId, paramMap)
	if err != nil {
		log.Printf("[CRITAL]%s create clickhouse data backup strategy failed, reason:%+v", logId, err)
		return err
	}
	d.SetId(instanceId)

	return resourceTencentCloudClickhouseBackupStrategyUpdate(d, meta)
}

func resourceTencentCloudClickhouseBackupStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_backup_strategy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()
	_ = d.Set("instance_id", instanceId)

	backUpSchedule, err := service.DescribeBackUpScheduleById(ctx, instanceId)
	if err != nil {
		return err
	}

	if backUpSchedule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClickhouseBackUpSchedule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if backUpSchedule.MetaStrategy != nil {
		metaBackupStrategyMap := make(map[string]interface{})
		metaBackupStrategyMap["week_days"] = backUpSchedule.MetaStrategy.WeekDays
		metaBackupStrategyMap["execute_hour"] = backUpSchedule.MetaStrategy.ExecuteHour
		metaBackupStrategyMap["retain_days"] = backUpSchedule.MetaStrategy.RetainDays
		_ = d.Set("meta_backup_strategy", []interface{}{metaBackupStrategyMap})
	}

	dataBackupStrategyMap := make(map[string]interface{})
	if backUpSchedule.DataStrategy != nil {
		dataBackupStrategyMap["week_days"] = backUpSchedule.DataStrategy.WeekDays
		dataBackupStrategyMap["execute_hour"] = backUpSchedule.DataStrategy.ExecuteHour
		dataBackupStrategyMap["retain_days"] = backUpSchedule.DataStrategy.RetainDays
	}

	if backUpSchedule.BackUpContents != nil {
		backUpTablesList := []interface{}{}
		for _, backUpContent := range backUpSchedule.BackUpContents {
			backUpContentMap := map[string]interface{}{}

			if backUpContent.Database != nil {
				backUpContentMap["database"] = backUpContent.Database
			}

			if backUpContent.Table != nil {
				backUpContentMap["table"] = backUpContent.Table
			}

			if backUpContent.TotalBytes != nil {
				backUpContentMap["total_bytes"] = backUpContent.TotalBytes
			}

			if backUpContent.VCluster != nil {
				backUpContentMap["v_cluster"] = backUpContent.VCluster
			}

			if backUpContent.Ips != nil {
				backUpContentMap["ips"] = backUpContent.Ips
			}

			if backUpContent.ZooPath != nil {
				backUpContentMap["zoo_path"] = backUpContent.ZooPath
			}

			if backUpContent.Rip != nil {
				backUpContentMap["rip"] = backUpContent.Rip
			}

			backUpTablesList = append(backUpTablesList, backUpContentMap)
		}

		dataBackupStrategyMap["back_up_tables"] = backUpTablesList
	}

	_ = d.Set("data_backup_strategy", []interface{}{dataBackupStrategyMap})

	return nil
}

func resourceTencentCloudClickhouseBackupStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_backup_strategy.update")()
	defer inconsistentCheck(d, meta)()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	backUpSchedule, err := service.DescribeBackUpScheduleById(ctx, instanceId)
	if err != nil {
		return err
	}
	if backUpSchedule == nil || backUpSchedule.MetaStrategy == nil || backUpSchedule.MetaStrategy.ScheduleId == nil {
		return fmt.Errorf("Can't fetch meta scheduleId!")
	}
	if backUpSchedule == nil || backUpSchedule.DataStrategy == nil || backUpSchedule.DataStrategy.ScheduleId == nil {
		return fmt.Errorf("Can't fetch data scheduleId!")
	}

	if d.HasChange("data_backup_strategy") {
		paramMap := make(map[string]interface{})
		dataBackupStrategyList := d.Get("data_backup_strategy").([]interface{})
		dataBackupStrategy := dataBackupStrategyList[0].(map[string]interface{})
		paramMap["operation_type"] = OPERATION_TYPE_UPDATE
		paramMap["schedule_type"] = SCHEDULE_TYPE_DATA
		dataScheduleId := *backUpSchedule.DataStrategy.ScheduleId
		paramMap["schedule_id"] = dataScheduleId
		if v, ok := dataBackupStrategy["week_days"]; ok {
			paramMap["week_days"] = v.(string)
		}
		if v, ok := dataBackupStrategy["execute_hour"]; ok {
			paramMap["execute_hour"] = v.(int)
		}
		if v, ok := dataBackupStrategy["retain_days"]; ok {
			paramMap["retain_days"] = v.(int)
		}
		if v, ok := dataBackupStrategy["back_up_tables"]; ok {
			paramMap["back_up_tables"] = v.([]interface{})
		}

		err = service.CreateBackUpSchedule(ctx, instanceId, paramMap)
		if err != nil {
			log.Printf("[CRITAL]%s create clickhouse data backup strategy failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("meta_backup_strategy") {
		if v, ok := d.GetOk("meta_backup_strategy"); ok {
			paramMap := make(map[string]interface{})
			metaBackupStrategyList := v.([]interface{})
			metaBackupStrategy := metaBackupStrategyList[0].(map[string]interface{})
			paramMap["operation_type"] = OPERATION_TYPE_UPDATE
			paramMap["schedule_type"] = SCHEDULE_TYPE_META
			metaScheduleId := *backUpSchedule.MetaStrategy.ScheduleId
			paramMap["schedule_id"] = metaScheduleId
			if v, ok := metaBackupStrategy["week_days"]; ok {
				paramMap["week_days"] = v.(string)
			}
			if v, ok := metaBackupStrategy["execute_hour"]; ok {
				paramMap["execute_hour"] = v.(int)
			}
			if v, ok := metaBackupStrategy["retain_days"]; ok {
				paramMap["retain_days"] = v.(int)
			}

			err = service.CreateBackUpSchedule(ctx, instanceId, paramMap)
			if err != nil {
				log.Printf("[CRITAL]%s create clickhouse meta backup strategy failed, reason:%+v", logId, err)
				return err
			}
		}

	}

	return resourceTencentCloudClickhouseBackupStrategyRead(d, meta)
}

func resourceTencentCloudClickhouseBackupStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_backup_strategy.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
