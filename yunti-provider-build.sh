python yunti-provider-process.py
go mod vendor
yq eval 'keys' yunti-provider-code.yaml | awk '{print $2}' | grep -v 'go.mod' | xargs -I {} goimports -w {}
yq eval 'keys' yunti-provider-code.yaml | awk '{print $2}' | grep -v 'go.mod' | xargs -I {} go fmt {}
