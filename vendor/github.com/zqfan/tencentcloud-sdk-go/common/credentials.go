package common

type Credential interface {
	GetCredentialParams() map[string]string
	GetSecretId() string
	GetSecretKey() string
}

type BasicCredential struct {
	secretId  string
	secretKey string
}

func NewBasicCredential(secretId, secretKey string) *BasicCredential {
	return &BasicCredential{
		secretId:  secretId,
		secretKey: secretKey,
	}
}

func (c *BasicCredential) GetSecretId() string {
	return c.secretId
}

func (c *BasicCredential) GetSecretKey() string {
	return c.secretKey
}

func (c *BasicCredential) GetCredentialParams() map[string]string {
	return map[string]string{
		"SecretId": c.secretId,
	}
}

type TokenCredential struct {
	secretId  string
	secretKey string
	token     string
}

func NewTokenCredential(secretId, secretKey, token string) *TokenCredential {
	return &TokenCredential{
		secretId:  secretId,
		secretKey: secretKey,
		token:     token,
	}
}

func (c *TokenCredential) GetSecretId() string {
	return c.secretId
}

func (c *TokenCredential) GetSecretKey() string {
	return c.secretKey
}

func (c *TokenCredential) GetCredentialParams() map[string]string {
	return map[string]string{
		"SecretId": c.secretId,
		"Token":    c.token,
	}
}
