*catatan
- setiap uji coba, restart servernya agar kode terimplementasi
- setiap nama variabel atau function yg ingin di import untuk digunakan file lain maka huruf pertama nama variabel tsb harus kapital
- catatan cek jika ada depencies yg kesatuan, misal kita instal echo/v4 maka middleware nya juga harus echo/v4/middleware jangan echo/middleware saja
  
- bedakan Param dan FormValue
  - param: variabel di url
  - FormValue : variabel dalam bentuk form di kirim di body 

ref
- https://www.youtube.com/playlist?list=PLO2Rv4lKm-K0vR95sfEXznno4421Z67OF

## Fundamental

1. insialisasi project
   ``go mod init namaProject``
   secara otomatis terbuat file  go.mod

2. install echo framework
   ``go get github.com/labstack/echo/v4@latest``
   ref: 
   - https://github.com/labstack/echo
   - https://echo.labstack.com/guide/

3. buat file server.go

4. masukkan kode berikut pada fie tsb
   ```go
    package main

    import (
    "net/http"
    "github.com/labstack/echo/v4" 
    )


    func main() {
    e := echo.New()

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, this is echo")
    })

    e.Logger.Fatal(e.Start(":1234"))
    }
   ```
5. uji coba. jalankan kode berikut di terminal
   ``go run server.go``

    coba buka di browser url berikut ``localhost:1234``

    jika telah terbuka akan memunculkan tulisan "Hello, this is echo"

## Using module and make structure folder

make structure folder
- projectFolder
  - config
  - controller
  - db
  - models
  - routes   

