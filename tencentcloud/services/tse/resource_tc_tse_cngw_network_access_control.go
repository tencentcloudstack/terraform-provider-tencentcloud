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

func ResourceTencentCloudTseCngwNetworkAccessControl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwNetworkAccessControlCreate,
		Read:   resourceTencentCloudTseCngwNetworkAccessControlRead,
		Update: resourceTencentCloudTseCngwNetworkAccessControlUpdate,
		Delete: resourceTencentCloudTseCngwNetworkAccessControlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "network id.",
			},

			"access_control": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "access control policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Access mode: `Whitelist`, `Blacklist`.",
						},
						"cidr_white_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "White list.",
						},
						"cidr_black_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Black list.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTseCngwNetworkAccessControlCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_network_access_control.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		gatewayId string
		groupId   string
		networkId string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
	}
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}
	if v, ok := d.GetOk("network_id"); ok {
		networkId = v.(string)
	}
	d.SetId(gatewayId + tccommon.FILED_SP + groupId + tccommon.FILED_SP + networkId)

	return resourceTencentCloudTseCngwNetworkAccessControlUpdate(d, meta)
}

func resourceTencentCloudTseCngwNetworkAccessControlRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_network_access_control.read")()
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
		log.Printf("[WARN]%s resource `TseCngwNetworkAccessControl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("gateway_id", gatewayId)
	_ = d.Set("group_id", groupId)
	_ = d.Set("network_id", networkId)

	if cngwNetwork.PublicNetwork != nil {
		internetConfig := cngwNetwork.PublicNetwork
		if internetConfig.AccessControl != nil {
			accessControlMap := map[string]interface{}{}

			accessControl := internetConfig.AccessControl
			if accessControl.Mode != nil {
				accessControlMap["mode"] = accessControl.Mode
			}
			if accessControl.Mode != nil {
				accessControlMap["cidr_white_list"] = accessControl.CidrWhiteList
			}
			if accessControl.Mode != nil {
				accessControlMap["cidr_black_list"] = accessControl.CidrBlackList
			}
			_ = d.Set("access_control", []interface{}{accessControlMap})
		}
	}

	return nil
}

func resourceTencentCloudTseCngwNetworkAccessControlUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_network_access_control.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := tse.NewModifyNetworkAccessStrategyRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	groupId := idSplit[1]
	networkId := idSplit[2]

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	cngwNetwork, err := service.DescribeTseCngwNetworkById(ctx, gatewayId, groupId, networkId)
	if err != nil {
		return err
	}
	if cngwNetwork == nil {
		return fmt.Errorf("[WARN]%s resource `TseCngwNetworkAccessControl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
	}

	request.GatewayId = helper.String(gatewayId)
	request.GroupId = helper.String(groupId)
	// The interface only supports public network
	request.NetworkType = helper.String("Open")
	request.Vip = cngwNetwork.PublicNetwork.Vip

	if d.HasChange("access_control") {
		if dMap, ok := helper.InterfacesHeadMap(d, "access_control"); ok {
			accessControl := tse.NetworkAccessControl{}
			if v, ok := dMap["mode"]; ok {
				accessControl.Mode = helper.String(v.(string))
			}
			if v, ok := dMap["cidr_white_list"]; ok {
				whitelist := v.([]interface{})
				accessControl.CidrWhiteList = helper.InterfacesStringsPoint(whitelist)
			}
			if v, ok := dMap["cidr_black_list"]; ok {
				blacklist := v.([]interface{})
				accessControl.CidrBlackList = helper.InterfacesStringsPoint(blacklist)
			}
			request.AccessControl = &accessControl
		}
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().ModifyNetworkAccessStrategy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwNetworkAccessStrategy failed, reason:%+v", logId, err)
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Open"}, 5*tccommon.ReadRetryTimeout, time.Second, service.TseCngwNetworkStateRefreshFunc(gatewayId, groupId, networkId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTseCngwNetworkAccessControlRead(d, meta)
}

func resourceTencentCloudTseCngwNetworkAccessControlDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_network_access_control.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
