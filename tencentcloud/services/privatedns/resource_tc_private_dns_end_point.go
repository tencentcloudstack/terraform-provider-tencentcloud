package privatedns

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatednsIntlv20201028 "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/privatedns/v20201028"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPrivateDnsEndPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPrivateDnsEndPointCreate,
		Read:   resourceTencentCloudPrivateDnsEndPointRead,
		Delete: resourceTencentCloudPrivateDnsEndPointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"end_point_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Endpoint name.",
			},

			"end_point_service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Endpoint service ID (namely, VPC endpoint service ID).",
			},

			"end_point_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Endpoint region, which should be consistent with the region of the endpoint service.",
			},

			"ip_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Number of endpoint IP addresses.",
			},

			"end_point_vip_set": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Vip list of endpoint.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudPrivateDnsEndPointCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_end_point.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = privatednsIntlv20201028.NewCreateEndPointRequest()
		response   = privatednsIntlv20201028.NewCreateEndPointResponse()
		endPointId string
	)

	if v, ok := d.GetOk("end_point_name"); ok {
		request.EndPointName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_service_id"); ok {
		request.EndPointServiceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_region"); ok {
		request.EndPointRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("ip_num"); ok {
		request.IpNum = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsIntlV20201028Client().CreateEndPointWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create private dns end point failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create private dns end point failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.EndPointId == nil {
		return fmt.Errorf("EndPointId is nil.")
	}

	endPointId = *response.Response.EndPointId
	d.SetId(endPointId)
	return resourceTencentCloudPrivateDnsEndPointRead(d, meta)
}

func resourceTencentCloudPrivateDnsEndPointRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_end_point.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = PrivatednsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		endPointId = d.Id()
	)

	respData, err := service.DescribePrivateDnsEndPointById(ctx, endPointId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_private_dns_end_point` [%s] not found, please check if it has been deleted.\n", logId, endPointId)
		return nil
	}

	if err := resourceTencentCloudPrivateDnsEndPointReadPreHandleResponse0(ctx, respData); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudPrivateDnsEndPointDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_end_point.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = privatednsIntlv20201028.NewDeleteEndPointRequest()
		endPointId = d.Id()
	)

	request.EndPointId = helper.String(endPointId)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsIntlV20201028Client().DeleteEndPointWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete private dns end point failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
