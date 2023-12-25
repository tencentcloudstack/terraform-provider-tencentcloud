package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcPeerConnectManager() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcPeerConnectManagerCreate,
		Read:   resourceTencentCloudVpcPeerConnectManagerRead,
		Update: resourceTencentCloudVpcPeerConnectManagerUpdate,
		Delete: resourceTencentCloudVpcPeerConnectManagerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"source_vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the local VPC.",
			},

			"peering_connection_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Peer connection name.",
			},

			"destination_vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the peer VPC.",
			},

			"destination_uin": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Peer user UIN.",
			},

			"destination_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Peer region.",
			},

			"bandwidth": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Bandwidth upper limit, unit Mbps.",
			},

			"type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Interworking type, VPC_PEER interworking between VPCs; VPC_BM_PEER interworking between VPC and BM Network.",
			},

			"charge_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Billing mode, daily peak value POSTPAID_BY_DAY_MAX, monthly value 95 POSTPAID_BY_MONTH_95.",
			},

			"qos_level": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Service classification PT, AU, AG.",
			},
		},
	}
}

func resourceTencentCloudVpcPeerConnectManagerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_peer_connect_manager.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request             = vpc.NewCreateVpcPeeringConnectionRequest()
		response            = vpc.NewCreateVpcPeeringConnectionResponse()
		peeringConnectionId string
	)
	if v, ok := d.GetOk("source_vpc_id"); ok {
		request.SourceVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("peering_connection_name"); ok {
		request.PeeringConnectionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_vpc_id"); ok {
		request.DestinationVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_uin"); ok {
		request.DestinationUin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_region"); ok {
		request.DestinationRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("bandwidth"); ok {
		request.Bandwidth = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("charge_type"); ok {
		request.ChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qos_level"); ok {
		request.QosLevel = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateVpcPeeringConnection(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc PeerConnectManager failed, reason:%+v", logId, err)
		return err
	}

	peeringConnectionId = *response.Response.PeeringConnectionId
	d.SetId(peeringConnectionId)

	return resourceTencentCloudVpcPeerConnectManagerRead(d, meta)
}

func resourceTencentCloudVpcPeerConnectManagerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_peer_connect_manager.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	peeringConnectionId := d.Id()

	PeerConnectManager, err := service.DescribeVpcPeerConnectManagerById(ctx, peeringConnectionId)
	if err != nil {
		return err
	}

	if PeerConnectManager == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcPeerConnectManager` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if PeerConnectManager.SourceVpcId != nil {
		_ = d.Set("source_vpc_id", PeerConnectManager.SourceVpcId)
	}

	if PeerConnectManager.PeeringConnectionName != nil {
		_ = d.Set("peering_connection_name", PeerConnectManager.PeeringConnectionName)
	}

	if PeerConnectManager.DestinationUin != nil {
		_ = d.Set("destination_uin", helper.Int64ToStr(*PeerConnectManager.DestinationUin))
	}

	if PeerConnectManager.DestinationRegion != nil {
		_ = d.Set("destination_region", common.ShortRegionNameParse(*PeerConnectManager.DestinationRegion))
	}

	if PeerConnectManager.Bandwidth != nil {
		_ = d.Set("bandwidth", PeerConnectManager.Bandwidth)
	}

	if PeerConnectManager.Type != nil {
		_ = d.Set("type", PeerConnectManager.Type)
	}

	if PeerConnectManager.ChargeType != nil {
		_ = d.Set("charge_type", PeerConnectManager.ChargeType)
	}

	if PeerConnectManager.QosLevel != nil {
		_ = d.Set("qos_level", PeerConnectManager.QosLevel)
	}

	return nil
}

func resourceTencentCloudVpcPeerConnectManagerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_peer_connect_manager.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vpc.NewModifyVpcPeeringConnectionRequest()

	peeringConnectionId := d.Id()

	request.PeeringConnectionId = &peeringConnectionId

	immutableArgs := []string{"source_vpc_id", "destination_vpc_id", "destination_uin", "destination_region", "type", "qos_level"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("peering_connection_name") {
		if v, ok := d.GetOk("peering_connection_name"); ok {
			request.PeeringConnectionName = helper.String(v.(string))
		}
	}

	if d.HasChange("bandwidth") {
		if v, ok := d.GetOkExists("bandwidth"); ok {
			request.Bandwidth = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("charge_type") {
		if v, ok := d.GetOk("charge_type"); ok {
			request.ChargeType = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyVpcPeeringConnection(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc PeerConnectManager failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVpcPeerConnectManagerRead(d, meta)
}

func resourceTencentCloudVpcPeerConnectManagerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_peer_connect_manager.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	peeringConnectionId := d.Id()

	if err := service.DeleteVpcPeerConnectManagerById(ctx, peeringConnectionId); err != nil {
		return err
	}

	return nil
}
