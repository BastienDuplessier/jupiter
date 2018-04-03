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

## Geolocation API

### Usage
Run with
```bash
mix run --no-halt
```
Or inside iex
```bash
iex -S mix
```

### Request
```bash
curl --request GET \
     --url "http://localhost:4000/jobs?latitude=48.8659387&longitude=2.34532&radius=10"
```

## Why Elixir ?
Elixir is a langage I really like. Many functionnal features (high order functions, lists/map easy to manipulate), a cool syntax, great tooling and awesome performance ! I've seen this as an opportunity to practice on my API skills with Elixir.

## Global design
I've used an umbrella design. Which means I separate the business side of the application (here, Geocoding) from the API part which will just get requests as inputs and format the response.
I think that's a great design because it avoid to pack all the logic in a web framework. It allow code to be clear about who do what. If there's a change to do about how we store/retrieve the jobs, I'll go in the Geocoding app. If I need to add a new endpoint, I'll go in the API.
