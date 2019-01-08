package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudCbsStorageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsStorageAttachmentCreate,
		Read:   resourceTencentCloudCbsStorageAttachmentRead,
		Delete: resourceTencentCloudCbsStorageAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"storage_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTencentCloudCbsStorageAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn

	storageId := d.Get("storage_id").(string)
	instanceId := d.Get("instance_id").(string)
	params := map[string]string{
		"Action":       "AttachCbsStorages",
		"storageIds.0": storageId,
		"uInstanceId":  instanceId,
	}

	response, err := client.SendRequest("cbs", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		Code     int
		Message  string
		CodeDesc string
		Detail   interface{}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("[ERROR] code=%s msg=%s", jsonresp.CodeDesc, jsonresp.Message)
	}
	attachStatus := jsonresp.Detail.(map[string]interface{})
	for _, detail := range attachStatus {
		taskDetail := detail.(map[string]interface{})
		taskCode := int(taskDetail["code"].(float64))
		if taskCode != 0 {
			return fmt.Errorf("[ERROR] code=%v msg=%s", taskCode, taskDetail["msg"].(string))
		}
	}
	d.SetId("att::" + instanceId + "::" + storageId)
	return nil
}

func resourceTencentCloudCbsStorageAttachmentRead(d *schema.ResourceData, m interface{}) error {
	instanceId, storageId, err := parseAttachmentId(d.Id())
	if err != nil {
		return err
	}

	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action":       "DescribeCbsStorages",
		"storageIds.0": storageId,
	}

	response, err := client.SendRequest("cbs", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		Code       int    `json:"code"`
		Message    string `json:"message"`
		StorageSet []struct {
			Attached      int    `json:"attached"`
			StorageStatus string `json:"storageStatus"`
			UInstanceId   string `json:"uInstanceId"`
		} `json:"storageSet"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("[ERROR] code=%v msg=%s", jsonresp.Code, jsonresp.Message)
	}
	// storage no longer exists
	if len(jsonresp.StorageSet) != 1 {
		d.SetId("")
		return nil
	}
	attached := jsonresp.StorageSet[0].Attached
	status := jsonresp.StorageSet[0].StorageStatus
	uInstanceId := jsonresp.StorageSet[0].UInstanceId
	// storage has been detached and is not attaching
	if attached != 1 && status != "attaching" {
		log.Printf("[DEBUG] storage=%v status=%v", storageId, status)
		d.SetId("")
	}
	// storage has been attached but belongs to another instance
	if attached == 1 && uInstanceId != instanceId {
		log.Printf("[DEBUG] storage=%s has been attached to instance=%s", storageId, uInstanceId)
		d.SetId("")
	}
	return nil
}

// attachment id is a format like: att::ins-m5vh60me::disk-ojhtwo3k
func parseAttachmentId(attachId string) (iid, sid string, err error) {
	ids := strings.Split(attachId, "::")
	if len(ids) != 3 {
		return "", "", fmt.Errorf("Invalid attachment ID: %v", attachId)
	}
	return ids[1], ids[2], nil
}

func resourceTencentCloudCbsStorageAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	instanceId, storageId, err := parseAttachmentId(d.Id())
	if err != nil {
		return err
	}

	client := m.(*TencentCloudClient).commonConn

	params := map[string]string{
		"Action":       "DetachCbsStorages",
		"storageIds.0": storageId,
	}

	response, err := client.SendRequest("cbs", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		Code    int
		Message string
		Detail  interface{}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf("[ERROR] code=%v msg=%s", jsonresp.Code, jsonresp.Message)
	}
	asyncTasks := jsonresp.Detail.(map[string]interface{})
	for _, detail := range asyncTasks {
		taskDetail := detail.(map[string]interface{})
		taskCode := int(taskDetail["code"].(float64))
		if taskCode != 0 {
			return fmt.Errorf("[ERROR] code=%v msg=%s", taskCode, taskDetail["msg"].(string))
		}
	}

	if _, err := waitInstanceReachTargetStatus(client, []string{instanceId}, "PENDING"); err != nil {
		return err
	}
	if _, err := waitInstanceReachOneOfTargetStatusList(client, []string{instanceId}, []string{"RUNNING", "STOPPED"}); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
