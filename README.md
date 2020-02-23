# Simple performance tracer

Initialization:
```go
jsoncfg := nptrace.NewJsonEncoderConfig(time.StampMicro, func(d time.Duration) []byte {
    return []byte(strconv.FormatInt(d.Nanoseconds(), 10))
})
cfg := nptrace.NewJsonEncoder(jsoncfg)
tracer := nptrace.NewTracer(cfg, traceWriter)
```

Middleware:
```go
func traceMiddleware(npTrace *nptrace.NPTrace) func(next http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
      ctx := r.Context()
      tracer := npTrace.New(ctx.Value("requestId").(string), strings.Trim(r.URL.Path, "/"))
      defer npTrace.Close(tracer)

      ctx = context.WithValue(ctx, infrastructure.Tracer, tracer)
      next.ServeHTTP(w, r.WithContext(ctx))
    }
    return http.HandlerFunc(fn)
  }
}
```

Usage:
```go
func Foo(ctx context.Context){
    tr := ctx.Value("Tracer").(*nptrace.Task).Start(name, args)
    defer ctx.Value("Tracer").(*nptrace.Task).Stop(tr)
    /*
        any work
    */
}
```

Output:
```json
{
  "id": "ArV23oYGbv-000001",
  "time": "Feb 23 00:14:42.496831",
  "trace": {
    "name": "api/user/login",
    "duration": 2869182,
    "args": [],
    "traces": [
      {
        "name": "loginUser",
        "duration": 2468308,
        "args": [],
        "traces": [
          {
            "name": "FindByEmail",
            "duration": 1876567,
            "args": [
              "SELECT * FROM users where email=$1"
            ],
            "traces": []
          }
        ]
      }
    ]
  }
}
```