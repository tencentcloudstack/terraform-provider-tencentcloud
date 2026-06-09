package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudTeoConfirmMultiPathGatewayOriginAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoConfirmMultiPathGatewayOriginAclCreate,
		Read:   resourceTencentCloudTeoConfirmMultiPathGatewayOriginAclRead,
		Update: resourceTencentCloudTeoConfirmMultiPathGatewayOriginAclUpdate,
		Delete: resourceTencentCloudTeoConfirmMultiPathGatewayOriginAclDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},

			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Gateway ID.",
			},

			"origin_acl_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Origin ACL version number to confirm.",
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
							Description: "Current active origin ACL info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entire_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Current origin IP segment details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "IPv4 segment list.",
												},
												"ipv6": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "IPv6 segment list.",
												},
											},
										},
									},
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Current version number.",
									},
									"is_planed": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether the update confirmation is completed.",
									},
								},
							},
						},
						"multi_path_gateway_next_origin_acl": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Next version origin ACL info when there is a pending update.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Next version number.",
									},
									"entire_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Next version origin IP segment details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "IPv4 segment list.",
												},
												"ipv6": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "IPv6 segment list.",
												},
											},
										},
									},
									"added_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Added IP segments compared to current ACL.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "IPv4 segment list.",
												},
												"ipv6": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "IPv6 segment list.",
												},
											},
										},
									},
									"removed_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Removed IP segments compared to current ACL.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "IPv4 segment list.",
												},
												"ipv6": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "IPv6 segment list.",
												},
											},
										},
									},
									"no_change_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Unchanged IP segments compared to current ACL.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "IPv4 segment list.",
												},
												"ipv6": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "IPv6 segment list.",
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
	}
}

func resourceTencentCloudTeoConfirmMultiPathGatewayOriginAclCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_confirm_multi_path_gateway_origin_acl.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		zoneId    string
		gatewayId string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
	}

	d.SetId(strings.Join([]string{zoneId, gatewayId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoConfirmMultiPathGatewayOriginAclUpdate(d, meta)
}

func resourceTencentCloudTeoConfirmMultiPathGatewayOriginAclRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_confirm_multi_path_gateway_origin_acl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	gatewayId := idSplit[1]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("gateway_id", gatewayId)

	respData, err := service.DescribeTeoConfirmMultiPathGatewayOriginAclById(ctx, zoneId, gatewayId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_confirm_multi_path_gateway_origin_acl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	aclInfoMap := make(map[string]interface{}, 0)

	if respData.MultiPathGatewayCurrentOriginACL != nil {
		currentAclMap := make(map[string]interface{}, 0)
		currentAcl := respData.MultiPathGatewayCurrentOriginACL

		if currentAcl.EntireAddresses != nil {
			currentAclMap["entire_addresses"] = flattenAddresses(currentAcl.EntireAddresses)
		}
		if currentAcl.Version != nil {
			currentAclMap["version"] = *currentAcl.Version
		}
		if currentAcl.IsPlaned != nil {
			currentAclMap["is_planed"] = *currentAcl.IsPlaned
		}

		aclInfoMap["multi_path_gateway_current_origin_acl"] = []interface{}{currentAclMap}
	}

	if respData.MultiPathGatewayNextOriginACL != nil {
		nextAclMap := make(map[string]interface{}, 0)
		nextAcl := respData.MultiPathGatewayNextOriginACL

		if nextAcl.Version != nil {
			nextAclMap["version"] = *nextAcl.Version
		}
		if nextAcl.EntireAddresses != nil {
			nextAclMap["entire_addresses"] = flattenAddresses(nextAcl.EntireAddresses)
		}
		if nextAcl.AddedAddresses != nil {
			nextAclMap["added_addresses"] = flattenAddresses(nextAcl.AddedAddresses)
		}
		if nextAcl.RemovedAddresses != nil {
			nextAclMap["removed_addresses"] = flattenAddresses(nextAcl.RemovedAddresses)
		}
		if nextAcl.NoChangeAddresses != nil {
			nextAclMap["no_change_addresses"] = flattenAddresses(nextAcl.NoChangeAddresses)
		}

		aclInfoMap["multi_path_gateway_next_origin_acl"] = []interface{}{nextAclMap}
	}

	_ = d.Set("multi_path_gateway_origin_acl_info", []interface{}{aclInfoMap})

	return nil
}

func resourceTencentCloudTeoConfirmMultiPathGatewayOriginAclUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_confirm_multi_path_gateway_origin_acl.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	gatewayId := idSplit[1]

	request := teov20220901.NewConfirmMultiPathGatewayOriginACLRequest()
	request.ZoneId = &zoneId
	request.GatewayId = &gatewayId

	if v, ok := d.GetOkExists("origin_acl_version"); ok {
		originAclVersion := int64(v.(int))
		request.OriginACLVersion = &originAclVersion
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ConfirmMultiPathGatewayOriginACLWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update teo confirm multi path gateway origin acl failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoConfirmMultiPathGatewayOriginAclRead(d, meta)
}

func resourceTencentCloudTeoConfirmMultiPathGatewayOriginAclDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_confirm_multi_path_gateway_origin_acl.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

// flattenAddresses flattens an Addresses struct into a map for Terraform state
func flattenAddresses(addr *teov20220901.Addresses) []interface{} {
	if addr == nil {
		return nil
	}

	m := make(map[string]interface{}, 0)

	if addr.IPv4 != nil {
		ipv4List := make([]string, 0, len(addr.IPv4))
		for _, v := range addr.IPv4 {
			if v != nil {
				ipv4List = append(ipv4List, *v)
			}
		}
		m["ipv4"] = ipv4List
	}

	if addr.IPv6 != nil {
		ipv6List := make([]string, 0, len(addr.IPv6))
		for _, v := range addr.IPv6 {
			if v != nil {
				ipv6List = append(ipv6List, *v)
			}
		}
		m["ipv6"] = ipv6List
	}

	return []interface{}{m}
}
