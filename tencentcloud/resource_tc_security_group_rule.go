package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSecurityGroupRuleCreate,
		Read:   resourceTencentCloudSecurityGroupRuleRead,
		Delete: resourceTencentCloudSecurityGroupRuleDelete,

		Schema: map[string]*schema.Schema{
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					value = strings.ToUpper(value)
					if value != "INGRESS" && value != "EGRESS" {
						errors = append(errors, fmt.Errorf("%s of rule, ingress (inbound) or egress (outbound) value:%v", k, value))
					}
					return
				},
			},
			"cidr_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					_, ip_err := validateIp(v, k)
					log.Printf("[DEBUG] validateIp ip_err:%v", ip_err)
					if len(ip_err) == 0 {
						return
					}
					_, cidr_err := validateCIDRNetworkAddress(v, k)
					log.Printf("[DEBUG] validateCIDRNetworkAddress cidr_err:%v", ip_err)
					if len(cidr_err) == 0 {
						return
					}
					errors = append(errors, fmt.Errorf("%s can be IP, or CIDR, otherwise it's invalid, value:%v", k, v))
					return
				},
			},
			"ip_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					value = strings.ToUpper(value)
					if value != "UDP" && value != "TCP" && value != "ICMP" {
						errors = append(errors, fmt.Errorf("%s support 'UDP', 'TCP', 'ICMP' and not configured means all protocols. But got %s", k, v))
					}
					return
				},
			},
			"port_range": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "example: 53、80,443、80-90",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					match, _ := regexp.MatchString("^(\\d{1,5},)*\\d{1,5}$|^\\d{1,5}\\-\\d{1,5}$", value)
					if !match {
						errors = append(errors, fmt.Errorf("%s example: 53、80,443、80-90, Not configured to represent all ports", k))
					}
					return
				},
			},
			"policy": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "accept" && value != "drop" {
						errors = append(errors, fmt.Errorf("Policy of rule, 'accept' or 'drop'"))
					}
					return
				},
			},
		},
	}
}

func resourceTencentCloudSecurityGroupRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":           "CreateSecurityGroupPolicy",
		"sgId":             d.Get("security_group_id").(string),
		"direction":        d.Get("type").(string),
		"index":            "-1",
		"policys.0.action": d.Get("policy").(string),
		"policys.0.cidrIp": d.Get("cidr_ip").(string),
	}
	if ip_protocol, ok := d.GetOk("ip_protocol"); ok {
		params["policys.0.ipProtocol"] = ip_protocol.(string)
	}
	if port_range, ok := d.GetOk("port_range"); ok {
		params["policys.0.portRange"] = port_range.(string)
	}

	log.Printf("[DEBUG] resource_tc_security_group_rule create params:%v", params)

	response, err := client.SendRequest("dfw", params)
	if err != nil {
		log.Printf("[ERROR] resource_tc_security_group_rule create client.SendRequest error:%v", err)
		return err
	}
	var jsonresp struct {
		Code     int    `json:tag"code"`
		Message  string `json:tag"message"`
		CodeDesc string `json:tag"codeDesc"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[ERROR] resource_tc_security_group_rule create json.Unmarshal error:%v", err)
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("resource_tc_security_group_rule create error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
	}

	rule := map[string]string{
		"sgId":       params["sgId"],
		"direction":  params["direction"],
		"action":     params["policys.0.action"],
		"cidrIp":     params["policys.0.cidrIp"],
		"ipProtocol": "ALL",
		"portRange":  "ALL",
	}
	if ipProtocol, ok := params["policys.0.ipProtocol"]; ok {
		rule["ipProtocol"] = ipProtocol
	}
	if portRange, ok := params["policys.0.portRange"]; ok {
		rule["portRange"] = portRange
	}

	uniqSecurityGroupRuleId := buildSecurityGroupRuleId(rule)
	log.Printf("[DEBUG] uniqSecurityGroupRuleId=%s", uniqSecurityGroupRuleId)
	d.SetId(uniqSecurityGroupRuleId)
	return nil
}

func resourceTencentCloudSecurityGroupRuleRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] resource_tc_security_group_rule read id:%v", d.Id())
	client := m.(*TencentCloudClient).commonConn
	rule, ok := parseSecurityGroupRuleId(d.Id())
	if ok == false {
		return fmt.Errorf("resource_tc_security_group_rule read error, id decode faild! id:%v", d.Id())
	}

	_, err := describeSecurityGroupRuleIndex(client, rule)
	if err != nil {
		if err == errSecurityGroupRuleNotFound {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("security_group_id", rule["sgId"])
	d.Set("type", rule["direction"])
	d.Set("cidr_ip", rule["cidrIp"])
	ipProtocol := strings.ToLower(rule["ipProtocol"])
	portRange := strings.ToLower(rule["portRange"])
	if ipProtocol != "all" {
		d.Set("ip_protocol", rule["ipProtocol"])
	}
	if portRange != "all" {
		d.Set("port_range", rule["portRange"])
	}
	d.Set("policy", rule["action"])

	return nil
}

func resourceTencentCloudSecurityGroupRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	rule, ok := parseSecurityGroupRuleId(d.Id())
	if ok == false {
		return fmt.Errorf("resource_tc_security_group_rule read error, id decode faild! id:%v", d.Id())
	}

	index, err := describeSecurityGroupRuleIndex(client, rule)
	if err != nil {
		return err
	}

	params := map[string]string{
		"Action":    "DeleteSecurityGroupPolicy",
		"sgId":      rule["sgId"],
		"direction": rule["direction"],
		"indexes.0": strconv.Itoa(index),
	}

	log.Printf("[DEBUG] resource_tc_security_group_rule delete params:%v", params)

	response, err := client.SendRequest("dfw", params)
	if err != nil {
		log.Printf("[ERROR] resource_tc_security_group_rule delete client.SendRequest error:%v", err)
		return err
	}
	var jsonresp struct {
		Code     int    `json:tag"code"`
		Message  string `json:tag"message"`
		CodeDesc string `json:tag"codeDesc"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[ERROR] resource_tc_security_group_rule delete json.Unmarshal error:%v", err)
		return err
	}

	if jsonresp.Code != 0 {
		log.Printf("[DEBUG] resource_tc_security_group_rule delete error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
		return errors.New(jsonresp.Message)
	}
	return nil
}

// Build an ID for a Security Group Rule
func buildSecurityGroupRuleId(rule map[string]string) (ruleId string) {
	log.Printf("[DEBUG] buildSecurityGroupRuleId before: %v", rule)
	var paramsArray []string
	for k, v := range rule {
		paramsArray = append(paramsArray, k+"="+v)
	}
	ruleId = strings.Join(paramsArray, "&")
	log.Printf("[DEBUG] buildSecurityGroupRuleId after: %v", ruleId)
	return
}

//Parse Security Group Rule ID
func parseSecurityGroupRuleId(ruleId string) (rule map[string]string, ok bool) {
	log.Printf("[DEBUG] parseSecurityGroupRuleId before: %v", ruleId)
	ok = true
	rule = map[string]string{}
	ruleQueryStrings := strings.Split(ruleId, "&")
	if len(ruleQueryStrings) == 0 {
		ok = false
		return
	}
	for _, str := range ruleQueryStrings {
		arr := strings.Split(str, "=")
		if len(arr) != 2 {
			ok = false
			return
		}
		rule[arr[0]] = arr[1]
	}
	log.Printf("[DEBUG] parseSecurityGroupRuleId after: %v", rule)
	return
}
