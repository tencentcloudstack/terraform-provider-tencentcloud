package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpnConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnConnectionCreate,
		Read:   resourceTencentCloudVpnConnectionRead,
		Update: resourceTencentCloudVpnConnectionUpdate,
		Delete: resourceTencentCloudVpnConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the VPN connection. The length of character is limited to 1-60.",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("is_ccn_type"); ok && v.(bool) {
						return true
					}
					return old == new
				},
				Description: "ID of the VPC. Required if vpn gateway is not in `CCN` type, and doesn't make sense for `CCN` vpn gateway.",
			},
			"is_ccn_type": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicate whether is ccn type. Modification of this field only impacts force new logic of `vpc_id`. If `is_ccn_type` is true, modification of `vpc_id` will be ignored.",
			},
			"customer_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the customer gateway.",
			},
			"vpn_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPN gateway.",
			},
			"pre_share_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pre-shared key of the VPN connection.",
			},
			"security_group_policy": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "SPD policy group, for example: {\"10.0.0.5/24\":[\"172.123.10.5/16\"]}, 10.0.0.5/24 is the vpc intranet segment, and 172.123.10.5/16 is the IDC network segment. " +
					"Users specify which network segments in the VPC can communicate with which network segments in your IDC.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local_cidr_block": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Local cidr block.",
						},
						"remote_cidr_block": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Remote cidr block list.",
						},
					},
				},
			},
			"ike_proto_encry_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     svcvpc.VPN_IKE_PROPO_ENCRY_ALGORITHM_3DESCBC,
				Description: "Proto encrypt algorithm of the IKE operation specification. Valid values: `3DES-CBC`, `AES-CBC-128`, `AES-CBC-192`, `AES-CBC-256`, `DES-CBC`, `SM4`, `AES128GCM128`, `AES192GCM128`, `AES256GCM128`,`AES128GCM128`, `AES192GCM128`, `AES256GCM128`. Default value is `3DES-CBC`.",
			},
			"ike_proto_authen_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     svcvpc.VPN_IKE_PROPO_AUTHEN_ALGORITHM_MD5,
				Description: "Proto authenticate algorithm of the IKE operation specification. Valid values: `MD5`, `SHA`, `SHA-256`. Default Value is `MD5`.",
			},
			"ike_exchange_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     svcvpc.VPN_IKE_EXCHANGE_MODE_MAIN,
				Description: "Exchange mode of the IKE operation specification. Valid values: `AGGRESSIVE`, `MAIN`. Default value is `MAIN`.",
			},
			"ike_local_identity": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     svcvpc.VPN_IKE_IDENTITY_ADDRESS,
				Description: "Local identity way of IKE operation specification. Valid values: `ADDRESS`, `FQDN`. Default value is `ADDRESS`.",
			},
			"ike_remote_identity": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     svcvpc.VPN_IKE_IDENTITY_ADDRESS,
				Description: "Remote identity way of IKE operation specification. Valid values: `ADDRESS`, `FQDN`. Default value is `ADDRESS`.",
			},
			"ike_local_address": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ike_local_fqdn_name"},
				Description:   "Local address of IKE operation specification, valid when ike_local_identity is `ADDRESS`, generally the value is `public_ip_address` of the related VPN gateway.",
			},
			"ike_remote_address": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ike_remote_fqdn_name"},
				Description:   "Remote address of IKE operation specification, valid when ike_remote_identity is `ADDRESS`, generally the value is `public_ip_address` of the related customer gateway.",
			},
			"ike_local_fqdn_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"ike_local_address"},
				Description:   "Local FQDN name of the IKE operation specification.",
			},
			"ike_remote_fqdn_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"ike_remote_address"},
				Description:   "Remote FQDN name of the IKE operation specification.",
			},
			"ike_dh_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     svcvpc.VPN_IKE_DH_GROUP_NAME_GROUP1,
				Description: "DH group name of the IKE operation specification. Valid values: `GROUP1`, `GROUP2`, `GROUP5`, `GROUP14`, `GROUP24`. Default value is `GROUP1`.",
			},
			"ike_sa_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      86400,
				ValidateFunc: tccommon.ValidateIntegerInRange(60, 604800),
				Description:  "SA lifetime of the IKE operation specification, unit is `second`. The value ranges from 60 to 604800. Default value is 86400 seconds.",
			},
			"ike_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "IKEV1",
				Description: "Version of the IKE operation specification, values: `IKEV1`, `IKEV2`. Default value is `IKEV1`.",
			},
			"ipsec_encrypt_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     svcvpc.VPN_IPSEC_ENCRY_ALGORITHM_3DESCBC,
				Description: "Encrypt algorithm of the IPSEC operation specification. Valid values: `3DES-CBC`, `AES-CBC-128`, `AES-CBC-192`, `AES-CBC-256`, `DES-CBC`, `SM4`, `NULL`, `AES128GCM128`, `AES192GCM128`, `AES256GCM128`. Default value is `3DES-CBC`.",
			},
			"ipsec_integrity_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     svcvpc.VPN_IPSEC_INTEGRITY_ALGORITHM_MD5,
				Description: "Integrity algorithm of the IPSEC operation specification. Valid values: `SHA1`, `MD5`, `SHA-256`. Default value is `MD5`.",
			},
			"ipsec_sa_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				ValidateFunc: tccommon.ValidateIntegerInRange(180, 604800),
				Description:  "SA lifetime of the IPSEC operation specification, unit is second. Valid value ranges: [180~604800]. Default value is 3600 seconds.",
			},
			"ipsec_pfs_dh_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NULL",
				Description: "PFS DH group. Valid value: `DH-GROUP1`, `DH-GROUP2`, `DH-GROUP5`, `DH-GROUP14`, `DH-GROUP24`, `NULL`. Default value is `NULL`.",
			},
			"ipsec_sa_lifetime_traffic": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1843200,
				ValidateFunc: tccommon.ValidateIntegerMin(2560),
				Description:  "SA lifetime of the IPSEC operation specification, unit is KB. The value should not be less then 2560. Default value is 1843200.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A list of tags used to associate different resources.",
			},
			"dpd_enable": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 1),
				Description:  "Specifies whether to enable DPD. Valid values: 0 (disable) and 1 (enable).",
			},
			"dpd_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(30, 60),
				Description:  "DPD timeout period.Valid value ranges: [30~60], Default: 30; unit: second. If the request is not responded within this period, the peer end is considered not exists. This parameter is valid when the value of DpdEnable is 1.",
			},
			"dpd_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(svcvpc.DPD_ACTIONS),
				Description:  "The action after DPD timeout. Valid values: clear (disconnect) and restart (try again). It is valid when DpdEnable is 1.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the VPN connection.",
			},
			"vpn_proto": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Vpn proto of the VPN connection.",
			},
			"encrypt_proto": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Encrypt proto of the VPN connection.",
			},
			"route_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(svcvpc.VPN_CONNECTION_ROUTE_TYPE),
				Description:  "Route type of the VPN connection. Valid value: `STATIC`, `StaticRoute`, `Policy`, `Bgp`.",
			},
			"negotiation_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The default negotiation type is `active`. Optional values: `active` (active negotiation), `passive` (passive negotiation), `flowTrigger` (traffic negotiation).",
			},
			// "route": {
			// 	Type:        schema.TypeList,
			// 	Optional:    true,
			// 	ForceNew:    true,
			// 	MaxItems:    1,
			// 	Description: "Create channel routing information.",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"destination_cidr_block": {
			// 				Type:        schema.TypeString,
			// 				Required:    true,
			// 				Description: "Destination IDC network segment.",
			// 			},
			// 			"priority": {
			// 				Type:        schema.TypeInt,
			// 				Optional:    true,
			// 				Description: "Priority. Optional value [0, 100].",
			// 			},
			// 		},
			// 	},
			// },
			"bgp_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "BGP config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tunnel_cidr": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "BGP tunnel segment.",
						},
						"local_bgp_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cloud BGP address. It must be allocated from within the BGP tunnel network segment.",
						},
						"remote_bgp_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User side BGP address. It must be allocated from within the BGP tunnel network segment.",
						},
					},
				},
			},
			"health_check_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "VPN channel health check configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"probe_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Detection mode, default is `NQA`, cannot be modified.",
						},
						"probe_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Detection interval, Tencent Cloud's interval between two health checks, range [1000-5000], Unit: ms.",
						},
						"probe_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Detection times, perform route switching after N consecutive health check failures, range [3-8], Unit: times.",
						},
						"probe_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Detection timeout, range [10-5000], Unit: ms.",
						},
					},
				},
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the connection. Valid value: `PENDING`, `AVAILABLE`, `DELETING`.",
			},
			"net_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Net status of the VPN connection. Valid value: `AVAILABLE`.",
			},
			"enable_health_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether intra-tunnel health checks are supported.",
			},
			"health_check_local_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Health check the address of this terminal.",
			},
			"health_check_remote_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Health check peer address.",
			},
		},
	}
}

func resourceTencentCloudVpnConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_connection.create")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	// pre check vpn gateway id
	has, gateway, err := service.DescribeVpngwById(ctx, d.Get("vpn_gateway_id").(string))
	if err != nil {
		return err
	}

	if !has {
		return fmt.Errorf("[CRITAL] vpn_gateway_id %s doesn't exist", d.Get("vpn_gateway_id").(string))
	}

	// create vpn connection
	request := vpc.NewCreateVpnConnectionRequest()
	request.VpnConnectionName = helper.String(d.Get("name").(string))
	if *gateway.Type != "CCN" {
		if _, ok := d.GetOk("vpc_id"); !ok {
			return fmt.Errorf("[CRITAL] vpc_id is required for this vpn connection which vpn gateway is in %s type", *gateway.Type)
		}
		request.VpcId = helper.String(d.Get("vpc_id").(string))
	} else {
		if _, ok := d.GetOk("vpc_id"); ok {
			return fmt.Errorf("[CRITAL] vpc_id doesn't make sense when vpn gateway is in CCN type")
		}
		request.VpcId = helper.String("")
	}

	request.VpnGatewayId = helper.String(d.Get("vpn_gateway_id").(string))
	request.CustomerGatewayId = helper.String(d.Get("customer_gateway_id").(string))
	request.PreShareKey = helper.String(d.Get("pre_share_key").(string))
	if v, ok := d.GetOk("dpd_enable"); ok {
		dpdEnable := v.(int)
		request.DpdEnable = helper.IntInt64(dpdEnable)
	}

	if v, ok := d.GetOk("dpd_action"); ok {
		request.DpdAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dpd_timeout"); ok {
		request.DpdTimeout = helper.String(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("route_type"); ok {
		request.RouteType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("negotiation_type"); ok {
		request.NegotiationType = helper.String(v.(string))
	}

	//set up SecurityPolicyDatabases
	if v, ok := d.GetOk("security_group_policy"); ok {
		for _, item := range v.(*schema.Set).List() {
			if dMap, ok := item.(map[string]interface{}); ok && dMap != nil {
				var sgp vpc.SecurityPolicyDatabase
				if v, ok := dMap["local_cidr_block"].(string); ok && v != "" {
					sgp.LocalCidrBlock = &v
				}

				if v, ok := dMap["remote_cidr_block"].(*schema.Set); ok {
					remoteCidrBlocks := v.List()
					for _, rcb := range remoteCidrBlocks {
						if v, ok := rcb.(string); ok && v != "" {
							sgp.RemoteCidrBlock = append(sgp.RemoteCidrBlock, &v)
						}
					}
				}

				request.SecurityPolicyDatabases = append(request.SecurityPolicyDatabases, &sgp)
			}
		}
	}

	//set up IKEOptionsSpecification
	var ikeOptionsSpecification vpc.IKEOptionsSpecification
	ikeOptionsSpecification.PropoEncryAlgorithm = helper.String(d.Get("ike_proto_encry_algorithm").(string))
	ikeOptionsSpecification.PropoAuthenAlgorithm = helper.String(d.Get("ike_proto_authen_algorithm").(string))
	ikeOptionsSpecification.ExchangeMode = helper.String(d.Get("ike_exchange_mode").(string))
	ikeOptionsSpecification.LocalIdentity = helper.String(d.Get("ike_local_identity").(string))
	ikeOptionsSpecification.RemoteIdentity = helper.String(d.Get("ike_remote_identity").(string))
	if *ikeOptionsSpecification.LocalIdentity == svcvpc.VPN_IKE_IDENTITY_ADDRESS {
		if v, ok := d.GetOk("ike_local_address"); ok {
			ikeOptionsSpecification.LocalAddress = helper.String(v.(string))
		} else {
			return fmt.Errorf("ike_local_address need to be set when ike_local_identity is `ADDRESS`.")
		}
	} else {
		if v, ok := d.GetOk("ike_local_fqdn_name"); ok {
			ikeOptionsSpecification.LocalFqdnName = helper.String(v.(string))
		} else {
			return fmt.Errorf("ike_local_fqdn_name need to be set when ike_local_identity is `FQDN`")
		}
	}

	if *ikeOptionsSpecification.LocalIdentity == svcvpc.VPN_IKE_IDENTITY_ADDRESS {
		if v, ok := d.GetOk("ike_remote_address"); ok {
			ikeOptionsSpecification.RemoteAddress = helper.String(v.(string))
		} else {
			return fmt.Errorf("ike_remote_address need to be set when ike_remote_identity is `ADDRESS`.")
		}
	} else {
		if v, ok := d.GetOk("ike_remote_fqdn_name"); ok {
			ikeOptionsSpecification.RemoteFqdnName = helper.String(v.(string))
		} else {
			return fmt.Errorf("ike_remote_fqdn_name need to be set when ike_remote_identity is `FQDN`")
		}
	}

	ikeOptionsSpecification.DhGroupName = helper.String(d.Get("ike_dh_group_name").(string))
	saLifetime := d.Get("ike_sa_lifetime_seconds").(int)
	saLifetime64 := uint64(saLifetime)
	ikeOptionsSpecification.IKESaLifetimeSeconds = &saLifetime64
	ikeOptionsSpecification.IKEVersion = helper.String(d.Get("ike_version").(string))
	request.IKEOptionsSpecification = &ikeOptionsSpecification

	//set up IPSECOptionsSpecification
	var ipsecOptionsSpecification vpc.IPSECOptionsSpecification
	ipsecOptionsSpecification.EncryptAlgorithm = helper.String(d.Get("ipsec_encrypt_algorithm").(string))
	ipsecOptionsSpecification.IntegrityAlgorith = helper.String(d.Get("ipsec_integrity_algorithm").(string))
	ipsecSaLifetimeSeconds := d.Get("ipsec_sa_lifetime_seconds").(int)
	ipsecSaLifetimeSeconds64 := uint64(ipsecSaLifetimeSeconds)
	ipsecOptionsSpecification.IPSECSaLifetimeSeconds = &ipsecSaLifetimeSeconds64
	ipsecOptionsSpecification.PfsDhGroup = helper.String(d.Get("ipsec_pfs_dh_group").(string))
	ipsecSaLifetimeTraffic := d.Get("ipsec_sa_lifetime_traffic").(int)
	ipsecSaLifetimeTraffic64 := uint64(ipsecSaLifetimeTraffic)
	ipsecOptionsSpecification.IPSECSaLifetimeTraffic = &ipsecSaLifetimeTraffic64
	request.IPSECOptionsSpecification = &ipsecOptionsSpecification
	if v, ok := d.GetOk("enable_health_check"); ok {
		request.EnableHealthCheck = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("health_check_local_ip"); ok {
		request.HealthCheckLocalIp = helper.String(v.(string))
	}

	if v, ok := d.GetOk("health_check_remote_ip"); ok {
		request.HealthCheckRemoteIp = helper.String(v.(string))
	}

	// if v, ok := d.GetOk("route"); ok {
	// 	for _, item := range v.([]interface{}) {
	// 		dMap := item.(map[string]interface{})
	// 		route := vpc.CreateVpnConnRoute{}
	// 		if v, ok := dMap["destination_cidr_block"]; ok {
	// 			route.DestinationCidrBlock = helper.String(v.(string))
	// 		}

	// 		if v, ok := dMap["priority"]; ok {
	// 			route.Priority = helper.IntUint64(v.(int))
	// 		}

	// 		request.Route = &route
	// 	}
	// }

	if v, ok := d.GetOk("bgp_config"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			bgpConfig := vpc.BgpConfig{}
			if v, ok := dMap["tunnel_cidr"]; ok {
				bgpConfig.TunnelCidr = helper.String(v.(string))
			}

			if v, ok := dMap["local_bgp_ip"]; ok {
				bgpConfig.LocalBgpIp = helper.String(v.(string))
			}

			if v, ok := dMap["remote_bgp_ip"]; ok {
				bgpConfig.RemoteBgpIp = helper.String(v.(string))
			}

			request.BgpConfig = &bgpConfig
		}
	}

	if v, ok := d.GetOk("health_check_config"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			healthCheckConfig := vpc.HealthCheckConfig{}
			if v, ok := dMap["probe_type"]; ok {
				healthCheckConfig.ProbeType = helper.String(v.(string))
			}

			if v, ok := dMap["probe_interval"]; ok {
				healthCheckConfig.ProbeInterval = helper.IntInt64(v.(int))
			}

			if v, ok := dMap["probe_threshold"]; ok {
				healthCheckConfig.ProbeThreshold = helper.IntInt64(v.(int))
			}

			if v, ok := dMap["probe_timeout"]; ok {
				healthCheckConfig.ProbeTimeout = helper.IntInt64(v.(int))
			}

			request.HealthCheckConfig = &healthCheckConfig
		}
	}

	var response *vpc.CreateVpnConnectionResponse
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateVpnConnection(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create VPN connection failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create VPN connection failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response.VpnConnection == nil {
		return fmt.Errorf("VpnConnection is nil.")
	}

	vpnConnectionId := ""
	//the response will return "" id
	if *response.Response.VpnConnection.VpnConnectionId == "" {
		idRequest := vpc.NewDescribeVpnConnectionsRequest()
		idRequest.Filters = make([]*vpc.Filter, 0, 3)
		params := make(map[string]string)
		if v, ok := d.GetOk("vpn_gateway_id"); ok {
			params["vpn-gateway-id"] = v.(string)
		}

		if v, ok := d.GetOk("vpc_id"); ok && *gateway.Type != "CCN" {
			params["vpc-id"] = v.(string)
		}

		if v, ok := d.GetOk("customer_gateway_id"); ok {
			params["customer-gateway-id"] = v.(string)
		}

		for k, v := range params {
			filter := &vpc.Filter{
				Name:   helper.String(k),
				Values: []*string{helper.String(v)},
			}

			idRequest.Filters = append(idRequest.Filters, filter)
		}

		offset := uint64(0)
		idRequest.Offset = &offset

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeVpnConnections(idRequest)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, idRequest.GetAction(), idRequest.ToJsonString(), e.Error())
				return tccommon.RetryError(e, tccommon.InternalError)
			} else {
				if len(result.Response.VpnConnectionSet) == 0 || *result.Response.VpnConnectionSet[0].VpnConnectionId == "" {
					return resource.RetryableError(fmt.Errorf("Id is creating, wait..."))
				} else {
					vpnConnectionId = *result.Response.VpnConnectionSet[0].VpnConnectionId
					return nil
				}
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s create VPN connection failed, reason:%s\n", logId, err.Error())
			return err
		}
	} else {
		vpnConnectionId = *response.Response.VpnConnection.VpnConnectionId
	}

	d.SetId(vpnConnectionId)
	// must wait for finishing creating connection
	statRequest := vpc.NewDescribeVpnConnectionsRequest()
	statRequest.VpnConnectionIds = []*string{&vpnConnectionId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeVpnConnections(statRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, statRequest.GetAction(), statRequest.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			//if not, quit
			if len(result.Response.VpnConnectionSet) != 1 {
				return resource.NonRetryableError(fmt.Errorf("creating error"))
			} else {
				if *result.Response.VpnConnectionSet[0].State == svcvpc.VPN_STATE_AVAILABLE {
					return nil
				} else {
					return resource.RetryableError(fmt.Errorf("State is not available: %s, wait for state to be AVAILABLE.", *result.Response.VpnConnectionSet[0].State))
				}
			}
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s create VPN connection failed, reason:%s\n", logId, err.Error())
		return err
	}

	//modify tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("vpc", "vpnx", region, vpnConnectionId)

		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpnConnectionRead(d, meta)
}

func resourceTencentCloudVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_connection.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	connectionId := d.Id()
	request := vpc.NewDescribeVpnConnectionsRequest()
	request.VpnConnectionIds = []*string{&connectionId}
	var response *vpc.DescribeVpnConnectionsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeVpnConnections(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && ee.Code == svcvpc.VPCNotFound {
				return nil
			}
			return tccommon.RetryError(e)
		}

		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read VPN connection failed, reason:%s\n", logId, err.Error())
		return err
	}
	if response == nil || response.Response == nil || len(response.Response.VpnConnectionSet) < 1 {
		d.SetId("")
		return nil
	}

	connection := response.Response.VpnConnectionSet[0]
	_ = d.Set("name", *connection.VpnConnectionName)
	_ = d.Set("create_time", *connection.CreatedTime)
	_ = d.Set("vpn_gateway_id", *connection.VpnGatewayId)

	// get vpngw type
	has, gateway, err := service.DescribeVpngwById(ctx, d.Get("vpn_gateway_id").(string))
	if err != nil {
		log.Printf("[CRITAL]%s read vpn_geteway failed, reason:%s\n", logId, err.Error())
		return err
	}
	if !has {
		return fmt.Errorf("[CRITAL] vpn_gateway_id %s doesn't exist", d.Get("vpn_gateway_id").(string))
	}

	if *gateway.Type != "CCN" {
		_ = d.Set("vpc_id", *connection.VpcId)
	}

	_ = d.Set("is_ccn_type", *gateway.Type == "CCN")
	_ = d.Set("customer_gateway_id", *connection.CustomerGatewayId)
	_ = d.Set("pre_share_key", *connection.PreShareKey)
	//set up SPD
	if *connection.RouteType != svcvpc.ROUTE_TYPE_STATIC_ROUTE {
		_ = d.Set("security_group_policy", svcvpc.FlattenVpnSPDList(connection.SecurityPolicyDatabaseSet))
	}

	//set up IKE
	_ = d.Set("ike_proto_encry_algorithm", *connection.IKEOptionsSpecification.PropoEncryAlgorithm)
	_ = d.Set("ike_proto_authen_algorithm", *connection.IKEOptionsSpecification.PropoAuthenAlgorithm)
	_ = d.Set("ike_exchange_mode", *connection.IKEOptionsSpecification.ExchangeMode)
	_ = d.Set("ike_local_identity", *connection.IKEOptionsSpecification.LocalIdentity)
	_ = d.Set("ike_remote_identity", *connection.IKEOptionsSpecification.RemoteIdentity)
	//optional
	if connection.IKEOptionsSpecification.LocalAddress != nil {
		_ = d.Set("ike_local_address", *connection.IKEOptionsSpecification.LocalAddress)
	}
	if connection.IKEOptionsSpecification.RemoteAddress != nil {
		_ = d.Set("ike_remote_address", *connection.IKEOptionsSpecification.RemoteAddress)
	}
	if connection.IKEOptionsSpecification.LocalFqdnName != nil {
		_ = d.Set("ike_local_fqdn_name", *connection.IKEOptionsSpecification.LocalFqdnName)
	}
	if connection.IKEOptionsSpecification.RemoteFqdnName != nil {
		_ = d.Set("ike_remote_fqdn_name", *connection.IKEOptionsSpecification.RemoteFqdnName)
	}
	_ = d.Set("ike_dh_group_name", *connection.IKEOptionsSpecification.DhGroupName)
	_ = d.Set("ike_sa_lifetime_seconds", int(*connection.IKEOptionsSpecification.IKESaLifetimeSeconds))
	_ = d.Set("ike_version", *connection.IKEOptionsSpecification.IKEVersion)

	//set up ipsec
	_ = d.Set("ipsec_encrypt_algorithm", *connection.IPSECOptionsSpecification.EncryptAlgorithm)
	_ = d.Set("ipsec_integrity_algorithm", *connection.IPSECOptionsSpecification.IntegrityAlgorith)
	_ = d.Set("ipsec_sa_lifetime_seconds", int(*connection.IPSECOptionsSpecification.IPSECSaLifetimeSeconds))
	_ = d.Set("ipsec_pfs_dh_group", *connection.IPSECOptionsSpecification.PfsDhGroup)
	_ = d.Set("ipsec_sa_lifetime_traffic", int(*connection.IPSECOptionsSpecification.IPSECSaLifetimeTraffic))

	//to be add
	_ = d.Set("state", *connection.State)
	_ = d.Set("net_status", *connection.NetStatus)
	_ = d.Set("vpn_proto", *connection.VpnProto)
	_ = d.Set("encrypt_proto", *connection.EncryptProto)
	_ = d.Set("route_type", *connection.RouteType)
	_ = d.Set("enable_health_check", *connection.EnableHealthCheck)
	_ = d.Set("health_check_local_ip", *connection.HealthCheckLocalIp)
	_ = d.Set("health_check_remote_ip", *connection.HealthCheckRemoteIp)
	// dpd
	_ = d.Set("dpd_enable", *connection.DpdEnable)
	if *connection.DpdTimeout != "" {
		dpdTimeoutInt, err := strconv.Atoi(*connection.DpdTimeout)
		if err != nil {
			return err
		}
		_ = d.Set("dpd_timeout", dpdTimeoutInt)
	}

	if connection.NegotiationType != nil {
		_ = d.Set("negotiation_type", *connection.NegotiationType)
	}

	_ = d.Set("dpd_action", *connection.DpdAction)

	if connection.BgpConfig != nil {
		tmpList := make([]map[string]interface{}, 0)
		dMap := make(map[string]interface{})
		if connection.BgpConfig.TunnelCidr != nil {
			dMap["tunnel_cidr"] = *connection.BgpConfig.TunnelCidr
		}

		if connection.BgpConfig.LocalBgpIp != nil {
			dMap["local_bgp_ip"] = *connection.BgpConfig.LocalBgpIp
		}

		if connection.BgpConfig.RemoteBgpIp != nil {
			dMap["remote_bgp_ip"] = *connection.BgpConfig.RemoteBgpIp
		}

		tmpList = append(tmpList, dMap)
		_ = d.Set("bgp_config", tmpList)
	}

	if connection.HealthCheckConfig != nil {
		tmpList := make([]map[string]interface{}, 0)
		dMap := make(map[string]interface{})
		if connection.HealthCheckConfig.ProbeType != nil {
			dMap["probe_type"] = *connection.HealthCheckConfig.ProbeType
		}

		if connection.HealthCheckConfig.ProbeInterval != nil {
			dMap["probe_interval"] = *connection.HealthCheckConfig.ProbeInterval
		}

		if connection.HealthCheckConfig.ProbeThreshold != nil {
			dMap["probe_threshold"] = *connection.HealthCheckConfig.ProbeThreshold
		}

		if connection.HealthCheckConfig.ProbeTimeout != nil {
			dMap["probe_timeout"] = *connection.HealthCheckConfig.ProbeTimeout
		}

		tmpList = append(tmpList, dMap)
		_ = d.Set("health_check_config", tmpList)
	}

	//tags
	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "vpnx", region, connectionId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_connection.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	d.Partial(true)
	connectionId := d.Id()
	request := vpc.NewModifyVpnConnectionAttributeRequest()
	request.VpnConnectionId = &connectionId
	changeFlag := false
	if d.HasChange("name") {
		request.VpnConnectionName = helper.String(d.Get("name").(string))
		changeFlag = true
	}
	if d.HasChange("pre_share_key") {
		request.PreShareKey = helper.String(d.Get("pre_share_key").(string))
		changeFlag = true
	}
	//healthcheck
	if d.HasChange("enable_health_check") {
		request.EnableHealthCheck = helper.Bool(d.Get("enable_health_check").(bool))
		changeFlag = true
	}
	if d.HasChange("health_check_local_ip") {
		request.HealthCheckLocalIp = helper.String(d.Get("health_check_local_ip").(string))
		changeFlag = true
	}
	if d.HasChange("health_check_remote_ip") {
		request.HealthCheckRemoteIp = helper.String(d.Get("health_check_remote_ip").(string))
		changeFlag = true
	}

	//set up  SecurityPolicyDatabases
	if d.HasChange("security_group_policy") {
		if v, ok := d.GetOk("security_group_policy"); ok {
			sgps := v.(*schema.Set).List()
			request.SecurityPolicyDatabases = make([]*vpc.SecurityPolicyDatabase, 0, len(sgps))
			for _, v := range sgps {
				m := v.(map[string]interface{})
				var sgp vpc.SecurityPolicyDatabase
				local := m["local_cidr_block"].(string)
				sgp.LocalCidrBlock = &local
				// list
				remoteCidrBlocks := m["remote_cidr_block"].(*schema.Set).List()
				for _, vv := range remoteCidrBlocks {
					remoteCidrBlock := vv.(string)
					sgp.RemoteCidrBlock = append(sgp.RemoteCidrBlock, &remoteCidrBlock)
				}
				request.SecurityPolicyDatabases = append(request.SecurityPolicyDatabases, &sgp)
			}
			changeFlag = true
		}
	}

	if d.HasChange("dpd_enable") {
		request.DpdEnable = helper.IntInt64(d.Get("dpd_enable").(int))
		changeFlag = true
	}
	if d.HasChange("dpd_timeout") {
		if v, ok := d.GetOk("dpd_timeout"); ok {
			request.DpdTimeout = helper.String(strconv.Itoa(v.(int)))
			changeFlag = true
		}
	}
	if d.HasChange("dpd_action") {
		if v, ok := d.GetOk("dpd_action"); ok {
			request.DpdAction = helper.String(v.(string))
			changeFlag = true
		}
	}

	ikeChangeKeySet := map[string]bool{
		"ike_proto_encry_algorithm":  false,
		"ike_proto_authen_algorithm": false,
		"ike_exchange_mode":          false,
		"ike_local_identity":         false,
		"ike_remote_identity":        false,
		"ike_local_address":          false,
		"ike_local_fqdn_name":        false,
		"ike_remote_address":         false,
		"ike_remote_fqdn_name":       false,
		"ike_sa_lifetime_seconds":    false,
		"ike_dh_group_name":          false,
		"ike_version":                false,
	}
	ikeChangeFlag := false
	for key := range ikeChangeKeySet {
		if d.HasChange(key) {
			ikeChangeFlag = true
			ikeChangeKeySet[key] = true
		}
	}
	if ikeChangeFlag {
		//set up IKEOptionsSpecification
		var ikeOptionsSpecification vpc.IKEOptionsSpecification
		ikeOptionsSpecification.PropoEncryAlgorithm = helper.String(d.Get("ike_proto_encry_algorithm").(string))
		ikeOptionsSpecification.PropoAuthenAlgorithm = helper.String(d.Get("ike_proto_authen_algorithm").(string))
		ikeOptionsSpecification.ExchangeMode = helper.String(d.Get("ike_exchange_mode").(string))
		ikeOptionsSpecification.LocalIdentity = helper.String(d.Get("ike_local_identity").(string))
		ikeOptionsSpecification.RemoteIdentity = helper.String(d.Get("ike_remote_identity").(string))
		if *ikeOptionsSpecification.LocalIdentity == svcvpc.VPN_IKE_IDENTITY_ADDRESS {
			if v, ok := d.GetOk("ike_local_address"); ok {
				ikeOptionsSpecification.LocalAddress = helper.String(v.(string))
			} else {
				return fmt.Errorf("ike_local_address need to be set when ike_local_identity is `ADDRESS`.")
			}
		} else {
			if v, ok := d.GetOk("ike_local_fqdn_name"); ok {
				ikeOptionsSpecification.LocalFqdnName = helper.String(v.(string))
			} else {
				return fmt.Errorf("ike_local_fqdn_name need to be set when ike_local_identity is `FQDN`")
			}
		}
		if *ikeOptionsSpecification.LocalIdentity == svcvpc.VPN_IKE_IDENTITY_ADDRESS {
			if v, ok := d.GetOk("ike_remote_address"); ok {
				ikeOptionsSpecification.RemoteAddress = helper.String(v.(string))
			} else {
				return fmt.Errorf("ike_remote_address need to be set when ike_remote_identity is `ADDRESS`.")
			}
		} else {
			if v, ok := d.GetOk("ike_remote_fqdn_name"); ok {
				ikeOptionsSpecification.RemoteFqdnName = helper.String(v.(string))
			} else {
				return fmt.Errorf("ike_remote_fqdn_name need to be set when ike_remote_identity is `FQDN`")
			}
		}

		ikeOptionsSpecification.DhGroupName = helper.String(d.Get("ike_dh_group_name").(string))
		saLifetime := d.Get("ike_sa_lifetime_seconds").(int)
		saLifetime64 := uint64(saLifetime)
		ikeOptionsSpecification.IKESaLifetimeSeconds = &saLifetime64
		ikeOptionsSpecification.IKEVersion = helper.String(d.Get("ike_version").(string))
		request.IKEOptionsSpecification = &ikeOptionsSpecification
		changeFlag = true
	}
	//set up IPSECOptionsSpecification
	ipsecChangeKeySet := map[string]bool{
		"ipsec_encrypt_algorithm":   false,
		"ipsec_integrity_algorithm": false,
		"ipsec_sa_lifetime_seconds": false,
		"ipsec_pfs_dh_group":        false,
		"ipsec_sa_lifetime_traffic": false}
	ipsecChangeFlag := false
	for key := range ipsecChangeKeySet {
		if d.HasChange(key) {
			ipsecChangeFlag = true
			ipsecChangeKeySet[key] = true
		}
	}
	if ipsecChangeFlag {
		var ipsecOptionsSpecification vpc.IPSECOptionsSpecification
		ipsecOptionsSpecification.EncryptAlgorithm = helper.String(d.Get("ipsec_encrypt_algorithm").(string))
		ipsecOptionsSpecification.IntegrityAlgorith = helper.String(d.Get("ipsec_integrity_algorithm").(string))
		ipsecSaLifetimeSeconds := d.Get("ipsec_sa_lifetime_seconds").(int)
		ipsecSaLifetimeSeconds64 := uint64(ipsecSaLifetimeSeconds)
		ipsecOptionsSpecification.IPSECSaLifetimeSeconds = &ipsecSaLifetimeSeconds64
		ipsecOptionsSpecification.PfsDhGroup = helper.String(d.Get("ipsec_pfs_dh_group").(string))
		ipsecSaLifetimeTraffic := d.Get("ipsec_sa_lifetime_traffic").(int)
		ipsecSaLifetimeTraffic64 := uint64(ipsecSaLifetimeTraffic)
		ipsecOptionsSpecification.IPSECSaLifetimeTraffic = &ipsecSaLifetimeTraffic64
		request.IPSECOptionsSpecification = &ipsecOptionsSpecification
		changeFlag = true
	}

	if d.HasChange("negotiation_type") {
		if v, ok := d.GetOk("negotiation_type"); ok {
			request.NegotiationType = helper.String(v.(string))
		}
	}

	if d.HasChange("health_check_config") {
		if v, ok := d.GetOk("health_check_config"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				healthCheckConfig := vpc.HealthCheckConfig{}
				if v, ok := dMap["probe_type"]; ok {
					healthCheckConfig.ProbeType = helper.String(v.(string))
				}

				if v, ok := dMap["probe_interval"]; ok {
					healthCheckConfig.ProbeInterval = helper.IntInt64(v.(int))
				}

				if v, ok := dMap["probe_threshold"]; ok {
					healthCheckConfig.ProbeThreshold = helper.IntInt64(v.(int))
				}

				if v, ok := dMap["probe_timeout"]; ok {
					healthCheckConfig.ProbeTimeout = helper.IntInt64(v.(int))
				}

				request.HealthCheckConfig = &healthCheckConfig
			}

			changeFlag = true
		}
	}

	if changeFlag {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyVpnConnectionAttribute(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify VPN connection failed, reason:%s\n", logId, err.Error())
			return err
		}
	}
	time.Sleep(3 * time.Minute)

	//tag
	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("vpc", "vpnx", region, connectionId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}
	d.Partial(false)

	return resourceTencentCloudVpnConnectionRead(d, meta)
}

func resourceTencentCloudVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_connection.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	connectionId := d.Id()
	vpnGatewayId := d.Get("vpn_gateway_id").(string)
	//test when tunneling exists, delete may cause a fault, see if sdk returns that error or not
	request := vpc.NewDeleteVpnConnectionRequest()
	request.VpnConnectionId = &connectionId
	request.VpnGatewayId = &vpnGatewayId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteVpnConnection(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.GetCode() == "UnsupportedOperation.InvalidState" {
					return resource.RetryableError(fmt.Errorf("state is not ready, wait to be `AVAILABLE`."))
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete VPN connection failed, reason:%s\n", logId, err.Error())
		return err
	}
	//to get the status of vpn connection
	statRequest := vpc.NewDescribeVpnConnectionsRequest()
	statRequest.VpnConnectionIds = []*string{&connectionId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeVpnConnections(statRequest)
		if e != nil {
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return tccommon.RetryError(e)
			}
			if ee.Code == svcvpc.VPCNotFound {
				log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return nil
			} else {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
		} else {
			//if not, quit
			if len(result.Response.VpnConnectionSet) == 0 {
				return nil
			}
			//else consider delete fail
			return resource.RetryableError(fmt.Errorf("describe retry"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete VPN connection failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
