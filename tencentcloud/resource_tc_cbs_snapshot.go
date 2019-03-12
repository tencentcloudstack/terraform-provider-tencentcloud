package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/client"
)

const (
	ecSnapshotNotExistError  = 16003
	ecSnapshotStatusError    = 16004
	ecSnapshotLifeStateError = 16033
)

var (
	errSnapshotNotFound = errors.New("snapshot not found")
)

type snapshotInfo struct {
	DiskType       string
	Percent        int
	StorageId      string
	StorageSize    int
	SnapshotName   string
	SnapshotStatus string
}

func resourceTencentCloudCbsSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsSnapshotCreate,
		Read:   resourceTencentCloudCbsSnapshotRead,
		Update: resourceTencentCloudCbsSnapshotUpdate,
		Delete: resourceTencentCloudCbsSnapshotDelete,

		Schema: map[string]*schema.Schema{
			"snapshot_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
			},
			"storage_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snapshot_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"pecent": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func modifySnapshot(snapshotId string, snapshotName string, client *client.Client) error {
	params := map[string]string{
		"Action":       "ModifySnapshot",
		"snapshotId":   snapshotId,
		"snapshotName": snapshotName,
	}

	response, err := client.SendRequest("snapshot", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		Code     int
		Message  string
		CodeDesc string
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf(
			"ModifySnapshot error, code:%v, message: %v, codeDesc: %v.",
			jsonresp.Code,
			jsonresp.Message,
			jsonresp.CodeDesc,
		)
	}

	log.Printf("[DEBUG] ModifySnapshot, new snapshotName: %#v.", snapshotName)
	return nil
}

func deleteSnapshot(snapshotId string, client *client.Client) *resource.RetryError {
	params := map[string]string{
		"Action":        "DeleteSnapshot",
		"snapshotIds.0": snapshotId,
	}

	response, err := client.SendRequest("snapshot", params)
	if err != nil {
		return resource.NonRetryableError(err)
	}

	var jsonresp struct {
		Code     int
		Message  string
		CodeDesc string
		Detail   struct {
			Result map[string]struct {
				Msg  string
				Code int
			}
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return resource.NonRetryableError(err)
	}
	if jsonresp.Code != 0 {
		return resource.NonRetryableError(fmt.Errorf(
			"DeleteSnapshot error, code:%v, message: %v, codeDesc: %v.",
			jsonresp.Code,
			jsonresp.Message,
			jsonresp.CodeDesc,
		))
	}
	code := jsonresp.Detail.Result[snapshotId].Code
	msg := jsonresp.Detail.Result[snapshotId].Msg

	if code == ecSnapshotNotExistError || code == 0 {
		return nil
	} else if code == ecSnapshotStatusError || code == ecSnapshotLifeStateError {
		return resource.RetryableError(fmt.Errorf("snapshot status error, please retry later"))
	} else {
		return resource.NonRetryableError(fmt.Errorf("DeleteSnapshot failed, inner code:%v, message: %v", code, msg))
	}
}

func waitingSnapshotReady(snapshotId string, client *client.Client) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		snapshotInfo, _, err := describeSnapshot(snapshotId, client)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if snapshotInfo.SnapshotStatus == "creating" {
			return resource.RetryableError(fmt.Errorf("waiting snapshot ready"))
		}
		return nil
	})
}

func describeSnapshot(snapshotId string, client *client.Client) (*snapshotInfo, bool, error) {
	var jsonresp struct {
		Code        int
		Message     string
		CodeDesc    string
		SnapshotSet []snapshotInfo
	}
	params := map[string]string{
		"Action":        "DescribeSnapshots",
		"snapshotIds.0": snapshotId,
	}
	response, err := client.SendRequest("snapshot", params)
	canRetryError := false
	if err != nil {
		return nil, canRetryError, err
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return nil, canRetryError, err
	}
	if jsonresp.Code != 0 {
		return nil, canRetryError, fmt.Errorf(
			"DescribeSnapshots error, code:%v, message: %v, codeDesc: %v.",
			jsonresp.Code,
			jsonresp.Message,
			jsonresp.CodeDesc,
		)
	}

	if len(jsonresp.SnapshotSet) == 0 {
		canRetryError = true
		return nil, canRetryError, errSnapshotNotFound

	}

	snapshot := jsonresp.SnapshotSet[0]
	return &snapshot, canRetryError, nil
}

func resourceTencentCloudCbsSnapshotCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	snapshotName := d.Get("snapshot_name").(string)
	params := map[string]string{
		"Action":       "CreateSnapshot",
		"storageId":    d.Get("storage_id").(string),
		"snapshotName": snapshotName,
	}

	response, err := client.SendRequest("snapshot", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		Code       int
		Message    string
		CodeDesc   string
		SnapshotId string
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Code != 0 {
		return fmt.Errorf(
			"CreateSnapshot error, code:%v, message:%v, codeDesc:%v.",
			jsonresp.Code,
			jsonresp.Message,
			jsonresp.CodeDesc,
		)
	}
	d.SetId(jsonresp.SnapshotId)
	return resourceTencentCloudCbsSnapshotRead(d, m)
}

func resourceTencentCloudCbsSnapshotRead(d *schema.ResourceData, m interface{}) error {
	snapshot, _, err := describeSnapshot(d.Id(), m.(*TencentCloudClient).commonConn)
	if err != nil {
		if err == errSnapshotNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("disk_type", snapshot.DiskType)
	d.Set("percent", snapshot.Percent)
	d.Set("storage_size", snapshot.StorageSize)
	d.Set("storage_id", snapshot.StorageId)
	d.Set("snapshot_name", snapshot.SnapshotName)
	d.Set("snapshot_status", snapshot.SnapshotStatus)
	return nil
}

func resourceTencentCloudCbsSnapshotUpdate(d *schema.ResourceData, m interface{}) error {
	if !d.HasChange("snapshot_name") {
		return nil
	}

	_, n := d.GetChange("snapshot_name")
	snapshotName := n.(string)

	err := modifySnapshot(d.Id(), snapshotName, m.(*TencentCloudClient).commonConn)
	if err != nil {
		return err
	}

	return resourceTencentCloudCbsSnapshotRead(d, m)
}

func resourceTencentCloudCbsSnapshotDelete(d *schema.ResourceData, m interface{}) error {
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		return deleteSnapshot(d.Id(), m.(*TencentCloudClient).commonConn)
	})
	d.SetId("")
	return err
}
