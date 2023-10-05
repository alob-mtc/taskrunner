# taskrunner

This divides a large dataset into 'n' parallel execution units and processes the data in chunks of roughly equal size to optimize CPU usage while minimizing contention.

```go
_, err := taskrunner.Run(
  ctx,
  5, // worker count
  false, // order output
  tasks, // list to be processed
  func(ctx context.Context, chunkData []T, returnValue []T) error {}, // worker
)
```


TODO:
- [] Error propagation
- [] Return value ordering
