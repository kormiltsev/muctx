# MUCTX

Context with TRY implementation.

### Example

### 1.
```
muc := muctx.New()
muc.Lock()
defer muc.Unlock()
muc.Lock() // Deadlock
```
### 2.
```
muc := muctx.New()
muc.Lock()

muc.LockTry() // returns FALSE in 200 ms
```

### 3.
```
muc := muctx.New()
muc.Lock()

ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second))
defer cancel()

muc.LockTryCtx(ctx) // returns FALSE after ctx.Done()
```