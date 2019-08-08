package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudVpc() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.10.0. Please use 'tencentcloud_vpc_instances' instead.",
		Read:               dataSourceTencentCloudVpcRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_multicast": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceTencentCloudVpcRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action": "DescribeVpcEx",
		"offset": "0",
		"limit":  "100",
	}
	if cid, ok := d.GetOk("id"); ok {
		params["vpcId"] = cid.(string)
	}
	if name, ok := d.GetOk("name"); ok && name != "" {
		params["vpcName"] = name.(string)
	}

	log.Printf("[DEBUG] Reading TencentCloud VPC: %s", params)
	response, err := client.SendRequest("vpc", params)
	if err != nil {
		return err
	}

	var jsonresp struct {
		Code       int    `json:"code"`
		Message    string `json:"message"`
		TotalCount int    `json:"totalCount"`
		Data       []struct {
			UnVpcId     string `json:"unVpcId"`
			VpcName     string `json:"vpcName"`
			CidrBlock   string `json:"cidrBlock"`
			IsDefault   bool   `json:"isDefault"`
			IsMulticast bool   `json:"isMulticast"`
		} `json:"data"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("resource_tc_vpc got error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
	} else if jsonresp.TotalCount == 0 {
		log.Printf("[DEBUG] VPC Not Found, params=%s", params)
		return nil
	}

	if len(jsonresp.Data) == 0 {
		return fmt.Errorf("no matching VPC found")
	} else if len(jsonresp.Data) > 1 {
		return fmt.Errorf("multiple VPCs matched; use additional constraints to reduce matches to a single VPC")
	}

	vpc := jsonresp.Data[0]
	d.SetId(vpc.UnVpcId)
	d.Set("name", vpc.VpcName)
	d.Set("cidr_block", vpc.CidrBlock)
	d.Set("is_default", vpc.IsDefault)
	d.Set("is_multicast", vpc.IsMulticast)
	return nil
}
