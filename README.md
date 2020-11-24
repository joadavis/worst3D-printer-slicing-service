# Worst3D printer slicing service

Bringing the worst idea in 3D printing to the cloud

This is an experiment and excuse for me to try out different ways of creating a REST API for a service.
The idea is to use Golang and as simple of a framework as possible to allow submitting slicing jobs and retrieving status and results.

For a first pass, my intent is to try Gorilla as a framework.  As there are many other ways to create a REST API with Go, if I get a chance to implement more I'll put them off in decicated branches (might try Swagger for instance).  It may also be interesting to try this using a "serverless" architecture.

# Reference
There are a few sites I've used for reference in pulling this together.  None are an exact match for the result, but might be useful to anyone else who wants to try this out.
- https://golang.org/pkg/net/http/
- https://dzone.com/articles/how-to-write-a-http-rest-api-server-in-go-in-minut
- https://dev.to/moficodes/build-your-first-rest-api-with-go-2gcj
