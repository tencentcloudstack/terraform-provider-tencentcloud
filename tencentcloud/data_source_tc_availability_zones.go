package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	tencentCloudApiAvailibilityZoneStateAvailable   = "AVAILABLE"
	tencentCloudApiAvailibilityZoneStateUnavailable = "UNAVAILABLE"
)

func dataSourceTencentCloudAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAvailabilityZonesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"include_unavailable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			// Computed values.
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAvailabilityZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).commonConn

	params := map[string]string{
		"Version": "2017-03-12",
		"Action":  "DescribeZones",
	}

	log.Printf("[DEBUG] tencentcloud_instance_types - param: %v", params)
	response, err := client.SendRequest("cvm", params)
	if err != nil {
		return err
	}

	type Zone struct {
		Zone      string `json:"Zone"`
		ZoneName  string `json:"ZoneName"`
		ZoneId    string `json:"ZoneId"`
		ZoneState string `json:"ZoneState"`
	}
	var jsonresp struct {
		Response struct {
			Error struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			}
			RequestId string `json:"RequestId"`
			ZoneSet   []Zone
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Response.Error.Code != "" {
		return fmt.Errorf(
			"tencentcloud_availability_zones got error, code:%v, message:%v",
			jsonresp.Response.Error.Code,
			jsonresp.Response.Error.Message,
		)
	}

	var (
		resultZoneList []Zone
	)
	zoneList := jsonresp.Response.ZoneSet
	if len(zoneList) == 0 {
		return errors.New("No avalability zones found")
	}

	name, nameOk := d.GetOk("name")
	includeUnavailable, includeUnavailableOk := d.GetOk("include_unavailable")
	for _, zone := range zoneList {
		log.Printf(
			"[DEBUG] tencentcloud_availability_zones - Zone found id: %v, name:% v, description: %v, state: %v",
			zone.ZoneId,
			zone.Zone,
			zone.ZoneName,
			zone.ZoneState,
		)

		if zone.ZoneState == tencentCloudApiAvailibilityZoneStateUnavailable {
			if !includeUnavailableOk || !includeUnavailable.(bool) {
				continue
			}
		}

		if nameOk {
			zoneName := name.(string)
			if zone.Zone == zoneName {
				resultZoneList = append(resultZoneList, zone)
			}
			continue
		}
		resultZoneList = append(resultZoneList, zone)
	}

	if len(resultZoneList) == 0 {
		return errors.New("No avalability zones found")
	}

	var (
		result    []map[string]interface{}
		resultIds []string
	)

	for _, zone := range resultZoneList {
		m := make(map[string]interface{})
		m["id"] = zone.ZoneId
		m["name"] = zone.Zone
		m["description"] = zone.ZoneName
		m["state"] = zone.ZoneState
		result = append(result, m)
		resultIds = append(resultIds, zone.ZoneId)
	}
	id := dataResourceIdsHash(resultIds)
	d.SetId(id)
	log.Printf("[DEBUG] tencentcloud_availability_zones - instances[0]: %#v", result[0])
	if err := d.Set("zones", result); err != nil {
		return err
	}
	return nil
}
