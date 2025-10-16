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

func ResourceTencentCloudWedataWorkflowFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataWorkflowFolderCreate,
		Read:   resourceTencentCloudWedataWorkflowFolderRead,
		Update: resourceTencentCloudWedataWorkflowFolderUpdate,
		Delete: resourceTencentCloudWedataWorkflowFolderDelete,
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
				Description: "The absolute path of the parent folder, such as/abc/de, if it is the root directory, pass/.",
			},

			"folder_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the folder to create.",
			},
		},
	}
}

func resourceTencentCloudWedataWorkflowFolderCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow_folder.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		projectId string
		folderId  string
	)
	var (
		request  = wedatav20250806.NewCreateWorkflowFolderRequest()
		response = wedatav20250806.NewCreateWorkflowFolderResponse()
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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateWorkflowFolderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata workflow folder failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Data != nil && response.Response.Data.FolderId != nil {
		folderId = *response.Response.Data.FolderId
		d.SetId(projectId + tccommon.FILED_SP + folderId)
	}

	return resourceTencentCloudWedataWorkflowFolderRead(d, meta)
}

func resourceTencentCloudWedataWorkflowFolderRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow_folder.read")()
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

	var (
		respData []*wedatav20250806.WorkflowFolder
		innerErr error
	)

	parentFolderPath := d.Get("parent_folder_path").(string)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		respData, innerErr = service.DescribeWedataWorkflowFolders(ctx, projectId, folderId, parentFolderPath)
		if innerErr != nil {
			return resource.RetryableError(innerErr)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `wedata_workflow_folder` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = folderId
	return nil
}

func resourceTencentCloudWedataWorkflowFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow_folder.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"project_id", "parent_folder_path"}
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
	mutableArgs := []string{"folder_name"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wedatav20250806.NewUpdateWorkflowFolderRequest()
		request.ProjectId = helper.String(projectId)
		request.FolderId = helper.String(folderId)

		if v, ok := d.GetOk("folder_name"); ok {
			request.FolderName = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateWorkflowFolderWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update wedata workflow folder failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = folderId
	return resourceTencentCloudWedataWorkflowFolderRead(d, meta)
}

func resourceTencentCloudWedataWorkflowFolderDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow_folder.delete")()
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
		request  = wedatav20250806.NewDeleteWorkflowFolderRequest()
		response = wedatav20250806.NewDeleteWorkflowFolderResponse()
	)

	request.ProjectId = helper.String(projectId)
	request.FolderId = helper.String(folderId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteWorkflowFolderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete wedata workflow folder failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = folderId
	return nil
}
