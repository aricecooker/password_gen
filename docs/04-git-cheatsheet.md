# Git 常用命令速查

> 配套 README 使用的速查手册

## 1. 配置

### 1.1 首次使用：设置用户信息

```bash
git config --global user.name "你的名字"
git config --global user.email "你的邮箱@example.com"
```

### 1.2 查看配置

```bash
git config --list
```

## 2. 仓库操作

### 2.1 初始化本地仓库

```bash
cd D:/password_gen
git init
```

### 2.2 关联远程仓库

```bash
git remote add origin https://github.com/用户名/仓库名.git
```

### 2.3 查看远程仓库

```bash
git remote -v
```

## 3. 提交与推送

### 3.1 提交三步

```bash
git add .                  # 把所有改动加入暂存区
git commit -m "提交说明"   # 提交到本地仓库
git push                   # 推送到远程
```

### 3.2 首次推送

```bash
git push -u origin main    # -u 记住上游分支，以后直接 git push
```

### 3.3 拉取远程更新

```bash
git pull
```

## 4. 分支

### 4.1 查看分支

```bash
git branch                 # 本地分支
git branch -a              # 包含远程
```

### 4.2 创建并切换

```bash
git checkout -b 新分支名
```

### 4.3 切换分支

```bash
git checkout 分支名
```

### 4.4 重命名当前分支

```bash
git branch -M 新名字       # -M 强制改名
```

## 5. 撤销与回退

### 5.1 撤销工作区修改

```bash
git checkout -- 文件名
```

### 5.2 取消暂存（已 add 但未 commit）

```bash
git reset HEAD 文件名
```

### 5.3 回退到上一次提交（危险，慎用）

```bash
git reset --hard HEAD^
```

## 6. 查看状态与历史

### 6.1 查看当前状态

```bash
git status
```

### 6.2 查看提交历史

```bash
git log                     # 详细
git log --oneline           # 简洁
```

### 6.3 查看某次提交的改动

```bash
git show 提交哈希
```

### 6.4 查看文件差异

```bash
git diff                    # 工作区 vs 暂存区
git diff --staged           # 暂存区 vs 上次提交
```

## 7. GitHub Token 认证

GitHub 已**不再支持密码推送**。需要用 Personal Access Token (PAT) 或 SSH。

### 7.1 生成 Token

1. 打开 https://github.com/settings/tokens
2. 点 "Generate new token" → 选 "Generate new token (classic)"
3. 勾选 `repo` 权限
4. 生成后**立即复制保存**（页面关闭后无法再看到）

### 7.2 方式 A：Windows 凭据管理器（推荐）

```bash
git config --global credential.helper manager
git push
```

第一次 push 时会弹窗让你输入：
- Username：GitHub 用户名
- Password：**粘贴 token**（不是真正的密码）

凭据会被 Windows 记住，下次不用再输。

### 7.3 方式 B：把 token 嵌入 URL（不推荐）

```bash
git remote set-url origin https://用户名:token@github.com/用户名/仓库名.git
```

⚠️ token 会保存在 `.git/config` 里，慎用。

### 7.4 方式 C：SSH

```bash
# 1. 生成 key
ssh-keygen -t ed25519 -C "你的邮箱@example.com"

# 2. 复制公钥内容
cat ~/.ssh/id_ed25519.pub
```

把公钥内容粘贴到 https://github.com/settings/keys。

```bash
# 3. 改用 SSH 协议
git remote set-url origin git@github.com:用户名/仓库名.git

# 4. 测试
ssh -T git@github.com

# 5. 推送
git push
```

## 8. .gitignore 模板（Go 项目）

```gitignore
# Binaries
*.exe
*.dll
*.so
*.dylib

# Go
*.test
*.out
coverage.txt
!go.sum

# IDE
.idea/
.vscode/

# OS
.DS_Store
Thumbs.db

# Build output
dist/
```

## 9. 紧急情况

### 9.1 撤销所有未提交改动

```bash
git checkout .
git clean -fd        # 删除未跟踪文件
```

### 9.2 找回丢失的提交

```bash
git reflog
```

`reflog` 记录所有 HEAD 移动历史，可找回"丢失"的提交哈希。

### 9.3 撤销已推送的 commit

```bash
git revert 提交哈希
git push
```

`revert` 会生成一个新的"反向"提交，**不会改写历史**。安全。

## 10. 提交说明规范

格式：`<类型>: <简短说明>`

| 类型 | 用途 |
|------|------|
| feat | 新功能 |
| fix | 修复 bug |
| docs | 文档变更 |
| style | 格式（不影响代码运行）|
| refactor | 重构 |
| test | 增加测试 |
| chore | 构建过程或辅助工具变动 |

示例：
```bash
git commit -m "feat: 添加批量生成密码"
git commit -m "fix: 修复长度边界判断"
git commit -m "docs: 更新 README 安装说明"
```
