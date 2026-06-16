# 测试环境变量配置

测试配置现在支持通过环境变量进行配置，如果环境变量不存在则使用默认值。

## 环境变量列表

### AWS S3 配置
- `AWS_CALLER_SECRET_ID` - AWS 访问密钥 ID
- `AWS_CALLER_SECRET_KEY` - AWS 秘密访问密钥
- `AWS_ROLE_ARN` - AWS 角色 ARN
- `AWS_REGION` - AWS 区域
- `AWS_ENDPOINT` - AWS 端点
- `AWS_BUCKET` - AWS S3 存储桶名称

### COS 配置
- `COS_CALLER_SECRET_ID` - COS 访问密钥 ID
- `COS_CALLER_SECRET_KEY` - COS 秘密访问密钥
- `COS_ROLE_ARN` - COS 角色 ARN
- `COS_REGION` - COS 区域
- `COS_ENDPOINT` - COS 端点
- `COS_BUCKET` - COS 存储桶名称
- `COS_UIN` - COS 用户 UIN

## 使用方法

### 1. 设置环境变量

```bash
# 设置 AWS 配置
export AWS_CALLER_SECRET_ID="your-aws-access-key"
export AWS_CALLER_SECRET_KEY="your-aws-secret-key"
export AWS_ROLE_ARN="arn:aws:iam::123456789012:role/your-role"
export AWS_REGION="ap-northeast-1"
export AWS_ENDPOINT="s3.amazonaws.com"
export AWS_BUCKET="your-bucket-name"

# 设置 COS 配置
export COS_CALLER_SECRET_ID="your-cos-access-key"
export COS_CALLER_SECRET_KEY="your-cos-secret-key"
export COS_ROLE_ARN="qcs::cam::uin/200043701471:roleName/test-role"
export COS_REGION="ap-tokyo"
export COS_ENDPOINT="cos.ap-tokyo.myqcloud.com"
export COS_BUCKET="your-cos-bucket"
export COS_UIN="200043701471"
```

### 2. 运行测试

```bash
# 使用环境变量运行测试
go test ./test/... -v

# 或者使用默认配置运行测试
go test ./test/... -v
```

### 3. 在 VS Code 中调试

在 VS Code 的调试配置中，你可以通过 `env` 字段设置环境变量：

```json
{
    "name": "Test All",
    "type": "go",
    "request": "launch",
    "mode": "test",
    "program": "${workspaceFolder}/common/test",
    "args": ["-test.v"],
    "env": {
        "AWS_CALLER_SECRET_ID": "your-aws-access-key",
        "AWS_CALLER_SECRET_KEY": "your-aws-secret-key",
        "COS_CALLER_SECRET_ID": "your-cos-access-key",
        "COS_CALLER_SECRET_KEY": "your-cos-secret-key"
    }
}
```

## 默认值

如果环境变量未设置，将使用以下默认值：

- AWS 配置使用测试用的示例值
- COS 配置使用测试用的示例值

## 安全注意事项

1. 不要将真实的凭证提交到版本控制系统
2. 使用 `.env` 文件时，确保将其添加到 `.gitignore`
3. 在生产环境中使用环境变量或安全的配置管理系统
4. 定期轮换访问密钥

## 配置验证

测试会自动验证配置的有效性，包括：
- 必需字段是否存在
- 格式是否正确
- 长度是否符合要求
