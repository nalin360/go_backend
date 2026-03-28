# Code Review — WealthArena Go Backend

**Date:** 28 March 2026  
**Reviewer:** AI Assistant  
**Scope:** Full codebase review (`cmd/`, `internal/`)

---

## Summary

| Area | Rating | Notes |
|------|--------|-------|
| Project Structure | ⭐⭐⭐⭐ | Clean layered architecture (handler → service → store) |
| Security | ⭐⭐ | Critical issues in auth & password handling |
| Naming & Conventions | ⭐⭐ | Several typos and inconsistent naming patterns |
| Error Handling | ⭐⭐⭐ | Good `CustomError` system, but misused in places |
| Code Quality | ⭐⭐⭐ | Functional but needs cleanup |

---

## 🔴 Critical Issues

### 1. Password hash exposed in API response

**File:** `internal/store/models.go:16`  
```go
PasswordHash string `json:"password_hash"`
```

The `Customer` struct serialises `password_hash` into every JSON response. When you call `GET /coustomer` or `POST /auth/register`, the hashed password is returned to the client.

**Fix:** Add `json:"-"` to exclude it, or create a separate response DTO:
```go
PasswordHash string `json:"-"` // never expose in JSON
```

> ⚠️ Since this file is **sqlc-generated**, you should NOT edit it directly. Instead, create a response struct that omits the field:
> ```go
> type CustomerResponse struct {
>     CustID      int64  `json:"cust_id"`
>     CustName    string `json:"cust_name"`
>     CustEmail   string `json:"cust_email"`
>     CustAddress string `json:"cust_address"`
>     IsAdmin     bool   `json:"is_admin"`
> }
> ```

---

### 2. JWT secret read at request-time with no validation

**File:** `internal/handlears/auth_handler.go:67`
```go
secretKey := os.Getenv("JWT_SECRET")
tokenString, err := token.SignedString([]byte(secretKey))
```

**Problems:**
- If `JWT_SECRET` is not set, `secretKey` is `""` → tokens are signed with an **empty key**. Anyone can forge tokens.
- Reading from `os.Getenv` on every request is wasteful; it should be loaded once at startup.

**Fix:** Load and validate at startup:
```go
// In main.go or config
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
    slog.Error("JWT_SECRET is not set")
    os.Exit(1)
}
```

Then pass it into the `AuthHandler` struct.

---

### 3. Request field named `password_hash` accepts plaintext

**File:** `internal/handlears/coustomer_handler.go:25`
```go
type CreateCustomerRequest struct {
    ...
    PasswordHash string `json:"password_hash"`
}
```

The client sends a **plaintext password**, but the JSON field is called `password_hash`. This is misleading — the client is NOT sending a hash. 

**Fix:**
```go
Password string `json:"password"`
```

---

## 🟡 Bugs & Logic Issues

### 4. `godotenv.Load()` called AFTER `env.GetString()`

**File:** `cmd/main.go:17-27`
```go
cfg := config{
    addr: ":8080",
    db: dbConfig{
        dsn: env.GetString("DATABASE_URL", ""),  // ← reads env HERE
    },
}

godotenv.Load()                      // ← .env loaded AFTER
dbURL := os.Getenv("DATABASE_URL")   // ← reads again
```

`env.GetString("DATABASE_URL", "")` runs **before** `godotenv.Load()`, so `cfg.db.dsn` will always be empty (unless the env var is set at the OS level). The `dbURL` variable gets the correct value but `cfg.db.dsn` doesn't. `cfg.db.dsn` is never used after this, making it dead code.

**Fix:** Move `godotenv.Load()` to the very first line of `main()`.

---

### 5. Error created before the error occurs

**File:** `internal/handlears/auth_handler.go:85,110`
```go
custErr := httputil.NewBadRequest(nil, "invalid request body")  // line 85
// ...
errs := httputil.NewBadRequest(err, "invalid request body")     // line 110
```

- **Line 85:** Creates an error object with `nil` before checking if decode failed. This caused the panic you hit earlier. Even with the nil-guard fix, the pattern is wrong — create the error only when you need it.
- **Line 110:** Creates a `BadRequest` error but then sends it as `StatusInternalServerError`. A DB failure is not a bad request.

**Fix:**
```go
// Don't pre-create errors. Create them inline:
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    jsons.Write(w, http.StatusBadRequest, httputil.NewBadRequest(err, "invalid request body"))
    return
}
```

---

### 6. `CreateCustomer` returns `200 OK` instead of `201 Created`

**File:** `internal/handlears/auth_handler.go:117`
```go
jsons.Write(w, http.StatusOK, coustomer)
```

A resource creation should return `201 Created`, as you correctly do in `CreateProduct`.

---

### 7. No input validation

**Files:** `auth_handler.go`, `handler.go`

There is **zero validation** on any input:
- Email format not checked
- Password length/strength not enforced
- Product name can be empty
- Price can be negative
- `// validation` comment exists but no code follows it

**Suggestion:** Add a validation helper or use a library like `go-playground/validator`.

---

## 🟡 Naming & Consistency Issues

### 8. Typo: `handlears` → `handlers`

