package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcACLCreate,
		Read:   resourceTencentCloudVpcACLRead,
		Update: resourceTencentCloudVpcACLUpdate,
		Delete: resourceTencentCloudVpcACLDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "ID of the VPC instance.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(0, 60),
				Description:  "Name of the network ACL.",
			},
			"ingress": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Ingress rules. A rule must match the following format: [action]#[cidr_ip]#[port]#[protocol]. The available value of 'action' is `ACCEPT` and `DROP`. The 'cidr_ip' must be an IP address network or segment. The 'port' valid format is `80`, `80,443`, `80-90` or `ALL`. The available value of 'protocol' is `TCP`, `UDP`, `ICMP` and `ALL`. When 'protocol' is `ICMP` or `ALL`, the 'port' must be `ALL`.",
			},
			"egress": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Egress rules. A rule must match the following format: [action]#[cidr_ip]#[port]#[protocol]. The available value of 'action' is `ACCEPT` and `DROP`. The 'cidr_ip' must be an IP address network or segment. The 'port' valid format is `80`, `80,443`, `80-90` or `ALL`. The available value of 'protocol' is `TCP`, `UDP`, `ICMP` and `ALL`. When 'protocol' is `ICMP` or `ALL`, the 'port' must be `ALL`.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the vpc acl.",
			},
			//compute
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of ACL.",
			},
		},
	}
}

func resourceTencentCloudVpcACLCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_acl.create")()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		vpcService = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

		ingress []VpcACLRule
		egress  []VpcACLRule
		vpcID   = d.Get("vpc_id").(string)
		name    = d.Get("name").(string)
		tags    map[string]string
	)

	if temp, ok := d.GetOk("ingress"); ok {
		ingressStrs := helper.InterfacesStrings(temp.([]interface{}))
		for _, ingressStr := range ingressStrs {
			liteRule, err := parseACLRule(ingressStr)
			if err != nil {
				return err
			}
			ingress = append(ingress, liteRule)
		}
	}
	if temp, ok := d.GetOk("egress"); ok {
		egressStrs := helper.InterfacesStrings(temp.([]interface{}))
		for _, egressStr := range egressStrs {
			liteRule, err := parseACLRule(egressStr)
			if err != nil {
				return err
			}
			egress = append(egress, liteRule)
		}
	}

	if temp := helper.GetTags(d, "tags"); len(temp) > 0 {
		tags = temp
	}

	aclID, err := vpcService.CreateVpcNetworkAcl(ctx, vpcID, name, tags)
	if err != nil {
		return err
	}

	d.SetId(aclID)

	client := meta.(*TencentCloudClient).apiV3Conn
	tagService := TagService{client: client}
	region := client.Region
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := BuildTagResourceName("vpc", "acl", region, aclID)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	if len(ingress) > 0 || len(egress) > 0 {
		err = vpcService.AttachRulesToACL(ctx, aclID, ingress, egress)
		if err != nil {
			return err
		}
	}
	return resourceTencentCloudVpcACLRead(d, meta)
}

func resourceTencentCloudVpcACLRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_acl.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		id      = d.Id()
	)

	info, has, err := service.DescribeNetWorkByACLID(ctx, id)
	if err != nil {
		return err
	}
	if has == 0 {
		log.Printf("[WARN]%s %s\n", logId, "ACL has been delete")
		d.SetId("")
		return nil
	}
	if has != 1 {
		errRet := fmt.Errorf("one acl_id read get %d ACL info", has)
		log.Printf("[CRITAL]%s %s\n", logId, errRet.Error())
		return errRet
	}

	_ = d.Set("vpc_id", info.VpcId)
	_ = d.Set("create_time", info.CreatedTime)
	_ = d.Set("name", info.NetworkAclName)
	egressList := make([]string, 0, len(info.EgressEntries))
	for i := range info.EgressEntries {
		// remove default rule
		if CheckIfDefaultRule(info.EgressEntries[i]) {
			continue
		}

		var (
			action    string
			cidrBlock string
			port      string
			protocol  string
		)

		if info.EgressEntries[i].Action != nil {
			action = *info.EgressEntries[i].Action
		}
		if info.EgressEntries[i].CidrBlock != nil {
			cidrBlock = *info.EgressEntries[i].CidrBlock
		}
		if info.EgressEntries[i].Port == nil || *info.EgressEntries[i].Port == "" {
			port = "ALL"
		} else {
			port = *info.EgressEntries[i].Port
		}
		if info.EgressEntries[i].Protocol != nil {
			protocol = *info.EgressEntries[i].Protocol
		}

		result := strings.Join([]string{
			action,
			cidrBlock,
			port,
			protocol,
		}, FILED_SP)

		egressList = append(egressList, strings.ToUpper(result))
	}

	ingressList := make([]string, 0, len(info.IngressEntries))
	for i := range info.IngressEntries {
		// remove default rule
		if CheckIfDefaultRule(info.IngressEntries[i]) {
			continue
		}

		var (
			action    string
			cidrBlock string
			port      string
			protocol  string
		)

		if info.IngressEntries[i].Action != nil {
			action = *info.IngressEntries[i].Action
		}
		if info.IngressEntries[i].CidrBlock != nil {
			cidrBlock = *info.IngressEntries[i].CidrBlock
		}
		if info.IngressEntries[i].Port == nil || *info.IngressEntries[i].Port == "" {
			port = "ALL"
		} else {
			port = *info.IngressEntries[i].Port
		}
		if info.IngressEntries[i].Protocol != nil {
			protocol = *info.IngressEntries[i].Protocol
		}

		result := strings.Join([]string{
			action,
			cidrBlock,
			port,
			protocol,
		}, FILED_SP)
		ingressList = append(ingressList, strings.ToUpper(result))
	}
	_ = d.Set("egress", egressList)
	_ = d.Set("ingress", ingressList)

	client := meta.(*TencentCloudClient).apiV3Conn
	tagService := TagService{client: client}
	region := client.Region
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "acl", region, id)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpcACLUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_acl.update")()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		id      = d.Id()

		name    *string
		ingress []VpcACLRule
		egress  []VpcACLRule
	)

	d.Partial(true)

	if d.HasChange("name") {
		name = helper.String(d.Get("name").(string))
		err := service.ModifyVpcNetworkAcl(ctx, &id, name)
		if err != nil {
			return nil
		}
	}

	if d.HasChange("ingress") {
		_, new := d.GetChange("ingress")
		if len(new.([]interface{})) == 0 {
			//del ingress
			ingress = []VpcACLRule{
				{
					protocol: "all",
					cidrIp:   "0.0.0.0/0",
					action:   "DROP",
				},
			}
		} else {
			newIngress := helper.InterfacesStrings(new.([]interface{}))
			for _, ingressStr := range newIngress {
				liteRule, err := parseACLRule(ingressStr)
				if err != nil {
					return err
				}
				ingress = append(ingress, liteRule)
			}
		}
	}

	if d.HasChange("egress") {
		_, new := d.GetChange("egress")
		if len(new.([]interface{})) == 0 {
			//del ingress
			egress = []VpcACLRule{
				{
					protocol: "all",
					cidrIp:   "0.0.0.0/0",
					action:   "DROP",
				},
			}
		} else {
			newIngress := helper.InterfacesStrings(new.([]interface{}))
			for _, egressStr := range newIngress {
				liteRule, err := parseACLRule(egressStr)
				if err != nil {
					return err
				}
				egress = append(egress, liteRule)
			}
		}
	}

	if d.HasChange("egress") || d.HasChange("ingress") {
		if err := service.ModifyNetWorkAclRules(ctx, id, ingress, egress); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		client := meta.(*TencentCloudClient).apiV3Conn
		tagService := TagService{client: client}
		region := client.Region

		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("vpc", "acl", region, id)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudVpcACLRead(d, meta)
}

func resourceTencentCloudVpcACLDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_acl.delete")()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		id      = d.Id()
	)

	err := service.DeleteAcl(ctx, id)
	if err != nil {
		return err
	}

	_, has, err := service.DescribeNetWorkByACLID(ctx, id)
	if err != nil {
		return err
	}

	if has > 0 {
		return fmt.Errorf("[CRITAL]%s delete network ACL : %s  failed\n", logId, id)
	}
	return nil
}

func CheckIfDefaultRule(aclEntry *vpc.NetworkAclEntry) bool {
	// remove default ipv6 rule
	if aclEntry.Protocol != nil && *aclEntry.Protocol == "all" &&
		aclEntry.Ipv6CidrBlock != nil && *aclEntry.Ipv6CidrBlock == "::/0" &&
		aclEntry.Action != nil && *aclEntry.Action == "Accept" {
		return true
	}
	// remove default cidr rule
	if aclEntry.Protocol != nil && *aclEntry.Protocol == "all" &&
		aclEntry.CidrBlock != nil && *aclEntry.CidrBlock == "0.0.0.0/0" &&
		aclEntry.Action != nil && *aclEntry.Action == "Drop" {
		return true
	}
	return false
}
