package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoOriginAcl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoOriginAclRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the site ID.",
			},

			"origin_acl_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Describes the binding relationship between the l7 acceleration domain/l4 proxy instance and the origin server IP range.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"l7_hosts": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "The list of L7 accelerated domains that enable the origin ACLs. This field is empty when origin protection is not enabled.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"l4_proxy_ids": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "The list of L4 proxy instances that enable the origin ACLs. This field is empty when origin protection is not enabled.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"current_origin_acl": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Currently effective origin ACLs. This field is empty when origin protection is not enabled.\nNote: This field may return null, which indicates a failure to obtain a valid value.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entire_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "IP range details.\nNote: This field may return null, which indicates a failure to obtain a valid value.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"i_pv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 subnet.",
													Deprecated:  "Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"i_pv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 subnet.",
													Deprecated:  "Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 subnet.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 subnet.",
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
										Description: "Version number.\nNote: This field may return null, which indicates a failure to obtain a valid value.",
									},
									"active_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version effective time in UTC+8, following the date and time format of the ISO 8601 standard.\nNote: This field may return null, which indicates a failure to obtain a valid value.",
									},
									"is_planed": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "This parameter is used to record whether \"I've upgraded to the lastest version\" is completed before the origin ACLs version is effective. valid values:.\n- true: specifies that the version is effective and the update to the latest version is confirmed.\n- false: when the version takes effect, the confirmation of updating to the latest origin ACLs are not completed. The IP range is forcibly updated to the latest version in the backend. When this parameter returns false, please confirm in time whether your origin server firewall configuration has been updated to the latest version to avoid origin-pull failure.\nNote: This field may return null, which indicates a failure to obtain a valid value.",
									},
								},
							},
						},
						"next_origin_acl": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "When the origin ACLs are updated, this field will be returned with the next version's origin IP range to take effect, including a comparison with the current origin IP range. This field is empty if not updated or origin protection is not enabled.\nNote: This field may return null, which indicates a failure to obtain a valid value.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version number.",
									},
									"planned_active_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version effective time, which adopts UTC+8 and follows the date and time format of the ISO 8601 standard.",
									},
									"entire_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "IP range details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"i_pv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 subnet.",
													Deprecated:  "Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"i_pv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 subnet.",
													Deprecated:  "Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 subnet.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 subnet.",
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
										Description: "The latest origin IP range newly-added compared with the origin IP range in CurrentOrginACL.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"i_pv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 subnet.",
													Deprecated:  "Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"i_pv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 subnet.",
													Deprecated:  "Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 subnet.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 subnet.",
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
										Description: "The latest origin IP range deleted compared with the origin IP range in CurrentOrginACL.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"i_pv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 subnet.",
													Deprecated:  "Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"i_pv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 subnet.",
													Deprecated:  "Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 subnet.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 subnet.",
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
										Description: "The latest origin IP range is unchanged compared with the origin IP range in CurrentOrginACL.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"i_pv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 subnet.",
													Deprecated:  "Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"i_pv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 subnet.",
													Deprecated:  "Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 subnet.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 subnet.",
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
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin protection status. Vaild values:\n- online: in effect;\n- offline: disabled;\n- updating: configuration deployment in progress.",
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

func dataSourceTencentCloudTeoOriginAclRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_origin_acl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
		zoneId = v.(string)
	}

	var respData *teov20220901.DescribeOriginACLResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoOriginAclByFilter(ctx, paramMap)
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
	if respData.OriginACLInfo != nil {
		if respData.OriginACLInfo.L7Hosts != nil {
			originACLInfoMap["l7_hosts"] = respData.OriginACLInfo.L7Hosts
		}

		if respData.OriginACLInfo.L4ProxyIds != nil {
			originACLInfoMap["l4_proxy_ids"] = respData.OriginACLInfo.L4ProxyIds
		}

		currentOriginACLMap := map[string]interface{}{}
		if respData.OriginACLInfo.CurrentOriginACL != nil {
			entireAddressesMap := map[string]interface{}{}
			if respData.OriginACLInfo.CurrentOriginACL.EntireAddresses != nil {
				if respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv4 != nil {
					entireAddressesMap["i_pv4"] = respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv4
					entireAddressesMap["ipv4"] = respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv4
				}

				if respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv6 != nil {
					entireAddressesMap["i_pv6"] = respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv6
					entireAddressesMap["ipv6"] = respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv6
				}

				currentOriginACLMap["entire_addresses"] = []interface{}{entireAddressesMap}
			}

			if respData.OriginACLInfo.CurrentOriginACL.Version != nil {
				currentOriginACLMap["version"] = respData.OriginACLInfo.CurrentOriginACL.Version
			}

			if respData.OriginACLInfo.CurrentOriginACL.ActiveTime != nil {
				currentOriginACLMap["active_time"] = respData.OriginACLInfo.CurrentOriginACL.ActiveTime
			}

			if respData.OriginACLInfo.CurrentOriginACL.IsPlaned != nil {
				currentOriginACLMap["is_planed"] = respData.OriginACLInfo.CurrentOriginACL.IsPlaned
			}

			originACLInfoMap["current_origin_acl"] = []interface{}{currentOriginACLMap}
		}

		nextOriginACLMap := map[string]interface{}{}
		if respData.OriginACLInfo.NextOriginACL != nil {
			if respData.OriginACLInfo.NextOriginACL.Version != nil {
				nextOriginACLMap["version"] = respData.OriginACLInfo.NextOriginACL.Version
			}

			if respData.OriginACLInfo.NextOriginACL.PlannedActiveTime != nil {
				nextOriginACLMap["planned_active_time"] = respData.OriginACLInfo.NextOriginACL.PlannedActiveTime
			}

			entireAddressesMap := map[string]interface{}{}
			if respData.OriginACLInfo.NextOriginACL.EntireAddresses != nil {
				if respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv4 != nil {
					entireAddressesMap["i_pv4"] = respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv4
					entireAddressesMap["ipv4"] = respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv4
				}

				if respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv6 != nil {
					entireAddressesMap["i_pv6"] = respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv6
					entireAddressesMap["ipv6"] = respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv6
				}

				nextOriginACLMap["entire_addresses"] = []interface{}{entireAddressesMap}
			}

			addedAddressesMap := map[string]interface{}{}
			if respData.OriginACLInfo.NextOriginACL.AddedAddresses != nil {
				if respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv4 != nil {
					addedAddressesMap["i_pv4"] = respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv4
					addedAddressesMap["ipv4"] = respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv4
				}

				if respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv6 != nil {
					addedAddressesMap["i_pv6"] = respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv6
					addedAddressesMap["ipv6"] = respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv6
				}

				nextOriginACLMap["added_addresses"] = []interface{}{addedAddressesMap}
			}

			removedAddressesMap := map[string]interface{}{}
			if respData.OriginACLInfo.NextOriginACL.RemovedAddresses != nil {
				if respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv4 != nil {
					removedAddressesMap["i_pv4"] = respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv4
					removedAddressesMap["ipv4"] = respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv4
				}

				if respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv6 != nil {
					removedAddressesMap["i_pv6"] = respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv6
					removedAddressesMap["ipv6"] = respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv6
				}

				nextOriginACLMap["removed_addresses"] = []interface{}{removedAddressesMap}
			}

			noChangeAddressesMap := map[string]interface{}{}
			if respData.OriginACLInfo.NextOriginACL.NoChangeAddresses != nil {
				if respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv4 != nil {
					noChangeAddressesMap["i_pv4"] = respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv4
					noChangeAddressesMap["ipv4"] = respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv4
				}

				if respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv6 != nil {
					noChangeAddressesMap["i_pv6"] = respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv6
					noChangeAddressesMap["ipv6"] = respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv6
				}

				nextOriginACLMap["no_change_addresses"] = []interface{}{noChangeAddressesMap}
			}

			originACLInfoMap["next_origin_acl"] = []interface{}{nextOriginACLMap}
		}

		if respData.OriginACLInfo.Status != nil {
			originACLInfoMap["status"] = respData.OriginACLInfo.Status
		}

		_ = d.Set("origin_acl_info", []interface{}{originACLInfoMap})
	}

	d.SetId(zoneId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), originACLInfoMap); e != nil {
			return e
		}
	}

	return nil
}
