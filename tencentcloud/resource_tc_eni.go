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
			},
			"ipv4_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"ipv4s"},
				ValidateFunc:  validateIntegerMin(1),
			},
			"ipv6s": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"ipv6_count"},
				Elem:          eniIpInputResource(),
				Set:           schema.HashResource(eniIpInputResource()),
			},
			"ipv6_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"ipv6s"},
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
			"ipv6_info": {
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
		ipv6s          []VpcEniIP
		ipv6Count      *int
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

	if raw, ok := d.GetOk("ipv6s"); ok {
		set := raw.(*schema.Set)
		ipv6s = make([]VpcEniIP, 0, set.Len())

		for _, v := range set.List() {
			m := v.(map[string]interface{})

			ipStr := m["ip"]
			ip := net.ParseIP(ipStr.(string))
			if ip == nil {
				return fmt.Errorf("ip %s is invalid", ipStr)
			}

			ipv6 := VpcEniIP{
				ip:      ip,
				primary: m["primary"].(bool),
			}

			if desc := m["description"].(string); desc != "" {
				ipv6.desc = stringToPointer(desc)
			}

			ipv6s = append(ipv6s, ipv6)
		}
	}

	if raw, ok := d.GetOk("ipv6_count"); ok {
		ipv6Count = common.IntPtr(raw.(int))
	}

	tags := getTags(d, "tags")

	vpcService := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := vpcService.CreateEni(ctx, name, vpcId, subnetId, desc, securityGroups, ipv4Count, ipv4s)
	if err != nil {
		return err
	}

	d.SetId(id)

	if len(ipv6s) > 0 || ipv6Count != nil {
		if err := vpcService.AssignIpv6ToEni(ctx, id, ipv6s, ipv6Count); err != nil {
			return err
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

	enis, err := service.DescribeEniById(ctx, id)
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

		if err := d.Set("security_groups", sgs); err != nil {
			panic(err)
		}
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

	ipv6s := make([]map[string]interface{}, 0, len(eni.Ipv6AddressSet))
	for _, ipv6 := range eni.Ipv6AddressSet {
		if nilFields := CheckNil(ipv6, map[string]string{
			"Address":     "ip",
			"Primary":     "primary",
			"Description": "description",
		}); len(nilFields) > 0 {
			return fmt.Errorf("eni ipv6 %v are nil", nilFields)
		}

		ipv6s = append(ipv6s, map[string]interface{}{
			"ip":          *ipv6.Address,
			"primary":     *ipv6.Primary,
			"description": *ipv6.Description,
		})
	}

	d.Set("ipv6_info", ipv6s)

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

			removeIpv4 := make([]string, 0, len(removeSet))
			for _, v := range removeSet {
				m := v.(map[string]interface{})
				if m["primary"].(bool) {
					return errors.New("primary ip can't be removed")
				}

				removeIpv4 = append(removeIpv4, m["ip"].(string))
			}

			addIpv4 := make([]VpcEniIP, 0, len(addSet))
			for _, v := range addSet {
				m := v.(map[string]interface{})

				ipStr := m["ip"].(string)
				ip := net.ParseIP(ipStr)
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

				addIpv4 = append(addIpv4, ipv4)
			}

			if len(removeIpv4) > 0 {
				if err := vpcService.UnAssignIpv4FromEni(ctx, id, removeIpv4); err != nil {
					return err
				}
			}

			if len(addIpv4) > 0 {
				if err := vpcService.AssignIpv4ToEni(ctx, id, addIpv4, nil); err != nil {
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
				if err := vpcService.AssignIpv4ToEni(ctx, id, nil, common.IntPtr(newCount-oldCount)); err != nil {
					return err
				}
			} else {
				removeCount := oldCount - newCount
				list := d.Get("ipv4_info").([]interface{})
				ipv4s := make([]string, 0, removeCount)
				for _, v := range list {
					if removeCount == 0 {
						break
					}
					m := v.(map[string]interface{})
					if m["primary"].(bool) {
						continue
					}
					ipv4s = append(ipv4s, m["ip"].(string))
					removeCount--
				}

				if err := vpcService.UnAssignIpv4FromEni(ctx, id, ipv4s); err != nil {
					return err
				}

				d.SetPartial("ipv4_count")
			}
		}
	}

	// if ipv6 set manually
	if _, ok := d.GetOk("ipv6s"); ok {
		if d.HasChange("ipv6s") {
			oldRaw, newRaw := d.GetChange("ipv6s")
			oldSet := oldRaw.(*schema.Set)
			newSet := newRaw.(*schema.Set)

			removeSet := oldSet.Difference(newSet).List()
			addSet := newSet.Difference(oldSet).List()

			removeIpv6 := make([]string, 0, len(removeSet))
			for _, v := range removeSet {
				m := v.(map[string]interface{})
				removeIpv6 = append(removeIpv6, m["ip"].(string))
			}

			addIpv6 := make([]VpcEniIP, 0, len(addSet))
			for _, v := range addSet {
				m := v.(map[string]interface{})

				ipStr := m["ip"].(string)
				ip := net.ParseIP(ipStr)
				if ip == nil {
					return fmt.Errorf("ip %s is invalid", ipStr)
				}

				ipv6 := VpcEniIP{
					ip:      ip,
					primary: m["primary"].(bool),
				}

				if desc, ok := m["description"]; ok {
					ipv6.desc = stringToPointer(desc.(string))
				}

				addIpv6 = append(addIpv6, ipv6)
			}

			if err := vpcService.UnAssignIpv6FromEni(ctx, id, removeIpv6); err != nil {
				return err
			}

			if err := vpcService.AssignIpv6ToEni(ctx, id, addIpv6, nil); err != nil {
				return err
			}

			d.SetPartial("ipv6s")
		}
	}

	if _, ok := d.GetOk("ipv6_count"); ok {
		if d.HasChange("ipv6_count") {
			oldRaw, newRaw := d.GetChange("ipv6_count")
			oldCount := oldRaw.(int)
			newCount := newRaw.(int)

			switch {
			case newCount > oldCount:
				if err := vpcService.AssignIpv6ToEni(ctx, id, nil, common.IntPtr(newCount-oldCount)); err != nil {
					return err
				}

			case newCount == 0:
				list := d.Get("ipv6_info").([]interface{})
				ipv6s := make([]string, 0, len(list))
				for _, ipv6 := range list {
					ipv6s = append(ipv6s, ipv6.(string))
				}

				if err := vpcService.UnAssignIpv6FromEni(ctx, id, ipv6s); err != nil {
					return err
				}

			default:
				removeCount := oldCount - newCount
				list := d.Get("ipv6s").(*schema.Set).List()
				ipv6s := make([]string, 0, removeCount)
				for i := 0; i < removeCount; i++ {
					ipv6s = append(ipv6s, list[i].(string))
				}

				if err := vpcService.UnAssignIpv6FromEni(ctx, id, ipv6s); err != nil {
					return err
				}
			}

			d.SetPartial("ipv6_count")
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
