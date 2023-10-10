/*
Provides a resource to create a rum release_file

Example Usage

```hcl
resource "tencentcloud_rum_release_file" "release_file" {
  project_id      = 123
  version         = "1.0"
  file_key        = "120000-last-1632921299138-index.js.map"
  file_name       = "index.js.map"
  file_hash 	  = "b148c43fd81d845ba7cc6907928ce430"
  release_file_id = 1
}
```

Import

rum release_file can be imported using the id, e.g.

```
terraform import tencentcloud_rum_release_file.release_file projectId#releaseFileId
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRumReleaseFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRumReleaseFileCreate,
		Read:   resourceTencentCloudRumReleaseFileRead,
		Delete: resourceTencentCloudRumReleaseFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Release File version.",
			},
			"file_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Release file unique key.",
			},
			"file_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Release file name.",
			},
			"file_hash": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Release file hash.",
			},
			"release_file_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Release file id.",
			},
		},
	}
}

func resourceTencentCloudRumReleaseFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_release_file.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = rum.NewCreateReleaseFileRequest()
		projectID     int
		releaseFileId int
	)
	if v, ok := d.GetOkExists("project_id"); ok {
		projectID = v.(int)
		request.ProjectID = helper.IntInt64(v.(int))
	}

	releaseFile := rum.ReleaseFile{}
	if v, ok := d.GetOkExists("version"); ok {
		releaseFile.Version = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("file_key"); ok {
		releaseFile.FileKey = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("file_name"); ok {
		releaseFile.FileName = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("file_hash"); ok {
		releaseFile.FileHash = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("release_file_id"); ok {
		releaseFileId = v.(int)
		releaseFile.ID = helper.IntInt64(v.(int))
	}
	request.Files = append(request.Files, &releaseFile)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().CreateReleaseFile(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create rum releaseFile failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strconv.Itoa(projectID) + FILED_SP + strconv.Itoa(releaseFileId))

	return resourceTencentCloudRumReleaseFileRead(d, meta)
}

func resourceTencentCloudRumReleaseFileRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_release_file.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectID, err := strconv.ParseInt(idSplit[0], 10, 64)
	if err != nil {
		return fmt.Errorf("id data format error,%s", d.Id())
	}
	releaseFileId, err := strconv.ParseInt(idSplit[1], 10, 64)
	if err != nil {
		return fmt.Errorf("id data format error,%s", d.Id())
	}

	releaseFile, err := service.DescribeRumReleaseFileById(ctx, projectID, releaseFileId)
	if err != nil {
		return err
	}

	if releaseFile == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RumReleaseFile` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectID)

	if releaseFile.Version != nil {
		_ = d.Set("version", releaseFile.Version)
	}

	if releaseFile.FileKey != nil {
		_ = d.Set("file_key", releaseFile.FileKey)
	}

	if releaseFile.FileName != nil {
		_ = d.Set("file_name", releaseFile.FileName)
	}

	if releaseFile.FileHash != nil {
		_ = d.Set("file_hash", releaseFile.FileHash)
	}

	if releaseFile.ID != nil {
		_ = d.Set("release_file_id", releaseFile.ID)
	}

	return nil
}

func resourceTencentCloudRumReleaseFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_release_file.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	// projectID, _ := strconv.ParseInt(idSplit[0], 10, 64)
	releaseFileId, _ := strconv.ParseInt(idSplit[1], 10, 64)

	if err := service.DeleteRumReleaseFileById(ctx, releaseFileId); err != nil {
		return err
	}

	return nil
}
