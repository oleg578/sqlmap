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
