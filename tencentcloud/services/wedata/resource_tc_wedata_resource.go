package wedata

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataResourceCreate,
		Read:   resourceTencentCloudWedataResourceRead,
		Update: resourceTencentCloudWedataResourceUpdate,
		Delete: resourceTencentCloudWedataResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},
			"file_path": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "For file path:/datastudio/resource/projectId/folderName; for folder path:/datastudio/resource/folderName.",
			},
			"file_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File name.",
			},
			"cos_bucket_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cos bucket name.",
			},
			"cos_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cos bucket region.",
			},
			"files_size": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File size.",
			},
			"resource_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Resource ID.",
			},
		},
	}
}

func resourceTencentCloudWedataResourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = wedata.NewCreateOrUpdateResourceRequest()
		response   = wedata.NewCreateOrUpdateResourceResponse()
		projectId  string
		filePath   string
		resourceId string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("file_path"); ok {
		request.FilePath = helper.String(v.(string))
		filePath = v.(string)
	}

	if v, ok := d.GetOk("file_name"); ok {
		request.Files = append(request.Files, helper.String(v.(string)))
	}

	if v, ok := d.GetOk("cos_bucket_name"); ok {
		request.CosBucketName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_region"); ok {
		request.CosRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("files_size"); ok {
		request.FilesSize = append(request.FilesSize, helper.String(v.(string)))
	}

	request.NewFile = helper.Bool(true)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().CreateOrUpdateResource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || len(result.Response.Data) == 0 {
			e = fmt.Errorf("wedata resource not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata resource failed, reason:%+v", logId, err)
		return err
	}

	resourceId = *response.Response.Data[0].ResourceId
	d.SetId(strings.Join([]string{projectId, filePath, resourceId}, tccommon.FILED_SP))

	return resourceTencentCloudWedataResourceRead(d, meta)
}

func resourceTencentCloudWedataResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	filePath := idSplit[1]
	resourceId := idSplit[2]

	resourceInfo, err := service.DescribeWedataResourceById(ctx, projectId, filePath, resourceId)
	if err != nil {
		return err
	}

	if resourceInfo == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataResource` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)
	_ = d.Set("file_path", filePath)

	if resourceInfo.ResourceId != nil {
		_ = d.Set("resource_id", resourceInfo.ResourceId)
	}

	if resourceInfo.Name != nil {
		_ = d.Set("file_name", resourceInfo.Name)
	}

	if resourceInfo.CosBucket != nil {
		_ = d.Set("cos_bucket_name", resourceInfo.CosBucket)
	}

	if resourceInfo.CosRegion != nil {
		_ = d.Set("cos_region", resourceInfo.CosRegion)
	}

	if resourceInfo.Size != nil {
		sizeStr := strconv.FormatInt(*resourceInfo.Size, 10)
		_ = d.Set("files_size", sizeStr)
	}

	return nil
}

func resourceTencentCloudWedataResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = wedata.NewCreateOrUpdateResourceRequest()
	)

	immutableArgs := []string{"file_path", "project_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	filePath := idSplit[1]

	request.ProjectId = &projectId
	request.FilePath = &filePath

	if v, ok := d.GetOk("file_name"); ok {
		request.Files = append(request.Files, helper.String(v.(string)))
	}

	if v, ok := d.GetOk("cos_bucket_name"); ok {
		request.CosBucketName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_region"); ok {
		request.CosRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("files_size"); ok {
		request.FilesSize = append(request.FilesSize, helper.String(v.(string)))
	}

	request.NewFile = helper.Bool(false)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().CreateOrUpdateResource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update wedata resource failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataResourceRead(d, meta)
}

func resourceTencentCloudWedataResourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	resourceId := idSplit[2]

	if err := service.DeleteWedataResourceById(ctx, projectId, resourceId); err != nil {
		return err
	}

	return nil
}
