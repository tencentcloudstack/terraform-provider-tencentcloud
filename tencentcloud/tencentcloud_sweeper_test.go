package tencentcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedClientForRegion(region string) (interface{}, error) {
	var secretId, secretKey string
	if secretId = os.Getenv("TENCENTCLOUD_SECRET_ID"); secretId == "" {
		return nil, fmt.Errorf("TENCENTCLOUD_SECRET_ID can not be empty")
	}
	if secretKey = os.Getenv("TENCENTCLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("TENCENTCLOUD_SECRET_KEY can not be empty")
	}

	config := Config{
		SecretId:  secretId,
		SecretKey: secretKey,
		Region:    region,
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}
	return client, nil
}
