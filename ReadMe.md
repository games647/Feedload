# Feedload

## Description

Small and simple go application that reads an RSS feed and downloads all episodes 
starting from the latest one. Rather than many shell or python scripts, this
application saves the episodes according to their episode title rather than their
file name on the server. 

Furthermore, it downloads the episodes in parallel with the aim of learning this 
programming language go.

## Features

* Starts 8 workers and downloads the files in parallel
* Progress bar for each worker
* Total progress bar
* Native binary

## Compiling

1. Download project
2. Download dependencies:
    1. `go get github.com/mmcdole/gofeed`
    2. `go get github.com/vbauerster/mpb`
3. Build executable `go build feedload.go`

## Running

* `go run feedload.go <RSS_URL>` (Without compiling and with golang installed)
* `./feedload <RSS_URL>` (*Nix)
* `feedload <RSS_URL>` (Windows)

## Example Output


>Downloading all episodes from PODCAST_NAME \
All                             0/339 [-------------------]  0s  0 % \
Episode 205: XXXX     1.8MiB/105.1MiB [->-----------------] 1m35s 2 % \
Episode 208: XXXX     1.6MiB/106.2MiB [->-----------------]  5s  2 % \
Episode 209: XXXX     1.6MiB/105.3MiB [->-----------------] 45s  2 % \
Episode 203: XXXX     1.7MiB/109.2MiB [->-----------------]  6s  2 % \
Episode 206: XXXX     1.7MiB/111.7MiB [->-----------------]  6s  2 % \
Episode 210: XXXX     1.7MiB/124.3MiB [>------------------]  7s  1 % \
Episode 204: XXXX     1.3MiB/108.5MiB [>------------------]  9s  1 % \
Episode 207: XXXX     1.2MiB/110.0MiB [>------------------] 10s  1 % \
