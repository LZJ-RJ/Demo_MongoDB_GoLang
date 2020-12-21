# Short URL Demo


## Before

* OS: Windows 10 Home
* Language: Golang version go1.15.6 windows/amd64
* DataBase: MongoDB (NoSQL) version v4.4.2
  - DB_NAME: admin
  - DB_USERNAME: adminrj
  - DB_PASSWORD:123456

## After
`In terminal`
* Enter `mongo.exe` to activate MongoDB for checking status of DB.
* At root path of this project, `/src/app/main`
  - input `go build` 
  - input `./main.exe`
  - input `go test`, if you want to test
* Input `mongo.exe` => `use admin` => `db.urls.drop()`, if you want to delete all this collection, `urls`, in MongoDB.
-----
`In browser`

* Input `/urls` to check list of urls.
* Input `/urls`, and then input data at the column, `destination`, its format must match url.
* Input `/urls`, and then input `alias`, if you want to customize url of `source`, optional.
* After creating data, you would get two urls, `Source`, `Destination`.
These two urls are the main purpose of the project.


-----------------------------
-----------------------------


## 前置

* 電腦作業系統: Windows 10 家用版
* 語言: Golang version go1.15.6 windows/amd64
* 資料庫:MongoDB (NoSQL) version v4.4.2
  - DB_NAME: admin
  - DB_USERNAME: adminrj
  - DB_PASSWORD: 123456

## 執行
(在終端機)
* 輸入`mongo.exe`或其他`啟動MongoDB的指令`來確認目前`DB`的狀況。
* 在`/src/app/main`目錄下
  - 輸入`go build`
  - 輸入`./main.exe`
  - 若要測試請輸入`go test`。
* 輸入`mongo.exe` => `use admin` => `db.urls.drop()`，若想刪除DataBase`admin`的所有Collections`urls`。

-----

(在瀏覽器)
* 輸入`/urls`後，可以到列表頁確認所有目前的URL。
* 輸入`/urls`後，在`destination`這個欄位，來輸入您要轉向的網址(要符合網址的格式)，或是透過`alias`輸入您想要客製的網址片段(非必填)。
* 輸入後便會新增url，該頁面會顯示`Source`和`Destination`這兩個網址。
這兩個網址為這專案的主要目的。 