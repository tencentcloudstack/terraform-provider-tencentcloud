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

func ResourceTencentCloudWedataResourceFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataResourceFolderCreate,
		Read:   resourceTencentCloudWedataResourceFolderRead,
		Update: resourceTencentCloudWedataResourceFolderUpdate,
		Delete: resourceTencentCloudWedataResourceFolderDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project id.",
			},

			"parent_folder_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Absolute path of parent folder, value example/wedata/test, root directory, please use/.",
			},

			"folder_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Folder name.",
			},
		},
	}
}

func resourceTencentCloudWedataResourceFolderCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_folder.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		projectId string
		folderId  string
	)
	var (
		request  = wedatav20250806.NewCreateResourceFolderRequest()
		response = wedatav20250806.NewCreateResourceFolderResponse()
	)

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(projectId)
	}

	if v, ok := d.GetOk("parent_folder_path"); ok {
		request.ParentFolderPath = helper.String(v.(string))
	}

	if v, ok := d.GetOk("folder_name"); ok {
		request.FolderName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateResourceFolderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata resource folder failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Data != nil && response.Response.Data.FolderId != nil {
		folderId = *response.Response.Data.FolderId
		d.SetId(strings.Join([]string{projectId, folderId}, tccommon.FILED_SP))
	}

	return resourceTencentCloudWedataResourceFolderRead(d, meta)
}

func resourceTencentCloudWedataResourceFolderRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_folder.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	folderId := idSplit[1]

	parentFolderPath := d.Get("parent_folder_path").(string)
	respData, err := service.DescribeWedataResourceFolderById(ctx, projectId, folderId, parentFolderPath)
	if err != nil {
		return err
	}

	if len(respData) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `wedata_resource_folder` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	folder := respData[0]
	if folder.FolderName != nil {
		_ = d.Set("folder_name", folder.FolderName)
	}

	_ = projectId
	_ = folderId
	return nil
}

func resourceTencentCloudWedataResourceFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_folder.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"parent_folder_path"}
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
	folderId := idSplit[1]

	needChange := false
	mutableArgs := []string{"project_id", "folder_id", "folder_name"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wedatav20250806.NewUpdateResourceFolderRequest()
		request.ProjectId = helper.String(projectId)
		request.FolderId = helper.String(folderId)

		if v, ok := d.GetOk("folder_name"); ok {
			request.FolderName = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateResourceFolderWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update wedata resource folder failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = projectId
	_ = folderId
	return resourceTencentCloudWedataResourceFolderRead(d, meta)
}

func resourceTencentCloudWedataResourceFolderDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_folder.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	folderId := idSplit[1]

	var (
		request  = wedatav20250806.NewDeleteResourceFolderRequest()
		response = wedatav20250806.NewDeleteResourceFolderResponse()
	)

	request.ProjectId = helper.String(projectId)
	request.FolderId = helper.String(folderId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteResourceFolderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete wedata resource folder failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = projectId
	_ = folderId
	return nil
}
