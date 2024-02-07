package tse

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTseCngwNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwNetworkCreate,
		Read:   resourceTencentCloudTseCngwNetworkRead,
		Update: resourceTencentCloudTseCngwNetworkUpdate,
		Delete: resourceTencentCloudTseCngwNetworkDelete,

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "gateway ID.",
			},

			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "gateway group ID.",
			},

			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "network id.",
			},

			"internet_address_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "internet type. Reference value:`IPV4` (default value), `IPV6`.",
			},
			"internet_pay_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "trade type of internet. Reference value:`BANDWIDTH` (default value), `TRAFFIC`.",
			},
			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "public network bandwidth.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "description of clb.",
			},
			"sla_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "specification type of clb. Default `shared` type when this parameter is empty, Note: input `shared` is not supported when creating. Reference value:`clb.c2.medium`, `clb.c3.small`, `clb.c3.medium`, `clb.c4.small`, `clb.c4.medium`, `clb.c4.large`, `clb.c4.xlarge`.",
			},
			"multi_zone_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether load balancing has multiple availability zones.",
			},
			"master_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "primary availability zone.",
			},
			"slave_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "alternate availability zone.",
			},
			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "clb vip.",
			},
		},
	}
}

func resourceTencentCloudTseCngwNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_network.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = tse.NewCreateCloudNativeAPIGatewayPublicNetworkRequest()
		response  = tse.NewCreateCloudNativeAPIGatewayPublicNetworkResponse()
		gatewayId string
		groupId   string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	internetConfig := tse.InternetConfig{}
	if v, ok := d.GetOk("internet_address_version"); ok {
		internetConfig.InternetAddressVersion = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_pay_mode"); ok {
		internetConfig.InternetPayMode = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		internetConfig.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("description"); ok {
		internetConfig.Description = helper.String(v.(string))
	}
	if v, ok := d.GetOk("sla_type"); ok {
		internetConfig.SlaType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("multi_zone_flag"); ok {
		internetConfig.MultiZoneFlag = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("master_zone_id"); ok {
		internetConfig.MasterZoneId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("slave_zone_id"); ok {
		internetConfig.SlaveZoneId = helper.String(v.(string))
	}
	request.InternetConfig = &internetConfig

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().CreateCloudNativeAPIGatewayPublicNetwork(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwNetwork failed, reason:%+v", logId, err)
		return err
	}

	networkId := *response.Response.Result.NetworkId
	d.SetId(gatewayId + tccommon.FILED_SP + groupId + tccommon.FILED_SP + networkId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Open"}, 5*tccommon.ReadRetryTimeout, time.Second, service.TseCngwNetworkStateRefreshFunc(gatewayId, groupId, networkId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTseCngwNetworkRead(d, meta)
}

func resourceTencentCloudTseCngwNetworkRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_network.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	groupId := idSplit[1]
	networkId := idSplit[2]

	cngwNetwork, err := service.DescribeTseCngwNetworkById(ctx, gatewayId, groupId, networkId)
	if err != nil {
		return err
	}

	if cngwNetwork == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwNetwork` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cngwNetwork.GatewayId != nil {
		_ = d.Set("gateway_id", cngwNetwork.GatewayId)
	}

	if cngwNetwork.GroupId != nil {
		_ = d.Set("group_id", cngwNetwork.GroupId)
	}

	_ = d.Set("network_id", networkId)

	if cngwNetwork.PublicNetwork != nil {
		internetConfig := cngwNetwork.PublicNetwork

		if internetConfig.NetType != nil {
			if *internetConfig.NetType == "Open" {
				_ = d.Set("internet_address_version", "IPV4")
			} else if *internetConfig.NetType == "Open-IPv6" {
				_ = d.Set("internet_address_version", "IPV6")
			}
		}

		if internetConfig.InternetMaxBandwidthOut != nil {
			_ = d.Set("internet_max_bandwidth_out", internetConfig.InternetMaxBandwidthOut)
		}

		if internetConfig.Description != nil {
			_ = d.Set("description", internetConfig.Description)
		}

		if internetConfig.SlaType != nil {
			_ = d.Set("sla_type", internetConfig.SlaType)
		}

		if internetConfig.MultiZoneFlag != nil {
			_ = d.Set("multi_zone_flag", internetConfig.MultiZoneFlag)
		}

		if internetConfig.MasterZoneId != nil {
			_ = d.Set("master_zone_id", internetConfig.MasterZoneId)
		}

		if internetConfig.SlaveZoneId != nil {
			_ = d.Set("slave_zone_id", internetConfig.SlaveZoneId)
		}

		if internetConfig.Vip != nil {
			_ = d.Set("vip", internetConfig.Vip)
		}
	}

	return nil
}

func resourceTencentCloudTseCngwNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_network.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := tse.NewModifyNetworkBasicInfoRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	groupId := idSplit[1]
	networkId := idSplit[2]

	request.GatewayId = &gatewayId
	request.GroupId = &groupId

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	cngwNetwork, e := service.DescribeTseCngwNetworkById(ctx, gatewayId, groupId, networkId)
	if e != nil {
		return e
	}
	if cngwNetwork == nil {
		return fmt.Errorf("[ERROR]%s resource `TseCngwNetwork` [%s] not found.\n", logId, d.Id())
	}
	request.Vip = cngwNetwork.PublicNetwork.Vip

	immutableArgs := []string{"internet_address_version", "internet_pay_mode", "sla_type", "multi_zone_flag", "master_zone_id", "slave_zone_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if v, ok := d.GetOk("internet_address_version"); ok {
		if v.(string) == "IPV4" {
			request.NetworkType = helper.String("Open")
		} else if v.(string) == "IPV6" {
			request.NetworkType = helper.String("Open-IPv6")
		}
	}

	if d.HasChange("internet_max_bandwidth_out") {
		if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
			request.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().ModifyNetworkBasicInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwNetwork failed, reason:%+v", logId, err)
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Open"}, 5*tccommon.ReadRetryTimeout, time.Second, service.TseCngwNetworkStateRefreshFunc(gatewayId, groupId, networkId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTseCngwNetworkRead(d, meta)
}

func resourceTencentCloudTseCngwNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_network.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	groupId := idSplit[1]
	networkId := idSplit[2]

	if err := service.DeleteTseCngwNetworkById(ctx, gatewayId, groupId, networkId); err != nil {
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Deleted"}, 5*tccommon.ReadRetryTimeout, time.Second, service.TseCngwNetworkStateRefreshFunc(gatewayId, groupId, networkId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
