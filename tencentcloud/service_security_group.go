package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/zqfan/tencentcloud-sdk-go/client"
)

var (
	errSecurityGroupRuleNotFound = errors.New("security group rule index not found")
)

func describeSecurityGroupRuleIndex(client *client.Client, rule map[string]string) (index int, err error) {
	if rule["sgId"] == "" {
		err = fmt.Errorf("describeSecurityGroupRuleIndex, sgId empty")
		return
	}
	params := map[string]string{
		"Action": "DescribeSecurityGroupPolicys",
		"sgId":   rule["sgId"],
	}
	var response string
	response, err = client.SendRequest("dfw", params)
	if err != nil {
		log.Printf("[ERROR] describeSecurityGroupRuleIndex client.SendRequest error:%v", err)
		return
	}

	var jsonresp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Ingress []struct {
				Index      int    `json:"index"`
				CidrIp     string `json:"cidrIp"`
				IpProtocol string `json:"ipProtocol"`
				PortRange  string `json:"portRange"`
				Action     string `json:"action"`
			}
			Egress []struct {
				Index      int    `json:"index"`
				CidrIp     string `json:"cidrIp"`
				IpProtocol string `json:"ipProtocol"`
				PortRange  string `json:"portRange"`
				Action     string `json:"action"`
			}
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[ERROR] describeSecurityGroupRules json.Unmarshal error:%v", err)
		return
	}

	if jsonresp.Code != 0 {
		log.Printf("[ERROR] describeSecurityGroupRuleIndex error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
		err = fmt.Errorf(jsonresp.Message)
		return
	}

	rules := jsonresp.Data
	_rule := rules.Ingress
	if rule["direction"] == "egress" {
		_rule = rules.Egress
	}

	exist := false

	for _, r := range _rule {

		r.IpProtocol = strings.ToLower(r.IpProtocol)
		r.PortRange = strings.ToLower(r.PortRange)
		r.Action = strings.ToLower(r.Action)
		rule["ipProtocol"] = strings.ToLower(rule["ipProtocol"])
		rule["portRange"] = strings.ToLower(rule["portRange"])
		rule["action"] = strings.ToLower(rule["action"])

		if r.CidrIp == rule["cidrIp"] && r.IpProtocol == rule["ipProtocol"] && r.PortRange == rule["portRange"] && r.Action == rule["action"] {
			exist = true
			index = r.Index
			break
		}
	}

	if exist == false {
		err = errSecurityGroupRuleNotFound
	}

	return
}

func getSecurityGroupAssociatedInstancesBySgId(client *client.Client, sgId string) (instanceIds []string, err error) {
	params := map[string]string{
		"Action": "DescribeInstancesOfSecurityGroup",
		"sgId":   sgId,
	}

	response, err := client.SendRequest("dfw", params)
	if err != nil {
		return
	}

	var jsonresp struct {
		Code       int    `json:"code"`
		Message    string `json:"message"`
		TotalCount int    `json:"totalCount"`
		Data       []struct {
			InstanceId string `json:"instanceId"`
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return
	}

	if jsonresp.Code != 0 {
		err = fmt.Errorf(jsonresp.Message)
		return
	}
	if jsonresp.TotalCount == 0 {
		return
	}
	items := jsonresp.Data
	for _, item := range items {
		instanceIds = append(instanceIds, item.InstanceId)
	}
	return
}
