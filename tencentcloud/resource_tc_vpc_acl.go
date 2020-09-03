/*
Provide a resource to create a VPC ACL instance.

Example Usage

```hcl
data "tencentcloud_vpc_instances" "default" {
}

resource "tencentcloud_vpc_acl" "foo" {
    vpc_id  = data.tencentcloud_vpc_instances.default.instance_list.0.vpc_id
    name  	= "test_acl_update"
	ingress = [
		"ACCEPT#192.168.1.0/24#800#TCP",
		"ACCEPT#192.168.1.0/24#800-900#TCP",
	]
	egress = [
    	"ACCEPT#192.168.1.0/24#800#TCP",
    	"ACCEPT#192.168.1.0/24#800-900#TCP",
	]
}
```

Import

Vpc ACL can be imported, e.g.

```
$ terraform import tencentcloud_vpc_acl.default acl-id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

	aclID, err := vpcService.CreateVpcNetworkAcl(ctx, vpcID, name)
	if err != nil {
		return err
	}

	d.SetId(aclID)
	err = vpcService.AttachRulesToACL(ctx, aclID, ingress, egress)
	if err != nil {
		return err
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
		if info.EgressEntries[i].Port == nil || *info.EgressEntries[i].Port == "" {
			continue
		}
		result := strings.Join([]string{
			*info.EgressEntries[i].Action,
			*info.EgressEntries[i].CidrBlock,
			*info.EgressEntries[i].Port,
			*info.EgressEntries[i].Protocol,
		}, FILED_SP)
		egressList = append(egressList, strings.ToUpper(result))
	}

	ingressList := make([]string, 0, len(info.IngressEntries))
	for i := range info.IngressEntries {
		if info.IngressEntries[i].Port == nil || *info.IngressEntries[i].Port == "" {
			continue
		}
		result := strings.Join([]string{
			*info.IngressEntries[i].Action,
			*info.IngressEntries[i].CidrBlock,
			*info.IngressEntries[i].Port,
			*info.IngressEntries[i].Protocol,
		}, FILED_SP)
		ingressList = append(ingressList, strings.ToUpper(result))
	}
	_ = d.Set("egress", egressList)
	_ = d.Set("ingress", ingressList)

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

		d.SetPartial("name")
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
		d.SetPartial("ingress")
		d.SetPartial("egress")
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
