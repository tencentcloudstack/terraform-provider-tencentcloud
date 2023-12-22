package sqlserver

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverBusinessIntelligenceFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverBusinessIntelligenceFileCreate,
		Read:   resourceTencentCloudSqlserverBusinessIntelligenceFileRead,
		Delete: resourceTencentCloudSqlserverBusinessIntelligenceFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},
			"file_url": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cos Url.",
			},
			"file_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "File Type FLAT - Flat File as Data Source, SSIS - ssis project package.",
			},
			"remark": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "remark.",
			},
		},
	}
}

func resourceTencentCloudSqlserverBusinessIntelligenceFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_business_intelligence_file.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		CreateBIFRequest  = sqlserver.NewCreateBusinessIntelligenceFileRequest()
		CreateBIFResponse = sqlserver.NewCreateBusinessIntelligenceFileResponse()
		instanceId        string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		CreateBIFRequest.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("file_url"); ok {
		CreateBIFRequest.FileURL = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_type"); ok {
		CreateBIFRequest.FileType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		CreateBIFRequest.Remark = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().CreateBusinessIntelligenceFile(CreateBIFRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, CreateBIFRequest.GetAction(), CreateBIFRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver CreateBusinessIntelligenceFile not exists")
			return resource.NonRetryableError(e)
		}

		CreateBIFResponse = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver businessIntelligenceFile failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{instanceId, *CreateBIFResponse.Response.FileTaskId}, tccommon.FILED_SP))

	return resourceTencentCloudSqlserverBusinessIntelligenceFileRead(d, meta)
}

func resourceTencentCloudSqlserverBusinessIntelligenceFileRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_business_intelligence_file.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	fileName := idSplit[1]

	businessIntelligenceFile, err := service.DescribeSqlserverBusinessIntelligenceFileById(ctx, instanceId, fileName)
	if err != nil {
		return err
	}

	if businessIntelligenceFile == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverBusinessIntelligenceFile` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if businessIntelligenceFile.InstanceId != nil {
		_ = d.Set("instance_id", businessIntelligenceFile.InstanceId)
	}

	if businessIntelligenceFile.FileURL != nil {
		_ = d.Set("file_url", businessIntelligenceFile.FileURL)
	}

	if businessIntelligenceFile.FileType != nil {
		_ = d.Set("file_type", businessIntelligenceFile.FileType)
	}

	if businessIntelligenceFile.Remark != nil {
		_ = d.Set("remark", businessIntelligenceFile.Remark)
	}

	return nil
}

func resourceTencentCloudSqlserverBusinessIntelligenceFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_business_intelligence_file.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	fileName := idSplit[1]

	if err := service.DeleteSqlserverBusinessIntelligenceFileById(ctx, instanceId, fileName); err != nil {
		return err
	}

	return nil
}
