# STLKER
A stock tracking chat bot that utilizes RESTful and gRPC microservices

## How it works?
The user sends a HTTP request to the control API's endpoint which then forwards it to the unary gRPC API. The gRPC API then sends a request to a third-party API, Alpha Vantage,
and sends back the result to the user in a similar manner. 

## Current Features
* Utilizes a unary gRPC microservice as well as a RESTful one
* JWT authentication and cookies
* Extensive use of GORM
* Redis Caching
* Test cases

## Future Features
* Front-end web application
* 

## Current endpoints
* /portfolios (Requires authentication)
* /portfolios/{name} (Requires authentication)
* /portfolios/{name}/{ticker} (Requires authentication)
* /stocks/more/{ticker}
* /stocks/{ticker}/{currency}
* /signup
* /login
* /logout
* /refresh
* /deleteuser (Requires authentication)
