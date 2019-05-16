package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cl "github.com/zqfan/tencentcloud-sdk-go/client"
)

func resourceTencentCloudSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSubnetCreate,
		Read:   resourceTencentCloudSubnetRead,
		Update: resourceTencentCloudSubnetUpdate,
		Delete: resourceTencentCloudSubnetDelete,

		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTencentCloudSubnetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":                 "CreateSubnet",
		"vpcId":                  d.Get("vpc_id").(string),
		"subnetSet.0.subnetName": d.Get("name").(string),
		"subnetSet.0.cidrBlock":  d.Get("cidr_block").(string),
		"subnetSet.0.zoneId":     d.Get("availability_zone").(string),
	}
	response, err := client.SendRequest("vpc", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		CodeDesc  string `json:"codeDesc"`
		SubnetSet []struct {
			UnSubnetId string `json:"unSubnetId"`
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("resource_tc_subnet got error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
	}
	subnet := jsonresp.SubnetSet[0]
	log.Printf("[DEBUG] UnSubnetId=%s", subnet.UnSubnetId)

	if routeTableId := d.Get("route_table_id").(string); routeTableId != "" {
		if err = bindRouteTableIdForSubnetId(client, d.Get("vpc_id").(string), subnet.UnSubnetId, routeTableId); err != nil {
			return err
		}
	}

	d.SetId(subnet.UnSubnetId)
	return nil
}

func bindRouteTableIdForSubnetId(client *cl.Client, vpcId, unSubnetId, routeTableId string) (err error) {
	params := map[string]string{
		"Action":       "AssociateRouteTable",
		"vpcId":        vpcId,
		"subnetId":     unSubnetId,
		"routeTableId": routeTableId,
	}
	response, err := client.SendRequest("vpc", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("AssociateRouteTable got error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
	}
	return
}

func getMyNewTableId(client *cl.Client, vpcId, oldTableId string) (routeTableId string, errRet error) {
	var jsonresp struct {
		Code       int                      `json:"code"`
		Message    string                   `json:"message"`
		TotalCount int                      `json:"totalCount"`
		Data       []map[string]interface{} `json:"data"`
	}
	params := map[string]string{
		"Action":       "DescribeRouteTable",
		"vpcId":        vpcId,
		"limit":        "1",
		"routeTableId": oldTableId,
	}

	response, err := client.SendRequest("vpc", params)
	if err != nil {
		errRet = err
		return
	}
	if err = json.Unmarshal([]byte(response), &jsonresp); err != nil {
		errRet = err
		return
	}

	if jsonresp.Code != 0 {
		errRet = fmt.Errorf("DescribeRouteTable got error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
		return
	}

	if len(jsonresp.Data) != 1 {
		errRet = fmt.Errorf("DescribeRouteTable got  %d table info.", len(jsonresp.Data))
		return
	}

	unRouteTableId, ok := jsonresp.Data[0]["unRouteTableId"]
	if !ok {
		errRet = fmt.Errorf("DescribeRouteTable got  nil unRouteTableId from table info.")
		return
	}
	routeTableId = unRouteTableId.(string)
	return
}

func getMyDefaultTableId(client *cl.Client, vpcId string) (routeTableId string, errRet error) {
	var (
		offset    int64 = 0
		limit     int64 = 50
		leftCount int64 = -1
	)
	var jsonresp struct {
		Code       int                      `json:"code"`
		Message    string                   `json:"message"`
		TotalCount int                      `json:"totalCount"`
		Data       []map[string]interface{} `json:"data"`
	}

again:
	params := map[string]string{
		"Action": "DescribeRouteTable",
		"vpcId":  vpcId,
		"offset": fmt.Sprintf("%d", offset),
		"limit":  fmt.Sprintf("%d", limit),
	}
	response, err := client.SendRequest("vpc", params)
	if err != nil {
		errRet = err
		return
	}
	if err = json.Unmarshal([]byte(response), &jsonresp); err != nil {
		errRet = err
		return
	}

	if jsonresp.Code != 0 {
		errRet = fmt.Errorf("DescribeRouteTable got error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
		return
	}
	for _, dataMap := range jsonresp.Data {
		unRouteTableId, uok := dataMap["unRouteTableId"]
		routeTableType, rok := dataMap["routeTableType"]
		if uok && rok {
			strUnRouteTableId, uok := unRouteTableId.(string)
			strRouteTableType := fmt.Sprintf("%v", routeTableType)
			if uok && rok && strRouteTableType == "1" {
				routeTableId = strUnRouteTableId
				return
			}
		}
	}
	if leftCount == -1 {
		leftCount = int64(jsonresp.TotalCount)
	}
	leftCount -= limit
	offset += limit
	if leftCount <= 0 {
		errRet = fmt.Errorf("no default route table was found.")
		return
	}
	goto again
}

func resourceTencentCloudSubnetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":   "DescribeSubnet",
		"vpcId":    d.Get("vpc_id").(string),
		"subnetId": d.Id(),
	}
	response, err := client.SendRequest("vpc", params)
	if err != nil {
		return err
	}

	var jsonresp struct {
		Code         int    `json:"code"`
		Message      string `json:"message"`
		CodeDesc     string `json:"codeDesc"`
		SubnetName   string `json:"subnetName"`
		CidrBlock    string `json:"cidrBlock"`
		RouteTableId string `json:"routeTableId"`
		ZoneId       int    `json:"zoneId"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		if jsonresp.CodeDesc == "InvalidSubnet.NotFound" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("resource_tc_subnet got error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
	}

	d.Set("cidr_block", jsonresp.CidrBlock)
	d.Set("name", jsonresp.SubnetName)

	if d.Get("route_table_id").(string) != "" {
		newTableId, err := getMyNewTableId(client, d.Get("vpc_id").(string), jsonresp.RouteTableId)
		if err != nil {
			return err
		}
		d.Set("route_table_id", newTableId)
	}
	return nil
}

func resourceTencentCloudSubnetUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":   "ModifySubnetAttribute",
		"vpcId":    d.Get("vpc_id").(string),
		"subnetId": d.Id(),
	}

	d.Partial(true)

	if d.HasChange("name") {
		params["subnetName"] = d.Get("name").(string)
		response, err := client.SendRequest("vpc", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return err
		}
		if jsonresp.Code != 0 {
			return fmt.Errorf("resource_tc_subnet got error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
		}
		d.SetPartial("name")
	}

	if d.HasChange("route_table_id") {
		routeTableId := d.Get("route_table_id").(string)

		if routeTableId == "" {
			//to default table id.
			defaultRouteTable, err := getMyDefaultTableId(client, d.Get("vpc_id").(string))
			if err != nil {
				return err
			}
			routeTableId = defaultRouteTable
		}

		if err := bindRouteTableIdForSubnetId(client, d.Get("vpc_id").(string), d.Id(), routeTableId); err != nil {
			return err
		}
		d.SetPartial("route_table_id")
	}

	d.Partial(false)

	return resourceTencentCloudSubnetRead(d, m)
}

func resourceTencentCloudSubnetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		params := map[string]string{
			"Action":   "DeleteSubnet",
			"vpcId":    d.Get("vpc_id").(string),
			"subnetId": d.Id(),
		}
		response, err := client.SendRequest("vpc", params)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("trying again while it is deleted."))
		}
		var jsonresp struct {
			Code     int    `json:"code"`
			Message  string `json:"message"`
			CodeDesc string `json:"codeDesc"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if jsonresp.CodeDesc == "InvalidSubnet.CannotDelete" {
			return resource.RetryableError(fmt.Errorf(jsonresp.Message))
		} else if jsonresp.CodeDesc == "InvalidVpc.NotFound" || jsonresp.CodeDesc == "InvalidSubnet.NotFound" {
			log.Printf("[DEBUG] Delete subnet faid failed, CodeDesc:%s, vpcId:%s, subnetId:%s", jsonresp.CodeDesc, params["vpcId"], params["subnetId"])
		} else if jsonresp.Code != 0 {
			return resource.NonRetryableError(fmt.Errorf(jsonresp.Message))
		}
		return nil
	})
}
