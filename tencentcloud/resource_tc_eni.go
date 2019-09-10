package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func eniIpInputResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"primary": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func eniIpOutputResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudEni() *schema.Resource {
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
				ValidateFunc: validateStringLengthInRange(0, 60),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validateStringLengthInRange(0, 60),
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"ipv4s": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"ipv4_count"},
				Elem:          eniIpInputResource(),
				Set:           schema.HashResource(eniIpInputResource()),
				MaxItems:      30,
			},
			"ipv4_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"ipv4s"},
				ValidateFunc:  validateIntegerInRange(1, 30),
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			// computed
			"mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_info": {
				Type:     schema.TypeList,
				Elem:     eniIpOutputResource(),
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudEniCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	name := d.Get("name").(string)
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)
	desc := d.Get("description").(string)

	var (
		securityGroups []string
		ipv4s          []VpcEniIP
		ipv4Count      *int
	)

	if raw, ok := d.GetOk("security_groups"); ok {
		securityGroups = expandStringList(raw.(*schema.Set).List())
	}

	if raw, ok := d.GetOk("ipv4s"); ok {
		set := raw.(*schema.Set)
		ipv4s = make([]VpcEniIP, 0, set.Len())

		for _, v := range set.List() {
			m := v.(map[string]interface{})

			ipStr := m["ip"]
			ip := net.ParseIP(ipStr.(string))
			if ip == nil {
				return fmt.Errorf("ip %s is invalid", ipStr)
			}

			ipv4 := VpcEniIP{
				ip:      ip,
				primary: m["primary"].(bool),
			}

			if desc := m["description"].(string); desc != "" {
				ipv4.desc = stringToPointer(desc)
			}

			ipv4s = append(ipv4s, ipv4)
		}
	}

	if raw, ok := d.GetOk("ipv4_count"); ok {
		ipv4Count = common.IntPtr(raw.(int))
	}

	if len(ipv4s) == 0 && ipv4Count == nil {
		return errors.New("ipv4s or ipv4_count must be set")
	}

	tags := getTags(d, "tags")

	vpcService := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	var (
		id  string
		err error
	)

	switch {
	case len(ipv4s) > 0 && len(ipv4s) <= 10:
		id, err = vpcService.CreateEni(ctx, name, vpcId, subnetId, desc, securityGroups, nil, ipv4s)
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
				primaryIpv4 := ipv4s[i]
				ipv4s = append(ipv4s[:i], ipv4s[i+1:]...)
				newIpv4s := make([]VpcEniIP, len(ipv4s)+1)
				newIpv4s[0] = primaryIpv4
				copy(newIpv4s[1:], ipv4s)
				ipv4s = newIpv4s
				break
			}
		}

		id, err = vpcService.CreateEni(ctx, name, vpcId, subnetId, desc, securityGroups, nil, ipv4s[:10])
		if err != nil {
			return err
		}

		d.SetId(id)

		ipv4s = ipv4s[10:]
		for len(ipv4s) > 10 {
			if err = vpcService.AssignIpv4ToEni(ctx, id, ipv4s[:10], nil); err != nil {
				return err
			}
			ipv4s = ipv4s[10:]
		}
		// assign last ipv4s
		if len(ipv4s) > 0 {
			if err = vpcService.AssignIpv4ToEni(ctx, id, ipv4s, nil); err != nil {
				return err
			}
		}

	case ipv4Count != nil && *ipv4Count <= 10:
		id, err = vpcService.CreateEni(ctx, name, vpcId, subnetId, desc, securityGroups, ipv4Count, nil)
		if err != nil {
			return err
		}

		d.SetId(id)

	case ipv4Count != nil:
		count := *ipv4Count

		id, err = vpcService.CreateEni(ctx, name, vpcId, subnetId, desc, securityGroups, common.IntPtr(10), nil)
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

	if tags != nil {
		tagService := TagService{client: m.(*TencentCloudClient).apiV3Conn}

		region := m.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:eni/%s", region, id)

		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudEniRead(d, m)
}

func resourceTencentCloudEniRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	enis, err := service.DescribeEniById(ctx, []string{id})
	if err != nil {
		return err
	}

	var eni *vpc.NetworkInterface
	for _, e := range enis {
		if e.NetworkInterfaceId == nil {
			return errors.New("eni id is nil")
		}

		if *e.NetworkInterfaceId == id {
			eni = e
			break
		}
	}

	if eni == nil {
		d.SetId("")
		return nil
	}

	if nilFields := CheckNil(eni, map[string]string{
		"NetworkInterfaceName":        "name",
		"NetworkInterfaceDescription": "description",
		"VpcId":                       "vpc id",
		"SubnetId":                    "subnet id",
		"MacAddress":                  "mac address",
		"State":                       "state",
		"CreatedTime":                 "create time",
		"Primary":                     "primary",
	}); len(nilFields) > 0 {
		return fmt.Errorf("eni %v are nil", nilFields)
	}

	d.Set("name", eni.NetworkInterfaceName)
	d.Set("vpc_id", eni.VpcId)
	d.Set("subnet_id", eni.SubnetId)
	d.Set("description", eni.NetworkInterfaceDescription)
	d.Set("mac", eni.MacAddress)
	d.Set("state", eni.State)
	d.Set("primary", eni.Primary)
	d.Set("create_time", eni.CreatedTime)

	if len(eni.GroupSet) > 0 {
		sgs := make([]string, 0, len(eni.GroupSet))
		for _, sg := range eni.GroupSet {
			sgs = append(sgs, *sg)
		}

		d.Set("security_groups", sgs)
	}

	if len(eni.PrivateIpAddressSet) == 0 {
		return errors.New("eni ipv4 is empty")
	}

	ipv4s := make([]map[string]interface{}, 0, len(eni.PrivateIpAddressSet))
	for _, ipv4 := range eni.PrivateIpAddressSet {
		if nilFields := CheckNil(ipv4, map[string]string{
			"PrivateIpAddress": "ip",
			"Primary":          "primary",
			"Description":      "description",
		}); len(nilFields) > 0 {
			return fmt.Errorf("eni ipv4 %v are nil", nilFields)
		}

		ipv4s = append(ipv4s, map[string]interface{}{
			"ip":          *ipv4.PrivateIpAddress,
			"primary":     *ipv4.Primary,
			"description": *ipv4.Description,
		})
	}
	d.Set("ipv4_info", ipv4s)

	_, manually := d.GetOk("ipv4s")
	_, count := d.GetOk("ipv4_count")
	if !manually && !count {
		// import mode
		d.Set("ipv4_count", len(ipv4s))
	}

	if len(eni.TagSet) > 0 {
		tags := make(map[string]string, len(eni.TagSet))
		for _, tag := range eni.TagSet {
			if tag.Key == nil {
				return errors.New("tag key is nil")
			}
			if tag.Value == nil {
				return errors.New("tag value is nil")
			}

			tags[*tag.Key] = *tag.Value
		}

		d.Set("tags", tags)
	}

	return nil
}

func resourceTencentCloudEniUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	d.Partial(true)

	vpcService := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	var (
		name        *string
		desc        *string
		sgs         []string
		updateAttrs []string
	)

	if d.HasChange("name") {
		updateAttrs = append(updateAttrs, "name")
		name = stringToPointer(d.Get("name").(string))
	}

	if d.HasChange("description") {
		updateAttrs = append(updateAttrs, "description")
		desc = stringToPointer(d.Get("description").(string))
	}

	if d.HasChange("security_groups") {
		updateAttrs = append(updateAttrs, "security_groups")
	}
	sgs = expandStringList(d.Get("security_groups").(*schema.Set).List())

	if len(updateAttrs) > 0 {
		if err := vpcService.ModifyEniAttribute(ctx, id, name, desc, sgs); err != nil {
			return err
		}

		for _, attr := range updateAttrs {
			d.SetPartial(attr)
		}
	}

	// if ipv4 set manually
	if _, ok := d.GetOk("ipv4s"); ok {
		if d.HasChange("ipv4s") {
			oldRaw, newRaw := d.GetChange("ipv4s")
			oldSet := oldRaw.(*schema.Set)
			newSet := newRaw.(*schema.Set)

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

					// not just update primary ipv4 description
					if modifyPrimaryIpv4.ip.String() != ipStr {
						return errors.New("can't change primary ipv4")
					}

					modifyPrimaryIpv4.desc = stringToPointer(m["description"].(string))
					continue
				}

				ipv4 := VpcEniIP{
					ip:      ip,
					primary: m["primary"].(bool),
				}

				if desc := m["description"].(string); desc != "" {
					ipv4.desc = stringToPointer(desc)
				}

				addIpv4 = append(addIpv4, ipv4)
			}

			if len(removeIpv4) > 0 {
				if len(removeIpv4) <= 10 {
					if err := vpcService.UnAssignIpv4FromEni(ctx, id, removeIpv4); err != nil {
						return err
					}
				} else {
					for len(removeIpv4) > 10 {
						if err := vpcService.UnAssignIpv4FromEni(ctx, id, removeIpv4[:10]); err != nil {
							return err
						}
						removeIpv4 = removeIpv4[10:]
					}
					// unassign last ipv4
					if len(removeIpv4) > 0 {
						if err := vpcService.UnAssignIpv4FromEni(ctx, id, removeIpv4); err != nil {
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
					for len(addIpv4) > 10 {
						if err := vpcService.AssignIpv4ToEni(ctx, id, addIpv4[:10], nil); err != nil {
							return err
						}
						addIpv4 = addIpv4[10:]
					}
					// assign last ipv4
					if len(addIpv4) > 0 {
						if err := vpcService.AssignIpv4ToEni(ctx, id, addIpv4, nil); err != nil {
							return err
						}
					}
				}
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

			d.SetPartial("ipv4s")
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
					for len(removeIpv4) > 10 {
						if err := vpcService.UnAssignIpv4FromEni(ctx, id, removeIpv4[:10]); err != nil {
							return err
						}
						removeIpv4 = removeIpv4[10:]
					}
					// unassign last ipv4
					if len(removeIpv4) > 0 {
						if err := vpcService.UnAssignIpv4FromEni(ctx, id, removeIpv4); err != nil {
							return err
						}
					}
				}

				d.SetPartial("ipv4_count")
			}
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := TagService{client: m.(*TencentCloudClient).apiV3Conn}

		region := m.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:eni/%s", region, id)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudEniRead(d, m)
}

func resourceTencentCloudEniDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteEni(ctx, id)
}
