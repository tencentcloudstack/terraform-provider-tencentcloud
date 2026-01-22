package igtm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIgtmInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIgtmInstanceCreate,
		Read:   resourceTencentCloudIgtmInstanceRead,
		Update: resourceTencentCloudIgtmInstanceUpdate,
		Delete: resourceTencentCloudIgtmInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Business domain.",
			},

			"access_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CUSTOM: Custom access domain\nSYSTEM: System access domain.",
			},

			"global_ttl": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Resolution effective time.",
			},

			"package_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Package type\nFREE: Free version\nSTANDARD: Standard version\nULTIMATE: Ultimate version.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name.",
			},

			"access_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access main domain.",
			},

			"access_sub_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access subdomain.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark.",
			},

			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Package resource ID.",
			},

			// computed
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudIgtmInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = igtmv20231024.NewCreateInstanceRequest()
		response   = igtmv20231024.NewCreateInstanceResponse()
		instanceId string
	)

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("access_type"); ok {
		request.AccessType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("global_ttl"); ok {
		request.GlobalTtl = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("package_type"); ok {
		request.PackageType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("access_domain"); ok {
		request.AccessDomain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("access_sub_domain"); ok {
		request.AccessSubDomain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_id"); ok {
		request.ResourceId = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().CreateInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create igtm instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create igtm instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.InstanceId == nil {
		return fmt.Errorf("InstanceId is nil.")
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)
	return resourceTencentCloudIgtmInstanceRead(d, meta)
}

func resourceTencentCloudIgtmInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	respData, err := service.DescribeIgtmInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_igtm_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Domain != nil {
		_ = d.Set("domain", respData.Domain)
	}

	if respData.AccessType != nil {
		_ = d.Set("access_type", respData.AccessType)
	}

	if respData.GlobalTtl != nil {
		_ = d.Set("global_ttl", respData.GlobalTtl)
	}

	if respData.PackageType != nil {
		_ = d.Set("package_type", respData.PackageType)
	}

	if respData.InstanceName != nil {
		_ = d.Set("instance_name", respData.InstanceName)
	}

	if respData.AccessDomain != nil {
		_ = d.Set("access_domain", respData.AccessDomain)
	}

	if respData.AccessSubDomain != nil {
		_ = d.Set("access_sub_domain", respData.AccessSubDomain)
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	if respData.ResourceId != nil {
		_ = d.Set("resource_id", respData.ResourceId)
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	return nil
}

func resourceTencentCloudIgtmInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"domain", "access_type", "global_ttl", "instance_name", "access_domain", "access_sub_domain", "remark"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := igtmv20231024.NewModifyInstanceConfigRequest()
		instanceConfig := igtmv20231024.InstanceConfig{}
		if v, ok := d.GetOk("domain"); ok && v.(string) != "" {
			instanceConfig.Domain = helper.String(v.(string))
		}

		if v, ok := d.GetOk("access_type"); ok && v.(string) != "" {
			instanceConfig.AccessType = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("global_ttl"); ok {
			instanceConfig.GlobalTtl = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("instance_name"); ok && v.(string) != "" {
			instanceConfig.InstanceName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("access_domain"); ok && v.(string) != "" {
			instanceConfig.AccessDomain = helper.String(v.(string))
		}

		if v, ok := d.GetOk("access_sub_domain"); ok && v.(string) != "" {
			instanceConfig.AccessSubDomain = helper.String(v.(string))
		}

		if v, ok := d.GetOk("remark"); ok && v.(string) != "" {
			instanceConfig.Remark = helper.String(v.(string))
		}

		request.InstanceConfig = &instanceConfig
		request.InstanceId = &instanceId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().ModifyInstanceConfigWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify igtm instance config failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudIgtmInstanceRead(d, meta)
}

func resourceTencentCloudIgtmInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return fmt.Errorf("tencentcloud igtm instance not supported delete, please contact the work order for processing")
}
