package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
				Computed: true,
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
	d.SetId(subnet.UnSubnetId)
	return nil
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
	d.Set("route_table_id", jsonresp.RouteTableId)
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
