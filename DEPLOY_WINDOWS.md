# WorldCup Mate — Windows 手动部署指南

## 部署目录

将前端构建产物复制到部署目录：

```
C:\app\worldcup-mate\frontend\     ← 复制 frontend\dist\ 下的全部文件
C:\app\worldcup-mate\backend\      ← 复制 backend\server.exe + backend\.env.production + backend\uploads\
```

---

## 步骤一：MySQL 准备

```powershell
# 打开 MySQL 命令行 (管理员)
"C:\Program Files\MySQL\MySQL Server 8.0\bin\mysql.exe" -u root -p

# 在 MySQL 中执行：
CREATE DATABASE IF NOT EXISTS worldcup_mate CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

## 步骤二：Redis 准备

确保 Redis 服务已安装并启动：

```powershell
# 检查 Redis 是否运行（默认 6379）
redis-cli ping
# 应返回 PONG
```

如果没有 Redis，可以从 https://redis.io/download/ 下载 Windows 版。

## 步骤三：部署后端

```powershell
# 1. 复制文件
mkdir C:\app\worldcup-mate\backend -Force
mkdir C:\app\worldcup-mate\backend\uploads\avatars -Force
copy backend\server.exe C:\app\worldcup-mate\backend\
copy backend\.env.production C:\app\worldcup-mate\backend\.env
copy -Recurse backend\uploads\* C:\app\worldcup-mate\backend\uploads\

# 2. 启动后端（新窗口）
cd C:\app\worldcup-mate\backend
.\server.exe
# 看到 "Server starting on :8080" 表示启动成功
```

## 步骤四：部署前端

```powershell
# 1. 复制前端构建产物
mkdir C:\app\worldcup-mate\frontend -Force
copy -Recurse frontend\dist\* C:\app\worldcup-mate\frontend\
```

## 步骤五：配置 Nginx

```powershell
# 1. 把生成的 nginx 配置复制到 nginx 目录
copy frontend\worldcup-mate-nginx.conf C:\app\Workspace\nginx-1.28.3\nginx-1.28.3\conf\worldcup-mate.conf
```

然后编辑 nginx 主配置文件 `C:\app\Workspace\nginx-1.28.3\nginx-1.28.3\conf\nginx.conf`，在 `http { }` 块中的最后（在最后一个 `}` 之前）添加一行：

```nginx
include worldcup-mate.conf;
```

## 步骤六：启动所有服务

**启动顺序很重要**：MySQL → Redis → 后端 → Nginx

### 1. 启动 MySQL（如果还没启动）
```powershell
net start MySQL80
# 或者
sc start MySQL80
```

### 2. 启动 Redis
```powershell
redis-server
# 或
net start Redis
```

### 3. 启动后端
```powershell
cd C:\app\worldcup-mate\backend
.\server.exe
```

### 4. 启动 Nginx
```powershell
cd C:\app\Workspace\nginx-1.28.3\nginx-1.28.3
start nginx.exe

# 检查 nginx 是否运行
tasklist /fi "imagename eq nginx.exe"
```

## 步骤七：验证部署

打开浏览器访问：

| 页面 | 地址 |
|------|------|
| 前端首页 | `http://localhost/` |
| API 测试 | `http://localhost/api/matches` |
| 球队列表 | `http://localhost/api/teams` |

## 常见问题

### nginx 报端口被占用
```powershell
netstat -ano | findstr :80
taskkill /PID <PID> /F
```

### 后端启动报 MySQL 连接失败
检查 MySQL 服务是否运行，以及 `.env` 中的 `MYSQL_DSN` 用户名密码是否正确。

### 后端启动报 CORS 错误
修改 `.env` 中 `CORS_ALLOWED_ORIGINS=http://localhost,http://你的实际域名`。

### nginx 配置修改后重载
```powershell
cd C:\app\Workspace\nginx-1.28.3\nginx-1.28.3
nginx -s reload
```

### 停止 nginx
```powershell
nginx -s stop
```