The package is named `handlears` (misspelled). This appears everywhere:
- Package name: `package handlears`
- Struct: `type handlears struct`
- Directory: `internal/handlears/`
- Imports: `wealtharena.in/api/internal/handlears`

> ⚠️ Renaming a package is a breaking change — do it now while the project is small. It'll be painful later.

### 9. Typo: `coustomer` → `customer`

Appears in:
- `CoustomerHandlears` struct
- `CoustomerHandler()` function
- Route path: `/coustomer`
- Variable names: `coustomer`, `coustomers`
- Comments: `// -- coustomer --`

### 10. Inconsistent constructor naming

| Handler | Constructor | Pattern |
|---------|------------|---------|
| Product | `ProductHandler()` | No `New` prefix |
| Customer | `CoustomerHandler()` | No `New` prefix |
| Auth | `NewAuthHandler()` | Has `New` prefix ✅ |

Pick one style. `New` prefix is idiomatic Go — use it everywhere.

### 11. Inconsistent struct visibility

| Struct | Exported? |
|--------|-----------|
| `handlears` (product) | ❌ unexported |
| `CoustomerHandlears` | ✅ exported |
| `AuthHandler` | ✅ exported |

Product handler struct is unexported while the others are exported.

### 12. Inconsistent error response format

**Product handlers** return:
```go
map[string]string{"error": err.Error()}
```

**Auth/Customer handlers** return:
```go
httputil.NewBadRequest(err, "message")  // CustomError struct
```

Your API clients will receive different error shapes depending on which endpoint they hit. Pick one format.

---

## 🟡 Architecture Issues

### 13. Monolithic `Service` interface

**File:** `internal/services/service.go:10-24`

The `Service` interface has **13 methods** spanning products, customers, search, and updates. As you add orders, auth, inventory etc., this will grow to 50+ methods.

**Suggestion:** Split by domain:
```go
type ProductService interface {
    ListProduct(ctx context.Context) ([]store.Product, error)
    CreateProduct(ctx context.Context, req store.CreateProductParams) (store.Product, error)
    GetProduct(ctx context.Context, id int64) (store.Product, error)
}

type CustomerService interface {
    CreateCustomer(ctx context.Context, req store.CreateCustomerParams) (store.Customer, error)
    GetCustomer(ctx context.Context, id int64) (store.Customer, error)
    // ...
}
```

### 14. Hardcoded pagination

**File:** `internal/services/service.go:53-56`
```go
return s.db.ListCustomers(ctx, store.ListCustomersParams{
    Limit:  10,
    Offset: 0,
})
```

Limit and offset are hardcoded to `10` and `0`. The handler never passes pagination params → users can only ever see the first 10 records.

### 15. `store.New()` dereference

**File:** `cmd/api.go:51,54`
```go
queries := store.New(app.db)          // returns *Queries
productsService := services.NewService(*queries)  // dereferences to value
```

`NewService` accepts `store.Queries` (value type), but `store.New()` returns `*store.Queries`. You dereference it with `*queries`, meaning the service works with a **copy**. This works for now because `Queries` just holds a `DBTX` interface, but it's fragile. Consider accepting `*store.Queries` in `NewService`.

### 16. Dependency instantiation inside route handler

**File:** `cmd/api.go:50-67`

Store, service, and handler creation all happen inside `mount()`. This makes testing impossible — you can't inject mocks. Move dependency creation to `main()` and pass them into `application`.

---

## 🟢 What's Good

- ✅ Clean layered architecture: **handler → service → store (sqlc)**
- ✅ Using `chi` router with sensible middleware stack
- ✅ Structured logging with `slog`
- ✅ `CustomError` type with proper `error` interface + `Unwrap()` support
- ✅ Password hashing with bcrypt in `CreateCustomer`
- ✅ JWT implementation with meaningful claims
- ✅ Proper use of `context.Context` throughout
- ✅ Server timeouts configured (read/write/idle)
- ✅ Using `pgxpool` for connection pooling

---

## 📋 Action Items (Priority Order)

| # | Priority | Task |
|---|----------|------|
| 1 | 🔴 Critical | Stop exposing `password_hash` in API responses |
| 2 | 🔴 Critical | Validate `JWT_SECRET` at startup, fail if empty |
| 3 | 🔴 Critical | Move `godotenv.Load()` before `env.GetString()` |
| 4 | 🟡 High | Rename `password_hash` → `password` in request structs |
| 5 | 🟡 High | Fix all typos (`handlears`, `coustomer`) while project is small |
| 6 | 🟡 High | Standardize error response format across all handlers |
| 7 | 🟡 Medium | Add input validation (email, password strength, required fields) |
| 8 | 🟡 Medium | Make constructor naming consistent (`New` prefix everywhere) |
| 9 | 🟢 Low | Split `Service` interface by domain |
| 10 | 🟢 Low | Add pagination params to list endpoints |
| 11 | 🟢 Low | Move dependency creation out of `mount()` |
| 12 | 🟢 Low | Add JWT auth middleware for protected routes |

---

*Review covers files as of 28 March 2026, 20:20 IST.*
