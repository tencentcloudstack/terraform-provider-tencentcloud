/*
Provides a resource to create a NAT forwarding.

Example Usage

```hcl
resource "tencentcloud_dnat" "foo" {
  vpc_id             = "vpc-asg3sfa3"
  nat_id             = "nat-2515tdg"
  protocol           = "tcp"
  elastic_ip         = "139.199.232.238"
  elastic_port       = 80
  private_ip         = "10.0.0.1"
  private_port       = 22
  description        = "test"
}
```

Import

NAT forwarding can be imported using the id, e.g.

```
$ terraform import tencentcloud_dnat.foo tcp://vpc-asg3sfa3:nat-1asg3t63@127.15.2.3:8080
```
*/
package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDnat() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnatCreate,
		Read:   resourceTencentCloudDnatRead,
		Update: resourceTencentCloudDnatUpdate,
		Delete: resourceTencentCloudDnatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the VPC.",
			},
			"nat_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the NAT gateway.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"TCP", "UDP"}),
				Description:  "Type of the network protocol, the available values are: `TCP` and `UDP`.",
			},
			"elastic_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIp,
				Description:  "Network address of the EIP.",
			},
			"elastic_port": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePort,
				Description:  "Port of the EIP.",
			},
			"private_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIp,
				Description:  "Network address of the backend service.",
			},
			"private_port": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePort,
				Description:  "Port of intranet.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the NAT forward.",
			},
		},
	}
}

func resourceTencentCloudDnatCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnat.create")()

	logId := getLogId(contextNil)
	request := vpc.NewCreateNatGatewayDestinationIpPortTranslationNatRuleRequest()
	var natForward vpc.DestinationIpPortTranslationNatRule
	natForward.IpProtocol = helper.String(d.Get("protocol").(string))
	natForward.PublicIpAddress = helper.String(d.Get("elastic_ip").(string))
	uePort, epErr := strconv.ParseInt(d.Get("elastic_port").(string), 10, 64)
	if epErr != nil {
		return fmt.Errorf("elastic port format error")
	}
	ePort := uint64(uePort)
	natForward.PublicPort = &ePort
	natForward.PrivateIpAddress = helper.String(d.Get("private_ip").(string))
	upPort, ppErr := strconv.ParseInt(d.Get("private_port").(string), 10, 64)
	if ppErr != nil {
		return fmt.Errorf("private port format error")
	}
	pPort := uint64(upPort)
	natForward.PrivatePort = &pPort
	description := ""
	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}
	natForward.Description = helper.String(description)
	natGatewayId := d.Get("nat_id").(string)
	request.NatGatewayId = helper.String(natGatewayId)

	request.DestinationIpPortTranslationNatRules = []*vpc.DestinationIpPortTranslationNatRule{&natForward}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateNatGatewayDestinationIpPortTranslationNatRule(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create NAT forward failed, reason:%s\n", logId, err.Error())
		return err
	}

	dnatId := buildDnatId(&natForward, d.Get("vpc_id").(string), natGatewayId)

	d.SetId(dnatId)
	return resourceTencentCloudDnatRead(d, meta)
}

func resourceTencentCloudDnatRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnat.read")()

	logId := getLogId(contextNil)
	_, params, e := parseDnatId(d.Id())
	if e != nil {
		return fmt.Errorf("[CRITAL]parse DNAT id fail, reason[%s]\n", e.Error())
	}
	request := vpc.NewDescribeNatGatewayDestinationIpPortTranslationNatRulesRequest()
	var response *vpc.DescribeNatGatewayDestinationIpPortTranslationNatRulesResponse
	request.Filters = make([]*vpc.Filter, 0, len(params))
	for k, v := range params {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, filter)
	}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeNatGatewayDestinationIpPortTranslationNatRules(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read DNAT failed, reason:%s\n", logId, err.Error())
		return err
	}
	if len(response.Response.NatGatewayDestinationIpPortTranslationNatRuleSet) < 1 {
		d.SetId("")
		return nil
	}

	dnat := response.Response.NatGatewayDestinationIpPortTranslationNatRuleSet[0]

	_ = d.Set("vpc_id", dnat.VpcId)
	_ = d.Set("nat_id", dnat.NatGatewayId)
	_ = d.Set("protocol", dnat.IpProtocol)
	_ = d.Set("elastic_ip", dnat.PublicIpAddress)
	_ = d.Set("elastic_port", strconv.Itoa(int(*dnat.PublicPort)))
	_ = d.Set("private_ip", dnat.PrivateIpAddress)
	_ = d.Set("private_port", strconv.Itoa(int(*dnat.PrivatePort)))
	_ = d.Set("description", dnat.Description)
	return nil
}

func resourceTencentCloudDnatUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnat.update")()

	logId := getLogId(contextNil)
	//only modify description
	if d.HasChange("description") {
		description := ""
		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		}
		natForward, params, e := parseDnatId(d.Id())
		if e != nil {
			return fmt.Errorf("[CRITAL]parse DNAT id fail, reason[%s]\n", e.Error())
		}
		//missing target port and ip
		srequest := vpc.NewDescribeNatGatewayDestinationIpPortTranslationNatRulesRequest()
		srequest.Filters = make([]*vpc.Filter, 0, len(params))
		for k, v := range params {
			filter := &vpc.Filter{
				Name:   helper.String(k),
				Values: []*string{helper.String(v)},
			}
			srequest.Filters = append(srequest.Filters, filter)
		}
		var sresponse *vpc.DescribeNatGatewayDestinationIpPortTranslationNatRulesResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeNatGatewayDestinationIpPortTranslationNatRules(srequest)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, srequest.GetAction(), srequest.ToJsonString(), e.Error())
				return retryError(e)
			}
			sresponse = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read DNAT failed, reason:%s\n", logId, err.Error())
			return err
		}
		if len(sresponse.Response.NatGatewayDestinationIpPortTranslationNatRuleSet) < 1 {
			return fmt.Errorf("modify error, forwarding rule not found")
		}
		target := sresponse.Response.NatGatewayDestinationIpPortTranslationNatRuleSet[0]
		natForward.PrivateIpAddress = target.PrivateIpAddress
		natForward.PrivatePort = target.PrivatePort
		natForward.Description = target.Description
		request := vpc.NewModifyNatGatewayDestinationIpPortTranslationNatRuleRequest()
		request.NatGatewayId = helper.String(params["nat-gateway-id"])
		request.SourceNatRule = natForward
		newNatForward := &vpc.DestinationIpPortTranslationNatRule{}
		newNatForward.PublicPort = natForward.PublicPort
		newNatForward.PublicIpAddress = natForward.PublicIpAddress
		newNatForward.PrivatePort = natForward.PrivatePort
		newNatForward.PrivateIpAddress = natForward.PrivateIpAddress
		newNatForward.IpProtocol = natForward.IpProtocol
		newNatForward.Description = helper.String(description)
		request.DestinationNatRule = newNatForward
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyNatGatewayDestinationIpPortTranslationNatRule(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s delete DNAT failed, reason:%s\n", logId, err.Error())
			return err
		}
	}
	return nil
}

func resourceTencentCloudDnatDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnat.delete")()

	logId := getLogId(contextNil)
	natForward, params, e := parseDnatId(d.Id())
	if e != nil {
		return fmt.Errorf("[CRITAL]parse DNAT id fail, reason[%s]\n", e.Error())
	}
	//missing target port and ip
	srequest := vpc.NewDescribeNatGatewayDestinationIpPortTranslationNatRulesRequest()
	srequest.Filters = make([]*vpc.Filter, 0, len(params))
	for k, v := range params {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		srequest.Filters = append(srequest.Filters, filter)
	}
	var sresponse *vpc.DescribeNatGatewayDestinationIpPortTranslationNatRulesResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeNatGatewayDestinationIpPortTranslationNatRules(srequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, srequest.GetAction(), srequest.ToJsonString(), e.Error())
			return retryError(e)
		}
		sresponse = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read DNAT failed, reason:%s\n", logId, err.Error())
		return err
	}
	if len(sresponse.Response.NatGatewayDestinationIpPortTranslationNatRuleSet) < 1 {
		return fmt.Errorf("delete error, forwarding rule not found")
	}
	target := sresponse.Response.NatGatewayDestinationIpPortTranslationNatRuleSet[0]
	natForward.PrivateIpAddress = target.PrivateIpAddress
	natForward.PrivatePort = target.PrivatePort
	request := vpc.NewDeleteNatGatewayDestinationIpPortTranslationNatRuleRequest()
	request.NatGatewayId = helper.String(params["nat-gateway-id"])
	request.DestinationIpPortTranslationNatRules = []*vpc.DestinationIpPortTranslationNatRule{natForward}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DeleteNatGatewayDestinationIpPortTranslationNatRule(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete DNAT failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}

// Build an id for a Forward Entry, eg "tcp://VpcId:NatId@127.15.2.3:8080"
func buildDnatId(entry *vpc.DestinationIpPortTranslationNatRule, vpcId string, natGatewayId string) (entryId string) {
	entryId = fmt.Sprintf("%s://%s:%s@%s:%d", *entry.IpProtocol, vpcId, natGatewayId, *entry.PublicIpAddress, *entry.PublicPort)
	log.Printf("[DEBUG] buildDnatId entryId=%s", entryId)
	return
}

//Parse Forward Entry id
func parseDnatId(entryId string) (entry *vpc.DestinationIpPortTranslationNatRule, params map[string]string, err error) {
	log.Printf("[DEBUG] parseDnatId entryId: %s", entryId)
	params = make(map[string]string)
	u, errors := url.Parse(entryId)
	if errors != nil {
		err = errors
		return
	}
	natId, _ := u.User.Password()
	host, port, _ := net.SplitHostPort(u.Host)
	entry = &vpc.DestinationIpPortTranslationNatRule{}
	params["nat-gateway-id"] = natId
	params["vpc-id"] = u.User.Username()
	entry.IpProtocol = helper.String(strings.ToUpper(u.Scheme))
	entry.PublicIpAddress = helper.String(host)

	portInt, err := strconv.Atoi(port)
	port64 := uint64(portInt)
	entry.PublicPort = &port64
	b, _ := json.Marshal(entry)
	params["public-ip-address"] = host
	params["public-port"] = port
	log.Printf("[DEBUG] parseDnatId result: %s", b)
	return
}
