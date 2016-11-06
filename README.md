Model templates with gorp orm for gin 
====================================

[![Build Status](https://travis-ci.org/yvesago/gin-model-template.svg?branch=master)](https://travis-ci.org/yvesago/gin-model-template)
[![Coverage Status](https://coveralls.io/repos/github/yvesago/gin-model-template/badge.svg?branch=master)](https://coveralls.io/github/yvesago/gin-model-template)

## Description

Data model templates with gorp orm for gin

## Usage

Create a "models" directory in your project

  cp gin-model-template/*.go  myproject/models/ 

``agent.go`` and ``user.go`` are tables templates. Feel free to rename files and fix ``XXX`` tags.
Test with ``go test``

``repo.go`` contains database parameters. 

In your main.go project import ``./models``

Sample :

```go
  package main
    
  import (
      "fmt"
      "github.com/gin-gonic/gin"
      . "./models"
  )
    
  func main() {
    r := gin.Default()
    
    r.Use(Database("test.sqlite3"))
  
    v1 := r.Group("api/v1")
    {
        v1.GET("/users", GetUsers)
        v1.GET("/users/:id", GetUser)
        v1.POST("/users", PostUser)
        v1.PUT("/users/:id", UpdateUser)
        v1.DELETE("/users/:id", DeleteUser)
        v1.OPTIONS("/users", Options)     // POST
        v1.OPTIONS("/users/:id", Options) // PUT, DELETE
     ...

```


## TODO

Sample with config

## Licence

MIT License

Copyright (c) 2016 Yves Agostini

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.


<yves+github@yvesago.net>
