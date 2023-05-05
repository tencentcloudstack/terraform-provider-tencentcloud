/*
Provides a resource to create a VPN connection.

Example Usage

```hcl
resource "tencentcloud_vpn_connection" "foo" {
  name                       = "vpn_connection_test"
  vpc_id                     = "vpc-dk8zmwuf"
  vpn_gateway_id             = "vpngw-8ccsnclt"
  customer_gateway_id        = "cgw-xfqag"
  pre_share_key              = "testt"
  ike_proto_encry_algorithm  = "3DES-CBC"
  ike_proto_authen_algorithm = "SHA"
  ike_local_identity         = "ADDRESS"
  ike_exchange_mode          = "AGGRESSIVE"
  ike_local_address          = "1.1.1.1"
  ike_remote_identity        = "ADDRESS"
  ike_remote_address         = "2.2.2.2"
  ike_dh_group_name          = "GROUP2"
  ike_sa_lifetime_seconds    = 86401
  ipsec_encrypt_algorithm    = "3DES-CBC"
  ipsec_integrity_algorithm  = "SHA1"
  ipsec_sa_lifetime_seconds  = 7200
  ipsec_pfs_dh_group         = "NULL"
  ipsec_sa_lifetime_traffic  = 2570

  security_group_policy {
    local_cidr_block  = "172.16.0.0/16"
    remote_cidr_block = ["2.2.2.0/26", ]
  }
  tags = {
    test = "testt"
  }
}
```

Import

VPN connection can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_connection.foo vpnx-nadifg3s
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpnConnection() *schema.Resource {
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
				ValidateFunc: validateStringLengthInRange(1, 60),
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
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Security group policy of the VPN connection.",
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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      VPN_IKE_PROPO_ENCRY_ALGORITHM_3DESCBC,
				ValidateFunc: validateAllowedStringValue(VPN_IKE_PROPO_ENCRY_ALGORITHM),
				Description:  "Proto encrypt algorithm of the IKE operation specification. Valid values: `3DES-CBC`, `AES-CBC-128`, `AES-CBC-128`, `AES-CBC-256`, `DES-CBC`. Default value is `3DES-CBC`.",
			},
			"ike_proto_authen_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      VPN_IKE_PROPO_AUTHEN_ALGORITHM_MD5,
				ValidateFunc: validateAllowedStringValue(VPN_IKE_PROPO_AUTHEN_ALGORITHM),
				Description:  "Proto authenticate algorithm of the IKE operation specification. Valid values: `MD5`, `SHA`, `SHA-256`. Default Value is `MD5`.",
			},
			"ike_exchange_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      VPN_IKE_EXCHANGE_MODE_MAIN,
				ValidateFunc: validateAllowedStringValue(VPN_IKE_EXCHANGE_MODE),
				Description:  "Exchange mode of the IKE operation specification. Valid values: `AGGRESSIVE`, `MAIN`. Default value is `MAIN`.",
			},
			"ike_local_identity": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      VPN_IKE_IDENTITY_ADDRESS,
				ValidateFunc: validateAllowedStringValue(VPN_IKE_IDENTITY),
				Description:  "Local identity way of IKE operation specification. Valid values: `ADDRESS`, `FQDN`. Default value is `ADDRESS`.",
			},
			"ike_remote_identity": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      VPN_IKE_IDENTITY_ADDRESS,
				ValidateFunc: validateAllowedStringValue(VPN_IKE_IDENTITY),
				Description:  "Remote identity way of IKE operation specification. Valid values: `ADDRESS`, `FQDN`. Default value is `ADDRESS`.",
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
				ConflictsWith: []string{"ike_local_address"},
				Description:   "Local FQDN name of the IKE operation specification.",
			},
			"ike_remote_fqdn_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ike_remote_address"},
				Description:   "Remote FQDN name of the IKE operation specification.",
			},
			"ike_dh_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      VPN_IKE_DH_GROUP_NAME_GROUP1,
				ValidateFunc: validateAllowedStringValue(VPN_IKE_DH_GROUP_NAME),
				Description:  "DH group name of the IKE operation specification. Valid values: `GROUP1`, `GROUP2`, `GROUP5`, `GROUP14`, `GROUP24`. Default value is `GROUP1`.",
			},
			"ike_sa_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      86400,
				ValidateFunc: validateIntegerInRange(60, 604800),
				Description:  "SA lifetime of the IKE operation specification, unit is `second`. The value ranges from 60 to 604800. Default value is 86400 seconds.",
			},
			"ike_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "IKEV1",
				Description: "Version of the IKE operation specification. Default value is `IKEV1`.",
			},
			"ipsec_encrypt_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      VPN_IPSEC_ENCRY_ALGORITHM_3DESCBC,
				ValidateFunc: validateAllowedStringValue(VPN_IPSEC_ENCRY_ALGORITHM),
				Description:  "Encrypt algorithm of the IPSEC operation specification. Valid values: `3DES-CBC`, `AES-CBC-128`, `AES-CBC-128`, `AES-CBC-256`, `DES-CBC`. Default value is `3DES-CBC`.",
			},
			"ipsec_integrity_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      VPN_IPSEC_INTEGRITY_ALGORITHM_MD5,
				ValidateFunc: validateAllowedStringValue(VPN_IPSEC_INTEGRITY_ALGORITHM),
				Description:  "Integrity algorithm of the IPSEC operation specification. Valid values: `SHA1`, `MD5`, `SHA-256`. Default value is `MD5`.",
			},
			"ipsec_sa_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				ValidateFunc: validateIntegerInRange(180, 604800),
				Description:  "SA lifetime of the IPSEC operation specification, unit is second. Valid value ranges: [180~604800]. Default value is 3600 seconds.",
			},
			"ipsec_pfs_dh_group": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NULL",
				ValidateFunc: validateAllowedStringValue(VPN_IPSEC_PFS_DH_GROUP_NAME),
				Description:  "PFS DH group. Valid value: `GROUP1`, `GROUP2`, `GROUP5`, `GROUP14`, `GROUP24`, `NULL`. Default value is `NULL`.",
			},
			"ipsec_sa_lifetime_traffic": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1843200,
				ValidateFunc: validateIntegerMin(2560),
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
				ValidateFunc: validateIntegerInRange(0, 1),
				Description:  "Specifies whether to enable DPD. Valid values: 0 (disable) and 1 (enable).",
			},
			"dpd_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(30, 60),
				Description:  "DPD timeout period.Valid value ranges: [30~60], Default: 30; unit: second. If the request is not responded within this period, the peer end is considered not exists. This parameter is valid when the value of DpdEnable is 1.",
			},
			"dpd_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue(DPD_ACTIONS),
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Route type of the VPN connection.",
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
	defer logElapsed("resource.tencentcloud_vpn_connection.create")()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
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
	//set up  SecurityPolicyDatabases

	sgps := d.Get("security_group_policy").(*schema.Set).List()
	if len(sgps) < 1 {
		return fmt.Errorf("Para `security_group_policy` should be set at least one.")
	}

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

	//set up IKEOptionsSpecification
	var ikeOptionsSpecification vpc.IKEOptionsSpecification
	ikeOptionsSpecification.PropoEncryAlgorithm = helper.String(d.Get("ike_proto_encry_algorithm").(string))
	ikeOptionsSpecification.PropoAuthenAlgorithm = helper.String(d.Get("ike_proto_authen_algorithm").(string))
	ikeOptionsSpecification.ExchangeMode = helper.String(d.Get("ike_exchange_mode").(string))
	ikeOptionsSpecification.LocalIdentity = helper.String(d.Get("ike_local_identity").(string))
	ikeOptionsSpecification.RemoteIdentity = helper.String(d.Get("ike_remote_identity").(string))
	if *ikeOptionsSpecification.LocalIdentity == VPN_IKE_IDENTITY_ADDRESS {
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
	if *ikeOptionsSpecification.LocalIdentity == VPN_IKE_IDENTITY_ADDRESS {
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
	var response *vpc.CreateVpnConnectionResponse
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateVpnConnection(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create VPN connection failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response.VpnConnection == nil {
		return fmt.Errorf("VPN connection id is nil")
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

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnConnections(idRequest)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, idRequest.GetAction(), idRequest.ToJsonString(), e.Error())
				return retryError(e, InternalError)
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
	}

	if vpnConnectionId == "" {
		return fmt.Errorf("VPN connection id is nil")
	}

	d.SetId(vpnConnectionId)
	// must wait for finishing creating connection
	statRequest := vpc.NewDescribeVpnConnectionsRequest()
	statRequest.VpnConnectionIds = []*string{&vpnConnectionId}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnConnections(statRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, statRequest.GetAction(), statRequest.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			//if not, quit
			if len(result.Response.VpnConnectionSet) != 1 {
				return resource.NonRetryableError(fmt.Errorf("creating error"))
			} else {
				if *result.Response.VpnConnectionSet[0].State == VPN_STATE_AVAILABLE {
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
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}

		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("vpc", "vpnx", region, vpnConnectionId)

		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpnConnectionRead(d, meta)
}

func resourceTencentCloudVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_connection.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	connectionId := d.Id()
	request := vpc.NewDescribeVpnConnectionsRequest()
	request.VpnConnectionIds = []*string{&connectionId}
	var response *vpc.DescribeVpnConnectionsResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnConnections(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && ee.Code == VPCNotFound {
				return nil
			}
			return retryError(e)
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
	_ = d.Set("security_group_policy", flattenVpnSPDList(connection.SecurityPolicyDatabaseSet))

	//set up IKE
	_ = d.Set("ike_proto_encry_algorithm", *connection.IKEOptionsSpecification.PropoEncryAlgorithm)
	_ = d.Set("ike_proto_authen_algorithm", *connection.IKEOptionsSpecification.PropoAuthenAlgorithm)
	_ = d.Set("ike_exchange_mode", *connection.IKEOptionsSpecification.ExchangeMode)
	_ = d.Set("ike_local_identity", *connection.IKEOptionsSpecification.LocalIdentity)
	_ = d.Set("ike_remote_idetity", *connection.IKEOptionsSpecification.RemoteIdentity)
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

	_ = d.Set("dpd_action", *connection.DpdAction)

	//tags
	tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "vpnx", region, connectionId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_connection.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
		sgps := d.Get("security_group_policy").(*schema.Set).List()
		if len(sgps) < 1 {
			return fmt.Errorf("Para `security_group_policy` should be set at least one.")
		}
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
		if *ikeOptionsSpecification.LocalIdentity == VPN_IKE_IDENTITY_ADDRESS {
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
		if *ikeOptionsSpecification.LocalIdentity == VPN_IKE_IDENTITY_ADDRESS {
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
	if changeFlag {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyVpnConnectionAttribute(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
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
		replaceTags, deleteTags := diffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := TagService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("vpc", "vpnx", region, connectionId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}
	d.Partial(false)

	return resourceTencentCloudVpnConnectionRead(d, meta)
}

func resourceTencentCloudVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_connection.delete")()

	logId := getLogId(contextNil)

	connectionId := d.Id()
	vpnGatewayId := d.Get("vpn_gateway_id").(string)
	//test when tunneling exists, delete may cause a fault, see if sdk returns that error or not
	request := vpc.NewDeleteVpnConnectionRequest()
	request.VpnConnectionId = &connectionId
	request.VpnGatewayId = &vpnGatewayId
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DeleteVpnConnection(request)
		if e != nil {
			if ee, ok := e.(*errors.TencentCloudSDKError); ok {
				if ee.GetCode() == "UnsupportedOperation.InvalidState" {
					return resource.RetryableError(fmt.Errorf("state is not ready, wait to be `AVAILABLE`."))
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
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
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnConnections(statRequest)
		if e != nil {
			ee, ok := e.(*errors.TencentCloudSDKError)
			if !ok {
				return retryError(e)
			}
			if ee.Code == VPCNotFound {
				log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return nil
			} else {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
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
