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

func ResourceTencentCloudWedataSqlFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataSqlFolderCreate,
		Read:   resourceTencentCloudWedataSqlFolderRead,
		Update: resourceTencentCloudWedataSqlFolderUpdate,
		Delete: resourceTencentCloudWedataSqlFolderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"folder_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Folder name.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID.",
			},

			"parent_folder_path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The parent folder path is /aaa/bbb/ccc. The path header must have a slash. To query the root directory, pass /.",
			},

			"access_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Permission range: SHARED, PRIVATE.",
			},

			// computed
			"folder_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Folder ID.",
			},

			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node path.",
			},
		},
	}
}

func resourceTencentCloudWedataSqlFolderCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_sql_folder.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = wedatav20250806.NewCreateSQLFolderRequest()
		response  = wedatav20250806.NewCreateSQLFolderResponse()
		projectId string
		folderId  string
	)

	if v, ok := d.GetOk("folder_name"); ok {
		request.FolderName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("parent_folder_path"); ok {
		request.ParentFolderPath = helper.String(v.(string))
	}

	if v, ok := d.GetOk("access_scope"); ok {
		request.AccessScope = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateSQLFolderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata sql folder failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata sql folder failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.FolderId == nil {
		return fmt.Errorf("FolderId is nil.")
	}

	folderId = *response.Response.Data.FolderId
	d.SetId(strings.Join([]string{projectId, folderId}, tccommon.FILED_SP))
	return resourceTencentCloudWedataSqlFolderRead(d, meta)
}

func resourceTencentCloudWedataSqlFolderRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_sql_folder.read")()
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

	respData, err := service.DescribeWedataGetSqlFolderById(ctx, projectId, folderId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_sql_folder` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("project_id", projectId)

	if respData.Name != nil {
		_ = d.Set("folder_name", *respData.Name)
	}

	if respData.ParentFolderPath != nil {
		if *respData.ParentFolderPath == "" {
			_ = d.Set("parent_folder_path", "/")
		} else {
			_ = d.Set("parent_folder_path", *respData.ParentFolderPath)
		}
	}

	if respData.AccessScope != nil {
		_ = d.Set("access_scope", *respData.AccessScope)
	}

	if respData.Id != nil {
		_ = d.Set("folder_id", *respData.Id)
	}

	if respData.Path != nil {
		_ = d.Set("path", *respData.Path)
	}

	return nil
}

func resourceTencentCloudWedataSqlFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_sql_folder.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	folderId := idSplit[1]

	needChange := false
	mutableArgs := []string{"folder_name", "access_scope"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wedatav20250806.NewUpdateSQLFolderRequest()
		if v, ok := d.GetOk("folder_name"); ok {
			request.FolderName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("access_scope"); ok {
			request.AccessScope = helper.String(v.(string))
		}

		request.ProjectId = helper.String(projectId)
		request.FolderId = helper.String(folderId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateSQLFolderWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("Update wedata sql folder failed, Response is nil."))
			}

			if !*result.Response.Data.Status {
				return resource.NonRetryableError(fmt.Errorf("Update wedata sql folder failed, Status is false."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update wedata sql folder failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudWedataSqlFolderRead(d, meta)
}

func resourceTencentCloudWedataSqlFolderDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_sql_folder.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wedatav20250806.NewDeleteSQLFolderRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	folderId := idSplit[1]

	request.ProjectId = helper.String(projectId)
	request.FolderId = helper.String(folderId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteSQLFolderWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata sql folder failed, Response is nil."))
		}

		if !*result.Response.Data.Status {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata sql folder failed, Status is false."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata sql folder failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
