# Litestream Demo

## 1. Start Services

```bash
docker compose up --build -d
```

## 2. Navigate to MinIO Console

```bash
open http://localhost:9001
```

### Sign In

See: `/env/minio.env`

### Create Litestream Bucket

Bucket name: `litestream`

## 3. Query and Create Items

Postman Collection: `postman_collection.json`

## 4. Verify Database Restore

### Delete the SQLite Database

Comment out the SQLite database `COPY` instruction in `/api/Dockerfile`

```dockerfile
# SQLite database:
# COPY items.db /var/api/data/
```

### Build and Restart the API Service

```bash
docker compose up --build -d
```

### Query Items

Postman Collection: `postman_collection.json`

The `Get Items` query should return the same items that were in the database at the time it was deleted.
