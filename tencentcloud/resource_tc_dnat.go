package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	vpc "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
)

func resourceTencentCloudDnat() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnatCreate,
		Read:   resourceTencentCloudDnatRead,
		Delete: resourceTencentCloudDnatDelete,

		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nat_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"tcp", "udp"}),
			},
			"elastic_ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIp,
			},
			"elastic_port": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePort,
			},
			"private_ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIp,
			},
			"private_port": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePort,
			},
		},
	}
}

func resourceTencentCloudDnatCreate(d *schema.ResourceData, meta interface{}) error {

	args := vpc.NewAddDnaptRuleRequest()

	args.VpcId = common.StringPtr(d.Get("vpc_id").(string))
	args.NatId = common.StringPtr(d.Get("nat_id").(string))
	args.Proto = common.StringPtr(d.Get("protocol").(string))
	args.Eip = common.StringPtr(d.Get("elastic_ip").(string))
	args.Eport = common.StringPtr(d.Get("elastic_port").(string))
	args.Pip = common.StringPtr(d.Get("private_ip").(string))
	args.Pport = common.StringPtr(d.Get("private_port").(string))

	conn := meta.(*TencentCloudClient).vpcConn
	response, err := conn.AddDnaptRule(args)
	b, _ := json.Marshal(response)
	log.Printf("[DEBUG] conn.AddDnaptRule response: %s", b)
	if _, ok := err.(*common.APIError); ok {
		return fmt.Errorf("conn.AddDnaptRule error: %v", err)
	}

	dnatId := buildDnatId(args)

	log.Printf("[DEBUG] resourceTencentCloudDnatCreate dnatId: %s", dnatId)

	d.SetId(dnatId)
	return nil
}

func resourceTencentCloudDnatRead(d *schema.ResourceData, meta interface{}) error {

	_entry, err := parseDnatId(d.Id())

	if err != nil {
		d.SetId("")
		return err
	}

	client := meta.(*TencentCloudClient)
	entry, err := client.DescribeDnat(_entry)

	if err == dnatNotFound {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	d.Set("vpc_id", entry.VpcId)
	d.Set("nat_id", entry.NatId)
	d.Set("ip_protocol", entry.Proto)
	d.Set("external_ip", entry.Eip)
	d.Set("external_port", strconv.Itoa(*entry.Eport)) //todo: delete strconv.Itoa
	d.Set("internal_ip", entry.Pip)
	d.Set("internal_port", entry.Pport)
	return nil
}

func resourceTencentCloudDnatDelete(d *schema.ResourceData, meta interface{}) error {

	_entry, err := parseDnatId(d.Id())

	if err != nil {
		d.SetId("")
		return err
	}

	args := vpc.NewDeleteDnaptRuleRequest()
	args.VpcId = _entry.UniqVpcId
	args.NatId = _entry.UniqNatId
	args.DnatList = []*vpc.DnaptRuleInput{
		&vpc.DnaptRuleInput{
			Eip:   _entry.Eip,
			Eport: common.StringPtr(strconv.Itoa(*_entry.Eport)), //todo: delete strconv.Itoa
			Proto: _entry.Proto,
		},
	}

	conn := meta.(*TencentCloudClient).vpcConn
	response, err := conn.DeleteDnaptRule(args)
	b, _ := json.Marshal(response)
	log.Printf("[DEBUG] conn.DeleteDnaptRule response: %s", b)
	if _, ok := err.(*common.APIError); ok {
		return fmt.Errorf("conn.DeleteDnaptRule error: %v", err)
	}
	return nil
}

// Build an ID for a Forward Entry, eg "tcp://VpcId:NatId@127.15.2.3:8080"
func buildDnatId(entry *vpc.AddDnaptRuleRequest) (entryId string) {
	log.Printf("[DEBUG] args=%v", entry)
	entryId = fmt.Sprintf("%s://%s:%s@%s:%s", *entry.Proto, *entry.VpcId, *entry.NatId, *entry.Eip, *entry.Eport)
	log.Printf("[DEBUG] buildDnatId entryId=%s", entryId)
	return
}

//Parse Forward Entry ID
func parseDnatId(entryId string) (entry *vpc.DnaptRule, err error) {
	log.Printf("[DEBUG] parseDnatId entryId: %s", entryId)
	u, errors := url.Parse(entryId)
	if errors != nil {
		err = errors
		return
	}
	host, port, _ := net.SplitHostPort(u.Host)
	_port, _ := strconv.Atoi(port)
	natId, _ := u.User.Password()
	entry = &vpc.DnaptRule{}
	entry.UniqVpcId = common.StringPtr(u.User.Username())
	entry.UniqNatId = common.StringPtr(natId)
	entry.Proto = common.StringPtr(u.Scheme)
	entry.Eip = common.StringPtr(host)
	entry.Eport = common.IntPtr(_port) //todo: delete strconv.Atoi
	b, _ := json.Marshal(entry)
	log.Printf("[DEBUG] parseDnatId result: %s", b)
	return
}
