package dlc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcUserVpcConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcUserVpcConnectionCreate,
		Read:   resourceTencentCloudDlcUserVpcConnectionRead,
		Delete: resourceTencentCloudDlcUserVpcConnectionDelete,
		Schema: map[string]*schema.Schema{
			"user_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User vpc ID.",
			},

			"user_subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User subnet ID.",
			},

			"user_vpc_endpoint_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User vpc endpoint name.",
			},

			"engine_network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Engine network ID.",
			},

			"user_vpc_endpoint_vip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Manually specify VIP, if not filled in, an IP address under the subnet will be automatically assigned.",
			},

			// computed
			"user_vpc_endpoint_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User endpoint ID.",
			},
		},
	}
}

func resourceTencentCloudDlcUserVpcConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_user_vpc_connection.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request           = dlcv20210125.NewCreateUserVpcConnectionRequest()
		response          = dlcv20210125.NewCreateUserVpcConnectionResponse()
		engineNetworkId   string
		userVpcEndpointId string
	)

	if v, ok := d.GetOk("user_vpc_id"); ok {
		request.UserVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_subnet_id"); ok {
		request.UserSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_vpc_endpoint_name"); ok {
		request.UserVpcEndpointName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_network_id"); ok {
		request.EngineNetworkId = helper.String(v.(string))
		engineNetworkId = v.(string)
	}

	if v, ok := d.GetOk("user_vpc_endpoint_vip"); ok {
		request.UserVpcEndpointVip = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().CreateUserVpcConnectionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc user vpc connection failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc user vpc connection failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.UserVpcEndpointId == nil {
		return fmt.Errorf("UserVpcEndpointId is nil.")
	}

	userVpcEndpointId = *response.Response.UserVpcEndpointId
	d.SetId(strings.Join([]string{engineNetworkId, userVpcEndpointId}, tccommon.FILED_SP))
	return resourceTencentCloudDlcUserVpcConnectionRead(d, meta)
}

func resourceTencentCloudDlcUserVpcConnectionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_user_vpc_connection.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	engineNetworkId := idSplit[0]
	userVpcEndpointId := idSplit[1]

	respData, err := service.DescribeDlcUserVpcConnectionById(ctx, engineNetworkId, userVpcEndpointId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dlc_user_vpc_connection` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.UserVpcId != nil {
		_ = d.Set("user_vpc_id", respData.UserVpcId)
	}

	if respData.UserVpcEndpointName != nil {
		_ = d.Set("user_vpc_endpoint_name", respData.UserVpcEndpointName)
	}

	if respData.EngineNetworkId != nil {
		_ = d.Set("engine_network_id", respData.EngineNetworkId)
	}

	if respData.UserVpcEndpointId != nil {
		_ = d.Set("user_vpc_endpoint_id", respData.UserVpcEndpointId)
	}

	return nil
}

func resourceTencentCloudDlcUserVpcConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_user_vpc_connection.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = dlcv20210125.NewDeleteUserVpcConnectionRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	engineNetworkId := idSplit[0]
	userVpcEndpointId := idSplit[1]

	request.EngineNetworkId = helper.String(engineNetworkId)
	request.UserVpcEndpointId = helper.String(userVpcEndpointId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DeleteUserVpcConnectionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dlc user vpc connection failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
