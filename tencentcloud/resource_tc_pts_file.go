/*
Provides a resource to create a pts file

Example Usage

```hcl
resource "tencentcloud_pts_file" "file" {
  file_id = &lt;nil&gt;
  project_id = &lt;nil&gt;
  kind = &lt;nil&gt;
  name = &lt;nil&gt;
  size = &lt;nil&gt;
  type = &lt;nil&gt;
  line_count = &lt;nil&gt;
  head_lines = &lt;nil&gt;
  tail_lines = &lt;nil&gt;
  header_in_file = &lt;nil&gt;
  header_columns = &lt;nil&gt;
  file_infos {
		name = &lt;nil&gt;
		size = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		file_id = &lt;nil&gt;

  }
}
```

Import

pts file can be imported using the id, e.g.

```
terraform import tencentcloud_pts_file.file file_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudPtsFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsFileCreate,
		Read:   resourceTencentCloudPtsFileRead,
		Update: resourceTencentCloudPtsFileUpdate,
		Delete: resourceTencentCloudPtsFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"file_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File id.",
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project id.",
			},

			"kind": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "File kind, parameter file-1, protocol file-2, request file-3.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File name.",
			},

			"size": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "File size.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File type, folder-folder.",
			},

			"line_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Line count.",
			},

			"head_lines": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The first few lines of data.",
			},

			"tail_lines": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The last few lines of data.",
			},

			"header_in_file": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether the header is in the file.",
			},

			"header_columns": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Meter head.",
			},

			"file_infos": {
				Optional:    true,
				Type:        schema.TypeList,
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
	defer logElapsed("resource.tencentcloud_pts_file.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = pts.NewCreateFileRequest()
		response  = pts.NewCreateFileResponse()
		projectId string
	)
	if v, ok := d.GetOk("file_id"); ok {
		request.FileId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("kind"); ok {
		request.Kind = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("size"); ok {
		request.Size = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("line_count"); ok {
		request.LineCount = helper.IntUint64(v.(int))
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

	if v, ok := d.GetOkExists("header_in_file"); ok {
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
				fileInfo.Size = helper.IntUint64(v.(int))
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().CreateFile(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create pts file failed, reason:%+v", logId, err)
		return err
	}

	projectId = *response.Response.ProjectId
	d.SetId(strings.Join([]string{projectId}, FILED_SP))

	return resourceTencentCloudPtsFileRead(d, meta)
}

func resourceTencentCloudPtsFileRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_file.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	fileIds := idSplit[1]

	file, err := service.DescribePtsFileById(ctx, projectId, fileIds)
	if err != nil {
		return err
	}

	if file == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PtsFile` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
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

			if file.FileInfos.Name != nil {
				fileInfosMap["name"] = file.FileInfos.Name
			}

			if file.FileInfos.Size != nil {
				fileInfosMap["size"] = file.FileInfos.Size
			}

			if file.FileInfos.Type != nil {
				fileInfosMap["type"] = file.FileInfos.Type
			}

			if file.FileInfos.UpdatedAt != nil {
				fileInfosMap["updated_at"] = file.FileInfos.UpdatedAt
			}

			if file.FileInfos.FileId != nil {
				fileInfosMap["file_id"] = file.FileInfos.FileId
			}

			fileInfosList = append(fileInfosList, fileInfosMap)
		}

		_ = d.Set("file_infos", fileInfosList)

	}

	return nil
}

func resourceTencentCloudPtsFileUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_file.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := pts.NewRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	fileIds := idSplit[1]

	request.ProjectId = &projectId
	request.FileIds = &fileIds

	immutableArgs := []string{"file_id", "project_id", "kind", "name", "size", "type", "line_count", "head_lines", "tail_lines", "header_in_file", "header_columns", "file_infos"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update pts file failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPtsFileRead(d, meta)
}

func resourceTencentCloudPtsFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_file.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	fileIds := idSplit[1]

	if err := service.DeletePtsFileById(ctx, projectId, fileIds); err != nil {
		return err
	}

	return nil
}
