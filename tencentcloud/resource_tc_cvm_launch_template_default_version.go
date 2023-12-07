package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				ForceNew:    true,
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

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	launchTemplateId := d.Get("launch_template_id").(string)
	defaultVersion := d.Get("default_version").(int)
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := service.ModifyLaunchTemplateDefaultVersion(ctx, launchTemplateId, defaultVersion)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(launchTemplateId)

	return resourceTencentCloudCvmLaunchTemplateDefaultVersionRead(d, meta)
}

func resourceTencentCloudCvmLaunchTemplateDefaultVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_default_version.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	launchTemplateId := d.Id()
	launchTemplateVersions, err := service.DescribeLaunchTemplateVersions(ctx, launchTemplateId)
	if err != nil {
		return err
	}

	if len(launchTemplateVersions) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmLaunchTemplateDefaultVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	for _, launchTemplateVersion := range launchTemplateVersions {
		if *launchTemplateVersion.IsDefaultVersion {
			_ = d.Set("default_version", launchTemplateVersion.LaunchTemplateVersion)
			break
		}
	}
	_ = d.Set("launch_template_id", d.Id())

	return nil
}

func resourceTencentCloudCvmLaunchTemplateDefaultVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_default_version.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	if d.HasChange("default_version") {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := service.ModifyLaunchTemplateDefaultVersion(ctx, d.Id(), d.Get("default_version").(int))
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudCvmLaunchTemplateDefaultVersionRead(d, meta)
}

func resourceTencentCloudCvmLaunchTemplateDefaultVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_default_version.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
