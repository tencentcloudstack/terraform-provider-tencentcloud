package dts

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDtsSyncJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsSyncJobCreate,
		Read:   resourceTencentCloudDtsSyncJobRead,
		Update: resourceTencentCloudDtsSyncJobUpdate,
		Delete: resourceTencentCloudDtsSyncJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"pay_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Billing type. Valid values: `PrePay` (subscription, monthly/yearly billing), `PostPay` (pay-as-you-go).",
			},

			"src_database_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source database type, such as `mysql`, `mariadb`, `percona`, `postgresql`, `cynosdbmysql` (TDSQL-C MySQL), `tdpg` (TDSQL for PostgreSQL), `tdsqlmysql`, `tdstore` (TDSQL TDStore), etc.",
			},

			"src_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region where the source database resides, such as `ap-guangzhou`.",
			},

			"dst_database_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Destination database type, such as `mysql`, `mariadb`, `percona`, `cynosdbmysql` (TDSQL-C MySQL), `tdpg` (TDSQL for PostgreSQL), `tdsqlmysql`, `kafka`, `tdstore` (TDSQL TDStore), etc.",
			},

			"dst_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region where the destination database resides, such as `ap-guangzhou`.",
			},

			"specification": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Sync job specification. `Standard` indicates the standard edition; currently only `Standard` is supported.",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Tag information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"auto_renew": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Auto-renewal flag. Only takes effect when `pay_mode` is `PrePay`. Valid values: `1` (enable auto-renewal), `0` (disable auto-renewal, default).",
			},

			"instance_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Sync link specification, such as `micro`, `small`, `medium`, `large`. Default is `medium`.",
			},

			"job_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Sync job name.",
			},

			"existed_job_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The existing sync job ID used to create a similar job.",
			},

			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sync job ID.",
			},
		},
	}
}

func resourceTencentCloudDtsSyncJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_job.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = dts.NewCreateSyncJobRequest()
		response *dts.CreateSyncJobResponse
		jobId    string
	)

	if v, ok := d.GetOk("pay_mode"); ok {
		request.PayMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_database_type"); ok {
		request.SrcDatabaseType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_region"); ok {
		request.SrcRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_database_type"); ok {
		request.DstDatabaseType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_region"); ok {
		request.DstRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("specification"); ok {
		request.Specification = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tagItem := dts.TagItem{}
			if v, ok := dMap["tag_key"]; ok {
				tagItem.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				tagItem.TagValue = helper.String(v.(string))
			}
			request.Tags = append(request.Tags, &tagItem)
		}
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request.AutoRenew = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("instance_class"); ok {
		request.InstanceClass = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_name"); ok {
		request.JobName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("existed_job_id"); ok {
		request.ExistedJobId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDtsClient().CreateSyncJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dts syncJob failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dts syncJob failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.JobIds == nil || len(response.Response.JobIds) == 0 {
		return fmt.Errorf("JobIds is nil.")
	}

	jobId = *response.Response.JobIds[0]
	d.SetId(jobId)
	return resourceTencentCloudDtsSyncJobRead(d, meta)
}

func resourceTencentCloudDtsSyncJobRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_job.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		syncJobId = d.Id()
	)

	syncJob, err := service.DescribeDtsSyncJob(ctx, helper.String(syncJobId))
	if err != nil {
		return err
	}

	if syncJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `tencentcloud_dts_sync_job` %s does not exist", syncJobId)
	}

	if syncJob.PayMode != nil {
		_ = d.Set("pay_mode", syncJob.PayMode)
	}

	if syncJob.SrcDatabaseType != nil {
		_ = d.Set("src_database_type", syncJob.SrcDatabaseType)
	}

	if syncJob.SrcRegion != nil {
		_ = d.Set("src_region", syncJob.SrcRegion)
	}

	if syncJob.DstDatabaseType != nil {
		_ = d.Set("dst_database_type", syncJob.DstDatabaseType)
	}

	if syncJob.DstRegion != nil {
		_ = d.Set("dst_region", syncJob.DstRegion)
	}

	if syncJob.Specification != nil {
		_ = d.Set("specification", syncJob.Specification)
	}

	if syncJob.Tags != nil {
		tagsList := []interface{}{}
		for _, tags := range syncJob.Tags {
			tagsMap := map[string]interface{}{}
			if tags.TagKey != nil {
				tagsMap["tag_key"] = tags.TagKey
			}
			if tags.TagValue != nil {
				tagsMap["tag_value"] = tags.TagValue
			}

			tagsList = append(tagsList, tagsMap)
		}
		_ = d.Set("tags", tagsList)
	}

	if syncJob.AutoRenew != nil {
		_ = d.Set("auto_renew", syncJob.AutoRenew)
	}

	if syncJob.InstanceClass != nil {
		_ = d.Set("instance_class", syncJob.InstanceClass)
	}

	if syncJob.JobName != nil {
		_ = d.Set("job_name", syncJob.JobName)
	}

	if syncJob.JobId != nil {
		_ = d.Set("job_id", syncJob.JobId)
	}

	return nil
}

func resourceTencentCloudDtsSyncJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_job.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		syncJobId = d.Id()
	)

	if d.HasChange("instance_class") {
		request := dts.NewResizeSyncJobRequest()
		if v, ok := d.GetOk("instance_class"); ok {
			request.NewInstanceClass = helper.String(v.(string))
		}

		request.JobId = &syncJobId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDtsClient().ResizeSyncJob(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Resize syncJob failed, Response is nil."))
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s resize dts syncJob failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudDtsSyncJobRead(d, meta)
}

func resourceTencentCloudDtsSyncJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_job.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		syncJobId = d.Id()
	)

	if err := service.IsolateDtsSyncJobById(ctx, syncJobId); err != nil {
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Isolated", "Stopped", "NotBilledByInternational", "NotBilled"}, 2*tccommon.ReadRetryTimeout, time.Second, service.DtsSyncJobConfigIsolateStateRefreshFunc(d.Id(), []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	syncConfig, e := service.DescribeDtsSyncConfigById(ctx, syncJobId)
	if e != nil {
		return e
	}

	if syncConfig != nil && syncConfig.TradeStatus != nil && (*syncConfig.TradeStatus == "NotBilledByInternational" || *syncConfig.TradeStatus == "NotBilled") {
		return nil
	}

	if err := service.DestroyDtsSyncJobById(ctx, syncJobId); err != nil {
		return err
	}

	conf = tccommon.BuildStateChangeConf([]string{}, []string{"Offlined"}, 2*tccommon.ReadRetryTimeout, time.Second, service.DtsSyncJobConfigDeleteStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
