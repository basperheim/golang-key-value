# Golang key-value store Redis alternative

Example use:

```bash
go run main.go
```

## Requests

### Set new key-value pair

```json
POST /set
{
  "key": "hello",
  "value": {
    "hello": "world"
  }
}
```

200 Resp:

```
ok
```

### Get key-value pair

```
GET /get?key=hello
```

200 Resp:

```json
{
  "createdAt": "2024-04-04T21:02:15.145488+08:00",
  "key": "hello",
  "updatedAt": "2024-04-04T21:02:15.145488+08:00",
  "value": {
    "hello": "world"
  }
}
```

404 Resp:

```json
Key not found
```

### Delete key-value pair

```
DELETE /delete?key=hello
```

200 Resp:

```json
{
  "key": "hello"
}
```

404 Resp:

```json
Key not found
```

## NodeJS Axios example

```js
const axios = require("axios");

const data = {
  key: "testJSON",
  value: { hello: "world" }, // JSON value
};

axios
  .post("http://localhost:8080/set", data)
  .then((response) => {
    console.log(response.data); // Log the response from the server

    axios
      .get(`http://localhost:8080/get?key=${data.key}`)
      .then((response) => {
        console.log(response.data); // Log the response from the server

        axios
          .delete(`http://localhost:8080/delete?key=${data.key}`)
          .then((response) => {
            console.log(response.data); // Log the response from the server
          })
          .catch((error) => {
            console.error("Error:", error);
          });
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  })
  .catch((error) => {
    console.error("Error:", error.code);
    console.error(error.message);
  });
```

## Persistent storage possibilities

Here's a summary of various ways to implement persistent storage in Go, along with minimal code examples for each approach:

File-based Storage:

Store data as files on disk.
Use encoding packages like encoding/json or encoding/gob for serialization.

### Minimal Example:

```go
// Write data to file
ioutil.WriteFile("data.json", jsonData, 0644)

// Read data from file
ioutil.ReadFile("data.json")
```

### Database Systems

Use database systems for structured data storage and retrieval.
Options include SQL databases (e.g., PostgreSQL, MySQL) or NoSQL databases (e.g., MongoDB, Redis).

#### Minimal Example (using SQLite):

```go
// Open SQLite database
db, err := sql.Open("sqlite3", "data.db")

// Execute SQL query
db.Exec("CREATE TABLE IF NOT EXISTS data (key TEXT PRIMARY KEY, value TEXT)")
```

#### Key-Value Stores:

Use key-value stores for simple data storage and retrieval.
Options include embedded databases (e.g., BoltDB) or standalone key-value store servers (e.g., Redis).

#### Minimal Example (using BoltDB):

```go
Copy code
// Open BoltDB database
db, err := bolt.Open("data.db", 0644, nil)
// Write data to database
db.Update(func(tx \*bolt.Tx) error {
bucket, err := tx.CreateBucketIfNotExists([]byte("data"))
// Store key-value pair in bucket
return bucket.Put([]byte("key"), []byte("value"))
})
```

### ORM Libraries:

Use Object-Relational Mapping (ORM) libraries to map Go structs to database tables.

ORM libraries provide high-level abstractions for database interactions.

Options include GORM, XORM, and SQLBoiler.

#### Minimal Example (using GORM):

```go
// Define model struct
type Data struct {
  Key string
  Value string
}
// Auto-migrate the schema
db.AutoMigrate(&Data{})
// Create new record
db.Create(&Data{Key: "key", Value: "value"})
```

### Cloud Storage Services:

Utilize cloud storage services like Amazon S3 or Google Cloud Storage for scalable and durable storage.

Interact with these services using Go SDKs or APIs.

#### Minimal Example (using AWS S3):

```go
// Upload data to S3 bucket
\_, err := svc.PutObject(&s3.PutObjectInput{
  Body: strings.NewReader("data"),
  Bucket: aws.String("myBucket"),
  Key: aws.String("myKey"),
})
```

Each of these approaches has its own advantages and trade-offs in terms of simplicity, performance, scalability, and maintenance. Choose the one that best fits your requirements and constraints.
