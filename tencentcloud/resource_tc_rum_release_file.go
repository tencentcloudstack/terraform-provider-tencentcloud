/*
Provides a resource to create a rum release_file

Example Usage

```hcl
resource "tencentcloud_rum_release_file" "release_file" {
  project_i_d = 123
  files {
		version = "1.0"
		file_key = "120000-last-1632921299138-index.js.map"
		file_name = "index.js.map"
		file_hash = "b148c43fd81d845ba7cc6907928ce430"
		i_d = 1

  }
}
```

Import

rum release_file can be imported using the id, e.g.

```
terraform import tencentcloud_rum_release_file.release_file release_file_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudRumReleaseFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRumReleaseFileCreate,
		Read:   resourceTencentCloudRumReleaseFileRead,
		Update: resourceTencentCloudRumReleaseFileUpdate,
		Delete: resourceTencentCloudRumReleaseFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_i_d": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"files": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "File list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Release File version.",
						},
						"file_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Release file unique key.",
						},
						"file_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Release file name.",
						},
						"file_hash": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Release file hash.",
						},
						"i_d": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Release file id.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudRumReleaseFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_release_file.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = rum.NewCreateReleaseFileRequest()
		response  = rum.NewCreateReleaseFileResponse()
		projectID int
	)
	if v, ok := d.GetOkExists("project_i_d"); ok {
		projectID = v.(int64)
		request.ProjectID = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("files"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			releaseFile := rum.ReleaseFile{}
			if v, ok := dMap["version"]; ok {
				releaseFile.Version = helper.String(v.(string))
			}
			if v, ok := dMap["file_key"]; ok {
				releaseFile.FileKey = helper.String(v.(string))
			}
			if v, ok := dMap["file_name"]; ok {
				releaseFile.FileName = helper.String(v.(string))
			}
			if v, ok := dMap["file_hash"]; ok {
				releaseFile.FileHash = helper.String(v.(string))
			}
			if v, ok := dMap["i_d"]; ok {
				releaseFile.ID = helper.IntInt64(v.(int))
			}
			request.Files = append(request.Files, &releaseFile)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().CreateReleaseFile(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create rum releaseFile failed, reason:%+v", logId, err)
		return err
	}

	projectID = *response.Response.ProjectID
	d.SetId(helper.Int64ToStr(projectID))

	return resourceTencentCloudRumReleaseFileRead(d, meta)
}

func resourceTencentCloudRumReleaseFileRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_release_file.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	releaseFileId := d.Id()

	releaseFile, err := service.DescribeRumReleaseFileById(ctx, projectID)
	if err != nil {
		return err
	}

	if releaseFile == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RumReleaseFile` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if releaseFile.ProjectID != nil {
		_ = d.Set("project_i_d", releaseFile.ProjectID)
	}

	if releaseFile.Files != nil {
		filesList := []interface{}{}
		for _, files := range releaseFile.Files {
			filesMap := map[string]interface{}{}

			if releaseFile.Files.Version != nil {
				filesMap["version"] = releaseFile.Files.Version
			}

			if releaseFile.Files.FileKey != nil {
				filesMap["file_key"] = releaseFile.Files.FileKey
			}

			if releaseFile.Files.FileName != nil {
				filesMap["file_name"] = releaseFile.Files.FileName
			}

			if releaseFile.Files.FileHash != nil {
				filesMap["file_hash"] = releaseFile.Files.FileHash
			}

			if releaseFile.Files.ID != nil {
				filesMap["i_d"] = releaseFile.Files.ID
			}

			filesList = append(filesList, filesMap)
		}

		_ = d.Set("files", filesList)

	}

	return nil
}

func resourceTencentCloudRumReleaseFileUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_release_file.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"project_i_d", "files"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudRumReleaseFileRead(d, meta)
}

func resourceTencentCloudRumReleaseFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_release_file.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}
	releaseFileId := d.Id()

	if err := service.DeleteRumReleaseFileById(ctx, projectID); err != nil {
		return err
	}

	return nil
}
