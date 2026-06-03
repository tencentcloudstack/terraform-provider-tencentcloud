package postgresql

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresAuditService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresAuditServiceCreate,
		Read:   resourceTencentCloudPostgresAuditServiceRead,
		Update: resourceTencentCloudPostgresAuditServiceUpdate,
		Delete: resourceTencentCloudPostgresAuditServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "PostgreSQL instance ID.",
			},
			"log_expire_day": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Log retention days. Valid values: 7, 30, 90, 180, 365, 1095, 1825.",
			},
			"hot_log_expire_day": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Hot log retention days. Valid values: 7, 30, 90, 180, 365, 1095, 1825.",
			},
			"audit_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Audit type. Valid values: complex (fine-grained audit), simple (fast audit).",
			},
			"audit_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Audit status. Values: ON, OFF.",
			},
			"cold_log_expire_day": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Cold log retention days.",
			},
			"hot_log_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Hot log size in MB.",
			},
			"cold_log_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Cold log size in MB.",
			},
		},
	}
}
func resourceTencentCloudPostgresAuditServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgres_audit_service.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = postgresql.NewOpenAuditServiceRequest()
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOkExists("log_expire_day"); ok {
		request.LogExpireDay = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("hot_log_expire_day"); ok {
		request.HotLogExpireDay = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("audit_type"); ok {
		request.AuditType = helper.String(v.(string))
	}

	request.Product = helper.String("postgres")
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		err := service.OpenAuditService(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create postgres audit service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(instanceId)
	return resourceTencentCloudPostgresAuditServiceRead(d, meta)
}
func resourceTencentCloudPostgresAuditServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgres_audit_service.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	instanceId := d.Id()

	request := postgresql.NewDescribeAuditInstanceListRequest()
	request.Product = helper.String("postgres")
	request.AuditSwitch = helper.Uint64(1)
	request.Limit = helper.Uint64(100)
	request.Offset = helper.Uint64(0)
	request.Filters = []*postgresql.Filter{
		{
			Name:   helper.String("InstanceId"),
			Values: []*string{helper.String(instanceId)},
		},
	}

	var auditInfo *postgresql.AuditInstanceInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		items, err := service.DescribeAuditInstanceList(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}

		for _, item := range items {
			if item.InstanceId != nil && *item.InstanceId == instanceId {
				auditInfo = item
				break
			}
		}

		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	if auditInfo == nil {
		log.Printf("[WARN]%s resource `tencentcloud_postgres_audit_service` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if auditInfo.AuditStatus != nil {
		_ = d.Set("audit_status", auditInfo.AuditStatus)
	}

	if auditInfo.LogExpireDay != nil {
		_ = d.Set("log_expire_day", auditInfo.LogExpireDay)
	}

	if auditInfo.HotLogExpireDay != nil {
		_ = d.Set("hot_log_expire_day", auditInfo.HotLogExpireDay)
	}

	if auditInfo.ColdLogExpireDay != nil {
		_ = d.Set("cold_log_expire_day", auditInfo.ColdLogExpireDay)
	}

	if auditInfo.HotLogSize != nil {
		_ = d.Set("hot_log_size", auditInfo.HotLogSize)
	}

	if auditInfo.ColdLogSize != nil {
		_ = d.Set("cold_log_size", auditInfo.ColdLogSize)
	}

	return nil
}
func resourceTencentCloudPostgresAuditServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgres_audit_service.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	instanceId := d.Id()

	needChange := false
	mutableArgs := []string{"log_expire_day", "hot_log_expire_day", "audit_type"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := postgresql.NewModifyAuditServiceRequest()
		request.InstanceId = helper.String(instanceId)

		if v, ok := d.GetOkExists("log_expire_day"); ok {
			request.LogExpireDay = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("hot_log_expire_day"); ok {
			request.HotLogExpireDay = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("audit_type"); ok {
			request.AuditType = helper.String(v.(string))
		}

		request.Product = helper.String("postgres")
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err := service.ModifyAuditService(ctx, request)
			if err != nil {
				return tccommon.RetryError(err)
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update postgres audit service failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudPostgresAuditServiceRead(d, meta)
}
func resourceTencentCloudPostgresAuditServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgres_audit_service.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = postgresql.NewCloseAuditServiceRequest()
	)

	instanceId := d.Id()
	request.InstanceId = helper.String(instanceId)
	request.Product = helper.String("postgres")

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		err := service.CloseAuditService(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete postgres audit service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
