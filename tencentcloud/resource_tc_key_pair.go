package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKeyPairCreate,
		Read:   resourceTencentCloudKeyPairRead,
		Delete: resourceTencentCloudKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateKeyPairName,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return strings.TrimSpace(v.(string))
					default:
						return ""
					}
				},
			},
		},
	}
}

func resourceTencentCloudKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Version":   "2017-03-12",
		"ProjectId": "0", // TODO only support default projectId in v1.0
	}

	keyName := d.Get("key_name").(string)

	if publicKey, ok := d.GetOk("public_key"); ok {
		params["Action"] = "ImportKeyPair"
		params["KeyName"] = keyName
		params["PublicKey"] = publicKey.(string)

		response, err := client.SendRequest("cvm", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Response struct {
				Error struct {
					Code    string `json:"Code"`
					Message string `json:"Message"`
				}
				KeyId     string `json:"KeyId"`
				RequestId string
			}
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return err
		}
		if jsonresp.Response.Error.Code != "" {
			return fmt.Errorf(
				"tencentcloud_instance got error, code:%v, message:%v",
				jsonresp.Response.Error.Code,
				jsonresp.Response.Error.Message,
			)
		}
		id := jsonresp.Response.KeyId
		d.SetId(id)

	} else {
		params["Action"] = "CreateKeyPair"
		params["KeyName"] = keyName
		response, err := client.SendRequest("cvm", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Response struct {
				Error struct {
					Code    string `json:"Code"`
					Message string `json:"Message"`
				}
				KeyPair struct {
					KeyId      string `json:"KeyId"`
					KeyName    string `json:"KeyName"`
					ProjectId  int    `json:"ProjectId"`
					PublicKey  string `json:"PublicKey"`
					PrivateKey string `json:"PrivateKey"`
				} `json:"KeyPair"`
				RequestId string
			}
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return err
		}
		if jsonresp.Response.Error.Code != "" {
			return fmt.Errorf(
				"tencentcloud_instance got error, code:%v, message:%v",
				jsonresp.Response.Error.Code,
				jsonresp.Response.Error.Message,
			)
		}

		id := jsonresp.Response.KeyPair.KeyId
		d.SetId(id)
	}

	return resourceTencentCloudKeyPairRead(d, meta)
}

func resourceTencentCloudKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	client := meta.(*TencentCloudClient).commonConn
	keyName, _, err := findKeyPairById(client, id)
	if err != nil {
		if err == errKeyPairNotFound {
			d.SetId("")
			return nil
		}
		return err
	}

	if err := d.Set("key_name", keyName); err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("[DEBUG] tencentcloud_key_pair - deleting key pair:% v", id)
	params := map[string]string{
		"Version":  "2017-03-12",
		"Action":   "DeleteKeyPairs",
		"KeyIds.0": id,
	}
	client := meta.(*TencentCloudClient).commonConn
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, bindedInstanceIds, err := findKeyPairById(client, id)
		if err != nil && err != errKeyPairNotFound {
			return resource.NonRetryableError(err)
		}
		if len(bindedInstanceIds) > 0 {
			var stillUnbinedInstanceIds []string
			for _, insId := range bindedInstanceIds {
				if err := unbindKeyPiar(client, insId, id); err != nil {
					stillUnbinedInstanceIds = append(stillUnbinedInstanceIds, insId)
				}
			}
			if len(stillUnbinedInstanceIds) > 0 {
				s := fmt.Sprintf(
					"key pair: %v need to unbind with instanceIds: %v",
					id,
					stillUnbinedInstanceIds,
				)
				log.Printf("[DEBUG] %v", s)
				err = fmt.Errorf(s)
				return resource.RetryableError(err)
			}
		}

		response, err := client.SendRequest("cvm", params)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		var jsonresp struct {
			Response struct {
				Error struct {
					Code    string `json:"Code"`
					Message string `json:"Message"`
				}
				RequestId string
			}
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if retryable(jsonresp.Response.Error.Code, jsonresp.Response.Error.Message) {
			return resource.RetryableError(fmt.Errorf(jsonresp.Response.Error.Message))
		}
		if jsonresp.Response.Error.Code != "" {
			err = fmt.Errorf(
				"tencentcloud_key_pair got error, code:%v, message:%v",
				jsonresp.Response.Error.Code,
				jsonresp.Response.Error.Message,
			)
			return resource.NonRetryableError(err)
		}
		return nil
	})
}
