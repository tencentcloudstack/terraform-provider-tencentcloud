/*
Provides a resource to create a cvm launch_template_default_version

Example Usage

```hcl
resource "tencentcloud_cvm_launch_template_default_version" "launch_template_default_version" {
  launch_template_id = "lt-34vaef8fe"
  default_version = 2
}
```

Import

cvm launch_template_default_version can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_launch_template_default_version.launch_template_default_version launch_template_default_version_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"log"
)

func resourceTencentCloudCvmLaunchTemplateDefaultVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmLaunchTemplateDefaultVersionCreate,
		Read:   resourceTencentCloudCvmLaunchTemplateDefaultVersionRead,
		Update: resourceTencentCloudCvmLaunchTemplateDefaultVersionUpdate,
		Delete: resourceTencentCloudCvmLaunchTemplateDefaultVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"launch_template_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance launch template ID.",
			},

			"default_version": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The number of the version that you want to set as the default version.",
			},
		},
	}
}

func resourceTencentCloudCvmLaunchTemplateDefaultVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_default_version.create")()
	defer inconsistentCheck(d, meta)()

	var launchTemplateId string
	if v, ok := d.GetOk("launch_template_id"); ok {
		launchTemplateId = v.(string)
	}

	d.SetId(launchTemplateId)

	return resourceTencentCloudCvmLaunchTemplateDefaultVersionUpdate(d, meta)
}

func resourceTencentCloudCvmLaunchTemplateDefaultVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_default_version.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	launchTemplateDefaultVersionId := d.Id()

	launchTemplateDefaultVersion, err := service.DescribeCvmLaunchTemplateDefaultVersionById(ctx, launchTemplateId)
	if err != nil {
		return err
	}

	if launchTemplateDefaultVersion == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmLaunchTemplateDefaultVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if launchTemplateDefaultVersion.LaunchTemplateId != nil {
		_ = d.Set("launch_template_id", launchTemplateDefaultVersion.LaunchTemplateId)
	}

	if launchTemplateDefaultVersion.DefaultVersion != nil {
		_ = d.Set("default_version", launchTemplateDefaultVersion.DefaultVersion)
	}

	return nil
}

func resourceTencentCloudCvmLaunchTemplateDefaultVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_default_version.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cvm.NewModifyLaunchTemplateDefaultVersionRequest()

	launchTemplateDefaultVersionId := d.Id()

	request.LaunchTemplateId = &launchTemplateId

	immutableArgs := []string{"launch_template_id", "default_version"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyLaunchTemplateDefaultVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cvm launchTemplateDefaultVersion failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCvmLaunchTemplateDefaultVersionRead(d, meta)
}

func resourceTencentCloudCvmLaunchTemplateDefaultVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_default_version.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
