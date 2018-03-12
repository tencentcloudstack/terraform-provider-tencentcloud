package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudRouteTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRouteTableRead,

		Schema: map[string]*schema.Schema{
			"route_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 60 {
						errors = append(errors, fmt.Errorf("%s Cannot be less than 2 characters and cannot be longer than 60 characters", k))
					}
					return
				},
			},
			"subnet_num": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hub": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"create_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTencentCloudRouteTableRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":       "DescribeRouteTable",
		"routeTableId": d.Get("route_table_id").(string),
	}
	if _, ok := d.GetOk("name"); ok {
		params["routeTableName"] = d.Get("name").(string)
	}

	log.Printf("[DEBUG] data_source_tc_route_table read params:%v", params)

	response, err := client.SendRequest("vpc", params)
	if err != nil {
		log.Printf("[DEBUG] data_source_tc_route_table read client.SendRequest error:%v", err)
		return err
	}

	var jsonresp struct {
		Code       int    `json:"code"`
		Message    string `json:"message"`
		CodeDesc   string `json:"codeDesc"`
		TotalCount int    `json:"totalCount"`
		Data       []struct {
			UnVpcId              string `json:"UnVpcId"`
			UnRouteTableId       string `json:"unRouteTableId"`
			RouteTableName       string `json:"routeTableName"`
			RouteTableCreateTime string `json:"routeTableCreateTime"`
			SubnetNum            int    `json:"subnetNum"`
			RouteTableSet        []struct {
				DestinationCidrBlock string `json:"destinationCidrBlock"`
				NextType             int    `json:"nextType"`
				UnNextHub            string `json:"unNextHub"`
				Description          string `json:"description"`
			} `json:"routeSet"`
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[DEBUG] data_source_tc_route_table read json.Unmarshal error:%v", err)
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("data_source_tc_route_table read error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
	} else if jsonresp.TotalCount <= 0 || len(jsonresp.Data) <= 0 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	rt := jsonresp.Data[0]

	log.Printf("[DEBUG] data_source_tc_route_table read result route_table_id:%v", rt.UnRouteTableId)

	d.SetId(rt.UnRouteTableId)
	d.Set("route_table_id", rt.UnRouteTableId)
	d.Set("vpc_id", rt.UnVpcId)
	d.Set("name", rt.RouteTableName)
	d.Set("subnet_num", rt.SubnetNum)
	d.Set("create_time", rt.RouteTableCreateTime)

	routes := make([]map[string]interface{}, 0, len(rt.RouteTableSet))
	for _, r := range rt.RouteTableSet {
		if strings.ToUpper(r.UnNextHub) == "LOCAL" {
			continue
		}
		m := make(map[string]interface{})
		for vgwKey, vgwType := range nextTypes {
			if vgwType == r.NextType {
				m["next_type"] = vgwKey
				break
			}
		}
		m["next_hub"] = r.UnNextHub
		m["cidr_block"] = r.DestinationCidrBlock
		m["description"] = r.Description
		routes = append(routes, m)
	}

	if err := d.Set("routes", routes); err != nil {
		log.Printf("[DEBUG] data_source_tc_route_table read d.Set routes error:%v", err)
		return err
	}

	return nil
}
