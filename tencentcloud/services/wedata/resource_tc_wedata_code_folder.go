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

func ResourceTencentCloudWedataCodeFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataCodeFolderCreate,
		Read:   resourceTencentCloudWedataCodeFolderRead,
		Update: resourceTencentCloudWedataCodeFolderUpdate,
		Delete: resourceTencentCloudWedataCodeFolderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID.",
			},

			"folder_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Folder name.",
			},

			"parent_folder_path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Parent folder path, for example /aaa/bbb/ccc, path header must start with a slash, root directory pass /.",
			},

			// computed
			"folder_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Folder ID.",
			},

			"access_scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Permission range: SHARED, PRIVATE.",
			},

			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type. folder, script.",
			},

			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node path.",
			},
		},
	}
}

func resourceTencentCloudWedataCodeFolderCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_folder.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = wedatav20250806.NewCreateCodeFolderRequest()
		response  = wedatav20250806.NewCreateCodeFolderResponse()
		projectId string
		folderId  string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("folder_name"); ok {
		request.FolderName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parent_folder_path"); ok {
		request.ParentFolderPath = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateCodeFolderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("create wedata code folder failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata code folder failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.FolderId == nil {
		return fmt.Errorf("FolderId is nil.")
	}

	folderId = *response.Response.Data.FolderId
	d.SetId(strings.Join([]string{projectId, folderId}, tccommon.FILED_SP))
	return resourceTencentCloudWedataCodeFolderRead(d, meta)
}

func resourceTencentCloudWedataCodeFolderRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_folder.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	folderId := idSplit[1]

	respData, err := service.DescribeWedataGetCodeFolderById(ctx, projectId, folderId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_code_folder` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("project_id", projectId)

	if respData.Title != nil {
		_ = d.Set("folder_name", *respData.Title)
	}

	if respData.ParentFolderPath != nil {
		if *respData.ParentFolderPath == "" {
			_ = d.Set("parent_folder_path", "/")
		} else {
			_ = d.Set("parent_folder_path", *respData.ParentFolderPath)
		}
	}

	if respData.Id != nil {
		_ = d.Set("folder_id", *respData.Id)
	}

	if respData.AccessScope != nil {
		_ = d.Set("access_scope", *respData.AccessScope)
	}

	if respData.Type != nil {
		_ = d.Set("type", *respData.Type)
	}

	if respData.Path != nil {
		_ = d.Set("path", *respData.Path)
	}

	return nil
}

func resourceTencentCloudWedataCodeFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_folder.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wedatav20250806.NewUpdateCodeFolderRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	folderId := idSplit[1]

	if d.HasChange("folder_name") {
		if v, ok := d.GetOk("folder_name"); ok {
			request.FolderName = helper.String(v.(string))
		}

		request.ProjectId = &projectId
		request.FolderId = &folderId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateCodeFolderWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("Update wedata code folder failed, Response is nil."))
			}

			if !*result.Response.Data.Status {
				return resource.NonRetryableError(fmt.Errorf("Update wedata code folder failed, Status is false."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update wedata code folder failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudWedataCodeFolderRead(d, meta)
}

func resourceTencentCloudWedataCodeFolderDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_folder.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wedatav20250806.NewDeleteCodeFolderRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	folderId := idSplit[1]

	request.ProjectId = &projectId
	request.FolderId = &folderId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteCodeFolderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata code folder failed, Response is nil."))
		}

		if !*result.Response.Data.Status {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata code folder failed, Status is false."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata code folder failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
