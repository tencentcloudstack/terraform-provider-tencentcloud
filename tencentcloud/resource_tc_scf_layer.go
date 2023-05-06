/*
Provide a resource to create a SCF layer.

Example Usage

```hcl
resource "tencentcloud_scf_layer" "foo" {
  layer_name = "foo"
  compatible_runtimes = ["Python3.6"]
  content {
    cos_bucket_name = "test-bucket"
    cos_object_name = "/foo.zip"
    cos_bucket_region = "ap-guangzhou"
  }
  description = "foo"
  license_info = "foo"
}
```
Import

Scf layer can be imported, e.g.

```
$ terraform import tencentcloud_scf_layer.layer layerId#layerVersion
```

*/
package tencentcloud

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func LayerContent() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// cos code
		"cos_bucket_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cos bucket name of the SCF layer, such as `cos-1234567890`, conflict with `zip_file`.",
		},
		"cos_object_name": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringSuffix(".zip", ".jar"),
			Description:  "Cos object name of the SCF layer, should have suffix `.zip` or `.jar`, conflict with `zip_file`.",
		},
		"cos_bucket_region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cos bucket region of the SCF layer, conflict with `zip_file`.",
		},
		// zip upload
		"zip_file": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Zip file of the SCF layer, conflict with `cos_bucket_name`, `cos_object_name`, `cos_bucket_region`.",
		},
	}
}

func resourceTencentCloudScfLayer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfLayerCreate,
		Read:   resourceTencentCloudScfLayerRead,
		Update: resourceTencentCloudScfLayerUpdate,
		Delete: resourceTencentCloudScfLayerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"layer_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of layer.",
			},
			"compatible_runtimes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "The name of runtime",
				},
				Description: "The compatible runtimes of layer.",
			},
			"content": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The source code of layer.",
				Elem: &schema.Resource{
					Schema: LayerContent(),
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of layer.",
			},
			"license_info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The license info of layer.",
			},

			//compute
			"layer_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of layer.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of layer.",
			},
			"code_sha_256": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The code type of layer.",
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The download location url of layer.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of layer.",
			},
		},
	}
}

func resourceTencentCloudScfLayerCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_layer.create")()
	logId := getLogId(contextNil)
	var (
		scfService   = ScfService{client: meta.(*TencentCloudClient).apiV3Conn}
		request      = scf.NewPublishLayerVersionRequest()
		layerName    *string
		layerVersion *int64
	)

	if v, ok := d.GetOk("layer_name"); ok {
		layerName = helper.String(v.(string))
		request.LayerName = layerName
	}

	if v, ok := d.GetOk("compatible_runtimes"); ok {
		items := make([]*string, 0, 10)
		for _, item := range v.([]interface{}) {
			items = append(items, helper.String(item.(string)))
		}
		request.CompatibleRuntimes = items
	}

	if v, ok := d.GetOk("content"); ok {
		items := v.([]interface{})
		if len(items) != 1 {
			return fmt.Errorf("need only one content.")
		}
		item := items[0].(map[string]interface{})
		var content = scf.Code{}
		if item["cos_bucket_name"] != "" {
			content.CosBucketName = helper.String(item["cos_bucket_name"].(string))
		}
		if item["cos_object_name"] != "" {
			content.CosObjectName = helper.String(item["cos_object_name"].(string))
		}
		if item["cos_bucket_region"] != "" {
			content.CosBucketRegion = helper.String(item["cos_bucket_region"].(string))
		}
		if item["zip_file"] != "" {
			path, err := homedir.Expand(item["zip_file"].(string))
			if err != nil {
				return fmt.Errorf("zip file (%s) homedir expand error: %s", item["zip_file"].(string), err.Error())
			}
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("zip file (%s) open error: %s", path, err.Error())
			}
			defer file.Close()
			body, err := ioutil.ReadAll(file)
			if err != nil {
				return fmt.Errorf("zip file (%s) read error: %s", path, err.Error())
			}

			zipContent := base64.StdEncoding.EncodeToString(body)
			content.ZipFile = &zipContent
		}
		request.Content = &content
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("license_info"); ok {
		request.LicenseInfo = helper.String(v.(string))
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := scfService.client.UseScfClient().PublishLayerVersion(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err, InternalError)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		layerVersion = response.Response.LayerVersion
		return nil
	}); err != nil {
		return err
	}

	d.SetId(*layerName + FILED_SP + helper.Int64ToStr(*layerVersion))
	return resourceTencentCloudScfLayerRead(d, meta)
}

func resourceTencentCloudScfLayerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTencentCloudScfLayerRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_layer.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	//ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		scfService    = ScfService{client: meta.(*TencentCloudClient).apiV3Conn}
		layerRequest  = scf.NewGetLayerVersionRequest()
		layerResponse = scf.GetLayerVersionResponse{}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("scf layer id is borken, id is %s", d.Id())
	}
	layerName := idSplit[0]
	layerVersion := idSplit[1]

	layerRequest.LayerName = &layerName
	layerRequest.LayerVersion = helper.Int64(helper.StrToInt64(layerVersion))

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(layerRequest.GetAction())
		response, err := scfService.client.UseScfClient().GetLayerVersion(layerRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, layerRequest.GetAction(), layerRequest.ToJsonString(), err.Error())
			return retryError(err, InternalError)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, layerRequest.GetAction(), layerRequest.ToJsonString(), response.ToJsonString())
		layerResponse = *response
		return nil
	}); err != nil {
		return err
	}

	var errs []error
	errs = append(errs,
		d.Set("layer_name", layerResponse.Response.LayerName),
		d.Set("layer_version", layerResponse.Response.LayerVersion),
		d.Set("location", layerResponse.Response.Location),
		d.Set("create_time", layerResponse.Response.AddTime),
		d.Set("description", layerResponse.Response.Description),
		d.Set("license_info", layerResponse.Response.LicenseInfo),
		d.Set("status", layerResponse.Response.Status),
		d.Set("code_sha_256", layerResponse.Response.CodeSha256),
	)

	var runtimes = make([]interface{}, 0, 100)
	for _, runtime := range layerResponse.Response.CompatibleRuntimes {
		runtimes = append(runtimes, runtime)
	}

	errs = append(errs, d.Set("compatible_runtimes", runtimes))

	if len(errs) > 0 {
		return errs[0]
	} else {
		return nil
	}
}

func resourceTencentCloudScfLayerDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_layer.delete")()

	logId := getLogId(contextNil)

	var (
		scfService = ScfService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = scf.NewDeleteLayerVersionRequest()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("scf layer id is borken, id is %s", d.Id())
	}
	layerName := idSplit[0]
	layerVersion := idSplit[1]

	request.LayerName = &layerName
	request.LayerVersion = helper.Int64(helper.StrToInt64(layerVersion))

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err := scfService.client.UseScfClient().DeleteLayerVersion(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
