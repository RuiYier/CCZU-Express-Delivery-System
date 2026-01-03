# PackChann - 校园快递收发系统后端文档

PackChann 是一个基于 Go (Gin) + PostgreSQL 开发的校园快递管理系统后端。支持用户注册登录、包裹入库出库、寄件服务以及状态追踪等功能。

## 🛠 技术栈

- **语言**: Go (Golang)
- **Web 框架**: Gin
- **数据库**: PostgreSQL
- **ORM**: GORM
- **认证**: JWT (JSON Web Token) + 数据库双重验证
- **配置管理**: Viper
- **ID 生成**: Snowflake (雪花算法)

## 🚀 快速开始

### 1. 环境要求

- Go 1.20+
- PostgreSQL 12+

### 2. 配置文件

在 `config/config.toml` 中配置数据库和 JWT 信息：

```toml
[database]
host = "127.0.0.1"
port = 5432
user = "postgres"
password = "your_password"
dbname = "packchann"

[jwt]
secret = "your_secret_key"
expiration_hours = 24
```

### 3. 运行

```bash
go mod tidy
go run .
```

服务默认运行在 `:8088` 端口。

---

## 📡 API 接口文档

### 通用说明

- **Content-Type**: `application/json`
- **认证方式**: 受保护接口需要在 Header 中携带 Token。
  - `Authorization: Bearer <your_access_token>`

### 1. 用户认证 (Unprotected)

#### 1.1 用户注册

- **URL**: `/register`
- **Method**: `POST`
- **描述**: 注册新用户，成功后自动返回 Token。

**请求参数**:

```json
{
  "user_name": "张三",
  "password": "password123",
  "student_id": "20210001",
  "phone": "13800000001",
  "address": "南区宿舍1号楼",
  "role": "user" // 可选，默认为 "user"，管理员注册需填写 "admin"
}
```

**响应**:

```json
{
    "user": { ... },
    "access_token": "eyJ...",
    "refresh_token": "eyJ..."
}
```

#### 1.2 用户登录

- **URL**: `/login`
- **Method**: `POST`
- **描述**: 使用学号和密码登录。

**请求参数**:

```json
{
  "student_id": "20210001",
  "password": "password123"
}
```

#### 1.3 健康检查

- **URL**: `/ping`
- **Method**: `GET`
- **响应**: `{"message": "pong"}`

---

### 2. 业务接口 (Protected)

**注意**: 以下所有接口均需在 Header 中携带 `Authorization`。

#### 2.1 包裹入库 (Check In)

- **URL**: `/packCheckIn`
- **Method**: `POST`
- **描述**: 快递员或管理员将包裹录入系统。

**请求参数**:

```json
{
  "pack_id": 10001, // 快递单号
  "user_id": 1, // 收件人用户ID
  "shelf_code": 101 // 货架号
}
```

**响应**:

```json
{
    "message": "Pack checked in successfully",
    "pack": {
        "pack_id": 10001,
        "pack_status": "pending",
        "pickup_code": "101-1767", // 生成的取件码
        ...
    }
}
```

#### 2.2 包裹出库 (Check Out)

- **URL**: `/packCheckout`
- **Method**: `POST`
- **描述**: 用户取件出库。

**请求参数**:

```json
{
  "pack_id": 10001,
  "user_id": 1
}
```

#### 2.3 寄件 (Mail Pack)

- **URL**: `/mailPack`
- **Method**: `POST`
- **描述**: 用户发起寄件请求。系统会自动生成唯一的 Pack ID (Snowflake)。

**请求参数**:

```json
{
  "shipping_address": "北京海淀区...",
  "recipient": "李四",
  "reciving_address": "上海浦东新区...",
  "shipper_phone": "13800000001",
  "recipient_phone": "13900000002" // 必须是系统内存在的用户手机号
}
```

**响应**:

```json
{
    "message": "Mail pack created successfully",
    "pack": {
        "pack_id": 265495629717835776,
        "pack_status": "in_transit",
        ...
    }
}
```

#### 2.4 取消寄件

- **URL**: `/cancelMail`
- **Method**: `POST`
- **描述**: 取消尚未发货的寄件请求。

**请求参数**:

```json
{
  "pack_id": 265495629717835776,
  "user_id": 1
}
```

#### 2.5 更新包裹状态

- **URL**: `/updatePackStatus`
- **Method**: `POST`
- **描述**: 管理员更新包裹状态（如：arrived, shipped, cancelled）。

**请求参数**:

```json
{
  "pack_id": 10001,
  "pack_status": "arrived"
}
```

#### 2.6 获取包裹详情

- **URL**: `/getPackDetails/:pack_id`
- **Method**: `GET`
- **描述**: 根据包裹 ID 获取详细信息。

#### 2.7 获取用户所有包裹

- **URL**: `/allPacks/:user_id`
- **Method**: `GET`
- **描述**: 获取指定用户的所有包裹列表。

#### 2.8 更新用户信息

- **URL**: `/updateUserInfo`
- **Method**: `POST`
- **描述**: 更新用户的基本信息（不包含密码）。

**请求参数**:

```json
{
  "user_id": 1,
  "user_name": "张三丰",
  "address": "新的地址"
}
```

### 3. 管理员接口 (Admin Only)

#### 3.1 获取所有用户

- **URL**: `/admin/users`
- **Method**: `GET`
- **描述**: 获取系统中所有用户的列表。需要管理员权限。

#### 3.2 获取所有包裹

- **URL**: `/admin/packs`
- **Method**: `GET`
- **描述**: 获取系统中所有包裹的列表。需要管理员权限。
- **Query 参数**:
  - `status` (可选): 按包裹状态筛选 (e.g., `pending`, `arrived`, `shipped`)

#### 3.3 更新包裹信息

- **URL**: `/admin/pack`
- **Method**: `PUT`
- **描述**: 管理员更新包裹的任意可修改字段。

**请求参数**:

```json
{
  "pack_id": 10001,
  "pack_status": "arrived",
  "pickup_code": "101-1234",
  "user_id": 2
}
```

---

## 🗄 数据库设计

### Users 表

存储用户信息。

- `user_id` (PK): 用户唯一标识
- `student_id`: 学号 (Unique)
- `password_hash`: 加密后的密码
- `phone`: 手机号
- `address`: 地址

### Packs 表

存储包裹信息。

- `pack_id` (PK): 包裹 ID (入库时为单号，寄件时为 Snowflake ID)
- `user_id`: 关联用户
- `pack_status`: 状态 (pending, checked_out, in_transit, cancelled, arrived, shipped)
- `pickup_code`: 取件码 (货架号-时间戳后四位)
- `check_in_time`: 入库时间
- `check_out_time`: 出库时间

### UserTokens 表

存储 JWT Token，用于验证 Token 的有效性和实现登出/吊销功能。

- `access_token`
- `refresh_token`
- `expires_at`

---

## 🔒 安全特性

1.  **密码加密**: 使用 `bcrypt` 对用户密码进行哈希存储。
2.  **Token 验证**: 使用 JWT 进行身份认证，并结合数据库 `UserTokens` 表验证 Token 是否被吊销或过期。
3.  **中间件保护**: `AuthMiddleware` 拦截所有受保护路由，确保请求合法。
