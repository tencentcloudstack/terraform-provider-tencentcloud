package tmp

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpInstanceCreate,
		Read:   resourceTencentCloudMonitorTmpInstanceRead,
		Update: resourceTencentCloudMonitorTmpInstanceUpdate,
		Delete: resourceTencentCloudMonitorTmpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Vpc Id.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet Id.",
			},

			"data_retention_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Data retention time(in days). Value range: 15, 30, 45, 90, 180, 360, 720.",
			},

			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Available zone.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},

			"ipv4_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance IPv4 address.",
			},

			"remote_write": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Prometheus remote write address.",
			},

			"api_root_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Prometheus HTTP API root address.",
			},

			"proxy_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Proxy address.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = monitor.NewCreatePrometheusMultiTenantInstancePostPayModeRequest()
		response *monitor.CreatePrometheusMultiTenantInstancePostPayModeResponse
	)

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("data_retention_time"); ok {
		request.DataRetentionTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreatePrometheusMultiTenantInstancePostPayMode(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create monitor tmpInstance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor tmpInstance failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.InstanceId == nil {
		return fmt.Errorf("InstanceId is nil.")
	}

	tmpInstanceId := *response.Response.InstanceId
	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeMonitorTmpInstance(ctx, tmpInstanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if *instance.InstanceStatus == 2 {
			return nil
		}

		if *instance.InstanceStatus == 3 {
			return resource.NonRetryableError(fmt.Errorf("tmpInstance status is %v, operate failed.", *instance.InstanceStatus))
		}

		return resource.RetryableError(fmt.Errorf("tmpInstance status is %v, retry...", *instance.InstanceStatus))
	})

	if err != nil {
		return err
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::monitor:%s:uin/:prom-instance/%s", region, tmpInstanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	d.SetId(tmpInstanceId)
	return resourceTencentCloudMonitorTmpInstanceRead(d, meta)
}

func resourceTencentCloudMonitorTmpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmpInstance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	tmpInstanceId := d.Id()
	tmpInstance, err := service.DescribeMonitorTmpInstance(ctx, tmpInstanceId)
	if err != nil {
		return err
	}

	if tmpInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tmpInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if tmpInstance.InstanceName != nil {
		_ = d.Set("instance_name", tmpInstance.InstanceName)
	}

	if tmpInstance.VpcId != nil {
		_ = d.Set("vpc_id", tmpInstance.VpcId)
	}

	if tmpInstance.SubnetId != nil {
		_ = d.Set("subnet_id", tmpInstance.SubnetId)
	}

	if tmpInstance.DataRetentionTime != nil {
		_ = d.Set("data_retention_time", tmpInstance.DataRetentionTime)
	}

	if tmpInstance.Zone != nil {
		_ = d.Set("zone", tmpInstance.Zone)
	}

	if tmpInstance.IPv4Address != nil {
		_ = d.Set("ipv4_address", tmpInstance.IPv4Address)
	}

	if tmpInstance.RemoteWrite != nil {
		_ = d.Set("remote_write", tmpInstance.RemoteWrite)
	}

	if tmpInstance.ApiRootPath != nil {
		_ = d.Set("api_root_path", tmpInstance.ApiRootPath)
	}

	if tmpInstance.ProxyAddress != nil {
		_ = d.Set("proxy_address", tmpInstance.ProxyAddress)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "monitor", "prom-instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudMonitorTmpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = monitor.NewModifyPrometheusInstanceAttributesRequest()
	)

	if d.HasChange("vpc_id") {
		return fmt.Errorf("`vpc_id` do not support change now.")
	}

	if d.HasChange("subnet_id") {
		return fmt.Errorf("`subnet_id` do not support change now.")
	}

	if d.HasChange("zone") {
		return fmt.Errorf("`zone` do not support change now.")
	}

	request.InstanceId = helper.String(d.Id())
	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if d.HasChange("data_retention_time") {
		if v, ok := d.GetOk("data_retention_time"); ok {
			request.DataRetentionTime = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().ModifyPrometheusInstanceAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("monitor", "prom-instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudMonitorTmpInstanceRead(d, meta)
}

func resourceTencentCloudMonitorTmpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		tmpInstanceId = d.Id()
	)

	if err := service.IsolateMonitorTmpInstanceById(ctx, tmpInstanceId); err != nil {
		return err
	}

	err := resource.Retry(1*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeMonitorTmpInstance(ctx, tmpInstanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if *instance.InstanceStatus == 6 {
			return nil
		}

		if *instance.InstanceStatus == 3 {
			return resource.NonRetryableError(fmt.Errorf("tmpInstance status is %v, operate failed.", *instance.InstanceStatus))
		}

		return resource.RetryableError(fmt.Errorf("tmpInstance status is %v, retry...", *instance.InstanceStatus))
	})

	if err != nil {
		return err
	}

	if err := service.DeleteMonitorTmpInstanceById(ctx, tmpInstanceId); err != nil {
		return err
	}

	return nil
}
