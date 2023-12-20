package cos

import (
	"fmt"

	"github.com/beevik/etree"
)

func validateACLBody(v interface{}, k string) (ws []string, errors []error) {
	missingUrl := "https://cloud.tencent.com/document/product/436/7737"
	aclUrl := "https://cloud.tencent.com/document/product/436/30752#.E6.93.8D.E4.BD.9C-permission"

	value := v.(string)
	xmlDoc := etree.NewDocument()
	if err := xmlDoc.ReadFromString(value); err != nil {
		errors = append(errors, fmt.Errorf("[CRITAL]read acl_body xml from string error: %v\n", err))
		return
	}

	rawRoot := xmlDoc.SelectElement("AccessControlPolicy")
	if rawRoot == nil {
		errors = append(errors, fmt.Errorf("[CRITAL]missing the AccessControlPolicy element, please refer to the official document: %v\n", missingUrl))
		return
	}

	if len(rawRoot.FindElements("//Owner")) == 0 {
		errors = append(errors, fmt.Errorf("[CRITAL]missing the Owner element, please refer to the official document: %v\n", missingUrl))
	} else {
		if len(rawRoot.FindElements("//Owner/ID")) == 0 {
			errors = append(errors, fmt.Errorf("[CRITAL]missing the Owner/ID element, please refer to the official document: %v\n", missingUrl))
		}
	}

	if len(rawRoot.FindElements("//Grant")) == 0 {
		errors = append(errors, fmt.Errorf("[CRITAL]missing the Grant element, please refer to the official document: %v\n", missingUrl))
	}

	foundT := false
	grantees := rawRoot.FindElements("//Grantee")
	if len(grantees) == 0 {
		errors = append(errors, fmt.Errorf("[CRITAL]missing the Grant/Grantee element, please refer to the official document: %v\n", missingUrl))
	}

	for _, grantee := range grantees {
		var aType string
		foundT = false
		for _, validType := range COSACLGranteeTypeSeq {
			aType = grantee.SelectAttrValue("type", "unknown")
			if aType == validType {
				foundT = true
				break
			}
		}
		if !foundT {
			errors = append(errors, fmt.Errorf("[CRITAL]the Grantee type[%s] is not a valid type, please refer to the official document: %v\n", aType, aclUrl))
		}
	}

	foundP := false
	permissions := rawRoot.FindElements("//Permission")
	if len(permissions) == 0 {
		errors = append(errors, fmt.Errorf("[CRITAL]missing the Grant/Permission element, please refer to the official document: %v\n", missingUrl))
	}

	for _, permission := range permissions {
		var aPermisson string
		foundP = false
		for _, validPermission := range COSACLPermissionSeq {
			aPermisson = permission.Text()
			if aPermisson == validPermission {
				foundP = true
				break
			}
		}
		if !foundP {
			errors = append(errors, fmt.Errorf("[CRITAL]the Grant Permission[%s] is not a valid type, please refer to the official document: %v\n", aPermisson, aclUrl))
		}
	}

	return
}
