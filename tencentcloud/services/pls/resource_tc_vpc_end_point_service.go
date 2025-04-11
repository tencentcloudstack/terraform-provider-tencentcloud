package pls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcEndPointService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcEndPointServiceCreate,
		Read:   resourceTencentCloudVpcEndPointServiceRead,
		Update: resourceTencentCloudVpcEndPointServiceUpdate,
		Delete: resourceTencentCloudVpcEndPointServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of vpc instance.",
			},

			"end_point_service_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of end point service.",
			},

			"auto_accept_flag": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically accept.",
			},

			"service_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Id of service instance, like lb-xxx.",
			},

			"service_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Type of service instance, like `CLB`, `CDB`, `CRS`, default is `CLB`.",
			},

			"service_owner": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "APPID.",
			},

			"service_vip": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "VIP of backend service.",
			},

			"end_point_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Count of end point.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create Time.",
			},

			"cdc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CDC instance ID.",
			},
		},
	}
}

func resourceTencentCloudVpcEndPointServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point_service.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		request           = vpc.NewCreateVpcEndPointServiceRequest()
		response          = vpc.NewCreateVpcEndPointServiceResponse()
		endPointServiceId string
	)

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_service_name"); ok {
		request.EndPointServiceName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_accept_flag"); ok {
		request.AutoAcceptFlag = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("service_instance_id"); ok {
		request.ServiceInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_type"); ok {
		request.ServiceType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateVpcEndPointService(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.RetryableError(fmt.Errorf("Create vpc endPointService failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc endPointService failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.EndPointService == nil || response.Response.EndPointService.EndPointServiceId == nil {
		return fmt.Errorf("EndPointServiceId is nil.")
	}

	endPointServiceId = *response.Response.EndPointService.EndPointServiceId
	d.SetId(endPointServiceId)

	return resourceTencentCloudVpcEndPointServiceRead(d, meta)
}

func resourceTencentCloudVpcEndPointServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point_service.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service           = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		endPointServiceId = d.Id()
	)

	endPointService, err := service.DescribeVpcEndPointServiceById(ctx, endPointServiceId)
	if err != nil {
		return err
	}

	if endPointService == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if endPointService.VpcId != nil {
		_ = d.Set("vpc_id", endPointService.VpcId)
	}

	if endPointService.ServiceName != nil {
		_ = d.Set("end_point_service_name", endPointService.ServiceName)
	}

	if endPointService.AutoAcceptFlag != nil {
		_ = d.Set("auto_accept_flag", endPointService.AutoAcceptFlag)
	}

	if endPointService.ServiceInstanceId != nil {
		_ = d.Set("service_instance_id", endPointService.ServiceInstanceId)
	}

	if endPointService.ServiceType != nil {
		_ = d.Set("service_type", endPointService.ServiceType)
	}

	if endPointService.ServiceOwner != nil {
		_ = d.Set("service_owner", endPointService.ServiceOwner)
	}

	if endPointService.ServiceVip != nil {
		_ = d.Set("service_vip", endPointService.ServiceVip)
	}

	if endPointService.EndPointCount != nil {
		_ = d.Set("end_point_count", endPointService.EndPointCount)
	}

	if endPointService.CreateTime != nil {
		_ = d.Set("create_time", endPointService.CreateTime)
	}

	if endPointService.CdcId != nil {
		_ = d.Set("cdc_id", endPointService.CdcId)
	}

	return nil
}

func resourceTencentCloudVpcEndPointServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point_service.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		request           = vpc.NewModifyVpcEndPointServiceAttributeRequest()
		endPointServiceId = d.Id()
	)

	unsupportedUpdateFields := []string{
		"vpc_id",
		"service_type",
	}

	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_vpc_end_point_service update on %s is not support yet", field)
		}
	}

	request.EndPointServiceId = &endPointServiceId

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if d.HasChange("end_point_service_name") {
		if v, ok := d.GetOk("end_point_service_name"); ok {
			request.EndPointServiceName = helper.String(v.(string))
		}
	}

	if d.HasChange("auto_accept_flag") {
		if v, ok := d.GetOkExists("auto_accept_flag"); ok {
			request.AutoAcceptFlag = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("service_instance_id") {
		if v, ok := d.GetOk("service_instance_id"); ok {
			request.ServiceInstanceId = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyVpcEndPointServiceAttribute(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc endPointService failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVpcEndPointServiceRead(d, meta)
}

func resourceTencentCloudVpcEndPointServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point_service.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service           = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		endPointServiceId = d.Id()
	)

	if err := service.DeleteVpcEndPointServiceById(ctx, endPointServiceId); err != nil {
		return nil
	}

	return nil
}
