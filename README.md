# http-sniffer
This package parses your access log file entries using standard w3 httpd access log format. It keeps track of Global Request and Request Size totals as well as their respective
running averages. It also publishes more in depth statistics at a configurable interval. Alerts are updated based on your monitoring flags (Interval and Request threshhold)

## Usage

HTTP-SNIFFER monitors acess logs

### Synopsis

HTTP-SNIFFER monitors acess logs allowing 
	flexibility with cinfigurable options such as 
	request threshhold, statistics capturing interval, and alerting

```
run [flags]
```

### Options

```
  -h, --help                    help for run
  -f, --log-file string         path to access log file (default "/tmp/access.log")
  -m, --monitor-interval int    how often req/s threshold should be checked (default 2)
  -t, --monitor-threshold int   req/s (default 10)
  -s, --stats-interval int      how often stats should be captured (default 10)
  -w, --web                     view output in web UI
```

Example running and parsing the test files included in this project
```linux
$ go run cmd/main.go run -f ./test_files/access_log.log 
```

Sample ouput
<img src="http://i66.tinypic.com/jac094.png" border="0">

Please submit here:
https://app.greenhouse.io/tests/4c4cee4a387e4510e421b3585d7523b8