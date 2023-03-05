package main

import ( 
  "restapi/routes"
  "restapi/db"
)
  
func main() {

  db.Init()

  e := routes.Init()

  e.Logger.Fatal(e.Start(":1234"))
}