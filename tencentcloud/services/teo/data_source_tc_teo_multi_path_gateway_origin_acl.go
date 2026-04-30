package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoMultiPathGatewayOriginAcl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoMultiPathGatewayOriginAclRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone ID.",
			},

			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Gateway ID.",
			},

			"multi_path_gateway_origin_acl_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Multi-path gateway origin ACL info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"multi_path_gateway_current_origin_acl": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Currently effective origin ACLs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entire_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "IP CIDR details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 CIDR list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 CIDR list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version number.",
									},
									"is_planed": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether the update to the latest origin IP CIDR has been confirmed.",
									},
								},
							},
						},
						"multi_path_gateway_next_origin_acl": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Next version origin ACLs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version number.",
									},
									"entire_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "IP CIDR details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 CIDR list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 CIDR list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"added_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Added IP CIDRs compared to current.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 CIDR list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 CIDR list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"removed_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Removed IP CIDRs compared to current.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 CIDR list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 CIDR list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"no_change_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Unchanged IP CIDRs compared to current.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 CIDR list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 CIDR list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoMultiPathGatewayOriginAclRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_multi_path_gateway_origin_acl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(nil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId    string
		gatewayId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
		zoneId = v.(string)
	}
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
		gatewayId = v.(string)
	}

	var respData *teov20220901.DescribeMultiPathGatewayOriginACLResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoMultiPathGatewayOriginAclByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	originACLInfoMap := map[string]interface{}{}
	if respData.MultiPathGatewayOriginACLInfo != nil {
		currentOriginACLMap := map[string]interface{}{}
		if respData.MultiPathGatewayOriginACLInfo.MultiPathGatewayCurrentOriginACL != nil {
			currentACL := respData.MultiPathGatewayOriginACLInfo.MultiPathGatewayCurrentOriginACL
			entireAddressesMap := map[string]interface{}{}
			if currentACL.EntireAddresses != nil {
				if currentACL.EntireAddresses.IPv4 != nil {
					entireAddressesMap["ipv4"] = currentACL.EntireAddresses.IPv4
				}
				if currentACL.EntireAddresses.IPv6 != nil {
					entireAddressesMap["ipv6"] = currentACL.EntireAddresses.IPv6
				}
				currentOriginACLMap["entire_addresses"] = []interface{}{entireAddressesMap}
			}

			if currentACL.Version != nil {
				currentOriginACLMap["version"] = helper.Int64ToStr(*currentACL.Version)
			}

			if currentACL.IsPlaned != nil {
				currentOriginACLMap["is_planed"] = currentACL.IsPlaned
			}

			originACLInfoMap["multi_path_gateway_current_origin_acl"] = []interface{}{currentOriginACLMap}
		}

		nextOriginACLMap := map[string]interface{}{}
		if respData.MultiPathGatewayOriginACLInfo.MultiPathGatewayNextOriginACL != nil {
			nextACL := respData.MultiPathGatewayOriginACLInfo.MultiPathGatewayNextOriginACL
			if nextACL.Version != nil {
				nextOriginACLMap["version"] = helper.Int64ToStr(*nextACL.Version)
			}

			entireAddressesMap := map[string]interface{}{}
			if nextACL.EntireAddresses != nil {
				if nextACL.EntireAddresses.IPv4 != nil {
					entireAddressesMap["ipv4"] = nextACL.EntireAddresses.IPv4
				}
				if nextACL.EntireAddresses.IPv6 != nil {
					entireAddressesMap["ipv6"] = nextACL.EntireAddresses.IPv6
				}
				nextOriginACLMap["entire_addresses"] = []interface{}{entireAddressesMap}
			}

			addedAddressesMap := map[string]interface{}{}
			if nextACL.AddedAddresses != nil {
				if nextACL.AddedAddresses.IPv4 != nil {
					addedAddressesMap["ipv4"] = nextACL.AddedAddresses.IPv4
				}
				if nextACL.AddedAddresses.IPv6 != nil {
					addedAddressesMap["ipv6"] = nextACL.AddedAddresses.IPv6
				}
				nextOriginACLMap["added_addresses"] = []interface{}{addedAddressesMap}
			}

			removedAddressesMap := map[string]interface{}{}
			if nextACL.RemovedAddresses != nil {
				if nextACL.RemovedAddresses.IPv4 != nil {
					removedAddressesMap["ipv4"] = nextACL.RemovedAddresses.IPv4
				}
				if nextACL.RemovedAddresses.IPv6 != nil {
					removedAddressesMap["ipv6"] = nextACL.RemovedAddresses.IPv6
				}
				nextOriginACLMap["removed_addresses"] = []interface{}{removedAddressesMap}
			}

			noChangeAddressesMap := map[string]interface{}{}
			if nextACL.NoChangeAddresses != nil {
				if nextACL.NoChangeAddresses.IPv4 != nil {
					noChangeAddressesMap["ipv4"] = nextACL.NoChangeAddresses.IPv4
				}
				if nextACL.NoChangeAddresses.IPv6 != nil {
					noChangeAddressesMap["ipv6"] = nextACL.NoChangeAddresses.IPv6
				}
				nextOriginACLMap["no_change_addresses"] = []interface{}{noChangeAddressesMap}
			}

			originACLInfoMap["multi_path_gateway_next_origin_acl"] = []interface{}{nextOriginACLMap}
		}

		_ = d.Set("multi_path_gateway_origin_acl_info", []interface{}{originACLInfoMap})
	}

	d.SetId(zoneId + tccommon.FILED_SP + gatewayId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), originACLInfoMap); e != nil {
			return e
		}
	}

	return nil
}