1. kode pada server.go kita pindahkan ke routes/routes.go
   ```go
    package routes

    import (
        "net/http"
        "github.com/labstack/echo/v4" 
    )


    func Init() *echo.Echo {
    e := echo.New()

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, this is echo")
    })

    return e
    }
   ```
   catatan:
   - nama package nya menyesuaikan dengan foldernya yaitu ./routes/
   - ``package main`` digunakan jika tidak didalam folder atau berada di folder utama project
   - ``*echo.Echo {`` artinya return nilai

2. selanjutnya pada file server.go 
   ```go
    package main

    import ( 
        "restapi/routes"
    )
    
    func main() {
        e := routes.Init()

        e.Logger.Fatal(e.Start(":1234"))
    }
   ```

   catatan: 
   - kita import file routes moduleName/packageName ``restapi/routes``
   - pada func main kita memanggil fungsi init dari file routes
   - define port 1234

3. run ulang.

## Koneksi database

1. instal gonfig
   - dependencies untuk manage configuration file
   - go get github.com/tkanos/gonfig
   - ref: https://github.com/tkanos/gonfig

2. di folder ./config buat file
   - config.go
   - config.json

3. isi file tersebut dengan
   - config.json
        ```json
        {
            "DB_USERNAME" : "root",
            "DB_PASSWORD" : "",
            "DB_PORT" : "3306",
            "DB_HOST" : "127.0.0.1",
            "DB_NAME" : "test"
        } 
        ```
   - config.go

        ```go
        package config

        import "github.com/tkanos/gonfig"

        // menampung data
        type Configuration struct {
            DB_USERNAME string
            DB_PASSWORD string
            DB_PORT string
            DB_HOST string
            DB_NAME string
        }

        func GetConfig() Configuration {
            conf := Configuration{}

            //insert variabel conf <- data config.json
            gonfig.GetConf("config/config.json" , &conf)
            return conf
        }
        ```
4. install mysql driver
   - https://github.com/go-sql-driver/mysql
   - go get -u github.com/go-sql-driver/mysql

5. buat connection
   - pada folder db buat file db.go dan isi dengan kode berikut

        ```go
        package db

        import (
            "database/sql"
            "restapi/config"

            _ "github.com/go-sql-driver/mysql"
        )

        var db *sql.DB
        var err error

        func Init(){
            conf := config.GetConfig()

            connectionString := conf.DB_USERNAME + ":" + conf.DB_PASSWORD + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME

            db, err = sql.Open("mysql", connectionString)
            
            if err != nil {
                panic("connectionString error ..")
            }

            err = db.Ping()
            if err != nil {
                panic("DSN Invalid ..")
            }
        }

        func CreateCon() *sql.DB {
            return db
        }
        ```
   - pada file server.go, tambahkan kode berikut

        ```go 
        ...
        import ( 
            ...
            "restapi/db"
        )
        
        func main() {

            db.Init()
            ...
        }
        ``` 
6. buat table pegawai
   - id int
   - nama varchar
   - alamat varchar
   - telepon varchar

7. silahkan run ulang


## Method GET 

1. pada folder models buat file ``response.go`` dan isi dengan kode berikut

    ```go
    package models

    type Response struct {
        Status int `json:"status"`
        Message string `json:"message"`
        Data interface{} `json:"data"`
    }
    ```

2. masih di folder models buat file ``m_pegawai.go`` isi dengan kode berikut

    ```go
    package models

    import (
        "net/http"
        "restapi/db"
    )

    type Pegawai struct {
        Id      int     `json:"id"`
        Nama    string  `json"nama"`
        Alamat  string  `json:"alamat"`
        Telepon string  `json:"telepon"`
    }

    func FetchAllPegawai() (Response, error) {
        var obj Pegawai
        var arrobj []Pegawai
        var res Response

        con := db.CreateCon()

        sqlStatement := "SELECT * FROM pegawai"

        rows, err := con.Query(sqlStatement)
        defer rows.Close()

        // kondisi jika sql query error
        if err != nil {
            return res, err
        }

        for rows.Next() {
            // kondisi jika ada data
            err = rows.Scan(&obj.Id, &obj.Nama, &obj.Alamat, &obj.Telepon)
            if err != nil {
                return res, err
            }

            // insert 
            arrobj = append(arrobj, obj)
        }
    
        res.Status = http.StatusOK
        res.Message = "Success"
        res.Data = arrobj

        // return value Response, error
        return res, nil 
    }
    ```

3. pada file controller buat file c_pegawai.go dan isi kode berikut:

    ```go
    package controller

    import (
        "net/http"
        "restapi/models"
        "github.com/labstack/echo/v4o"
    )

    func GetAllPegawai(c echo.Context) error {
        result, err := models.FetchAllPegawai()

        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error() } )
        }

        return c.JSON(http.StatusOK, result)
    }
    ```

4. di file routes.go tambahkan kode berikut

    ```go 
    ..
    import (
    ...
    "restapi/controller" 
    ...
    )

    func Init() *echo.Echo { 
    ...
    e.GET("/pegawai", controller.GetAllPegawai)
    ...
    }
    ```

5. uji coba di browser atau di postman url berikut localhost:1234/pegawai. jika berhasil akan menampilkan data berikut

    ```json
    {
        "status":200,
        "message":"Success",
        "data":[
        {
            "id":1,
            "Nama":"Aswin",
            "alamat":"Jakarta",
            "telepon":"8233434343"
        },
    ...]
    }

    ```

## Metod Post 

1. pada file ``m_pegawai`` tambahkan kode berikut dibawah ``func FetchAllPegawai`` 

    ```go
    ...

    func StorePegawai(nama string, alamat string, telepon string) (Response, error) {
        var res Response

        con := db.CreateCon()

        sqlStatement := "INSERT pegawai (nama, alamat, telepon) VALUES (?, ?, ?)"

        stmt, err := con.Prepare(sqlStatement)
        if err != nil {
            return res, err
        }

        result, err := stmt.Exec(nama, alamat, telepon)
        if err != nil {
            return res, err
        }

        lastInsertedId, err := result.LastInsertId()
        if err != nil {
            return res,err
        }

        res.Status = http.StatusOK
        res.Message = "Success"
        res.Data = map[string]int64{
            "last_inserted_id" : lastInsertedId,
        }

        return res, nil
    }
    ```

2. pada file c_pegawai.go tambahkan kode berikut dibawah ``func GetAllPegawai``

    ```go
    ...
    func InsertPegawai(c echo.Context) error {
        nama := c.FormValue("nama")
        alamat := c.FormValue("alamat")
        telepon := c.FormValue("telepon")

        result, err := models.StorePegawai(nama, alamat, telepon)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err.Error())
        }

        return c.JSON(http.StatusOK, result)
    }
    ```

3. pada file ``routes.go`` tambahkan kode berikut

    ```go
    func Init() *echo.Echo {
    ...
    e.POST("/pegawai", controller.InsertPegawai)
    ...
    }
    ```
4. silahkan uji coba di postman dengan memilih method POST, url : localhost:1234/pegawai dan masukkan ``params dan value``  berikut kemudian send
   - nama : adhy
   - alamat : jepang
   - telepon : 8343424

5. jikka berhasil maka akan muncul data berikut

    ```json
    {
        "status": 200,
        "message": "Success",
        "data": {
            "last_inserted_id": 5
        }
    }
    ```

## Method PUT

1. pada file ``m_pegawai.go`` tambahkan kode berikut di bawah ``func StorePegawai``

    ```go
    func UpdatePegawai(id int, nama string, alamat string, telepon string) (Response, error) {
        var res Response

        con := db.CreateCon()

        sqlStatement := "UPDATE pegawai SET nama = ?, alamat = ?, telepon = ? WHERE id = ?"

        stmt, err := con.Prepare(sqlStatement)

        result, err := stmt.Exec(nama, alamat, telepon, id)
        if err != nil {
            return res, err
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            return res, err
        }

        res.Status = http.StatusOK
        res.Message = "Success"
        res.Data = map[string]int64{
            "rows_affected" : rowsAffected,
        }

        return res, nil
    }
    ```

2. pada file controller tambahkan kode berikut

    ```go
    import (
        "strconv" 
        ...
    )
    ...
    func UpdatePegawai(c echo.Context) error {
        id := c.FormValue("id")
        nama := c.FormValue("nama")
        alamat := c.FormValue("alamat")
        telepon := c.FormValue("telepon")

        // convert tpye id from string to int
        conv_id, err := strconv.Atoi(id)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err.Error())
        }

        result, err := models.UpdatePegawai(conv_id, nama, alamat, telepon)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err.Error())
        }

        return c.JSON(http.StatusOK, result)
    }
    ```

3. pada file routes.go tambahkan kode berikut

    ```go
    func Init() *echo.Echo {
    ...
    e.PUT("/pegawai", controller.UpdatePegawai)
    ...
    }
    ```

4. uji coba url yg sama ``localhost:1234/pegawai``, pada postman pilih method put, kemudian pilh body form-data kemudian masukkan data key-value berikut
   - id : 5
   - nama : asri
   - alamat : bandung
   - telepon : 8435454
   - kemudian send

5. jika berhasil akan mendatkan pesan berikut

    ```json
    {
        "status": 200,
        "message": "Success",
        "data": {
            "rows_affected": 1
        }
    }
    ```

## Method Delete

1. pada file ``m_pegawai.go`` tambahkan kode berikut dibawah ``func UpdatePegawai()``

    ```go
    ...
    func DeletePegawai(id int) (Response, error) {
        var res Response

        con := db.CreateCon()

        sqlStatement := "DELETE FROM pegawai WHERE id = ?"

        stmt, err := con.Prepare(sqlStatement)
        if err != nil {
            return res, nil
        }

        result, err := stmt.Exec(id)
        if err != nil {
            return res, nil
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            return res, nil
        }

        res.Status = http.StatusOK
        res.Message = "Success"
        res.Data = map[string]int64{
            "rows_affected": rowsAffected,
        }

        return res, nil
    }
    ```

2. pada file ``c_pegawai.go`` tambahkan kode berikut dibawah ``func UpdatePegawai()``

    ```go
    func DeletePegawai(c echo.Context) error {
        id := c.FormValue("id")

        conv_id, err := strconv.Atoi(id)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err.Error() )
        }

        result, err := models.DeletePegawai(conv_id)

        return c.JSON(http.StatusOK, result)
    }
    ```

3. pada file ``routes.go`` tambahkan kode berikut:

    ```go
    func Init() *echo.Echo {
    ...
    e.DELETE("/pegawai", controller.DeletePegawai)
    ...
    }

4. uji coba dengan url yg sama ``localhost:1234/pegawai``, pada postman pilih method ``delete``, kemudian pilh body form-data kemudian masukkan data key-value berikut
   - id : 5 

5. jika berhasil akan mendatkan pesan berikut

    ```json
    {
        "status": 200,
        "message": "Success",
        "data": {
            "rows_affected": 1
        }
    }

## Implementasi JWT to secure API

### hashing

ref
- https://pkg.go.dev/golang.org/x/crypto/bcrypt

1. buat table users
   - id int
   - username varchar
   - password varchar

2. import function to generate hash password
   - go get golang.org/x/crypto/bcrypt

3. buat folder helper dan file paswordHelper.go dan masukkan kode berikut

    ```go
    package helpers

    import (
        "golang.org/x/crypto/bcrypt"
    )

    func HashPassword(password string) (string, error) {
        bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost )
        return string(bytes), err
    }

    func CheckPasswordHash(password, hash string) (bool, error) {
        err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password) )
        if err != nil {
            return false, err
        }

        return true, nil
    }
    ```

4. buat controller ``c_login.go`` dan masukkan kode berikut

    ```go
    package controller

    import (
        "net/http"
        "restapi/helpers"
        "github.com/labstack/echo/v4"
    )

    func GenerateHashPassword(c echo.Context) error {
        password := c.Param("password")

        hash, _ := helpers.HashPassword(password)

        return c.JSON(http.StatusOK, hash)
    }
    ```

5. pada file ``routes.go`` tambahkan kode berikut:

    ```go
    func Init() *echo.Echo {
    ...
    e.GET("/generate-hash/:password", controller.GenerateHashPassword)
    ...
    }

6. uji coba
   - silahkan ke postman dan masukkan dan gunakan url berikut 
   - localhost:1234/generate-hash/inipassword
   - variabel dari url terakhir adalah password yg anda ingin hashing
   - jika berhasil akan menampilkan hasil hashing contoh : "$2a$10$cL7YMjClsjNGe1RZS7QwkeeY/KfgIzOFnsK1IWZLqU6QI6dnR6z9a"

7. silahkan coba isi tabel users dengan username dan hashpassword


### authenticator

1. buat file model ``m_login.go`` dan masukkan kode berikut

    ```go
    package models

    import ( 
        "database/sql"
        "fmt"

        "restapi/db"
        "restapi/helpers"
    )

    type User struct {
        Id int `json:"id"`
        Username string `json:"username"`
    }

    func CheckLogin(username, password string) (bool, error) {
        var obj User
        var pwd string
        con := db.CreateCon()

        // search user
        sqlStatement := "SELECT * FROM users where username=?"

        err := con.QueryRow(sqlStatement, username).Scan(
            &obj.Id, &obj.Username, &pwd,
        )

        if err == sql.ErrNoRows {
            fmt.Println("Username not found")
            return false, err
        }

        if err != nil {
            fmt.Println("Query error")
            return false, err
        }

        match, err := helpers.CheckPasswordHash(password, pwd)
        if !match {
            fmt.Println("Hash and password doesnt match")
            return false, err
        }

        //if match
        return true, nil
    }
    ```

2. pada file controller ``c_login.go`` tambahkan kode berikut

    ```go
    ...
    func CheckLogin(c echo.Context) error {
        username := c.FormValue("username")
        password := c.FormValue("password")
        
        res, err := models.CheckLogin(username, password)
        if err != nil{
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "messages": err.Error(),
            })
        }

        if !res {
            return echo.ErrUnauthorized
        }

        return c.String(http.StatusOK, "berhasil login")
    }
    ```

3. pada file ``routes.go`` tambahkan kode berikut:

    ```go
    func Init() *echo.Echo {
    ...
    e.POST("/login", controller.CheckLogin)
    ...
    }

