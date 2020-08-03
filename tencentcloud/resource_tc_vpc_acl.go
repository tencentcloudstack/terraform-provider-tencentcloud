/*
Provide a resource to create a VPC ACL instance.

Example Usage

```hcl
data "tencentcloud_vpc_instances" "foo" {
}

resource "resource_vpc_acl" "default" {
    vpc_id            	= data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
    network_acl_name  	= "test_acl"
	ingress [
		"ACCEPT#192.168.1.0/24#80#TCP",
		"ACCEPT#192.168.1.0/24#80-90#TCP",
	]
	egress [
		"ACCEPT#192.168.1.0/24#80#TCP",
		"ACCEPT#192.168.1.0/24#80-90#TCP",
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
			"acl_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "`ID` of the network ACL instance.",
			},
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
		vpcID      string
		name       string
		ingress    []VpcACLRule
		egress     []VpcACLRule
	)

	if temp, ok := d.GetOk("vpc_id"); ok {
		vpcID = temp.(string)
	}
	if temp, ok := d.GetOk("name"); ok {
		name = temp.(string)
	}
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

	err = vpcService.AttachRulesToACL(ctx, aclID, ingress, egress)
	if err != nil {
		return err
	}

	d.SetId(aclID)

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

	vpcID, createTime, has, err := service.DescribeNetWorkByACLID(ctx, id)
	if err != nil {
		return err
	}
	if has == 0 {
		log.Printf("[WARN]%s %s\n", logId, "acl has been delete")
		d.SetId("")
		return nil
	}
	if has != 1 {
		errRet := fmt.Errorf("one acl_id read get %d acl info", has)
		log.Printf("[CRITAL]%s %s\n", logId, errRet.Error())
		return errRet
	}

	_ = d.Set("acl_id", id)
	_ = d.Set("vpc_id", vpcID)
	_ = d.Set("create_time", createTime)
	return nil
}

func resourceTencentCloudVpcACLUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_acl.update")()

	const (
		DeleteAll int8 = iota
		DeleteIngress
		DeleteEgress
	)

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		id      = d.Id()

		name          *string
		ingress       []VpcACLRule
		egress        []VpcACLRule
		deleteIngress bool
		deleteEgress  bool
	)

	if d.HasChange("name") {
		name = helper.String(d.Get("name").(string))
		err := service.ModifyVpcNetworkAcl(ctx, &id, name)
		if err != nil {
			return nil
		}
	}
	if d.HasChange("ingress") {
		if raw, ok := d.GetOk("ingress"); ok {
			oldIngress := helper.InterfacesStrings(raw.([]interface{}))
			for _, ingressStr := range oldIngress {
				liteRule, err := parseACLRule(ingressStr)
				if err != nil {
					return err
				}
				ingress = append(ingress, liteRule)
			}
		} else {
			old, _ := d.GetChange("ingress")
			oldIngress := helper.InterfacesStrings(old.([]interface{}))
			for _, ingressStr := range oldIngress {
				liteRule, err := parseACLRule(ingressStr)
				if err != nil {
					return err
				}
				ingress = append(ingress, liteRule)
			}

			deleteIngress = true
		}
	}

	if d.HasChange("egress") {
		if raw, ok := d.GetOk("egress"); ok {
			oldEgress := helper.InterfacesStrings(raw.([]interface{}))
			for _, egressStr := range oldEgress {
				liteRule, err := parseACLRule(egressStr)
				if err != nil {
					return err
				}
				egress = append(egress, liteRule)
			}
		} else {
			old, _ := d.GetChange("egress")
			oldEgress := helper.InterfacesStrings(old.([]interface{}))
			for _, egressStr := range oldEgress {
				liteRule, err := parseACLRule(egressStr)
				if err != nil {
					return err
				}
				egress = append(egress, liteRule)
			}

			deleteEgress = true
		}
	}

	d.Partial(true)

	if deleteIngress && deleteEgress {
		if err := service.DeleteACLRulesByChoose(ctx, id, nil, nil, DeleteAll); err != nil {
			return err
		}

		d.Partial(false)

		return resourceTencentCloudVpcACLRead(d, meta)
	}

	if deleteIngress {
		if err := service.DeleteACLRulesByChoose(ctx, id, ingress, nil, DeleteIngress); err != nil {
			return err
		}

		d.SetPartial("ingress")

		ingress = nil
	}

	if deleteEgress {
		if err := service.DeleteACLRulesByChoose(ctx, id, nil, egress, DeleteEgress); err != nil {
			return err
		}

		d.SetPartial("egress")

		egress = nil
	}

	if len(ingress) > 0 || len(egress) > 0 {
		if err := service.ModifyNetWorkAclRules(ctx, id, ingress, egress); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudVpcACLRead(d, meta)
}

func resourceTencentCloudVpcACLDelete(d *schema.ResourceData, meta interface{}) (err error) {
	defer logElapsed("resource.tencentcloud_vpc_acl.delete")()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		id      = d.Id()
	)

	err = service.DeleteAcl(ctx, id)
	if err != nil {
		return err
	}

	_, _, has, err := service.DescribeNetWorkByACLID(ctx, id)
	if err != nil {
		return err
	}

	if has > 0 {
		return fmt.Errorf("[CRITAL]%s delete network acl : %s  failed\n", logId, id)
	}

	return
}
