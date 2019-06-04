package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSecurityGroupCreate,
		Read:   resourceTencentCloudSecurityGroupRead,
		Update: resourceTencentCloudSecurityGroupUpdate,
		Delete: resourceTencentCloudSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(2, 100),
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTencentCloudSecurityGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action": "CreateSecurityGroup",
	}
	if _, ok := d.GetOk("name"); ok {
		params["sgName"] = d.Get("name").(string)
	}
	if _, ok := d.GetOk("description"); ok {
		params["sgRemark"] = d.Get("description").(string)
	}
	if _, ok := d.GetOk("project_id"); ok {
		params["projectId"] = strconv.Itoa(d.Get("project_id").(int))
	}

	log.Printf("[DEBUG] resource_tc_security_group create params:%v", params)

	response, err := client.SendRequest("dfw", params)
	if err != nil {
		log.Printf("[ERROR] resource_tc_security_group create client.SendRequest error:%v", err)
		return err
	}
	var jsonresp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			SgId string `json:"sgId"`
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[ERROR] resource_tc_security_group create json.Unmarshal error:%v", err)
		return err
	}
	if jsonresp.Code != 0 {
		log.Printf("[ERROR] resource_tc_security_group create error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
		return errors.New(jsonresp.Message)
	}
	log.Printf("[DEBUG] SgId=%s", jsonresp.Data.SgId)
	d.SetId(jsonresp.Data.SgId)
	return nil
}

func resourceTencentCloudSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action": "DescribeSecurityGroupEx",
		"sgId":   d.Id(),
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
				SgName    string `json:"sgName"`
				SgRemark  string `json:"sgRemark"`
				ProjectId int    `json:"projectId"`
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
		d.SetId("")
		return nil
	}
	sg := jsonresp.Data.Detail[0]
	d.Set("name", sg.SgName)
	d.Set("description", sg.SgRemark)
	d.Set("project_id", sg.ProjectId)
	return nil
}

func resourceTencentCloudSecurityGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action": "ModifySecurityGroupAttributes",
		"sgId":   d.Id(),
	}

	d.Partial(true)

	attributeUpdate := false

	if d.HasChange("name") {
		params["sgName"] = d.Get("name").(string)
		d.SetPartial("name")
		attributeUpdate = true
	}

	if d.HasChange("description") {
		params["sgRemark"] = d.Get("description").(string)
		d.SetPartial("description")
		attributeUpdate = true
	}

	if attributeUpdate {

		log.Printf("[DEBUG] resource_tc_security_group update params:%v", params)

		response, err := client.SendRequest("dfw", params)
		if err != nil {
			log.Printf("[ERROR] resource_tc_security_group update client.SendRequest error:%v", err)
			return err
		}
		var jsonresp struct {
			Code     int    `json:"code"`
			Message  string `json:"message"`
			CodeDesc string `json:"codeDesc"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			log.Printf("[ERROR] resource_tc_security_group update json.Unmarshal error:%v", err)
			return err
		}
		if jsonresp.Code != 0 {
			log.Printf("[ERROR] resource_tc_security_group update error, code:%v, message:%v, CodeDesc:%v", jsonresp.Code, jsonresp.Message, jsonresp.CodeDesc)
			return errors.New(jsonresp.Message)
		}
	}

	d.Partial(false)

	return resourceTencentCloudSecurityGroupRead(d, m)
}

func resourceTencentCloudSecurityGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	var sgId = d.Id()

	params := map[string]string{
		"Action": "DeleteSecurityGroup",
		"sgId":   sgId,
	}

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		var (
			err         error
			instanceIds []string
			response    string
		)
		instanceIds, err = getSecurityGroupAssociatedInstancesBySgId(client, sgId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(instanceIds) > 0 {
			s := fmt.Sprintf(
				"security group: %v still bind with instanceIds: %v",
				sgId,
				instanceIds,
			)
			log.Printf("[DEBUG] %v", s)
			err = fmt.Errorf(s)
			return resource.RetryableError(err)
		}

		response, err = client.SendRequest("dfw", params)
		if err != nil {
			log.Printf("[DEBUG] resource_tc_route_table delete client.SendRequest error:%v", err)
			return resource.NonRetryableError(err)
		}
		var jsonresp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			log.Printf("[ERROR] resource_tc_security_group delete client.SendRequest error:%v", err)
			return resource.NonRetryableError(err)
		}

		//The security group does not exist, doc: https://intl.cloud.tencent.com/document/api/213/1362
		if jsonresp.Code == 7001 {
			return nil
		} else if jsonresp.Code != 0 {
			log.Printf("[ERROR] resource_tc_security_group delete error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
			err = errors.New(jsonresp.Message)
			return resource.NonRetryableError(err)
		}

		return nil
	})
}
