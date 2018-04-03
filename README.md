# Jupiter

## Regroup by continent

### Usage
While at the root of the project
```bash
go run scripts/regroup_by_continent.go
```

### Dependencies
```bash
go get github.com/olekukonko/tablewriter
```

### Why go ?
Many peole say that Go is awesome, I wanted to give it a try. Plus, Go has good performances and that's important for this. `go run` makes easy to run Go code as a script.


### Scaling
If we need to be more performant and handle hundred of thousand job offers, I'd say we could split the input in batches (like 10k) and then feed them to concurrent processes (here, Goroutines). Then another process will wait all responses to be computed to display the result.
In the case of we have 1k offers every second, I'd recommand the usage of a database that support transactions. Like explained above, each concurent process will compute its batch of data. At the end, it will write its result into the database using transactions to be sure it avoid any race condition.


