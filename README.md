# feishu2hexo

<div align="center">

![Golang](https://img.shields.io/github/go-mod/go-version/Mars160/feishu2hexo?color=00ADD8&logo=go)
![Release](https://img.shields.io/github/v/release/Mars160/feishu2hexo?color=orange&logo=github)
![License](https://img.shields.io/github/license/Mars160/feishu2hexo?color=blue)
![Tests](https://github.com/Mars160/feishu2hexo/actions/workflows/unittest.yaml/badge.svg)
![Docker](https://img.shields.io/badge/Docker-feishu2hexo-2496ed?logo=docker&logoColor=white)
![Render](https://img.shields.io/badge/Render-feishu2hexo-4cfac9?logo=render&logoColor=white)

**A powerful tool to download Feishu/LarkSuite documents as Markdown files**

[快速开始](#快速开始) • [使用文档](#使用文档) • [API 参考](#api-参考) • [贡献指南](#贡献指南)

</div>

## 📋 目录

- [功能特性](#功能特性)
- [快速开始](#快速开始)
- [安装方式](#安装方式)
- [配置说明](#配置说明)
- [使用文档](#使用文档)
- [进阶用法](#进阶用法)
- [架构设计](#架构设计)
- [常见问题](#常见问题)
- [贡献指南](#贡献指南)
- [致谢](#致谢)

## ✨ 功能特性

- 🚀 **多格式支持** - 支持标准 Markdown、Hexo、Hugo 格式输出
- 📁 **批量下载** - 支持整个文件夹或知识库的批量下载
- 🖼️ **图片处理** - 自动下载并本地化文档中的所有图片
- 🔄 **格式转换** - 智能转换 Feishu 文档块为标准 Markdown 语法
- 🌐 **跨平台** - 支持 Windows、macOS、Linux
- 🐳 **容器化** - 提供 Docker 镜像，开箱即用
- 🌍 **在线服务** - 提供在线版本，无需安装

## 🚀 快速开始

### 1. 获取 API 凭证

前往 [飞书开放平台](https://open.feishu.cn/app) 创建应用并获取：

- **App ID** - 应用标识
- **App Secret** - 应用密钥

### 2. 安装工具

#### 方式一：下载预编译二进制文件（推荐）

从 [Releases 页面](https://github.com/Mars160/feishu2hexo/releases) 下载对应平台的可执行文件。

#### 方式二：使用 Go 安装

```bash
go install github.com/Mars160/feishu2hexo/cmd/feishu2hexo@latest
```

#### 方式三：使用 Docker

```bash
docker run -it --rm -p 8080:8080 \
  -e FEISHU_APP_ID=<your_app_id> \
  -e FEISHU_APP_SECRET=<your_app_secret> \
  wwMars160/feishu2hexo
```

### 3. 配置凭证

```bash
feishu2hexo config --appId <your_app_id> --appSecret <your_app_secret>
```

### 4. 下载文档

```bash
# 下载单个文档
feishu2hexo dl "https://your-domain.feishu.cn/docx/xxxxx"

# 下载为 Hexo 博客格式
feishu2hexo hexo -o posts/ -t "技术,教程" "https://your-domain.feishu.cn/docx/xxxxx"

# 下载为 Hugo 博客格式
feishu2hexo hugo -o content/posts/ "https://your-domain.feishu.cn/docx/xxxxx"
```

## 📦 安装方式

### 预编译二进制文件

| 平台    | 架构  | 下载链接                                                                                |
| ------- | ----- | --------------------------------------------------------------------------------------- |
| Windows | x64   | [feishu2hexo-windows-amd64.exe](https://github.com/Mars160/feishu2hexo/releases/latest) |
| macOS   | x64   | [feishu2hexo-darwin-amd64](https://github.com/Mars160/feishu2hexo/releases/latest)      |
| macOS   | ARM64 | [feishu2hexo-darwin-arm64](https://github.com/Mars160/feishu2hexo/releases/latest)      |
| Linux   | x64   | [feishu2hexo-linux-amd64](https://github.com/Mars160/feishu2hexo/releases/latest)       |

### 从源码编译

```bash
git clone https://github.com/Mars160/feishu2hexo.git
cd feishu2hexo
make build
```

## ⚙️ 配置说明

### 必需权限

在飞书开放平台中，需要开通以下权限：

| 权限名称                               | 权限代码                       | 说明                   |
| -------------------------------------- | ------------------------------ | ---------------------- |
| 查看新版文档                           | `docx:document:readonly`       | 获取文档基本信息和内容 |
| 下载云文档中的图片和附件               | `docs:document.media:download` | 下载文档中的图片       |
| 查看、评论、编辑和管理云空间中所有文件 | `drive:file:readonly`          | 访问文件夹             |
| 查看知识库                             | `wiki:wiki:readonly`           | 访问知识库             |

### 配置文件

配置文件位置：`~/.config/feishu2hexo/config.json`

```json
{
  "feishu": {
    "app_id": "your_app_id",
    "app_secret": "your_app_secret"
  },
  "output": {
    "image_dir": "static",
    "title_as_filename": false,
    "use_html_tags": false,
    "skip_img_download": false
  }
}
```

## 📖 使用文档

### 命令行界面

#### 全局选项

```bash
feishu2hexo [global options] command [command options] [arguments...]

GLOBAL OPTIONS:
  --help, -h     显示帮助信息
  --version, -v  显示版本信息
```

#### config 命令 - 配置管理

```bash
feishu2hexo config [options...]

OPTIONS:
  --appId value      设置 App ID
  --appSecret value  设置 App Secret
```

示例：

```bash
# 设置配置
feishu2hexo config --appId cli_xxx --appSecret xxx

# 查看当前配置
feishu2hexo config
```

#### download 命令 - 下载文档

```bash
feishu2hexo download [options...] <url>

OPTIONS:
  --output value, -o value  指定输出目录 (默认: "./")
  --dump                    导出 API 响应的 JSON 数据
  --batch                   批量下载文件夹内所有文档
  --wiki                    下载知识库内所有文档
```

示例：

```bash
# 下载单个文档
feishu2hexo dl "https://domain.feishu.cn/docx/xxxxx"

# 批量下载文件夹
feishu2hexo dl --batch -o output/ "https://domain.feishu.cn/drive/folder/xxxxx"

# 下载整个知识库
feishu2hexo dl --wiki -o output/ "https://domain.feishu.cn/wiki/settings/xxxxx"
```

#### hexo 命令 - 转换为 Hexo 格式

```bash
feishu2hexo hexo [options...] <url>

OPTIONS:
  --output value, -o value  指定输出目录 (默认: "./")
  --tags value, -t value    设置文章标签，用逗号分隔 (默认: "论文,算法")
  --dump                    导出 API 响应的 JSON 数据
  --batch                   批量下载文件夹内所有文档
  --wiki                    下载知识库内所有文档
```

生成的 Hexo 文件格式：

```markdown
---
title: 文档标题
date: 2024-01-01 12:00:00
updated: 2024-01-01 12:00:00
tags: [标签1, 标签2]
categories: []
---

文档内容...
```

#### hugo 命令 - 转换为 Hugo 格式

```bash
feishu2hexo hugo [options...] <url>

OPTIONS:
  --output value, -o value  指定输出目录 (默认: "./")
  --tags value, -t value    设置文章标签，用逗号分隔 (默认: "论文,算法")
  --dump                    导出 API 响应的 JSON 数据
  --batch                   批量下载文件夹内所有文档
  --wiki                    下载知识库内所有文档
```

Hugo 模式会自动检测 Hugo 项目根目录，图片默认保存到 `static/post_imgs/`。

### Web 界面

#### 使用 Docker

```yaml
# docker-compose.yml
version: "3"
services:
  feishu2hexo:
    image: wwMars160/feishu2hexo
    environment:
      FEISHU_APP_ID: <your_app_id>
      FEISHU_APP_SECRET: <your_app_secret>
      GIN_MODE: release
    ports:
      - "8080:8080"
```

启动服务：

```bash
docker-compose up -d
```

访问 http://localhost:8080 使用 Web 界面。

#### 在线版本

访问 https://feishu2hexo.onrender.com/ 使用在线版本（仅供测试使用）。

## 🎯 进阶用法

### 批量处理脚本

```bash
#!/bin/bash
# 批量下载多个文档

urls=(
  "https://domain.feishu.cn/docx/doc1"
  "https://domain.feishu.cn/docx/doc2"
  "https://domain.feishu.cn/docx/doc3"
)

for url in "${urls[@]}"; do
  echo "Processing: $url"
  feishu2hexo hexo -o posts/ "$url"
done
```

### 集成到 CI/CD

```yaml
# GitHub Actions 示例
name: Convert Feishu Docs
on:
  push:
    paths:
      - "docs-list.txt"

jobs:
  convert:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: "1.21"

      - name: Install feishu2hexo
        run: go install github.com/Mars160/feishu2hexo/cmd/feishu2hexo@latest

      - name: Convert documents
        env:
          FEISHU_APP_ID: ${{ secrets.FEISHU_APP_ID }}
          FEISHU_APP_SECRET: ${{ secrets.FEISHU_APP_SECRET }}
        run: |
          feishu2hexo config --appId $FEISHU_APP_ID --appSecret $FEISHU_APP_SECRET
          while read -r url; do
            feishu2hexo hexo -o content/posts/ "$url"
          done < docs-list.txt
```

## 🏗️ 架构设计

```
feishu2hexo/
├── cmd/              # CLI 命令实现
│   ├── main.go       # 主入口和命令路由
│   ├── config.go     # 配置命令
│   ├── download.go   # 下载命令
│   ├── hexo.go       # Hexo 转换命令
│   └── hugo.go       # Hugo 转换命令
├── core/             # 核心业务逻辑
│   ├── client.go     # Feishu API 客户端
│   ├── config.go     # 配置管理
│   └── parser.go     # 文档解析器
├── utils/            # 工具函数
│   ├── common.go     # 通用工具
│   └── url.go        # URL 解析
└── web/              # Web 服务
    ├── main.go       # Web 服务器
    └── download.go   # Web 下载接口
```

### 核心流程

1. **URL 解析** - 提取文档类型和 Token
2. **API 调用** - 通过 Feishu SDK 获取文档内容
3. **内容解析** - 将文档块转换为 Markdown
4. **资源处理** - 下载并本地化图片等资源
5. **格式化输出** - 根据指定格式生成最终文件

## ❓ 常见问题

### Q: 为什么出现权限错误？

A: 请确保已在飞书开放平台开通了所有必需权限，并且应用已发布。

### Q: 批量下载速度很慢？

A: 工具已内置速率限制（4 QPS），这是 Feishu API 的限制。如需更快速度，请申请提高 API 限额。

### Q: 图片下载失败？

A: 检查网络连接和防火墙设置，确保能够访问 Feishu 的 CDN。

### Q: Docker 版本无法使用批量下载？

A: Docker 版本目前仅支持单个文档下载，批量功能请使用 CLI 版本。

### Q: 如何迁移旧版文档？

A: 旧版 Feishu 文档已不再支持，请使用 [v1_support](https://github.com/Mars160/feishu2hexo/tree/v1_support) 分支。

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 开发环境搭建

1. Fork 本仓库
2. 克隆到本地：

   ```bash
   git clone https://github.com/your-username/feishu2hexo.git
   cd feishu2hexo
   ```

3. 安装依赖：

   ```bash
   go mod download
   ```

4. 运行测试：
   ```bash
   make test
   ```

### 提交代码

1. 创建功能分支：

   ```bash
   git checkout -b feature/your-feature
   ```

2. 提交更改：

   ```bash
   git commit -m "feat: add new feature"
   ```

3. 推送分支：

   ```bash
   git push origin feature/your-feature
   ```

4. 创建 Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 运行 `make format` 格式化代码
- 添加适当的测试用例
- 更新相关文档

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [chyroc/lark](https://github.com/chyroc/lark) - Feishu Go SDK
- [88250/lute](https://github.com/88250/lute) - Markdown 引擎
- 所有贡献者和用户的支持

## 📞 联系方式

- 提交 Issue：[GitHub Issues](https://github.com/Mars160/feishu2hexo/issues)
- 功能建议：[GitHub Discussions](https://github.com/Mars160/feishu2hexo/discussions)
- 邮箱：your-email@example.com

---

<div align="center">

**如果这个项目对你有帮助，请给它一个 ⭐️**

Made with ❤️ by [Wsine](https://github.com/Wsine)

</div>
