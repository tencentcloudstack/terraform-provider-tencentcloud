package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"strings"

	"github.com/athom/goset"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/zqfan/tencentcloud-sdk-go/client"
)

var (
	errKeyPairNotFound = fmt.Errorf("tencentcloud_key_pair not found")
)

func bindKeyPiar(client *client.Client, instanceId string, keyId string) error {
	if err := operateKeyPiar(client, instanceId, keyId, "AssociateInstancesKeyPairs", waitForKeyPairBinded); err != nil {
		return err
	}
	return nil
}

func unbindKeyPiar(client *client.Client, instanceId string, keyId string) error {
	if err := operateKeyPiar(client, instanceId, keyId, "DisassociateInstancesKeyPairs", waitForKeyPairUnbinded); err != nil {
		return err
	}
	return nil
}

func operateKeyPiar(client *client.Client, instanceId string, keyId string, action string, wait func(client *client.Client, instanceId string, keyId string) error) error {
	params := map[string]string{
		"Version":       "2017-03-12",
		"Action":        action,
		"InstanceIds.0": instanceId,
		"KeyIds.0":      keyId,
	}
	return operateInstanceBetweenStopAndStart(client, params, instanceId, func() error {
		return wait(client, instanceId, keyId)
	})
}

func waitForKeyPairUnbinded(client *client.Client, instanceId string, keyId string) error {
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] waitForKeyPairUnbinded: keyId: %v, instanceId: %v", keyId, instanceId)
		_, associatedInstanceIds, err := findKeyPairById(client, keyId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(associatedInstanceIds) == 0 {
			return nil
		}
		if goset.IsIncluded(associatedInstanceIds, instanceId) {
			err := fmt.Errorf("key pair: %v still bind in instance: %v", keyId, instanceId)
			return resource.RetryableError(err)
		}
		return nil
	})
}

func waitForKeyPairBinded(client *client.Client, instanceId string, keyId string) error {
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] waitForKeyPairBinded: keyId: %v, instanceId: %v", keyId, instanceId)
		_, associatedInstanceIds, err := findKeyPairById(client, keyId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if goset.IsIncluded(associatedInstanceIds, instanceId) {
			return nil
		}
		err = fmt.Errorf("key pair: %v not bind in instance: %v yet, retry", keyId, instanceId)
		return resource.RetryableError(err)
	})
}

func findKeyPairById(client *client.Client, id string) (keyName string, associatedInstanceIds []string, err error) {
	params := map[string]string{
		"Version":  "2017-03-12",
		"Action":   "DescribeKeyPairs",
		"KeyIds.0": id,
	}
	var response string
	response, err = client.SendRequest("cvm", params)
	if err != nil {
		return
	}
	var jsonresp struct {
		Response struct {
			Error struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			}
			TotalCount int `json:"TotalCount"`
			KeyPairSet []struct {
				KeyId                 string    `json:"KeyId"`
				KeyName               string    `json:"KeyName"`
				Description           string    `json:"Description"`
				PublicKey             string    `json:"PublicKey"`
				AssociatedInstanceIds []string  `json:"AssociatedInstanceIds"`
				CreateTime            time.Time `json:"CreateTime"`
			} `json:"KeyPairSet"`
			RequestId string
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return
	}
	if jsonresp.Response.Error.Code != "" {
		err = fmt.Errorf(
			"tencentcloud_key_pair got error, code:%v, message:%v",
			jsonresp.Response.Error.Code,
			jsonresp.Response.Error.Message,
		)
		return
	}

	kpSet := jsonresp.Response.KeyPairSet
	if len(kpSet) == 0 {
		err = errKeyPairNotFound
		return
	}
	keyName = kpSet[0].KeyName
	associatedInstanceIds = kpSet[0].AssociatedInstanceIds
	return
}

func errAlreadyStopped(err error, instanceId string) bool {
	w := fmt.Sprintf("`%s` which is in the state of `STOPPED`", instanceId)
	return strings.Contains(err.Error(), w)
}
