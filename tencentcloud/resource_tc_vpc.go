package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcCreate,
		Read:   resourceTencentCloudVpcRead,
		Update: resourceTencentCloudVpcUpdate,
		Delete: resourceTencentCloudVpcDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr_block": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"is_default": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_multicast": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudVpcCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":    "CreateVpc",
		"vpcName":   d.Get("name").(string),
		"cidrBlock": d.Get("cidr_block").(string),
	}
	response, err := client.SendRequest("vpc", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		Code      int    `json:tag"code"`
		Message   string `json:tag"message"`
		UniqVpcId string `json:tag"uniqVpcId"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("resource_tc_vpc create error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
	}
	log.Printf("[DEBUG] UniqVpcId=%v", jsonresp.UniqVpcId)
	d.SetId(jsonresp.UniqVpcId)
	return nil
}

func resourceTencentCloudVpcRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action": "DescribeVpcEx",
		"vpcId":  d.Id(),
	}
	response, err := client.SendRequest("vpc", params)
	if err != nil {
		d.SetId("")
		return err
	}

	var jsonresp struct {
		Code       int    `json:tag"code"`
		Message    string `json:tag"message"`
		TotalCount int    `json:tag"totalCount"`
		Data       []struct {
			VpcName     string `json:tag"vpcName"`
			CidrBlock   string `json:tag"cidrBlock"`
			IsDefault   bool   `json:tag"isDefault"`
			IsMulticast bool   `json:tag"isMulticast"`
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("resource_tc_vpc read error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
	} else if jsonresp.TotalCount == 0 {
		d.SetId("")
		return nil
	}

	vpc := jsonresp.Data[0]
	d.Set("name", vpc.VpcName)
	d.Set("cidr_block", vpc.CidrBlock)
	d.Set("is_default", vpc.IsDefault)
	d.Set("is_multicast", vpc.IsMulticast)
	return nil
}

func resourceTencentCloudVpcUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action": "ModifyVpcAttribute",
	}
	d.Partial(true)
	params["vpcId"] = d.Id()
	attributeUpdate := false
	if d.HasChange("name") {
		d.SetPartial("name")
		params["vpcName"] = d.Get("name").(string)
		attributeUpdate = true
	}
	if attributeUpdate {
		response, err := client.SendRequest("vpc", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Code    int    `json:tag"code"`
			Message string `json:tag"message"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return err
		}
		if jsonresp.Code != 0 {
			return fmt.Errorf("resource_tc_vpc update error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
		}
	}
	d.Partial(false)
	return resourceTencentCloudVpcRead(d, m)
}

func resourceTencentCloudVpcDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		params := map[string]string{
			"Action": "DeleteVpc",
			"vpcId":  d.Id(),
		}
		response, err := client.SendRequest("vpc", params)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("trying again while it is deleted."))
		}
		var jsonresp struct {
			Code     int    `json:tag"code"`
			Message  string `json:tag"message"`
			CodeDesc string `json:tag"codeDesc"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if jsonresp.CodeDesc == "InvalidVpc.CannotDelete" {
			return resource.RetryableError(fmt.Errorf(jsonresp.Message))
		} else if jsonresp.CodeDesc == "InvalidVpc.NotFound" {
			log.Printf("[DEBUG] Delete vpc faid failed, CodeDesc:InvalidVpc.NotFound, vpcId:%s", params["vpcId"])
		} else if jsonresp.Code != 0 {
			return resource.NonRetryableError(fmt.Errorf(jsonresp.Message))
		}
		return nil
	})
}
