# Short URL Demo

## 前置

* 電腦作業系統/OS為 Windows 10 家用版
* 使用Golang Version go1.15.6 windows/amd64
* 使用MongoDB Version v4.4.2
* 創建DataBase admin
* 若資料庫無資料，可先跑測試去創建資料，若無可從網頁的`urls`網址進入列表頁，並開始流程(輸入文字)
* DataBase admin的連線資訊：使用者帳號為`adminrj`，使用者密碼為`123456`

## 流程

* 可開啟`terminal`輸入`mongo.exe`或其他`啟動MongoDB的指令`來確認目前`DB`的狀況。
* 準備好開始後，透過`terminal`進入到此專案的`/src/app/main`目錄下，並輸入指令`go build`(為了保險起見)，再來，若要執行請輸入`./main.exe`，若要測試請輸入`go test`。

-----

(接下來步驟請透過網頁)
* 到 `/urls` 列表頁可以確認所有目前的列表。
* 查看完畢後，可以到一樣的頁面`/urls`，透過`destination`輸入您要轉的網址(要符合網址的格式)，或是透過`alias`輸入您想要客製的網址片段(非必填)。
* 輸入後便會新增`URL`，該頁面會顯示`您的URL`和客製後的`ShortURL`，可點選直接使用，或是按下`按鈕`返回到列表頁面`/urls`。 
* 若想刪除DataBase`admin`的Collections`urls`，可直接從`terminal`啟動MongoDB後，輸入指令`use admin` => `db.urls.drop()`，如此便會把所有`urls`的Collections刪除。