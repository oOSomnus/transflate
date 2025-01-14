
## Redis 数据结构

### 键名格式

- 每个用户的任务集合以 `tasks:<username>` 的格式存储。
- 单个任务以 `task:<username>:<taskId>` 的格式存储。

### 哈希字段

#### 字段说明

- **`status`**: 任务状态，存储为整型字符串。
- **`filename`**: 文件名，表示与任务关联的文件。
- **`link`**: 下载链接，可根据需求更新。

#### Redis 示例数据

```txt
Key:    task:john:task123
Fields:
    status: "1"
    filename: "example.txt"
    link: "http://example.com/task123"
```

---

## 接口方法

### 1. `SetTaskState`

#### 功能

设置任务状态。如果任务不存在则创建；如果存在，则更新状态并保持原文件名。

#### 方法签名

```go
SetTaskState(ctx context.Context, username, taskId string, status int, filename string, ttl time.Duration) error
```

#### 参数

- **`username`**: 用户名，用于生成任务的 Redis 键。
- **`taskId`**: 任务的唯一标识。
- **`status`**: 任务状态。
- **`filename`**: 任务的文件名。
- **`ttl`**: 键的过期时间。

#### 示例

```go
err := repository.SetTaskState(ctx, "john", "task123", 1, "example.txt", 24*time.Hour)
```

#### Redis 操作

- 如果任务不存在：
  - 使用 `HSET` 创建任务数据。
  - 使用 `SADD` 将任务 ID 添加到用户的任务集合。
- 如果任务存在：
  - 使用 `HSET` 更新 `status`。
- 设置过期时间：`EXPIRE`。

---

### 2. `GetTaskState`

#### 功能

获取任务的状态和文件名。

#### 方法签名

```go
GetTaskState(ctx context.Context, username, taskId string) (int, string, error)
```

#### 参数

- **`username`**: 用户名。
- **`taskId`**: 任务的唯一标识。

#### 示例

```go
status, filename, err := repository.GetTaskState(ctx, "john", "task123")
```

#### Redis 操作

- 使用 `HMGET` 获取 `status` 和 `filename` 字段。

---

### 3. `FetchAllTask`

#### 功能

获取用户的所有任务，包括状态、文件名和链接。

#### 方法签名

```go
FetchAllTask(ctx context.Context, username string) (map[string]map[string]interface{}, error)
```

#### 参数

- **`username`**: 用户名。

#### 示例

```go
tasks, err := repository.FetchAllTask(ctx, "john")
```

#### 返回

- 返回包含任务 ID 及其状态、文件名、链接的嵌套字典。

#### Redis 操作

- 使用 `SMEMBERS` 获取用户任务集合。
- 使用 `HGETALL` 获取每个任务的详细数据。

---

### 4. `UpdateTaskLink`

#### 功能

更新任务的下载链接。

#### 方法签名

```go
UpdateTaskLink(ctx context.Context, username, taskId, link string) error
```

#### 参数

- **`username`**: 用户名。
- **`taskId`**: 任务的唯一标识。
- **`link`**: 新的下载链接。

#### 示例

```go
err := repository.UpdateTaskLink(ctx, "john", "task123", "http://example.com/new_link")
```

#### Redis 操作

- 使用 `HSET` 更新 `link` 字段。

---

## Redis 数据操作对照表

| 方法               | Redis 操作               | 描述               |
|------------------|------------------------|------------------|
| `SetTaskState`   | `HSET` + `EXPIRE`      | 创建/更新任务状态并设置过期时间 |
| `GetTaskState`   | `HMGET`                | 获取任务状态和文件名       |
| `FetchAllTask`   | `SMEMBERS` + `HGETALL` | 获取所有任务详细数据       |
| `UpdateTaskLink` | `HSET`                 | 更新任务的下载链接        |

---