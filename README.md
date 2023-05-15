# property-search-aws

This is just a basic scaffold to get one started with a property search API.
Full-text search is provided by RediSearch run in a container.
A Python script is provided to load sample data into RediSearch.
A GO app provides the search API.

---

Primer for RediSearch:  
https://www.youtube.com/watch?v=infTV4ifNZY&list=PLratyGi2ixLsqd3SRcsJticE9yt5LDX3R

---

## Quick Start with Docker

1. Run: `docker-compose up`
2. Open your browser with a valid search URL. e.g. http://localhost:3000/search?q=9005 

Data will automatically be loaded from "MessySampleData.txt"

---

## API Search Reference

All searches use the HTTP `GET` method. Just set the appropriate query params as needed:

 - `q`: "query"; the term or terms to query; terms must be separated by `+`; e.g. http://localhost:3000/search?q=9005+9006
 - `l`: limit for paging; min is 25 and max is 100; e.g. http://localhost:3000/search?q=owen&l=30
 - `o`: offset for paging

A health check URL is provided at `/`. e.g. GET http://localhost:3000/ returns HTTP Status 200 with `{"message":"OK"}`

Read the Go code in [goapp/server.go](goapp/server.go) for additional reference.
