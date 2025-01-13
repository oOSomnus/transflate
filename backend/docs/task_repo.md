# Task Repository - Redis 数据结构描述与方法用例

## 概述

该 `TaskRepository` 接口使用 **Redis** 的哈希结构 (`Hash`) 来管理与用户任务相关的数据。任务存储以用户为单位，每个用户对应一个
Redis 哈希键，哈希键内包含此用户的所有任务数据。任务字段以 JSON 格式序列化，存储在 Redis 哈希的值中。

---

## Redis 数据结构

### 键名格式

Redis 键名采用以下格式，每个用户的任务数据保存在以用户名开头的键中：

```txt
task:<username>
```

- `<username>` 是用户的名字，用于标识不同用户。

### 哈希键值对 (Field)

#### 键 (Field Key)

- 每个任务的 ID，标记为 `taskId`，唯一标识每个任务。

#### 值 (Field Value)

- 每个任务的详细信息存储为 JSON 格式，包含以下字段：
    - **`status`**: 表示任务的状态，存储为字符串（可以转换为整数）。
    - **`link`**: 任务的关联链接。

#### Redis 中的示例数据

```txt
Key:    task:john
Fields:
    task123: {"status": "1", "link": "http://example.com/task123"}
    task124: {"status": "2", "link": "http://example.com/task124"}
```

---

## 接口方法

### 1. `SetTaskState`

将任务的状态设置到 Redis，并以 JSON 格式保存任务数据。

#### 方法签名

```go
SetTaskState(ctx context.Context, username, taskId string, state int, ttl time.Duration) error
```

#### 参数说明

- **`username`**: 用户名，用于生成 Redis 中的哈希键名。
- **`taskId`**: 任务的唯一标识。
- **`state`**: 任务状态，整型（会转为字符串存储）。
- **`ttl`**: 哈希键的过期时间。

#### Redis 操作

- Redis 哈希 (`HSET`)：
    - 键名：`task:<username>`
    - 字段：`<taskId>`
    - 值：`{"status": "<state>", "link": ""}`
- 设置 TTL (`EXPIRE`)：
    - 键名：`task:<username>`

#### 使用示例

```go
err := repository.SetTaskState(ctx, "john", "task123", 1, 24*time.Hour)
```

#### Redis 数据结果

```txt
Key:    task:john
Field:  task123
Value:  {"status": "1", "link": ""}
TTL:    24小时
```

---

### 2. `GetTaskState`

从 Redis 中获取任务的状态和关联链接。

#### 方法签名

```go
GetTaskState(ctx context.Context, username, taskId string) (int, string, error)
```

#### 参数说明

- **`username`**: 用户名，用于获取哈希键名。
- **`taskId`**: 任务的唯一标识。

#### Redis 操作

- Redis 哈希 (`HGET`)：
    - 键名：`task:<username>`
    - 字段：`<taskId>`

#### 返回值

- **`int`**: 任务的状态。
- **`string`**: 任务的链接。
- **`error`**: 操作失败时返回错误。

#### 使用示例

```go
state, link, err := repository.GetTaskState(ctx, "john", "task123")
```

#### 返回结果

```plaintext
state: 1
link: "http://example.com/task123"
error: nil (无错误)
```

---

### 3. `FetchAllTask`

获取用户的所有任务及其对应的状态和链接。

#### 方法签名

```go
FetchAllTask(ctx context.Context, username string) (map[string]map[string]interface{}, error)
```

#### 参数说明

- **`username`**: 用户名，用于获取 Redis 中的哈希键名。

#### Redis 操作

- Redis 哈希 (`HGETALL`)：
    - 键名：`task:<username>`

#### 返回值

- **`map[string]map[string]interface{}`**:
    - 外层 `map` 的键是任务 ID。
    - 内层 `map` 包含任务的字段（`status` 和 `link`）。
- **`error`**: 操作失败时的错误。

#### 使用示例

```go
tasks, err := repository.FetchAllTask(ctx, "john")
```

#### 返回结果

```go
map[string]map[string]interface{}{
    "task123": {
        "status": "1",
        "link": "http://example.com/task123",
    },
    "task124": {
        "status": "2",
        "link": "http://example.com/task124",
    },
}, nil
```

---

### 4. `UpdateTaskLink`

更新指定任务的关联链接。

#### 方法签名

```go
UpdateTaskLink(ctx context.Context, username, taskId, link string) error
```

#### 参数说明

- **`username`**: 用户名，用于获取 Redis 中的哈希键名。
- **`taskId`**: 任务的唯一标识。
- **`link`**: 新的关联链接。

#### Redis 操作

- 读取并更新 Redis 哈希 (`HGET` + `HSET`)：
    - 键名：`task:<username>`
    - 字段：`<taskId>`

#### 使用示例

```go
err := repository.UpdateTaskLink(ctx, "john", "task123", "http://example.com/new_task_link")
```

#### Redis 数据结果

```txt
Key:    task:john
Field:  task123
Value:  {"status": "1", "link": "http://example.com/new_task_link"}
```

---

## 小结

### Redis 数据结构

- **键**:
  ```txt
  task:<username>
  ```
- **哈希字段**:
    - 任务 ID (`taskId`) 作为字段。
    - 字段值存储为 JSON，包含任务状态和链接。

### 方法调用对应的 Redis 操作

| 方法               | Redis 操作          | 描述            |
|------------------|-------------------|---------------|
| `SetTaskState`   | `HSET` + `EXPIRE` | 设置任务状态并设置过期时间 |
| `GetTaskState`   | `HGET`            | 获取任务状态和链接     |
| `FetchAllTask`   | `HGETALL`         | 获取所有任务及其状态和链接 |
| `UpdateTaskLink` | `HGET` + `HSET`   | 更新任务的链接       |

---
