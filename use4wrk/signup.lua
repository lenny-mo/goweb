-- post.lua
wrk.method = "POST"
wrk.body   = '{"username":"foo","password":"bar", "email":"boo@gmail.com", "gender":1}'
wrk.headers["Content-Type"] = "application/json"

