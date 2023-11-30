package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsMigrateService() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDtsMigrateServiceRead,
		Create: resourceTencentCloudDtsMigrateServiceCreate,
		Update: resourceTencentCloudDtsMigrateServiceUpdate,
		Delete: resourceTencentCloudDtsMigrateServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"src_database_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "source database type, optional value is mysql/redis/percona/mongodb/postgresql/sqlserver/mariadb.",
			},

			"dst_database_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "destination database type, optional value is mysql/redis/percona/mongodb/postgresql/sqlserver/mariadb.",
			},

			"src_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "source region.",
			},

			"dst_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "destination region.",
			},

			"instance_class": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance class, optional value is small/medium/large/xlarge/2xlarge.",
			},

			"job_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "job name.",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
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
			}, //tags
		},
	}
}

func resourceTencentCloudDtsMigrateServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_service.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		request  = dts.NewCreateMigrationServiceRequest()
		response *dts.CreateMigrationServiceResponse
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		service  = DtsService{client: meta.(*TencentCloudClient).apiV3Conn}
		jobId    string
	)

	if v, ok := d.GetOk("src_database_type"); ok {
		request.SrcDatabaseType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_database_type"); ok {
		request.DstDatabaseType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_region"); ok {
		request.SrcRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_region"); ok {
		request.DstRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_class"); ok {
		request.InstanceClass = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_name"); ok {
		request.JobName = helper.String(v.(string))
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().CreateMigrationService(request)
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
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}
	if response.Response == nil || response.Response.JobIds == nil {
		return fmt.Errorf("%s create dts migrateJob failed, response is nil!", logId)
	}

	jobId = *response.Response.JobIds[0]
	// wait created
	if err = service.PollingMigrateJobStatusUntil(ctx, jobId, DTSJobStatus, []string{"created"}); err != nil {
		return err
	}

	d.SetId(jobId)
	return resourceTencentCloudDtsMigrateServiceRead(d, meta)
}

func resourceTencentCloudDtsMigrateServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_service.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = DtsService{client: meta.(*TencentCloudClient).apiV3Conn}
		jobId   = d.Id()
	)

	migrateJob, err := service.DescribeDtsMigrateJobById(ctx, jobId)

	if err != nil {
		return err
	}

	if migrateJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `migrateJob` %s does not exist", jobId)
	}

	if migrateJob.SrcInfo != nil {
		srcInfo := migrateJob.SrcInfo
		if srcInfo.DatabaseType != nil {
			_ = d.Set("src_database_type", srcInfo.DatabaseType)
		}

		if srcInfo.Region != nil {
			_ = d.Set("src_region", srcInfo.Region)
		}
	}

	if migrateJob.DstInfo != nil {
		destInfo := migrateJob.DstInfo
		if destInfo.DatabaseType != nil {
			_ = d.Set("dst_database_type", destInfo.DatabaseType)
		}

		if destInfo.Region != nil {
			_ = d.Set("dst_region", destInfo.Region)
		}
	}

	if migrateJob.TradeInfo != nil && migrateJob.TradeInfo.InstanceClass != nil {
		_ = d.Set("instance_class", migrateJob.TradeInfo.InstanceClass)
	}

	if migrateJob.JobName != nil {
		_ = d.Set("job_name", migrateJob.JobName)
	}

	if migrateJob.Tags != nil {
		tagsList := []interface{}{}
		for _, tags := range migrateJob.Tags {
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

	return nil
}

func resourceTencentCloudDtsMigrateServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_service.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dts.NewModifyMigrateJobSpecRequest()
	request.JobId = helper.String(d.Id())

	if v, ok := d.GetOk("instance_class"); ok {
		request.NewInstanceClass = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().ModifyMigrateJobSpec(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateService failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDtsMigrateServiceRead(d, meta)
}

func resourceTencentCloudDtsMigrateServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_service.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	migrateJobId := d.Id()

	if err := service.DeleteDtsMigrateServiceById(ctx, migrateJobId); err != nil {
		return err
	}

	return nil
}
