## 1. Implementation
- [x] 1.1 Create GitHub workflow file `.github/workflows/chinese-char-check.yml`
- [x] 1.2 Configure workflow to trigger on PR events (opened, synchronize, edited)
- [x] 1.3 Implement script to detect Chinese characters in changed files
- [x] 1.4 Limit validation scope to `/tencentcloud/services/` and `/website/` directories
- [x] 1.5 Configure workflow to check only PR diff, not entire repository
- [x] 1.6 Set up proper error reporting when Chinese characters are found
- [x] 1.7 Test workflow with sample PR containing Chinese characters
- [x] 1.8 Test workflow with sample PR containing only English content
- [x] 1.9 Document workflow behavior and troubleshooting in PR template or contributing guide