4. uji coba
   - silahkan ke postman pilih methode post
   - masukkan url localhost:1234/login
   - pilih body dan form-data 
     - input key value username dan password dengan data sembarang
   - jika username dan password salah akan memunculkan pesan berikut

        ```json
        {
            "messages": "sql: no rows in result set"
        }
        ```
   - selanjutnya masukkan user dan password yang benar, jika berhasil akan memunculkan pesan "berhasil login" 

### Generate token with jwt

ref 
- https://echo.labstack.com/cookbook/jwt/
- https://github.com/golang-jwt/jwt

1. instal jwt
   - go get -u github.com/golang-jwt/jwt/v5

2. pada file controller ``c_login.go`` disable kode ``return c.String(http.StatusOK, ...`` dan tambahkan kode berikut ini dibawahnya

    ```go
    func CheckLogin(c echo.Context) error {
        ....
        // return c.String(http.StatusOK, "berhasil login")
        
        token := jwt.New(jwt.SigningMethodHS256)

        claims := token.Claims.(jwt.MapClaims)
        claims["username"] = username
        claims["level"] = "application"
        claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

        t, err := token.SignedString([]byte("secret"))
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "messages": err.Error(),
            })
        }

        return c.JSON(http.StatusOK, map[string]string{
            "token": t,
        })

    }
    ```

3. run ulang dan uji coba login ulang, di url ``localhost:1234/login`` menggunakan user dan pass yang benar
   - jika berhasil maka bukan lagi teks ``login berhasil`` yang muncul tapi token

    ```json
    {
        "token": "eyJhbGciO..."
    }
    ```

### protect route url API wiht echo middleware

ref
- https://echo.labstack.com/middleware/

1. buat folder middleware dan file ``middleware.go``

2. install echo middleware
   - go get github.com/labstack/echo/v4/middleware

3. pada file ``middleware.go`` masukkan kode berikut

    ```go
    package middleware

    import (
        "github.com/labstack/echo/middleware"
    )

    var isAuthenticated = middleware.JWTWithConfig(middleware.JWTConfig{
        SigningKey: []byte("secret"),
    })
    ```

4. pada file ``routes.go`` import file middleware.go dan variabel IsAuthenticated di route getall pegawai

    ```go
    import ( 
        "restapi/middleware"  
    )

    func Init() *echo.Echo {
    ...
    e.GET("/pegawai", controller.GetAllPegawai, middleware.IsAuthenticated)
    ...
    }

5. selanjutnya uji coba akses url ``localhost:1234/pegawai`` dengan method get. jika muncul seperti berikut

    ```json
    {"message":"missing or malformed jwt"}
    ```
   - maka middleware sukses dipasangkan, selanjutnya utk mengakses data pegawai kita perlu token

6. nah, coba copy token dari fungsi generate token sebelumnya kemudian pilih variabel auth di postman dengan type Bearer Token dan pastekan token tersebut. 
   - selanjutnya akses kembali url ``localhost:1234/pegawai`` maka data telah muncul

7. jika testingnya telah berhasil maka implementasikan middleware di route data lainnya

    ```go
    e.GET("/pegawai", controller.GetAllPegawai, middleware.IsAuthenticated)

    e.POST("/pegawai", controller.InsertPegawai, middleware.IsAuthenticated)

    e.PUT("/pegawai", controller.UpdatePegawai, middleware.IsAuthenticated)

    e.DELETE("/pegawai", controller.DeletePegawai, middleware.IsAuthenticated)
    ```

## validasi user input

ref
- https://github.com/go-playground/validator

1. install dependencies
   - go get github.com/go-playground/validator/v10

### Multivalidation

2. buat file controller ``c_validation.go``  dan masukkan kode berikut

    ```go
    package controller

    import(
        // alias path
        "net/http"
        validator "github.com/go-playground/validator/v10"
        "github.com/labstack/echo"
    )
    
    type Customer struct {
        Nama 	string 	`validate:"required"`
        Email 	string 	`validate:"required,email"`
        Alamat 	string 	`validate:"required"`
        Umur	int 	`validate:"gte=17,lte=35"`
    }

    
    func TestStructValidation(c echo.Context) error {
        v := validator.New()

        //variabel inputan uji coba
        cust := Customer{
            Nama: "Asw",
            Email: "asw",
            Alamat: "",
            Umur: 15, 
        }

        err := v.Struct(cust)
        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{
                "message": err.Error(),
            })
        }

        return c.JSON(http.StatusOK, map[string]string{
            "message": "Success",
        })
    }

    ``` 
  - gte : greater than equal, lte : less then equal
  - stuct text validate.. dapat dilihat di dokumentasi playground validator
  - validasi variabel sekaligus dengan struct, 
  - validasi variabel juga bisa satu2 tanpa struct
  - err := v.Var()

3. pada file ``routes.go`` tambahkan kode berikut

    ```go
    import ( 
        "restapi/middleware"  
    )

    func Init() *echo.Echo {
    ...
    e.GET("test-struct-validation", controller.TestStructValidation)
    ...
    }

