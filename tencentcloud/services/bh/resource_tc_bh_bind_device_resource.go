package bh

import (
	"context"
	"fmt"
	"log"

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
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Device ID set.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bindable bastion host service ID.",
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Network domain ID.",
			},

			"manage_dimension": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "K8S cluster managed account dimension. 1-cluster, 2-namespace, 3-workload.",
			},

			"manage_account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "K8S cluster managed account ID.",
			},

			"manage_account": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "K8S cluster managed account name.",
			},

			"manage_kubeconfig": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "K8S cluster managed account kubeconfig credential.",
			},

			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "K8S cluster managed namespace.",
			},

			"workload": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "K8S cluster managed workload.",
			},
		},
	}
}

// classifyBhDeviceIdsByKind queries the devices by the given device id set and splits them into
// normal ids and k8s ids. A device whose Kind is 12 or 13 is treated as a K8s id, otherwise it is
// treated as a normal id.
func classifyBhDeviceIdsByKind(ctx context.Context, meta interface{}, deviceIds []uint64) (normalIds []uint64, k8sIds []uint64, err error) {
	if len(deviceIds) == 0 {
		return nil, nil, nil
	}

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := bhv20230418.NewDescribeDevicesRequest()
	for i := range deviceIds {
		id := deviceIds[i]
		request.IdSet = append(request.IdSet, &id)
	}

	request.Limit = helper.Uint64(100)
	kindMap := make(map[uint64]uint64, len(deviceIds))
	var offset uint64 = 0
	for {
		request.Offset = helper.Uint64(offset)
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
			return nil, nil, reqErr
		}

		deviceSet := response.Response.DeviceSet
		for _, item := range deviceSet {
			if item != nil && item.Id != nil && item.Kind != nil {
				kindMap[*item.Id] = *item.Kind
			}
		}

		if len(deviceSet) < 100 {
			break
		}

		offset += 100
	}

	for _, id := range deviceIds {
		if kind, ok := kindMap[id]; ok && (kind == 12 || kind == 13) {
			k8sIds = append(k8sIds, id)
		} else {
			normalIds = append(normalIds, id)
		}
	}

	return normalIds, k8sIds, nil
}

// bindBhDeviceResourceByKind calls the BindDeviceResource interface according to the device kind.
// Normal ids can be bound/unbound in a single batch call, while k8s ids (Kind 12 or 13) must be
// processed one by one. All the binding parameters except DeviceIdSet are copied from baseRequest.
func bindBhDeviceResourceByKind(ctx context.Context, meta interface{}, baseRequest *bhv20230418.BindDeviceResourceRequest, normalIds, k8sIds []uint64) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	doBind := func(ids []uint64) error {
		request := bhv20230418.NewBindDeviceResourceRequest()
		request.ResourceId = baseRequest.ResourceId
		request.DomainId = baseRequest.DomainId
		request.ManageDimension = baseRequest.ManageDimension
		request.ManageAccountId = baseRequest.ManageAccountId
		request.ManageAccount = baseRequest.ManageAccount
		request.ManageKubeconfig = baseRequest.ManageKubeconfig
		request.Namespace = baseRequest.Namespace
		request.Workload = baseRequest.Workload
		for i := range ids {
			id := ids[i]
			request.DeviceIdSet = append(request.DeviceIdSet, &id)
		}

		return resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
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
	}

	// Normal ids can be bound/unbound in a single batch call.
	if len(normalIds) > 0 {
		if err := doBind(normalIds); err != nil {
			return err
		}
	}

	// K8s ids (Kind 12 or 13) must be bound/unbound one by one.
	for _, id := range k8sIds {
		if err := doBind([]uint64{id}); err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudBhBindDeviceResourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = bhv20230418.NewBindDeviceResourceRequest()
		resourceId string
		deviceIds  []uint64
	)

	if v, ok := d.GetOk("device_id_set"); ok {
		for _, item := range v.(*schema.Set).List() {
			deviceIds = append(deviceIds, uint64(item.(int)))
		}
	}

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

	normalIds, k8sIds, err := classifyBhDeviceIdsByKind(ctx, meta, deviceIds)
	if err != nil {
		log.Printf("[CRITAL]%s classify bh devices by kind failed, reason:%+v", logId, err)
		return err
	}

	if err := bindBhDeviceResourceByKind(ctx, meta, request, normalIds, k8sIds); err != nil {
		log.Printf("[CRITAL]%s bind bh device resource failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(resourceId)
	return resourceTencentCloudBhBindDeviceResourceRead(d, meta)
}

func resourceTencentCloudBhBindDeviceResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = bhv20230418.NewDescribeDevicesRequest()
		response   = bhv20230418.NewDescribeDevicesResponse()
		resourceId = d.Id()
	)

	request.ResourceIdSet = helper.Strings([]string{resourceId})
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

	_ = d.Set("resource_id", resourceId)

	tmpList := make([]interface{}, 0, len(response.Response.DeviceSet))
	for _, item := range response.Response.DeviceSet {
		if item.Id != nil {
			tmpList = append(tmpList, item.Id)
		}

		if item.DomainId != nil {
			_ = d.Set("domain_id", item.DomainId)
		}

		if item.ManageDimension != nil {
			_ = d.Set("manage_dimension", int(*item.ManageDimension))
		}

		if item.ManageAccountId != nil {
			_ = d.Set("manage_account_id", int(*item.ManageAccountId))
		}

		if item.Namespace != nil {
			_ = d.Set("namespace", item.Namespace)
		}

		if item.Workload != nil {
			_ = d.Set("workload", item.Workload)
		}
	}

	_ = d.Set("device_id_set", tmpList)

	return nil
}

func resourceTencentCloudBhBindDeviceResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		resourceId = d.Id()
	)

	if d.HasChange("device_id_set") {
		oldInterface, newInterface := d.GetChange("device_id_set")
		olds := oldInterface.(*schema.Set)
		news := newInterface.(*schema.Set)
		remove := helper.InterfacesIntegers(olds.Difference(news).List())
		add := helper.InterfacesIntegers(news.Difference(olds).List())
		if len(remove) > 0 {
			removeIds := make([]uint64, 0, len(remove))
			for _, item := range remove {
				removeIds = append(removeIds, uint64(item))
			}

			normalIds, k8sIds, err := classifyBhDeviceIdsByKind(ctx, meta, removeIds)
			if err != nil {
				log.Printf("[CRITAL]%s classify bh devices by kind failed, reason:%+v", logId, err)
				return err
			}

			request := bhv20230418.NewBindDeviceResourceRequest()
			request.ResourceId = helper.String("")
			if err := bindBhDeviceResourceByKind(ctx, meta, request, normalIds, k8sIds); err != nil {
				log.Printf("[CRITAL]%s operate dasb bindDeviceResource failed, reason:%+v", logId, err)
				return err
			}
		}

		if len(add) > 0 {
			addIds := make([]uint64, 0, len(add))
			for _, item := range add {
				addIds = append(addIds, uint64(item))
			}

			request := bhv20230418.NewBindDeviceResourceRequest()
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

			request.ResourceId = helper.String(resourceId)
			normalIds, k8sIds, err := classifyBhDeviceIdsByKind(ctx, meta, addIds)
			if err != nil {
				log.Printf("[CRITAL]%s classify bh devices by kind failed, reason:%+v", logId, err)
				return err
			}

			if err := bindBhDeviceResourceByKind(ctx, meta, request, normalIds, k8sIds); err != nil {
				log.Printf("[CRITAL]%s operate dasb bindDeviceResource failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudBhBindDeviceResourceRead(d, meta)
}

func resourceTencentCloudBhBindDeviceResourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_resource.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = bhv20230418.NewBindDeviceResourceRequest()
		deviceIds []uint64
	)

	if v, ok := d.GetOk("device_id_set"); ok {
		deviceIdSetSet := v.(*schema.Set).List()
		for i := range deviceIdSetSet {
			deviceIds = append(deviceIds, uint64(deviceIdSetSet[i].(int)))
		}
	}

	request.ResourceId = helper.String("")
	normalIds, k8sIds, err := classifyBhDeviceIdsByKind(ctx, meta, deviceIds)
	if err != nil {
		log.Printf("[CRITAL]%s classify bh devices by kind failed, reason:%+v", logId, err)
		return err
	}

	if err := bindBhDeviceResourceByKind(ctx, meta, request, normalIds, k8sIds); err != nil {
		log.Printf("[CRITAL]%s delete bh bind device resource failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
