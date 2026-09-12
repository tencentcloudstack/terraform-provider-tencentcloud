# Chinese Character Check Workflow

## Overview

The Chinese Character Check workflow automatically validates that Pull Requests do not introduce Chinese characters in public-facing code and documentation. This ensures the Terraform provider maintains international project standards.

## Scope

The workflow checks files in the following directories:
- `tencentcloud/services/` - Provider implementation code
- `website/` - Documentation and website content

Files outside these directories are not validated and may contain Chinese characters without affecting the CI check.

## How It Works

1. **Trigger**: The workflow runs automatically on PR events:
   - When a PR is opened
   - When a PR is synchronized (new commits pushed)
   - When a PR is edited

2. **File Detection**: Only files changed in the PR within the scope directories are checked

3. **Character Detection**: Uses Unicode ranges to detect:
   - Chinese characters (CJK Unified Ideographs): `\u4e00-\u9fff`
   - CJK Extension A: `\u3400-\u4dbf`
   - CJK Compatibility Ideographs: `\uf900-\ufaff`
   - Chinese punctuation and symbols: `\u3000-\u303f`
   - Full-width characters: `\uff00-\uffef`

4. **Reporting**: If Chinese characters are found:
   - The CI check fails
   - Detailed report shows file names and line numbers
   - Clear instructions provided for resolution

## Example Violations

The following would cause the check to fail:

```go
// ❌ Chinese comments
func CreateInstance() {
    // 创建实例
    description := "实例描述"
}
```

```markdown
<!-- ❌ Chinese documentation -->
# 资源配置

这是一个配置示例。
```

## Resolution

When the check fails:

1. **Locate the violations**: Check the CI output for specific files and line numbers
2. **Remove Chinese characters**: Replace with English equivalents
3. **Common replacements**:
   - Comments: Translate to English
   - Documentation: Rewrite in English
   - String literals: Use English descriptions
   - Variable names: Use English naming

## Manual Testing

You can manually check for Chinese characters using:

```bash
# Search for Chinese characters in scope directories
grep -rP '[\u4e00-\u9fff\u3400-\u4dbf\uf900-\ufaff\u3000-\u303f\uff00-\uffef]' tencentcloud/services/ website/

# Test the detection script
./scripts/test-chinese-check.sh
```

## Troubleshooting

### False Positives
If the check incorrectly flags content:
1. Verify the file is actually in scope (`tencentcloud/services/` or `website/`)
2. Check if the flagged content contains full-width characters or Chinese punctuation
3. Replace with standard ASCII equivalents

### Check Not Running
If the workflow doesn't trigger:
1. Ensure your PR modifies files in the scope directories
2. Check that the workflow file exists: `.github/workflows/chinese-char-check.yml`
3. Verify PR events are configured correctly

### Performance Issues
The workflow is designed to complete within 2 minutes for typical PRs. If it takes longer:
1. Check if the PR contains very large files
2. Verify the file filtering is working correctly
3. Consider splitting large changes into smaller PRs

## Workflow Configuration

The workflow is configured in `.github/workflows/chinese-char-check.yml` with:
- **Permissions**: Read-only access to PR content
- **Triggers**: PR opened, synchronized, edited
- **Path filters**: Only runs when scope directories are modified
- **Python dependencies**: Uses built-in Python 3 for character detection

## Contributing

When contributing to this project:
1. Write all code comments in English
2. Use English for documentation
3. Ensure variable names and function names use English
4. Test locally before submitting PR if unsure

The Chinese character check helps maintain code quality and ensures the project remains accessible to the international developer community.