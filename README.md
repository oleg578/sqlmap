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

#### reflect-chan

Elapsed time: 74_203 ms
TotalAlloc = 113_988 MiB HeapAlloc = 2_849 MiB    StackInuse = 544 Kb     Sys = 4_051_259 Kb        Frees = 2_192_703_342      NumGC = 355


#### reflect

Elapsed time: 72_446 ms
TotalAlloc = 117_737 MiB HeapAlloc = 5_201 MiB    StackInuse = 576 Kb     Sys = 8_926_631 Kb        Frees = 2_200_005_331      NumGC = 353


#### mapping

Elapsed time: 77_576 ms
TotalAlloc = 122_005 MiB HeapAlloc = 9_864 MiB    StackInuse = 544 Kb     Sys = 14_229_293 Kb       Frees = 2_151_872_625      NumGC = 347
