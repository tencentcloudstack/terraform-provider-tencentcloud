package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSubnetRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTencentCloudSubnetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	subnetParams := map[string]string{
		"Action":   "DescribeSubnet",
		"vpcId":    d.Get("vpc_id").(string),
		"subnetId": d.Get("subnet_id").(string),
	}
	response, err := client.SendRequest("vpc", subnetParams)
	if err != nil {
		return err
	}

	var jsonresp struct {
		Code             int    `json:"code"`
		Message          string `json:"message"`
		CodeDesc         string `json:"codeDesc"`
		SubnetName       string `json:"subnetName"`
		CidrBlock        string `json:"cidrBlock"`
		RouteTableId     string `json:"routeTableId"`
		ZoneId           int    `json:"zoneId"`
		AvailabilityZone string
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.CodeDesc == "InvalidSubnet.NotFound" {
		return fmt.Errorf("Subnet Not Found")
	} else if jsonresp.ZoneId == 0 {
		return fmt.Errorf("data_source_tc_subnet got error, ZoneId is empty")
	} else if jsonresp.Code != 0 {
		return fmt.Errorf("data_source_tc_subnet got error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
	}

	subnet := jsonresp
	subnetZoneId := strconv.Itoa(subnet.ZoneId)

	//[zone] transform [availability_zone], eg: 100002 --> ap-guangzhou-2
	params := map[string]string{
		"Version": "2017-03-12",
		"Action":  "DescribeZones",
	}
	response, err = client.SendRequest("cvm", params)
	if err != nil {
		return err
	}

	type Zone struct {
		Zone      string `json:"Zone"`
		ZoneName  string `json:"ZoneName"`
		ZoneId    string `json:"ZoneId"`
		ZoneState string `json:"ZoneState"`
	}
	var ZoneJsonresp struct {
		Response struct {
			Error struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			}
			RequestId string `json:"RequestId"`
			ZoneSet   []Zone
		}
	}
	err = json.Unmarshal([]byte(response), &ZoneJsonresp)
	if err != nil {
		return err
	}
	if ZoneJsonresp.Response.Error.Code != "" {
		return fmt.Errorf(
			"tencentcloud_availability_zones got error, code:%v, message:%v",
			ZoneJsonresp.Response.Error.Code,
			ZoneJsonresp.Response.Error.Message,
		)
	}

	if len(ZoneJsonresp.Response.ZoneSet) == 0 {
		return errors.New("No avalability zones found")
	}

	for _, zone := range ZoneJsonresp.Response.ZoneSet {
		log.Printf(
			"[DEBUG] tencentcloud_availability_zones - Zone found id: %v, name:% v, description: %v, state: %v",
			zone.ZoneId,
			zone.Zone,
			zone.ZoneName,
			zone.ZoneState,
		)

		if zone.ZoneId == subnetZoneId {
			subnet.AvailabilityZone = zone.Zone
		}
	}

	if subnet.AvailabilityZone == "" {
		return errors.New("No avalability zones found")
	}

	d.SetId(subnetParams["subnetId"])
	d.Set("cidr_block", subnet.CidrBlock)
	d.Set("name", subnet.SubnetName)
	d.Set("route_table_id", subnet.RouteTableId)
	d.Set("availability_zone", subnet.AvailabilityZone)
	return nil
}
