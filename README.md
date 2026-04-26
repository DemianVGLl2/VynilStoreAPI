# VynilStoreAPI

## Project description

REST API built with Go and the Gin Framework that manages the inventory of a vinyl record store, inspired by platforms like Merchbar. It supports token-based authentication, allowing multiple users to be logged in simultaneously. All data is stored in memory.

---

## Installation

**Requirements:** Go 1.21 or higher

```bash
git clone https://github.com/DemianVGLl2/VynilStoreAPI.git
cd VynilStoreAPI
go mod tidy
```

---

## How to run

```bash
go run .
```

The server starts at `http://localhost:8080`.

---

## Authentication

All endpoints except `/login` require a valid Bearer token in the `Authorization` header.

To get a token, log in with one of the available users:

| Username | Password  |
|----------|-----------|
| admin    | admin123  |
| user1    | pass1     |

---

## Endpoints

### POST `/login`
Authenticates a user and returns a session token.

```bash
curl -u admin:admin123 -X POST http://localhost:8080/login
```

```json
{
  "message": "Hi admin, welcome to the Store System",
  "token": "OjIE89GzFw"
}
```

---

### POST `/logout`
Revokes the current session token.

```bash
curl -H "Authorization: Bearer <TOKEN>" -X POST http://localhost:8080/logout
```

```json
{
  "message": "Bye admin, your token has been revoked"
}
```

---

### GET `/albums`
Returns all albums in the store.

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/albums
```

```json
[
  { "id": "1", "title": "Blue Train", "artist": "John Coltrane", "price": 56.99 },
  { "id": "2", "title": "Time Out", "artist": "Dave Brubeck", "price": 37.99 },
  { "id": "3", "title": "Flying Beagle", "artist": "Himiko Kikuchi", "price": 69.99 }
]
```

---

### GET `/albums/:id`
Returns a single album by ID. Returns 404 if not found.

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/albums/2
```

```json
{ "id": "2", "title": "Time Out", "artist": "Dave Brubeck", "price": 37.99 }
```

---

### POST `/createAlbum`
Adds a new album to the store. Returns 409 if the ID already exists.

```bash
curl -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -X POST \
     -d '{"id":"4","title":"Benny Golson New York Scene","artist":"Benny Golson","price":49.99}' \
     http://localhost:8080/createAlbum
```

```json
{ "id": "4", "title": "Benny Golson New York Scene", "artist": "Benny Golson", "price": 49.99 }
```

---

### GET `/status`
Returns the system status and the currently logged-in user.

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/status
```

```json
{
  "message": "Hi admin, the DPIP System is Up and Running",
  "time": "2026-04-26 14:32:01"
}
```

---

## Future work

1. **Database persistence** — replace the in-memory store with a database like PostgreSQL or SQLite so data survives server restarts.
2. **Token expiration** — add expiration times to tokens and implement refresh token support for more secure session management.
3. **Full CRUD for albums** — add `PUT /albums/:id` and `DELETE /albums/:id` endpoints to allow editing and removing albums from the store.