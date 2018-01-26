package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRouteTableCreate,
		Read:   resourceTencentCloudRouteTableRead,
		Update: resourceTencentCloudRouteTableUpdate,
		Delete: resourceTencentCloudRouteTableDelete,

		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTencentCloudRouteTableCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":         "CreateRouteTable",
		"vpcId":          d.Get("vpc_id").(string),
		"routeTableName": d.Get("name").(string),
	}
	response, err := client.SendRequest("vpc", params)
	if err != nil {
		log.Printf("[ERROR] resource_tc_route_table create json.Unmarshal error:%v", err)
		return err
	}
	var jsonresp struct {
		Code           int    `json:tag"code"`
		Message        string `json:tag"message"`
		CodeDesc       string `json:tag"codeDesc"`
		UnRouteTableId string `json:tag"unRouteTableId"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("resource_tc_route_table create error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
	}
	log.Printf("[DEBUG] UnRouteTableId=%s", jsonresp.UnRouteTableId)
	d.SetId(jsonresp.UnRouteTableId)
	return nil
}

func resourceTencentCloudRouteTableRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":       "DescribeRouteTable",
		"vpcId":        d.Get("vpc_id").(string),
		"routeTableId": d.Id(),
	}
	response, err := client.SendRequest("vpc", params)
	if err != nil {
		return err
	}

	var jsonresp struct {
		Code     int    `json:tag"code"`
		Message  string `json:tag"message"`
		CodeDesc string `json:tag"codeDesc"`
		Data     []struct {
			RouteTableName string `json:tag"routeTableName"`
		} `json:tag"data"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[ERROR] resource_tc_route_table read json.Unmarshal error:%v", err)
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("resource_tc_route_table read error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
	} else if jsonresp.CodeDesc == "InvalidRouteTableId.NotFound" {
		d.SetId("")
		return nil
	}

	d.Set("name", jsonresp.Data[0].RouteTableName)

	return nil
}

func resourceTencentCloudRouteTableUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":       "ModifyVPCRouteTable",
		"vpcId":        d.Get("vpc_id").(string),
		"routeTableId": d.Id(),
	}

	d.Partial(true)

	if d.HasChange("name") {
		params["routeTableName"] = d.Get("name").(string)
		response, err := client.SendRequest("vpc", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Code     int    `json:tag"code"`
			Message  string `json:tag"message"`
			CodeDesc string `json:tag"codeDesc"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			log.Printf("[ERROR] resource_tc_route_table update json.Unmarshal error:%v", err)
			return err
		}
		if jsonresp.Code != 0 {
			return fmt.Errorf("resource_tc_route_table update error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
		}
		d.SetPartial("name")
	}

	d.Partial(false)

	return resourceTencentCloudRouteTableRead(d, m)
}

func resourceTencentCloudRouteTableDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":       "DeleteRouteTable",
		"vpcId":        d.Get("vpc_id").(string),
		"routeTableId": d.Id(),
	}

	log.Printf("[DEBUG] resource_tc_route_table delete params:%v", params)

	response, err := client.SendRequest("vpc", params)
	if err != nil {
		log.Printf("[DEBUG] resource_tc_route_table delete client.SendRequest error:%v", err)
		return err
	}
	var jsonresp struct {
		Code     int    `json:tag"code"`
		Message  string `json:tag"message"`
		CodeDesc string `json:tag"codeDesc"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[DEBUG] resource_tc_route_table delete json.Unmarshal error:%v", err)
		return err
	}

	if jsonresp.Code != 0 {
		log.Printf("[DEBUG] resource_tc_route_table delete error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
		if jsonresp.CodeDesc != "InvalidRouteTableId.NotFound" {
			return fmt.Errorf(jsonresp.Message)
		}
	}

	return nil
}
