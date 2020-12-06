# GO-REST-API
![License: MIT](https://img.shields.io/badge/Language-Golang-blue.svg)
[![Build Status](https://travis-ci.org/edwinnduti/go-rest-api.svg?branch=master)](https://travis-ci.org/edwinnduti/go-rest-api)
![License: MIT](https://img.shields.io/badge/Database-MongoDB-lightgreen.svg)

A rest API made in Golang.

### Requirements
<ul>
	<li>GOLANG</li>
	<li>POSTMAN</li>
	<li>MONGO COMPASS GUI</li>
</ul>

To run it locally:
```bash
$ git clone https://github.com/edwinnduti/go-rest-api.git

$ cd go-rest-api

$ go install 

$ go run main.go

```

Available :

| function              |   path                    |   method  |
|   ----                |   ----                    |   ----    |
| Create user           |   /api					|	POST    |
| Get single user	    |   /api/{tablename}		|	GET     |
| Get All users         |   /api            		|	GET     |
| Delete single user	|   /api/{tablename}		|	DELETE  |
| update single user	|   /api/{tablename}		|	UPDATE  |
| Get QRcode alone 	    |   /api/{tablename}/qrcode |   GET     |   


Enjoy!
