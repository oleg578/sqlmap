# sql2json

## Description

The main function of this package convert result of sql request into JSON structure.

This allows for flexible use of database queries without creating new structures for each query every time.

When queries with a varying number of fields are executed, or when the fields are undefined,

we can receive a JSON response that is built based on the metadata of the query and response.


For example, when we call query like "SELECT `Column1`, `Column2` from `table_name`" we know

structure and can use it as

```golang
...

type SomeStruct struct {
	Column1: string
	Column2: int64
}
...
```

But, when we call data like "SELECT * from `table_name`", we can't create the structure in advance.

So, using RowsToJson function we can encode response into JSON.

### tests result for 10_000_000 rows

#### standard (structure way)
Elapsed Time = 3193<br>
Alloc = 2600 MiB        TotalAlloc = 5387 MiB   HeapAlloc = 2600 MiB    StackInuse = 544 Kb     Sys = 3356033 Kb        NumGC = 17
#### mapping
Elapsed Time = 31793<br>
Alloc = 8269 MiB        TotalAlloc = 24469 MiB  HeapAlloc = 8269 MiB    StackInuse = 608 Kb     Sys = 9293062 Kb        NumGC = 25
##### reflect
Elapsed time: 29106<br>
Alloc = 2871 MiB        TotalAlloc = 21040 MiB  HeapAlloc = 2871 MiB    StackInuse = 544 Kb     Sys = 4467661 Kb        NumGC = 32