4. uji coba jalankan url berikut ``localhost:1234/test-struct-validation`` di postman dengan method get, 
   - jika berhasil akan menampilkan text berikut

    ```json
    {
        "message": "Key: 'Customer.Email' Error:Field validation for 'Email' failed on the 'email' tag\nKey: 'Customer.Alamat' Error:Field validation for 'Alamat' failed on the 'required' tag\nKey: 'Customer.Umur' Error:Field validation for 'Umur' failed on the 'gte' tag"
    }
    ```

5. selanjutnya coba benarkan variabel inputan, dan send ulang. jika berhasil akan memunculkan notifkasi success

    ```go
    //variabel inputan uji coba
        cust := Customer{
            Nama: "Asw",
            Email: "asw@gmail.com",
            Alamat: "makassar",
            Umur: 17, 
        }
    ```

### one validation

1. tambahkan kode berikut pada file ``c_validation.go``

    ```go
    func TestVariabelValidation(c echo.Context) error {
        v := validator.New()

        email := "asw"

        err := v.Var(email, "required,email")
        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{
                "message": "Email not valid",
            })
        }

        return c.JSON(http.StatusOK, map[string]string{
            "message": "Success",
        })
    }

    ```

2. pada file ``routes.go`` tambahkan kode berikut

    ```go
    import ( 
        "restapi/middleware"  
    )

    func Init() *echo.Echo {
    ...
    e.GET("test-variable-validation", controller.TestVariabelValidation)
    ...
    }

3. uji coba papa url berikut ``localhost:1234/test-variable-validation`` di postman dengan method get

    ```json
    {
        "message": "Email not valid"
    }
    ```

4. silahkan perbaiki nilai dari variabel email dan coba ulangi lagi

### implementasi 

1. pada file model ``m_pegawai`` modifikasi kode pada struct menjadi seperti berikut

    ```go
    type Pegawai struct {
        Id 		int 	`json:"id"`
        Nama 	string 	`json"nama" validate:"required"`
        Alamat 	string 	`json:"alamat" validate:"required"`
        Telepon string 	`json:"telepon" validate:"required"`
    }
    ```

2. pada file yg sama di fungsi StorePegawai modifikasi kodenya dengan menambahkan kode validation

    ```go
    
    import ( 
        validator "github.com/go-playground/validator/v10"
    )

    func StorePegawai(nama string, alamat string, telepon string) (Response, error) {
        var res Response

        // validation
            v := validator.New()

            peg := Pegawai{
                Nama: nama,
                Alamat: alamat,
                Telepon: telepon,
            }

            err := v.Struct(peg)
            if err != nil {
                return res, err
            }
        ...
    ```

3. di file controller file ``c_pegawai`` modifikasi kode penampil pesan errornya agar lebih readable

    ```go
    func UpdatePegawai(c echo.Context) error {
        ...
        result, err := models.UpdatePegawai(conv_id, nama, alamat, telepon)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "message": err.Error(),
            } )
        } 
        ...
    }
    ```

4. uji coba di url ``localhost:1234/pegawai``
   - jangan lupa menggunakan auth token
   - coba buat form-data dengan key value salah satunya didisable misal cuman kirim nama saja
   - kemudian sent maka akan mendapati error seperti dibawah maka validasi user berhasil
    ```json
    "Key: 'Pegawai.Alamat' Error:Field validation for 'Alamat' failed on the 'required' tag\nKey: 'Pegawai.Telepon' Error:Field validation for 'Telepon' failed on the 'required' tag"
    ```
   - kemudian enable all nama, alamat, telepon lalu kirim. jika data masuk maka berhasil 
