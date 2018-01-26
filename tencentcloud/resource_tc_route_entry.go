package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

var nextTypes = map[string]int{
	"public_gateway":     0,
	"vpn_gateway":        1,
	"dc_gateway":         3,
	"peering_connection": 4,
	"sslvpn_gateway":     7,
	"nat_gateway":        8,
	"instance":           9,
}

func resourceTencentCloudRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRouteEntryCreate,
		Read:   resourceTencentCloudRouteEntryRead,
		Delete: resourceTencentCloudRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"next_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					_, ok := nextTypes[value]
					if !ok {
						var nextHubDesc []string
						for vgwKey, _ := range nextTypes {
							nextHubDesc = append(nextHubDesc, vgwKey)
						}
						errors = append(errors, fmt.Errorf("%s Only 1 of %s is allowed", k, strings.Join(nextHubDesc, ",")))
					}
					return
				},
			},
			"next_hub": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTencentCloudRouteEntryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	next_type := d.Get("next_type").(string)
	params := map[string]string{
		"Action":                          "CreateRoute",
		"vpcId":                           d.Get("vpc_id").(string),
		"routeTableId":                    d.Get("route_table_id").(string),
		"routeSet.0.destinationCidrBlock": d.Get("cidr_block").(string),
		"routeSet.0.nextType":             strconv.Itoa(nextTypes[next_type]),
		"routeSet.0.nextHub":              d.Get("next_hub").(string),
	}

	log.Printf("[DEBUG] resource_tc_route_entry create params:%v", params)

	response, err := client.SendRequest("vpc", params)
	if err != nil {
		log.Printf("[ERROR] resource_tc_route_entry create client.SendRequest error:%v", err)
		return err
	}
	var jsonresp struct {
		Code     int    `json:tag"code"`
		Message  string `json:tag"message"`
		CodeDesc string `json:tag"codeDesc"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[ERROR] resource_tc_route_entry create json.Unmarshal error:%v", err)
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("resource_tc_route_entry create error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
	}
	route := map[string]string{
		"vpcId":                params["vpcId"],
		"routeTableId":         params["routeTableId"],
		"destinationCidrBlock": params["routeSet.0.destinationCidrBlock"],
		"nextType":             params["routeSet.0.nextType"],
		"nextHub":              params["routeSet.0.nextHub"],
	}
	uniqRouteEntryId, ok := routeIdEncode(route)
	if ok == false {
		return errors.New("resource_tc_route_entry create error, Build route entry id fail")
	}
	log.Printf("[DEBUG] uniqRouteEntryId=%s", uniqRouteEntryId)
	d.SetId(uniqRouteEntryId)
	return nil
}

func resourceTencentCloudRouteEntryRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] resource_tc_route_entry read id:%v", d.Id())
	client := m.(*TencentCloudClient).commonConn
	_route, ok := routeIdDecode(d.Id())
	if ok == false {
		return fmt.Errorf("resource_tc_route_entry read error, id decode faild! id:%v", d.Id())
	}

	_route["Action"] = "DescribeRoutes"

	log.Printf("[DEBUG] resource_tc_route_entry read params:%v", _route)

	response, err := client.SendRequest("vpc", _route)
	if err != nil {
		log.Printf("[ERROR] resource_tc_route_entry read json.Unmarshal error:%v", err)
		return err
	}
	var jsonresp struct {
		Code     int    `json:tag"code"`
		Message  string `json:tag"message"`
		CodeDesc string `json:tag"codeDesc"`
		Data     struct {
			TotalNum int `json:tag"totalNum"`
			Data     []struct {
				DestinationCidrBlock string `json:tag"destinationCidrBlock"`
				NextType             int    `json:tag"nextType"`
				UnNextHub            string `json:tag"unNextHub"`
			}
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[DEBUG] resource_tc_route_entry read json.Unmarshal error:%v", err)
		return err
	}

	if jsonresp.Code != 0 {
		return fmt.Errorf("resource_tc_route_entry read error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
	} else if jsonresp.Data.TotalNum <= 0 || len(jsonresp.Data.Data) <= 0 {
		d.SetId("")
		return nil
	}

	nextType := -1
	for _, r := range jsonresp.Data.Data {
		if strings.ToUpper(r.UnNextHub) == "LOCAL" {
			continue
		}
		nextType = r.NextType
	}
	if nextType == -1 {
		d.SetId("")
		return nil
	}

	uniqNextType := ""
	for vgwKey, vgwType := range nextTypes {
		if vgwType == nextType {
			uniqNextType = vgwKey
			break
		}
	}
	if uniqNextType == "" {
		return errors.New("[ERROR] nextType is null")
	}

	d.Set("vpc_id", _route["vpcId"])
	d.Set("route_table_id", _route["routeTableId"])
	d.Set("cidr_block", _route["destinationCidrBlock"])
	d.Set("next_type", uniqNextType)
	d.Set("next_hub", _route["nextHub"])
	return nil
}

func resourceTencentCloudRouteEntryDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	routes, ok := routeIdDecode(d.Id())
	if ok == false {
		return fmt.Errorf("resource_tc_route_entry delete error, id decode faild! id:%v", d.Id())
	}
	params := map[string]string{
		"Action":                          "DeleteRoute",
		"vpcId":                           routes["vpcId"],
		"routeTableId":                    routes["routeTableId"],
		"routeSet.0.destinationCidrBlock": routes["destinationCidrBlock"],
		"routeSet.0.nextType":             routes["nextType"],
		"routeSet.0.nextHub":              routes["nextHub"],
	}

	log.Printf("[DEBUG] resource_tc_route_entry delete params:%v", params)

	response, err := client.SendRequest("vpc", params)
	if err != nil {
		log.Printf("[ERROR] resource_tc_route_entry delete client.SendRequest error:%v", err)
		return err
	}
	var jsonresp struct {
		Code     int    `json:tag"code"`
		Message  string `json:tag"message"`
		CodeDesc string `json:tag"codeDesc"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[ERROR] resource_tc_route_entry delete json.Unmarshal error:%v", err)
		return err
	}

	//jsonresp.Code == 28004 ==> Route table does not exist
	if jsonresp.Code != 0 && jsonresp.Code != 28004 {
		log.Printf("[DEBUG] resource_tc_route_entry delete error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
		return errors.New(jsonresp.Message)
	}
	return nil
}

// Build an ID for a route
func routeIdEncode(route map[string]string) (routeId string, ok bool) {
	vpcId, ok0 := route["vpcId"]
	rtbId, ok1 := route["routeTableId"]
	cidrBlock, ok2 := route["destinationCidrBlock"]
	nextType, ok3 := route["nextType"]
	nextHub, ok4 := route["nextHub"]
	if ok0 && ok1 && ok2 && ok3 && ok4 {
		return fmt.Sprintf("%v::%v::%v::%v::%v", vpcId, rtbId, cidrBlock, nextType, nextHub), true
	}
	return "", false
}

//Decompose a Route Id
func routeIdDecode(routeId string) (route map[string]string, ok bool) {
	route = map[string]string{}
	routeArray := strings.Split(routeId, "::")
	if len(routeArray) != 5 {
		return route, false
	}
	route["vpcId"] = routeArray[0]
	route["routeTableId"] = routeArray[1]
	route["destinationCidrBlock"] = routeArray[2]
	route["nextType"] = routeArray[3]
	route["nextHub"] = routeArray[4]
	log.Printf("[DEBUG] routeIdDecode routeId:%v", routeId)
	log.Printf("[DEBUG] routeIdDecode result route:%v", route)
	return route, true
}
