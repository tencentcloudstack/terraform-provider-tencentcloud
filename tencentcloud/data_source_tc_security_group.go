package tencentcloud

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(2, 100),
			},
			"create_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"be_associate_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceTencentCloudSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":    "DescribeSecurityGroupEx",
		"projectId": strconv.Itoa(projectId),
		"sgId":      d.Get("security_group_id").(string),
	}

	log.Printf("[DEBUG] resource_tc_security_group read params:%v", params)

	response, err := client.SendRequest("dfw", params)
	if err != nil {
		log.Printf("[ERROR] resource_tc_security_group read client.SendRequest error:%v", err)
		return err
	}

	var jsonresp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			TotalNum int `json:"totalNum"`
			Detail   []struct {
				SgId             string `json:"sgId"`
				SgName           string `json:"sgName"`
				SgRemark         string `json:"sgRemark"`
				BeAssociateCount string `json:"beAssociateCount"`
				CreateTime       string `json:"createTime"`
			}
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[ERROR] resource_tc_security_group read json.Unmarshal error:%v", err)
		return err
	}
	if jsonresp.Code != 0 {
		log.Printf("[ERROR] resource_tc_security_group read error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
		return errors.New(jsonresp.Message)
	} else if jsonresp.Data.TotalNum <= 0 || len(jsonresp.Data.Detail) <= 0 {
		return errors.New("Security group not found")
	}

	sg := jsonresp.Data.Detail[0]

	d.Set("security_group_id", sg.SgId)
	d.Set("name", sg.SgName)
	d.Set("description", sg.SgRemark)
	d.Set("create_time", sg.CreateTime)
	d.Set("be_associate_count", sg.BeAssociateCount)

	return nil
}
