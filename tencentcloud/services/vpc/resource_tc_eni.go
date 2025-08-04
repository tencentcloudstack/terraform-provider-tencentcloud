package vpc

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"errors"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func eniIpOutputResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Intranet IP.",
			},
			"primary": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the IP is primary.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the IP.",
			},
		},
	}
}

func ResourceTencentCloudEni() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEniCreate,
		Read:   resourceTencentCloudEniRead,
		Update: resourceTencentCloudEniUpdate,
		Delete: resourceTencentCloudEniDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 60),
				Description:  "Name of the ENI, maximum length 60.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the vpc.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the subnet within this vpc.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 60),
				Description:  "Description of the ENI, maximum length 60.",
			},
			"security_groups": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"orderly_security_groups"},
				Deprecated:    "It has been deprecated from version 1.82.15. Use `orderly_security_groups` instead.",
				Description:   "A set of security group IDs.",
			},
			"orderly_security_groups": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"security_groups"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "List of security group IDs.",
			},
			"ipv4s": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"ipv4_count"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Intranet IP.",
						},
						"primary": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates whether the IP is primary.",
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "",
							Description:  "Description of the IP, maximum length 25.",
							ValidateFunc: tccommon.ValidateStringLengthInRange(0, 25),
						},
					},
				},
				MaxItems:    30,
				Description: "Applying for intranet IPv4s collection, conflict with `ipv4_count`. When there are multiple ipv4s, can only be one primary IP, and the maximum length of the array is 30. Each element contains the following attributes:",
			},
			"ipv4_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"ipv4s"},
				ValidateFunc:  tccommon.ValidateIntegerInRange(1, 30),
				Description:   "The number of intranet IPv4s. When it is greater than 1, there is only one primary intranet IP. The others are auxiliary intranet IPs, which conflict with `ipv4s`.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the ENI.",
			},

			// computed
			"mac": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MAC address.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the ENI.",
			},
			"primary": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the IP is primary.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the ENI.",
			},
			"ipv4_info": {
				Type:        schema.TypeList,
				Elem:        eniIpOutputResource(),
				Computed:    true,
				Description: "An information list of IPv4s. Each element contains the following attributes:",
			},
			"cdc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CDC instance ID.",
			},
		},
	}
}

func resourceTencentCloudEniCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	name := d.Get("name").(string)
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)
	desc := d.Get("description").(string)

	var (
		securityGroups []string
		ipv4s          []VpcEniIP
		ipv4Count      *int
		tags           map[string]string
	)

	if raw, ok := d.GetOk("security_groups"); ok {
		securityGroups = helper.InterfacesStrings(raw.(*schema.Set).List())
	}

	if raw, ok := d.GetOk("orderly_security_groups"); ok {
		securityGroups = helper.InterfacesStrings(raw.([]interface{}))
	}

	if raw, ok := d.GetOk("ipv4s"); ok {
		set := raw.(*schema.Set)
		ipv4s = make([]VpcEniIP, 0, set.Len())
		var hasPrimary bool

		for _, v := range set.List() {
			m := v.(map[string]interface{})

			ipStr := m["ip"].(string)
			ip := net.ParseIP(ipStr)
			if ip == nil {
				return fmt.Errorf("ip %s is invalid", ipStr)
			}

			primary := m["primary"].(bool)

			switch {
			case !hasPrimary && primary:
				hasPrimary = true

			case hasPrimary && primary:
				return errors.New("only can have a primary ipv4")
			}

			ipv4 := VpcEniIP{
				ip:      ip,
				primary: primary,
			}

			ipv4.desc = helper.String(m["description"].(string))

			ipv4s = append(ipv4s, ipv4)
		}

		if !hasPrimary {
			return errors.New("need a primary ipv4")
		}
	}

	if raw, ok := d.GetOk("ipv4_count"); ok {
		ipv4Count = common.IntPtr(raw.(int))
	}

	if len(ipv4s) == 0 && ipv4Count == nil {
		return errors.New("ipv4s or ipv4_count must be set")
	}

	if raw := helper.GetTags(d, "tags"); len(raw) > 0 {
		tags = raw
	}

	client := m.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := VpcService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region

	var (
		id  string
		err error
	)

	switch {
	case len(ipv4s) > 0 && len(ipv4s) <= 10:
		id, err = vpcService.CreateEni(ctx, name, vpcId, subnetId, desc, securityGroups, nil, ipv4s, tags)
		if err != nil {
			return err
		}

		d.SetId(id)

	case len(ipv4s) > 0:
		// eni should create with primary ipv4
		for i := 0; i < len(ipv4s); i++ {
			if ipv4s[i].primary {
				if i < 10 {
					break
				}

				// move primary ip to the first
				ipv4s[0], ipv4s[i] = ipv4s[i], ipv4s[0]
				break
			}
		}

		ipv4ss := chunkEniIP(ipv4s)
		withPrimaryIpv4s := ipv4ss[0]

		id, err = vpcService.CreateEni(ctx, name, vpcId, subnetId, desc, securityGroups, nil, withPrimaryIpv4s, tags)
		if err != nil {
			return err
		}

		d.SetId(id)

		for _, ipv4s := range ipv4ss[1:] {
			if err = vpcService.AssignIpv4ToEni(ctx, id, ipv4s, nil); err != nil {
				return err
			}
		}

	case ipv4Count != nil && *ipv4Count <= 10:
		id, err = vpcService.CreateEni(ctx, name, vpcId, subnetId, desc, securityGroups, ipv4Count, nil, tags)
		if err != nil {
			return err
		}

		d.SetId(id)

	case ipv4Count != nil:
		count := *ipv4Count

		id, err = vpcService.CreateEni(ctx, name, vpcId, subnetId, desc, securityGroups, common.IntPtr(10), nil, tags)
		if err != nil {
			return err
		}

		d.SetId(id)

		count -= 10
		for count > 10 {
			if err = vpcService.AssignIpv4ToEni(ctx, id, nil, common.IntPtr(10)); err != nil {
				return err
			}
			count -= 10
		}
		// assign last ip
		if count > 0 {
			if err = vpcService.AssignIpv4ToEni(ctx, id, nil, &count); err != nil {
				return err
			}
		}
	}

	if len(tags) > 0 {
		resourceName := tccommon.BuildTagResourceName("vpc", "eni", region, id)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudEniRead(d, m)
}

func resourceTencentCloudEniRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni.read")()
	defer tccommon.InconsistentCheck(d, m)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := VpcService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	enis, err := service.DescribeEniById(ctx, []string{id})
	if err != nil {
		return err
	}

	if len(enis) < 1 {
		d.SetId("")
		return nil
	}

	eni := enis[0]

	_ = d.Set("name", eni.NetworkInterfaceName)
	_ = d.Set("vpc_id", eni.VpcId)
	_ = d.Set("subnet_id", eni.SubnetId)
	_ = d.Set("description", eni.NetworkInterfaceDescription)
	_ = d.Set("mac", eni.MacAddress)
	_ = d.Set("state", eni.State)
	_ = d.Set("primary", eni.Primary)
	_ = d.Set("create_time", eni.CreatedTime)

	sgs := make([]string, 0, len(eni.GroupSet))
	for _, sg := range eni.GroupSet {
		sgs = append(sgs, *sg)
	}
	_ = d.Set("security_groups", sgs)
	_ = d.Set("orderly_security_groups", sgs)

	ipv4s := make([]map[string]interface{}, 0, len(eni.PrivateIpAddressSet))
	for _, ipv4 := range eni.PrivateIpAddressSet {
		ipv4s = append(ipv4s, map[string]interface{}{
			"ip":          ipv4.PrivateIpAddress,
			"primary":     ipv4.Primary,
			"description": ipv4.Description,
		})
	}
	_ = d.Set("ipv4_info", ipv4s)

	_, manually := d.GetOk("ipv4s")
	_, count := d.GetOk("ipv4_count")
	if !manually && !count {
		// import mode
		_ = d.Set("ipv4_count", len(ipv4s))
	}

	if eni.CdcId != nil {
		_ = d.Set("cdc_id", eni.CdcId)
	}

	tags := make(map[string]string, len(eni.TagSet))
	for _, tag := range eni.TagSet {
		tags[*tag.Key] = *tag.Value
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudEniUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	d.Partial(true)

	client := m.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := VpcService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region

	var (
		name        *string
		desc        *string
		sgs         []string
		updateAttrs []string
	)

	if d.HasChange("name") {
		updateAttrs = append(updateAttrs, "name")
		name = helper.String(d.Get("name").(string))
	}

	if d.HasChange("description") {
		updateAttrs = append(updateAttrs, "description")
		desc = helper.String(d.Get("description").(string))
	}

	if d.HasChange("security_groups") {
		updateAttrs = append(updateAttrs, "security_groups")
		if v, ok := d.GetOk("security_groups"); ok {
			sgs = helper.InterfacesStrings(v.(*schema.Set).List())
		}
	}

	if d.HasChange("orderly_security_groups") {
		updateAttrs = append(updateAttrs, "orderly_security_groups")
		if v, ok := d.GetOk("orderly_security_groups"); ok {
			sgs = helper.InterfacesStrings(v.([]interface{}))
		}
	}

	if len(updateAttrs) > 0 {
		if err := vpcService.ModifyEniAttribute(ctx, id, name, desc, sgs); err != nil {
			return err
		}

	}

	// if ipv4 set manually
	if _, ok := d.GetOk("ipv4s"); ok && d.HasChange("ipv4s") {
		oldRaw, newRaw := d.GetChange("ipv4s")
		oldSet := oldRaw.(*schema.Set)
		newSet := newRaw.(*schema.Set)

		if newSet.Len() == 0 {
			return errors.New("can't remove all ipv4s")
		}

		removeSet := oldSet.Difference(newSet).List()
		addSet := newSet.Difference(oldSet).List()

		var modifyPrimaryIpv4 *VpcEniIP

		removeIpv4 := make([]string, 0, len(removeSet))
		for _, v := range removeSet {
			m := v.(map[string]interface{})
			if m["primary"].(bool) {
				// check if only modify primary description
				modifyPrimaryIpv4 = &VpcEniIP{
					ip: net.ParseIP(m["ip"].(string)),
				}
				continue
			}

			removeIpv4 = append(removeIpv4, m["ip"].(string))
		}

		addIpv4 := make([]VpcEniIP, 0, len(addSet))
		newPrimaryCount := 0
		for _, v := range addSet {
			m := v.(map[string]interface{})

			ipStr := m["ip"].(string)
			ip := net.ParseIP(ipStr)
			if ip == nil {
				return fmt.Errorf("ip %s is invalid", ipStr)
			}

			if m["primary"].(bool) {
				if modifyPrimaryIpv4 == nil {
					return errors.New("can't set more than one primary ipv4")
				}

				// if newPrimaryCount > 1, means new ipv4s have more than one primary ipv4,
				// if only one, maybe just update primary ipv4 description
				newPrimaryCount++
				if newPrimaryCount > 1 {
					return errors.New("can't set more than one primary ipv4")
				}

				// only can update primary ipv4 description
				if modifyPrimaryIpv4.ip.String() != ipStr {
					return errors.New("can't change primary ipv4")
				}

				modifyPrimaryIpv4.desc = helper.String(m["description"].(string))
				continue
			}

			ipv4 := VpcEniIP{
				ip:      ip,
				primary: m["primary"].(bool),
				desc:    helper.String(m["description"].(string)),
			}

			addIpv4 = append(addIpv4, ipv4)
		}

		if modifyPrimaryIpv4 != nil {
			// if desc is nil, means remove primary ipv4 but not add same primary ipv4,
			// that means not just update primary ipv4 description, user remove primary ipv4
			if modifyPrimaryIpv4.desc == nil {
				return errors.New("can't remove primary ipv4")
			}

			if err := vpcService.ModifyEniPrimaryIpv4Desc(ctx, id, modifyPrimaryIpv4.ip.String(), modifyPrimaryIpv4.desc); err != nil {
				return err
			}
		}

		if len(removeIpv4) > 0 {
			if len(removeIpv4) <= 10 {
				if err := vpcService.UnAssignIpv4FromEni(ctx, id, removeIpv4); err != nil {
					return err
				}
			} else {
				for _, remove := range chunkRemoveIpv4(removeIpv4) {
					if err := vpcService.UnAssignIpv4FromEni(ctx, id, remove); err != nil {
						return err
					}
				}
			}
		}

		if len(addIpv4) > 0 {
			if len(addIpv4) <= 10 {
				if err := vpcService.AssignIpv4ToEni(ctx, id, addIpv4, nil); err != nil {
					return err
				}
			} else {
				for _, add := range chunkEniIP(addIpv4) {
					if err := vpcService.AssignIpv4ToEni(ctx, id, add, nil); err != nil {
						return err
					}
				}
			}
		}

	}

	if _, ok := d.GetOk("ipv4_count"); ok {
		if d.HasChange("ipv4_count") {
			oldRaw, newRaw := d.GetChange("ipv4_count")
			oldCount := oldRaw.(int)
			newCount := newRaw.(int)

			if newCount > oldCount {
				count := newCount - oldCount

				if count <= 10 {
					if err := vpcService.AssignIpv4ToEni(ctx, id, nil, &count); err != nil {
						return err
					}
				} else {
					for count > 10 {
						if err := vpcService.AssignIpv4ToEni(ctx, id, nil, common.IntPtr(10)); err != nil {
							return err
						}
						count -= 10
					}
					// assign last ip
					if count > 0 {
						if err := vpcService.AssignIpv4ToEni(ctx, id, nil, &count); err != nil {
							return err
						}
					}
				}
			} else {
				removeCount := oldCount - newCount
				list := d.Get("ipv4_info").([]interface{})
				removeIpv4 := make([]string, 0, removeCount)
				for _, v := range list {
					if removeCount == 0 {
						break
					}
					m := v.(map[string]interface{})
					if m["primary"].(bool) {
						continue
					}
					removeIpv4 = append(removeIpv4, m["ip"].(string))
					removeCount--
				}

				if len(removeIpv4) <= 10 {
					if err := vpcService.UnAssignIpv4FromEni(ctx, id, removeIpv4); err != nil {
						return err
					}
				} else {
					for _, remove := range chunkRemoveIpv4(removeIpv4) {
						if err := vpcService.UnAssignIpv4FromEni(ctx, id, remove); err != nil {
							return err
						}
					}
				}

			}
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := tccommon.BuildTagResourceName("vpc", "eni", region, id)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudEniRead(d, m)
}

func resourceTencentCloudEniDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := VpcService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	return service.DeleteEni(ctx, id)
}

func chunkEniIP(ipv4s []VpcEniIP) [][]VpcEniIP {
	if len(ipv4s) <= 10 {
		return [][]VpcEniIP{ipv4s}
	}

	first := ipv4s[:10]
	return append([][]VpcEniIP{first}, chunkEniIP(ipv4s[10:])...)
}

func chunkRemoveIpv4(ss []string) [][]string {
	if len(ss) <= 10 {
		return [][]string{ss}
	}

	s := ss[:10]
	return append([][]string{s}, chunkRemoveIpv4(ss[10:])...)
}
