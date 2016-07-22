# kapacitor_utils

Useful TICKScripts and UDFs for Kapacitor

## How Do I Use It?

##### 1. Download the project

If you're just using TICKScripts you can simply `git clone` this repository.

If you want to build the UDFs, use `go get`

```
go get github.com/mpchadwick/kapacitor_utils
```

##### 2. Install the dependencies (if you want to build the UDFs)

```
cd $GOPATH/src/github.com/mpchadwick/kapacitor_utils
go get -d ./...
```

##### 3. Build the UDF you need

```
go install github.com/mpchadwick/kapacitor_utils/query_parser
```

## What Does It Do?

##### `query_parser`

Includes a UDF which is intended to ingest a stream sent to InfluxDb from [the logparser telegraf plugin](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/logparser). The UDF processes the URLs and outputs a stream of data which reports on query parameter usage in the URLs.

Understanding query parameter usage is valuable when for optimizing your page cache hit rate. For example, Varnish provides [an example of how to strip `gclid`](https://www.varnish-cache.org/trac/wiki/VCLExampleStripGoogleAdwordsGclidParameter) from URLs on their Wiki.

A sample TICKScript which makes use of the UDF is also included.

### More to come!
