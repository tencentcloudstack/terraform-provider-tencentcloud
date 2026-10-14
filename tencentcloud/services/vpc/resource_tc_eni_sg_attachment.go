package vpc

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func mergePrependSecurityGroups(prepend, current []string) []string {
	if len(prepend) == 0 {
		// no prepend requested, keep current unchanged
		out := make([]string, 0, len(current))
		out = append(out, current...)
		return out
	}

	seen := make(map[string]struct{}, len(prepend)+len(current))
	out := make([]string, 0, len(prepend)+len(current))

	// put prepend list first, keep its order
	for _, sg := range prepend {
		if sg == "" {
			continue
		}
		if _, ok := seen[sg]; ok {
			continue
		}
		seen[sg] = struct{}{}
		out = append(out, sg)
	}

	// then append current list excluding those already placed, keep relative order
	for _, sg := range current {
		if sg == "" {
			continue
		}
		if _, ok := seen[sg]; ok {
			continue
		}
		seen[sg] = struct{}{}
		out = append(out, sg)
	}

	return out
}

func ResourceTencentCloudEniSgAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEniSgAttachmentCreate,
		Read:   resourceTencentCloudEniSgAttachmentRead,
		Delete: resourceTencentCloudEniSgAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"network_interface_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems:    1,
				Description: "ENI instance ID. Such as:eni-pxir56ns. It Only support set one eni instance now.",
			},

			"security_group_ids": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"orderly_security_groups"},
				AtLeastOneOf:  []string{"security_group_ids", "orderly_security_groups", "prepend_orderly_security_groups"},
				Description: "Security group instance ID, for example:sg-33ocnj9n, can be obtained through DescribeSecurityGroups. There is a limit of 100 instances per request.",
			},

			"orderly_security_groups": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"security_group_ids"},
				AtLeastOneOf:  []string{"security_group_ids", "orderly_security_groups", "prepend_orderly_security_groups"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "List of security group IDs. This will replace the ENI security group list (full overwrite) and keep the order specified.",
			},

			"prepend_orderly_security_groups": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"security_group_ids", "orderly_security_groups"},
				AtLeastOneOf:  []string{"security_group_ids", "orderly_security_groups", "prepend_orderly_security_groups"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description: "List of security group IDs to be placed in front of the current ENI security group list. If it contains IDs that already exist on the ENI, their positions will be adjusted to the front according to this order; the remaining original security groups will keep their relative order after them.",
			},

			// internal: store original order for rollback on delete when using prepend_orderly_security_groups
			"original_orderly_security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Internal field: original security group order before applying prepend_orderly_security_groups.",
			},
		},
	}
}

func resourceTencentCloudEniSgAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni_sg_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		networkInterfaceId string
		securityGroups     []string
		prependGroups      []string
	)
	if v, ok := d.GetOk("network_interface_ids"); ok {
		networkInterfaceIdsSet := v.(*schema.Set).List()
		for i := range networkInterfaceIdsSet {
			networkInterfaceId = networkInterfaceIdsSet[i].(string)
		}
	}

	// keep order if user specifies it
	if raw, ok := d.GetOk("orderly_security_groups"); ok {
		securityGroups = helper.InterfacesStrings(raw.([]interface{}))
	} else if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroups = helper.InterfacesStrings(v.(*schema.Set).List())
	}

	// Use ModifyNetworkInterfaceAttribute (same path as `tencentcloud_eni`) to make order take effect.
	// AssociateNetworkInterfaceSecurityGroups does not guarantee order semantics.
	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if raw, ok := d.GetOk("prepend_orderly_security_groups"); ok {
		prependGroups = helper.InterfacesStrings(raw.([]interface{}))

		enis, err := service.DescribeEniById(ctx, []string{networkInterfaceId})
		if err != nil {
			log.Printf("[CRITAL]%s read eni before prepend failed, reason:%+v", logId, err)
			return err
		}
		if len(enis) < 1 || enis[0] == nil {
			// ENI not found
			d.SetId("")
			return nil
		}

		current := make([]string, 0, len(enis[0].GroupSet))
		for _, sg := range enis[0].GroupSet {
			if sg == nil {
				continue
			}
			current = append(current, *sg)
		}

		// persist original order for rollback on delete
		_ = d.Set("original_orderly_security_groups", current)

		securityGroups = mergePrependSecurityGroups(prependGroups, current)
	}

	if err := service.ModifyEniAttribute(ctx, networkInterfaceId, nil, nil, securityGroups); err != nil {
		log.Printf("[CRITAL]%s create vpc eniSgAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(networkInterfaceId)

	return resourceTencentCloudEniSgAttachmentRead(d, meta)
}

func resourceTencentCloudEniSgAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni_sg_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	networkInterfaceId := d.Id()

	enis, err := service.DescribeEniById(ctx, []string{networkInterfaceId})
	if err != nil {
		return err
	}

	if len(enis) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcEniSgAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	eni := enis[0]

	_ = d.Set("network_interface_ids", []string{networkInterfaceId})

	sgs := make([]string, 0, len(eni.GroupSet))
	for _, sg := range eni.GroupSet {
		if sg == nil {
			continue
		}
		sgs = append(sgs, *sg)
	}
	_ = d.Set("security_group_ids", sgs)
	_ = d.Set("orderly_security_groups", sgs)

	return nil
}

func resourceTencentCloudEniSgAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni_sg_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	networkInterfaceId := d.Id()

	// If prepend mode was used, rollback to original order to avoid removing pre-existing SGs.
	if _, ok := d.GetOk("prepend_orderly_security_groups"); ok {
		if v, ok2 := d.GetOk("original_orderly_security_groups"); ok2 {
			original := helper.InterfacesStrings(v.([]interface{}))
			if err := service.ModifyEniAttribute(ctx, networkInterfaceId, nil, nil, original); err != nil {
				return err
			}
			return nil
		}
		// fallback: if original not present, do nothing to avoid accidental detach
		return nil
	}

	// legacy delete behavior: detach the specified SGs
	var securityGroupIds []string
	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIds = helper.InterfacesStrings(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("orderly_security_groups"); ok {
		securityGroupIds = helper.InterfacesStrings(v.([]interface{}))
	}

	if err := service.DeleteVpcEniSgAttachmentById(ctx, networkInterfaceId, securityGroupIds); err != nil {
		return err
	}

	return nil
}
