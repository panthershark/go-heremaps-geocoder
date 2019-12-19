# go-heremaps-geocoder

Utility for geocoding addresses with the HERE maps api

This library takes a csv file as input, and writes json file as output including lat/lon. It writes json b/c javascript can simply parse it and languages like golang allow you to use structs to marshall. It could be changed to write csv or any other format if you like.

### CLI tool

The compiled CLI tool is checked into this repo so you don't need to install anything or run build to use it.. `./geocoder` is an executable so just run it. 


This following example runs it against the example input file.
```bash
geocoder --in=./example/basic.csv --out=./example/out.json --apikey="here-maps-api-key"
```

### Development

Setup go on your machine. Clone this repo into you go path. Run for development

```bash
go run main.go --in=./example/basic.csv --out=./example/out.json --apikey="here-maps-api-key"
```

### Input file format

The input file is a csv b/c most db engines easily export recordsets as csv. Example at [](./in.csv)

* column 1 - The id of the record in the db. An address will be associated with a row in a database. This is the id of the record which contains this address. This is passed back in the results file.
* column 2 - The address to geocode. e.g. `4200 North Lamar Blvd Austin TX` (There is a good taco place at this address)
