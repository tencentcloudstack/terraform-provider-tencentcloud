package pts

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPtsFile() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudPtsFileRead,
		Create: resourceTencentCloudPtsFileCreate,
		Update: resourceTencentCloudPtsFileUpdate,
		Delete: resourceTencentCloudPtsFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"file_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "File id.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project id.",
			},

			"kind": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "File kind, parameter file-1, protocol file-2, request file-3.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "File name.",
			},

			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "File size.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "File type, folder-folder.",
			},

			"line_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Line count.",
			},

			"head_lines": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The first few lines of data.",
			},

			"tail_lines": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The last few lines of data.",
			},

			"header_in_file": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the header is in the file.",
			},

			"header_columns": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Meter head.",
			},

			"file_infos": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Files in a folder.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "File size.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File type.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time.",
						},
						"file_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File id.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPtsFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_file.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = pts.NewCreateFileRequest()
		projectId string
		fileId    string
	)

	if v, ok := d.GetOk("file_id"); ok {
		fileId = v.(string)
		request.FileId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kind"); ok {
		request.Kind = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("size"); ok {
		request.Size = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("line_count"); ok {
		request.LineCount = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("head_lines"); ok {
		headLinesSet := v.(*schema.Set).List()
		for i := range headLinesSet {
			headLines := headLinesSet[i].(string)
			request.HeadLines = append(request.HeadLines, &headLines)
		}
	}

	if v, ok := d.GetOk("tail_lines"); ok {
		tailLinesSet := v.(*schema.Set).List()
		for i := range tailLinesSet {
			tailLines := tailLinesSet[i].(string)
			request.TailLines = append(request.TailLines, &tailLines)
		}
	}

	if v, _ := d.GetOk("header_in_file"); v != nil {
		request.HeaderInFile = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("header_columns"); ok {
		headerColumnsSet := v.(*schema.Set).List()
		for i := range headerColumnsSet {
			headerColumns := headerColumnsSet[i].(string)
			request.HeaderColumns = append(request.HeaderColumns, &headerColumns)
		}
	}

	if v, ok := d.GetOk("file_infos"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			fileInfo := pts.FileInfo{}
			if v, ok := dMap["name"]; ok {
				fileInfo.Name = helper.String(v.(string))
			}
			if v, ok := dMap["size"]; ok {
				fileInfo.Size = helper.Int64(int64(v.(int)))
			}
			if v, ok := dMap["type"]; ok {
				fileInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["updated_at"]; ok {
				fileInfo.UpdatedAt = helper.String(v.(string))
			}
			if v, ok := dMap["file_id"]; ok {
				fileInfo.FileId = helper.String(v.(string))
			}

			request.FileInfos = append(request.FileInfos, &fileInfo)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePtsClient().CreateFile(request)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation.DbRecordCreateFailed" {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create pts file failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(projectId + tccommon.FILED_SP + fileId)
	return resourceTencentCloudPtsFileRead(d, meta)
}

func resourceTencentCloudPtsFileRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_file.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	fileId := idSplit[1]

	file, err := service.DescribePtsFile(ctx, projectId, fileId)

	if err != nil {
		return err
	}

	if file == nil {
		d.SetId("")
		return fmt.Errorf("resource `file` %s does not exist", fileId)
	}

	if file.FileId != nil {
		_ = d.Set("file_id", file.FileId)
	}

	if file.ProjectId != nil {
		_ = d.Set("project_id", file.ProjectId)
	}

	if file.Kind != nil {
		_ = d.Set("kind", file.Kind)
	}

	if file.Name != nil {
		_ = d.Set("name", file.Name)
	}

	if file.Size != nil {
		_ = d.Set("size", file.Size)
	}

	if file.Type != nil {
		_ = d.Set("type", file.Type)
	}

	if file.LineCount != nil {
		_ = d.Set("line_count", file.LineCount)
	}

	if file.HeadLines != nil {
		_ = d.Set("head_lines", file.HeadLines)
	}

	if file.TailLines != nil {
		_ = d.Set("tail_lines", file.TailLines)
	}

	if file.HeaderInFile != nil {
		_ = d.Set("header_in_file", file.HeaderInFile)
	}

	if file.HeaderColumns != nil {
		_ = d.Set("header_columns", file.HeaderColumns)
	}

	if file.FileInfos != nil {
		fileInfosList := []interface{}{}
		for _, fileInfos := range file.FileInfos {
			fileInfosMap := map[string]interface{}{}
			if fileInfos.Name != nil {
				fileInfosMap["name"] = fileInfos.Name
			}
			if fileInfos.Size != nil {
				fileInfosMap["size"] = fileInfos.Size
			}
			if fileInfos.Type != nil {
				fileInfosMap["type"] = fileInfos.Type
			}
			if fileInfos.UpdatedAt != nil {
				fileInfosMap["updated_at"] = fileInfos.UpdatedAt
			}
			if fileInfos.FileId != nil {
				fileInfosMap["file_id"] = fileInfos.FileId
			}

			fileInfosList = append(fileInfosList, fileInfosMap)
		}
		_ = d.Set("file_infos", fileInfosList)
	}

	return nil
}

func resourceTencentCloudPtsFileUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_file.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	return resourceTencentCloudPtsFileRead(d, meta)
}

func resourceTencentCloudPtsFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_file.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	fileId := idSplit[1]

	if err := service.DeletePtsFileById(ctx, projectId, fileId); err != nil {
		return err
	}

	return nil
}
