package oceanus

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOceanusFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusFolderCreate,
		Read:   resourceTencentCloudOceanusFolderRead,
		Update: resourceTencentCloudOceanusFolderUpdate,
		Delete: resourceTencentCloudOceanusFolderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"folder_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "New file name.",
			},
			"parent_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Parent folder id.",
			},
			"folder_type": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Default:      0,
				Description:  "Folder type, 0: job folder, 1: resource folder. Default is 0.",
			},
			"work_space_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},
		},
	}
}

func resourceTencentCloudOceanusFolderCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_folder.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = oceanus.NewCreateFolderRequest()
		response    = oceanus.NewCreateFolderResponse()
		folderId    string
		folderType  string
		workSpaceId string
	)

	if v, ok := d.GetOk("folder_name"); ok {
		request.FolderName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parent_id"); ok {
		request.ParentId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("folder_type"); ok {
		request.FolderType = helper.IntInt64(v.(int))
		folderType = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		request.WorkSpaceId = helper.String(v.(string))
		workSpaceId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOceanusClient().CreateFolder(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("oceanus Folder not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create oceanus Folder failed, reason:%+v", logId, err)
		return err
	}

	folderId = *response.Response.FolderId
	d.SetId(strings.Join([]string{workSpaceId, folderId, folderType}, tccommon.FILED_SP))

	return resourceTencentCloudOceanusFolderRead(d, meta)
}

func resourceTencentCloudOceanusFolderRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_folder.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	workSpaceId := idSplit[0]
	folderId := idSplit[1]
	folderType := idSplit[2]

	Folder, err := service.DescribeOceanusFolderById(ctx, workSpaceId, folderId, folderType)
	if err != nil {
		return err
	}

	if Folder == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OceanusFolder` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if Folder.FolderName != nil {
		_ = d.Set("folder_name", Folder.FolderName)
	}

	if Folder.ParentId != nil {
		_ = d.Set("parent_id", Folder.ParentId)
	}

	if Folder.FolderType != nil {
		_ = d.Set("folder_type", Folder.FolderType)
	}

	if Folder.WorkSpaceId != nil {
		_ = d.Set("work_space_id", Folder.WorkSpaceId)
	}

	return nil
}

func resourceTencentCloudOceanusFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_folder.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = oceanus.NewModifyFolderRequest()
	)

	immutableArgs := []string{"folder_type", "work_space_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	workSpaceId := idSplit[0]
	folderId := idSplit[1]
	folderType := idSplit[2]

	request.WorkSpaceId = &workSpaceId
	request.SourceFolderId = &folderId
	folderTypeInt, _ := strconv.ParseInt(folderType, 10, 64)
	request.FolderType = &folderTypeInt

	if v, ok := d.GetOk("parent_id"); ok {
		request.TargetFolderId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("folder_name"); ok {
		request.FolderName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOceanusClient().ModifyFolder(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update oceanus Folder failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudOceanusFolderRead(d, meta)
}

func resourceTencentCloudOceanusFolderDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_folder.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	workSpaceId := idSplit[0]
	folderId := idSplit[1]
	folderType := idSplit[2]

	if err := service.DeleteOceanusFolderById(ctx, workSpaceId, folderId, folderType); err != nil {
		return err
	}

	return nil
}
