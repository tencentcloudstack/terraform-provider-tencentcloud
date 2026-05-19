package bh

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhBindDeviceResourceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhBindDeviceResourceConfigCreate,
		Read:   resourceTencentCloudBhBindDeviceResourceConfigRead,
		Update: resourceTencentCloudBhBindDeviceResourceConfigUpdate,
		Delete: resourceTencentCloudBhBindDeviceResourceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bastion host service ID.",
			},

			"device_id_set": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Asset ID collection.",
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Network Domain ID.",
			},

			"manage_dimension": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "K8S cluster managed account dimension. Valid values: `1`-Cluster, `2`-Namespace, `3`-Workload.",
			},

			"manage_account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "K8S cluster managed account ID.",
			},

			"manage_account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "K8S cluster managed account name.",
			},

			"manage_kubeconfig": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "K8S cluster managed account kubeconfig credential.",
			},

			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "K8S cluster managed namespace.",
			},

			"workload": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "K8S cluster managed workload.",
			},

			// computed
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network Domain name.",
			},
		},
	}
}

func resourceTencentCloudBhBindDeviceResourceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	resourceId := d.Get("resource_id").(string)
	d.SetId(resourceId)
	return resourceTencentCloudBhBindDeviceResourceConfigUpdate(d, meta)
}

func resourceTencentCloudBhBindDeviceResourceConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		resourceId = d.Id()
	)

	deviceSets, err := service.DescribeBhDeviceByResourceId(ctx, resourceId)
	if err != nil {
		return err
	}

	if len(deviceSets) == 0 {
		log.Printf("[WARN]%s resource `tencentcloud_bh_bind_device_resource_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("resource_id", resourceId)

	deviceIdList := make([]interface{}, 0, len(deviceSets))
	domainIdSet := false
	domainNameSet := false
	for _, item := range deviceSets {
		if item.Id != nil {
			deviceIdList = append(deviceIdList, int(*item.Id))
		}

		if !domainIdSet && item.DomainId != nil {
			_ = d.Set("domain_id", *item.DomainId)
			domainIdSet = true
		}

		if !domainNameSet && item.DomainName != nil {
			_ = d.Set("domain_name", *item.DomainName)
			domainNameSet = true
		}
	}

	_ = d.Set("device_id_set", deviceIdList)

	return nil
}

func resourceTencentCloudBhBindDeviceResourceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		resourceId = d.Id()
	)

	mutableArgs := []string{
		"device_id_set", "manage_dimension", "manage_account_id",
		"manage_account", "manage_kubeconfig", "namespace", "workload",
	}
	needChange := false
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if !needChange {
		return resourceTencentCloudBhBindDeviceResourceConfigRead(d, meta)
	}

	if d.HasChange("device_id_set") {
		oldVal, newVal := d.GetChange("device_id_set")
		olds := oldVal.(*schema.Set)
		news := newVal.(*schema.Set)
		remove := helper.InterfacesIntegers(olds.Difference(news).List())
		add := helper.InterfacesIntegers(news.Difference(olds).List())

		// Unbind removed devices by setting ResourceId to empty string
		if len(remove) > 0 {
			removeReq := bhv20230418.NewBindDeviceResourceRequest()
			for _, item := range remove {
				removeReq.DeviceIdSet = append(removeReq.DeviceIdSet, helper.IntUint64(item))
			}

			removeReq.ResourceId = helper.String("")
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().BindDeviceResourceWithContext(ctx, removeReq)
				if e != nil {
					return tccommon.RetryError(e)
				}

				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, removeReq.GetAction(), removeReq.ToJsonString(), result.ToJsonString())
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update bh bind device resource config (unbind) failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}

		// Bind newly added devices
		if len(add) > 0 {
			addReq := bhv20230418.NewBindDeviceResourceRequest()
			for _, item := range add {
				addReq.DeviceIdSet = append(addReq.DeviceIdSet, helper.IntUint64(item))
			}

			addReq.ResourceId = helper.String(resourceId)
			setBindDeviceResourceK8sParams(d, addReq)

			if v, ok := d.GetOk("domain_id"); ok {
				addReq.DomainId = helper.String(v.(string))
			}

			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().BindDeviceResourceWithContext(ctx, addReq)
				if e != nil {
					return tccommon.RetryError(e)
				}

				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, addReq.GetAction(), addReq.ToJsonString(), result.ToJsonString())
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update bh bind device resource config (bind) failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}
	} else {
		// Only K8S fields changed: re-bind all current devices with updated params
		request := bhv20230418.NewBindDeviceResourceRequest()
		if v, ok := d.GetOk("device_id_set"); ok {
			for _, item := range v.(*schema.Set).List() {
				request.DeviceIdSet = append(request.DeviceIdSet, helper.IntUint64(item.(int)))
			}
		}

		request.ResourceId = helper.String(resourceId)
		if v, ok := d.GetOk("domain_id"); ok {
			request.DomainId = helper.String(v.(string))
		}

		setBindDeviceResourceK8sParams(d, request)

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().BindDeviceResourceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bh bind device resource config failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBhBindDeviceResourceConfigRead(d, meta)
}

func resourceTencentCloudBhBindDeviceResourceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewBindDeviceResourceRequest()
	)

	if v, ok := d.GetOk("device_id_set"); ok {
		for _, item := range v.(*schema.Set).List() {
			request.DeviceIdSet = append(request.DeviceIdSet, helper.IntUint64(item.(int)))
		}
	}

	request.ResourceId = helper.String("")
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().BindDeviceResourceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bh bind device resource config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// setBindDeviceResourceK8sParams assigns K8S-related optional fields to the request.
func setBindDeviceResourceK8sParams(d *schema.ResourceData, request *bhv20230418.BindDeviceResourceRequest) {
	if v, ok := d.GetOk("manage_dimension"); ok {
		request.ManageDimension = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("manage_account_id"); ok {
		val := int64(v.(int))
		request.ManageAccountId = &val
	}

	if v, ok := d.GetOk("manage_account"); ok {
		request.ManageAccount = helper.String(v.(string))
	}

	if v, ok := d.GetOk("manage_kubeconfig"); ok {
		request.ManageKubeconfig = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
	}

	if v, ok := d.GetOk("workload"); ok {
		request.Workload = helper.String(v.(string))
	}
}
