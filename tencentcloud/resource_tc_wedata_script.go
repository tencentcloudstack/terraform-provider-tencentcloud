/*
Provides a resource to create a wedata script

Example Usage

```hcl
resource "tencentcloud_wedata_script" "example" {
  file_path           = "/datastudio/project/tf_example.sql"
  project_id          = "1470575647377821696"
  bucket_name         = "wedata-demo-1257305158"
  region              = "ap-guangzhou"
  file_extension_type = "sql"
}
```

Import

wedata script can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_script.example 1470575647377821696#/datastudio/project/tf_example.sql#4147824b-7ba2-432b-8a8b-7e747594c926
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWedataScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataScriptCreate,
		Read:   resourceTencentCloudWedataScriptRead,
		Update: resourceTencentCloudWedataScriptUpdate,
		Delete: resourceTencentCloudWedataScriptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cos file path:/datastudio/project/projectId/.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Project id.",
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cos bucket name.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cos region.",
			},
			"file_extension_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File Extension Type:jar, sql, zip, py, sh, txt, di, dg, pyspark, kjb, ktr, csv.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID.",
			},
		},
	}
}

func resourceTencentCloudWedataScriptCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_script.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = wedata.NewUploadContentRequest()
		response   = wedata.NewUploadContentResponse()
		projectId  string
		resourceId string
		filePath   string
		bucketName string
		region     string
	)

	scriptRequestInfo := wedata.ScriptRequestInfo{}

	if v, ok := d.GetOk("file_path"); ok {
		scriptRequestInfo.FilePath = helper.String(v.(string))
		filePath = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		scriptRequestInfo.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("bucket_name"); ok {
		scriptRequestInfo.BucketName = helper.String(v.(string))
		bucketName = v.(string)
	}

	if v, ok := d.GetOk("region"); ok {
		scriptRequestInfo.Region = helper.String(v.(string))
		region = v.(string)
	}

	if v, ok := d.GetOk("file_extension_type"); ok {
		scriptRequestInfo.FileExtensionType = helper.String(v.(string))
	}

	scriptRequestInfo.Operation = helper.String("create")
	tmpStr := fmt.Sprintf("%s|%s|%s", region, bucketName, filePath)
	ExtraInfoObj := map[string]string{
		"taskId": tmpStr,
	}
	extraInfoBytes, _ := json.Marshal(ExtraInfoObj)
	extraInfoStr := string(extraInfoBytes)
	scriptRequestInfo.ExtraInfo = helper.String(extraInfoStr)

	request.ScriptRequestInfo = &scriptRequestInfo
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().UploadContent(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata script failed, reason:%+v", logId, err)
		return err
	}

	resourceId = *response.Response.ScriptInfo.ResourceId
	d.SetId(strings.Join([]string{projectId, filePath, resourceId}, FILED_SP))

	return resourceTencentCloudWedataScriptRead(d, meta)
}

func resourceTencentCloudWedataScriptRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_script.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		service  = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
		fileInfo *wedata.UserFileInfo
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	filePath := idSplit[1]
	resourceId := idSplit[2]

	fileInfo, err := service.DescribeWedataScriptById(ctx, projectId, filePath)
	if err != nil {
		return err
	}

	if fileInfo == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataScript` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)
	_ = d.Set("file_path", filePath)
	_ = d.Set("resource_id", resourceId)

	if fileInfo.Bucket != nil {
		_ = d.Set("bucket_name", fileInfo.Bucket)
	}

	if fileInfo.Region != nil {
		_ = d.Set("region", fileInfo.Region)
	}

	if fileInfo.FileExtensionType != nil {
		_ = d.Set("file_extension_type", fileInfo.FileExtensionType)
	}

	return nil
}

func resourceTencentCloudWedataScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_script.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = wedata.NewUploadContentRequest()
		bucketName string
		region     string
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	filePath := idSplit[1]

	immutableArgs := []string{"file_path, project_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	scriptRequestInfo := wedata.ScriptRequestInfo{}
	scriptRequestInfo.ProjectId = &projectId
	scriptRequestInfo.FilePath = &filePath

	if v, ok := d.GetOk("bucket_name"); ok {
		scriptRequestInfo.BucketName = helper.String(v.(string))
		bucketName = v.(string)
	}

	if v, ok := d.GetOk("region"); ok {
		scriptRequestInfo.Region = helper.String(v.(string))
		region = v.(string)
	}

	if v, ok := d.GetOk("file_extension_type"); ok {
		scriptRequestInfo.FileExtensionType = helper.String(v.(string))
	}

	scriptRequestInfo.Operation = helper.String("create")
	tmpStr := fmt.Sprintf("%s|%s|%s", region, bucketName, filePath)
	ExtraInfoObj := map[string]string{
		"taskId": tmpStr,
	}
	extraInfoBytes, _ := json.Marshal(ExtraInfoObj)
	extraInfoStr := string(extraInfoBytes)
	scriptRequestInfo.ExtraInfo = helper.String(extraInfoStr)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().UploadContent(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update wedata script failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataScriptRead(d, meta)
}

func resourceTencentCloudWedataScriptDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_script.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	resourceId := idSplit[2]

	if err := service.DeleteWedataScriptById(ctx, projectId, resourceId); err != nil {
		return err
	}

	return nil
}
