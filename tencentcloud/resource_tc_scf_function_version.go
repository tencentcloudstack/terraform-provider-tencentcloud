/*
Provides a resource to create a scf function_version

Example Usage

```hcl
resource "tencentcloud_scf_function_version" "function_version" {
  function_name = "test_function"
  description = "test function"
  namespace = "test_namespace"
}
```

Import

scf function_version can be imported using the id, e.g.

```
terraform import tencentcloud_scf_function_version.function_version function_version_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudScfFunctionVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfFunctionVersionCreate,
		Read:   resourceTencentCloudScfFunctionVersionRead,
		Update: resourceTencentCloudScfFunctionVersionUpdate,
		Delete: resourceTencentCloudScfFunctionVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of the released function.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function description.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace.",
			},
		},
	}
}

func resourceTencentCloudScfFunctionVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_version.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = scf.NewPublishVersionRequest()
		response        = scf.NewPublishVersionResponse()
		functionName    string
		namespace       string
		functionVersion string
	)
	if v, ok := d.GetOk("function_name"); ok {
		functionName = v.(string)
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
		request.Namespace = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseScfClient().PublishVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create scf FunctionVersion failed, reason:%+v", logId, err)
		return err
	}

	functionName = *response.Response.FunctionName
	d.SetId(strings.Join([]string{functionName, namespace, functionVersion}, FILED_SP))

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Active"}, 120*readRetryTimeout, time.Second, service.ScfFunctionVersionStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudScfFunctionVersionRead(d, meta)
}

func resourceTencentCloudScfFunctionVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_version.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionName := idSplit[0]
	namespace := idSplit[1]
	functionVersion := idSplit[2]

	FunctionVersion, err := service.DescribeScfFunctionVersionById(ctx, functionName, namespace, functionVersion)
	if err != nil {
		return err
	}

	if FunctionVersion == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ScfFunctionVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if FunctionVersion.FunctionName != nil {
		_ = d.Set("function_name", FunctionVersion.FunctionName)
	}

	if FunctionVersion.Description != nil {
		_ = d.Set("description", FunctionVersion.Description)
	}

	if FunctionVersion.Namespace != nil {
		_ = d.Set("namespace", FunctionVersion.Namespace)
	}

	return nil
}

func resourceTencentCloudScfFunctionVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_version.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"function_name", "description", "namespace"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudScfFunctionVersionRead(d, meta)
}

func resourceTencentCloudScfFunctionVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_version.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionName := idSplit[0]
	namespace := idSplit[1]
	functionVersion := idSplit[2]

	if err := service.DeleteScfFunctionVersionById(ctx, functionName, namespace, functionVersion); err != nil {
		return err
	}

	return nil
}
