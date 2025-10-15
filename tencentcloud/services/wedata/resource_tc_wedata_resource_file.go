package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataResourceFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataResourceFileCreate,
		Read:   resourceTencentCloudWedataResourceFileRead,
		Update: resourceTencentCloudWedataResourceFileUpdate,
		Delete: resourceTencentCloudWedataResourceFileDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project id.",
			},

			"resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The resource file name should be consistent with the uploaded file name as much as possible.",
			},

			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cos bucket name, which can be obtained from the GetResourceCosPath interface.",
			},

			"cos_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cos bucket area corresponding to the BucketName bucket.",
			},

			"parent_folder_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path to upload resource files in the project, example value: /wedata/qxxxm/, root directory, please use/.",
			},

			"resource_file": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "- You can only choose one of the two methods of uploading a file and manually filling. If both are provided, the order of values is file> manual filling value\n-the manual filling value must be the existing cos path, /datastudio/resource/is a fixed prefix, projectId is the project ID, and a specific value needs to be passed in, parentFolderPath is the parent folder path, name is the file name, and examples of manual filling value values are: /datastudio/resource/projectId/parentFolderPath/name \n.",
			},

			"bundle_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "bundle client ID.",
			},

			"bundle_info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "bundle client information.",
			},
		},
	}
}

func resourceTencentCloudWedataResourceFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_file.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		projectId  string
		resourceId string
	)
	var (
		request  = wedatav20250806.NewCreateResourceFileRequest()
		response = wedatav20250806.NewCreateResourceFileResponse()
	)

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(projectId)
	}

	if v, ok := d.GetOk("resource_name"); ok {
		request.ResourceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bucket_name"); ok {
		request.BucketName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_region"); ok {
		request.CosRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parent_folder_path"); ok {
		request.ParentFolderPath = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_file"); ok {
		request.ResourceFile = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bundle_id"); ok {
		request.BundleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bundle_info"); ok {
		request.BundleInfo = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateResourceFileWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata resource file failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Data != nil && response.Response.Data.ResourceId != nil {
		resourceId = *response.Response.Data.ResourceId
		d.SetId(strings.Join([]string{projectId, resourceId}, tccommon.FILED_SP))

	}

	return resourceTencentCloudWedataResourceFileRead(d, meta)
}

func resourceTencentCloudWedataResourceFileRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_file.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	resourceId := idSplit[1]

	respData, err := service.DescribeWedataResourceFileById(ctx, projectId, resourceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `wedata_resource_file` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)
	if respData.ResourceName != nil {
		_ = d.Set("resource_name", respData.ResourceName)
	}

	if respData.BucketName != nil {
		_ = d.Set("bucket_name", respData.BucketName)
	}

	if respData.CosRegion != nil {
		_ = d.Set("cos_region", respData.CosRegion)
	}

	if respData.BundleId != nil {
		_ = d.Set("bundle_id", respData.BundleId)
	}

	if respData.BundleInfo != nil {
		_ = d.Set("bundle_info", respData.BundleInfo)
	}

	_ = projectId
	_ = resourceId
	return nil
}

func resourceTencentCloudWedataResourceFileUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_file.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"bucket_name", "cos_region", "parent_folder_path"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	resourceId := idSplit[1]

	needChange := false
	mutableArgs := []string{"resource_file", "resource_name", "bundle_id", "bundle_info"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wedatav20250806.NewUpdateResourceFileRequest()
		request.ProjectId = helper.String(projectId)
		request.ResourceId = helper.String(resourceId)

		if v, ok := d.GetOk("resource_file"); ok {
			request.ResourceFile = helper.String(v.(string))
		}

		if v, ok := d.GetOk("resource_name"); ok {
			request.ResourceName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("bundle_id"); ok {
			request.BundleId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("bundle_info"); ok {
			request.BundleInfo = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateResourceFileWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update wedata resource file failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = projectId
	_ = resourceId
	return resourceTencentCloudWedataResourceFileRead(d, meta)
}

func resourceTencentCloudWedataResourceFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_file.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	resourceId := idSplit[1]

	var (
		request  = wedatav20250806.NewDeleteResourceFileRequest()
		response = wedatav20250806.NewDeleteResourceFileResponse()
	)

	request.ProjectId = helper.String(projectId)
	request.ResourceId = helper.String(resourceId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteResourceFileWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete wedata resource file failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = projectId
	_ = resourceId
	return nil
}
