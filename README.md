
---

## Step 1.) Run the server

```bash
go run main.go
```

You should see:

```
⇨ http server started on [::]:8080
```

---

## Step 2.) Test with curl

```bash
# 1) Health
curl -s localhost:8080/healthz
# -> {"status":"ok"}

# 2) Hello (default)
curl -s localhost:8080/hello
# -> {"hello":"world"}

# 3) Hello with name
curl -s "localhost:8080/hello?name=Abhi"
# -> {"hello":"Abhi"}

# 4) Create user
curl -s -X POST localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ava","age":21}'
# -> {"id":1,"name":"Ava","age":21}

# 5) List users
curl -s localhost:8080/users
# -> [{"id":1,"name":"Ava","age":21}]

# 6) Filter users by min_age
curl -s "localhost:8080/users?min_age=22"
# -> []
```

---

