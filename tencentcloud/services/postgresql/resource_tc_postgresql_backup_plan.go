package postgresql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlBackupPlanCreate,
		Read:   resourceTencentCloudPostgresqlBackupPlanRead,
		Update: resourceTencentCloudPostgresqlBackupPlanUpdate,
		Delete: resourceTencentCloudPostgresqlBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID.",
			},

			"plan_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Backup plan name.",
			},

			"backup_period_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Backup period type, currently only supports month.",
			},

			"backup_period": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Backup dates, such as backing up on the 2nd of each month.",
			},

			"min_backup_start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The earliest time to start a backup.",
			},

			"max_backup_start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The latest time to start a backup.",
			},

			"base_backup_retention_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Data backup retention period in days. Value range: [0, 30000).",
			},

			"log_backup_retention_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Log backup retention period in days. Value range: 7-1830.",
			},

			"backup_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Backup method. Enumerated values: physical, logical, snapshot.",
			},

			// computed
			"plan_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Backup plan ID.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_backup_plan.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request      = postgres.NewCreateBackupPlanRequest()
		response     = postgres.NewCreateBackupPlanResponse()
		dbInstanceId string
		planId       string
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dbInstanceId = v.(string)
	}

	if v, ok := d.GetOk("plan_name"); ok {
		request.PlanName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_period_type"); ok {
		request.BackupPeriodType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_period"); ok {
		backupPeriodSet := v.(*schema.Set).List()
		for i := range backupPeriodSet {
			if backupPeriodSet[i] != nil {
				backupPeriod := backupPeriodSet[i].(string)
				request.BackupPeriod = append(request.BackupPeriod, &backupPeriod)
			}
		}
	}

	if v, ok := d.GetOk("min_backup_start_time"); ok {
		request.MinBackupStartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_backup_start_time"); ok {
		request.MaxBackupStartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("base_backup_retention_period"); ok {
		request.BaseBackupRetentionPeriod = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateBackupPlanWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create postgresql backup_plan failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create postgresql backup_plan failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[CRUD]%s create postgresql backup_plan, dbInstanceId=%s, d.Id()=%s", logId, dbInstanceId, d.Id())

	if response.Response.PlanId == nil || *response.Response.PlanId == "" {
		return fmt.Errorf("Create postgresql backup_plan failed, PlanId is nil or empty.")
	}

	planId = *response.Response.PlanId
	d.SetId(strings.Join([]string{dbInstanceId, planId}, tccommon.FILED_SP))
	return resourceTencentCloudPostgresqlBackupPlanUpdate(d, meta)
}

func resourceTencentCloudPostgresqlBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_backup_plan.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dbInstanceId := idSplit[0]
	planId := idSplit[1]

	backupPlan, err := service.DescribePostgresqlBackupPlanById(ctx, dbInstanceId, planId)
	if err != nil {
		return err
	}

	if backupPlan == nil {
		log.Printf("[WARN]%s resource `tencentcloud_postgresql_backup_plan` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("db_instance_id", dbInstanceId)

	if backupPlan.PlanName != nil {
		_ = d.Set("plan_name", backupPlan.PlanName)
	}

	if backupPlan.BackupPeriodType != nil {
		_ = d.Set("backup_period_type", backupPlan.BackupPeriodType)
	}

	if backupPlan.BackupPeriod != nil {
		var newJson interface{}
		if err := json.Unmarshal([]byte(*backupPlan.BackupPeriod), &newJson); err != nil {
			return fmt.Errorf("convert BackupPeriod from string to interface{} failed, reason:%+v", err)
		}

		_ = d.Set("backup_period", newJson)
	}

	if backupPlan.MinBackupStartTime != nil {
		_ = d.Set("min_backup_start_time", backupPlan.MinBackupStartTime)
	}

	if backupPlan.MaxBackupStartTime != nil {
		_ = d.Set("max_backup_start_time", backupPlan.MaxBackupStartTime)
	}

	if backupPlan.BaseBackupRetentionPeriod != nil {
		_ = d.Set("base_backup_retention_period", backupPlan.BaseBackupRetentionPeriod)
	}

	if backupPlan.LogBackupRetentionPeriod != nil {
		_ = d.Set("log_backup_retention_period", backupPlan.LogBackupRetentionPeriod)
	}

	if backupPlan.BackupMethod != nil {
		_ = d.Set("backup_method", backupPlan.BackupMethod)
	}

	if backupPlan.PlanId != nil {
		_ = d.Set("plan_id", backupPlan.PlanId)
	}

	return nil
}

func resourceTencentCloudPostgresqlBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_backup_plan.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dbInstanceId := idSplit[0]
	planId := idSplit[1]

	needChange := false
	mutableArgs := []string{"plan_name", "backup_period", "min_backup_start_time", "max_backup_start_time", "base_backup_retention_period", "log_backup_retention_period", "backup_method"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := postgres.NewModifyBackupPlanRequest()
		request.DBInstanceId = &dbInstanceId
		request.PlanId = &planId

		if v, ok := d.GetOk("plan_name"); ok {
			request.PlanName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("backup_period"); ok {
			backupPeriodSet := v.(*schema.Set).List()
			for i := range backupPeriodSet {
				if backupPeriodSet[i] != nil {
					backupPeriod := backupPeriodSet[i].(string)
					request.BackupPeriod = append(request.BackupPeriod, &backupPeriod)
				}
			}
		}

		if v, ok := d.GetOk("min_backup_start_time"); ok {
			request.MinBackupStartTime = helper.String(v.(string))
		}

		if v, ok := d.GetOk("max_backup_start_time"); ok {
			request.MaxBackupStartTime = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("base_backup_retention_period"); ok {
			request.BaseBackupRetentionPeriod = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("log_backup_retention_period"); ok {
			request.LogBackupRetentionPeriod = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("backup_method"); ok {
			request.BackupMethod = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyBackupPlanWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update postgresql backup plan failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudPostgresqlBackupPlanRead(d, meta)
}

func resourceTencentCloudPostgresqlBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_backup_plan.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = postgres.NewDeleteBackupPlanRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dbInstanceId := idSplit[0]
	planId := idSplit[1]

	request.DBInstanceId = &dbInstanceId
	request.PlanId = &planId

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DeleteBackupPlanWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete postgresql backup plan failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
