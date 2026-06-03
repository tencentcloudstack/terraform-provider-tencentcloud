package bh

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhBindDeviceResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhBindDeviceResourceCreate,
		Read:   resourceTencentCloudBhBindDeviceResourceRead,
		Update: resourceTencentCloudBhBindDeviceResourceUpdate,
		Delete: resourceTencentCloudBhBindDeviceResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"device_id_set": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Device ID set.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bindable bastion host service ID.",
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network domain ID.",
			},

			"manage_dimension": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "K8S cluster managed account dimension. 1-cluster, 2-namespace, 3-workload.",
			},

			"manage_account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "K8S cluster managed account ID.",
			},

			"manage_account": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
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
		},
	}
}

func resourceTencentCloudBhBindDeviceResourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewBindDeviceResourceRequest()
	)

	var deviceIdStrs []string
	if v, ok := d.GetOk("device_id_set"); ok {
		for _, item := range v.([]interface{}) {
			deviceId := uint64(item.(int))
			request.DeviceIdSet = append(request.DeviceIdSet, &deviceId)
			deviceIdStrs = append(deviceIdStrs, strconv.FormatUint(deviceId, 10))
		}
	}

	var resourceId string
	if v, ok := d.GetOk("resource_id"); ok {
		resourceId = v.(string)
		request.ResourceId = helper.String(resourceId)
	}

	if v, ok := d.GetOk("domain_id"); ok {
		request.DomainId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("manage_dimension"); ok {
		request.ManageDimension = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("manage_account_id"); ok {
		request.ManageAccountId = helper.IntInt64(v.(int))
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

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().BindDeviceResourceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Bind bh device resource failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s bind bh device resource failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join(deviceIdStrs, ",") + tccommon.FILED_SP + resourceId)
	return resourceTencentCloudBhBindDeviceResourceRead(d, meta)
}

func resourceTencentCloudBhBindDeviceResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewDescribeDevicesRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	deviceIdsStr := idSplit[0]
	resourceId := idSplit[1]

	deviceIdStrs := strings.Split(deviceIdsStr, ",")
	if len(deviceIdStrs) == 0 {
		return fmt.Errorf("device_id_set is empty in id: %s", d.Id())
	}

	// Use the first device ID to query
	firstDeviceId, err := strconv.ParseUint(deviceIdStrs[0], 10, 64)
	if err != nil {
		return fmt.Errorf("parse device id failed: %v", err)
	}

	request.IdSet = []*uint64{&firstDeviceId}

	var response *bhv20230418.DescribeDevicesResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().DescribeDevicesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe bh devices failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read bh bind device resource failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.DeviceSet == nil || len(response.Response.DeviceSet) == 0 {
		log.Printf("[WARN]%s resource `tencentcloud_bh_bind_device_resource` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	device := response.Response.DeviceSet[0]
	if device.Resource == nil {
		log.Printf("[WARN]%s resource `tencentcloud_bh_bind_device_resource` [%s] device is not bound to any service, removing from state.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	// Set device_id_set from the ID
	deviceIds := make([]int, 0, len(deviceIdStrs))
	for _, idStr := range deviceIdStrs {
		id, _ := strconv.Atoi(idStr)
		deviceIds = append(deviceIds, id)
	}
	_ = d.Set("device_id_set", deviceIds)

	_ = d.Set("resource_id", resourceId)

	if device.DomainId != nil {
		_ = d.Set("domain_id", device.DomainId)
	}

	if device.ManageDimension != nil {
		_ = d.Set("manage_dimension", int(*device.ManageDimension))
	}

	if device.ManageAccountId != nil {
		_ = d.Set("manage_account_id", int(*device.ManageAccountId))
	}

	if device.Namespace != nil {
		_ = d.Set("namespace", device.Namespace)
	}

	if device.Workload != nil {
		_ = d.Set("workload", device.Workload)
	}

	return nil
}

func resourceTencentCloudBhBindDeviceResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewBindDeviceResourceRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	deviceIdsStr := idSplit[0]
	deviceIdStrs := strings.Split(deviceIdsStr, ",")

	for _, idStr := range deviceIdStrs {
		deviceId, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return fmt.Errorf("parse device id failed: %v", err)
		}
		request.DeviceIdSet = append(request.DeviceIdSet, &deviceId)
	}

	if v, ok := d.GetOk("resource_id"); ok {
		request.ResourceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain_id"); ok {
		request.DomainId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("manage_dimension"); ok {
		request.ManageDimension = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("manage_account_id"); ok {
		request.ManageAccountId = helper.IntInt64(v.(int))
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

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().BindDeviceResourceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Update bh bind device resource failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update bh bind device resource failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Update the ID if resource_id changed
	if d.HasChange("resource_id") {
		newResourceId := d.Get("resource_id").(string)
		d.SetId(deviceIdsStr + tccommon.FILED_SP + newResourceId)
	}

	return resourceTencentCloudBhBindDeviceResourceRead(d, meta)
}

func resourceTencentCloudBhBindDeviceResourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewBindDeviceResourceRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	deviceIdsStr := idSplit[0]
	deviceIdStrs := strings.Split(deviceIdsStr, ",")

	for _, idStr := range deviceIdStrs {
		deviceId, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return fmt.Errorf("parse device id failed: %v", err)
		}
		request.DeviceIdSet = append(request.DeviceIdSet, &deviceId)
	}

	// Set ResourceId to empty string to unbind
	request.ResourceId = helper.String("")

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().BindDeviceResourceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Unbind bh device resource failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bh bind device resource failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
