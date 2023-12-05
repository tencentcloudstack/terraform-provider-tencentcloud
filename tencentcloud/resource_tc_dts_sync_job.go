package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsSyncJob() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDtsSyncJobRead,
		Create: resourceTencentCloudDtsSyncJobCreate,
		Delete: resourceTencentCloudDtsSyncJobDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"pay_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "pay mode, optional value is PrePay or PostPay.",
			},

			"src_database_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "source database type.",
			},

			"src_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "source region.",
			},

			"dst_database_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "destination database type.",
			},

			"dst_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "destination region.",
			},

			"specification": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "specification.",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag value.",
						},
					},
				},
			},

			"auto_renew": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "auto renew.",
			},

			"instance_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "instance class.",
			},

			"job_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "job name.",
			},

			"existed_job_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "existed job id.",
			},

			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "job id.",
			},
		},
	}
}

func resourceTencentCloudDtsSyncJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
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

	if v, ok := d.GetOk("auto_renew"); ok {
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().CreateSyncJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dts syncJob failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobIds[0]

	d.SetId(jobId)
	return resourceTencentCloudDtsSyncJobRead(d, meta)
}

func resourceTencentCloudDtsSyncJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	syncJobId := d.Id()

	syncJob, err := service.DescribeDtsSyncJob(ctx, helper.String(syncJobId))

	if err != nil {
		return err
	}

	if syncJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `syncJob` %s does not exist", syncJobId)
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

	if syncJob.JobName != nil {
		_ = d.Set("job_name", syncJob.JobName)
	}

	if syncJob.JobId != nil {
		_ = d.Set("job_id", syncJob.JobId)
	}

	return nil
}

func resourceTencentCloudDtsSyncJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	syncJobId := d.Id()

	if err := service.IsolateDtsSyncJobById(ctx, syncJobId); err != nil {
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{"Isolated", "Stopped"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobConfigIsolateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	if err := service.DestroyDtsSyncJobById(ctx, syncJobId); err != nil {
		return err
	}

	conf = BuildStateChangeConf([]string{}, []string{"Offlined"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobConfigDeleteStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
