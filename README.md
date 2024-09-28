# OpenZhiShu

OpenZhiShu是一個抽直屬的工具, 圍繞以下幾個核心理念設計:
- 明確的名單檔案, 非常低的使用門檻
- 以網頁構造的UI, 提供良好的跨平台支援性
- 簡單的設定檔, 具備高度可自訂性

---

## 使用方法

**NOTE:** 請至少閱讀本段落以確保程式執行預期內的行為, 如果還想要自訂畫面, 則請閱讀[Config](#config)段落  

### 準備階段

首先, 我們有一些必要的檔案要下載和編譯
1. 最簡單的方法是使用
    ```
    git clone https://github.com/OpenZhiShu/OpenZhiShu.git
    ```
    將這個倉庫下載到本地, 並使用
    ```
    go build
    ```
    將原始碼編譯為可執行檔

2. 若你的環境沒有且不想安裝git或go, 可以參考以下做法:
    - 沒有git:  
        從GitHub的Download ZIP功能直接下載最新版本, ~或是前往Releases頁面下載特定版本的壓縮檔~
    - 沒有go:  
        ~前往Releases頁面下載編譯完成的可執行檔, 並將其放入專案的目錄中~

準備完成後, 應該會有一個類似這樣的目錄結構:
```
├─ assets/
|  ├─ static/    # 裡面的檔案能透過`/static/paht/to/file`存取
|  └─ templates/
├─ pkg/
├─ .gitignore
├─ README.md
├─ main.go
├─ OpenZhiShu    # 編譯後產生的可執行檔, 在Windows環境則為OpenZhiShu.exe
├─ list.json     # 儲存名單的檔案
├─ config.json   # 設定檔, 修改前請參考Config段落
└─ results.json  # 儲存結果的檔案, 每次使用都會覆蓋它的內容
```

### 名單

名單應儲存在`list.json`中, 有兩個欄位`freshmen`代表新生和`seniors`代表學長姊  
這個範例應該能說明一切:

```json
{
    "freshmen": [
        {
            "number": 1,
            "name": "王小名"
        },
        {
            "number": 2,
            "name": "王大名"
        }
    ],
    "seniors": [
        {
            "number": 1,
            "name": "王中名"
        }
    ]
}
```

### 開始使用

只差最後一步了! 運行`OpenZhiShu`(或是OpenZhiShu.exe如果在Windows環境), 你應該會看到一個訊息
```
$ ./OpenZhiShu
choose a port to listen: 
```
在這裡, 你應該選擇一個通訊埠(port)使用, 在以下的範例中, 將以8080作為代表
```
$ ./OpenZhiShu
choose a port to listen: 8080
http://localhost:8080
```
到了這裡, 你就可以打開你最愛的瀏覽器並前往它提供的網址`http://localhost:8080`了!

### 主要頁面

- 主頁 `/`
- 準備 `/drawing`
- 結果 `/result/{number}`
- 設定畫面 `settings`
    - 下載結果  
      **NOTE:** 不要把它存到目錄中的`results.json`, 因為它會被新的結果覆蓋

## Config

`config.json`

- homepage: page
- drawing: page
- result: page

### 頁面

page

- title: string  
    網頁標題
- body_color: string  
    當視窗比例不符合時, "黑邊"的顏色
- ratio: string  
    框架的比例, 背景及所有元素都會在這個框架中
- elements: []element  
    一個元素的陣列, 其中每個元素會根據他的索引有一個獨特的id, 例如首個元素的id為`element-0`

### 元素

element

- type: string  
    指定這個元素的內容應該以什麼形式解析
- content: string  
    這個元素的內容, 見該類型的說明
- style: string  
    附加在這個元素上的css屬性, 
    - size: `width` and `height`  
        推薦使用百分比, 能在不同的視窗大小下保持比例
    - position: (`top` or `bottom`) and (`left` or `right`)  
        推薦使用百分比, 能在不同的視窗大小下保持比例
    - others: 例如圖片的`object-fit: cover;`, 文字的`color: lightgray;`
- appear: int (optional)  
    元素出現的時間, 以ms為單位, 數值應>=0
- hide: int (optional)  
    元素出現的時間, 以ms為單位, 數值應>=0
- link: string (optional)  
    設定連結網址, 例如: `https://github.com/Shiphan`, 或是在網站內移動, 例如回到首頁: `/`

### 各元素的說明

以下省略常見屬性, 僅列出該元素用法不同之處

- image
    - type: `image`
    - content: string  
        圖片的網址, 例如: `https://path.to.img`, 或在`./assets/static/`中的圖片, 例如: `/static/images/image.png`
- video
    - type: `video`
    - content: string  
        影片的網址, 例如: `https://path.to.video`, 或在`./assets/static/`中的影片, 例如: `/static/videos/video.mp4`
    - autoplay: bool
    - loop: bool
    - muted: bool
- text
    - type: `text`
    - content: string  
        顯示的文字內容
- variable  
    轉換成`to_type`的元素, 並取得變數後放入`content`
    - type: `variable`
    - content: string  
        變數的名稱
    - to_type: string
    - prefix: string  
        前綴, 加在變數值的前面
    - suffix: string  
        後綴, 加在變數值的後面

- input  
    輸入欄位, 詳見[說明](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input)
    - type: `input`
    - content: string  
        placeholder的值, 顯示在空的輸入欄
    - input_type
        type的值, [說明](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input#input_types)

- jump  
    一個按鈕, 用來讀取一個input元素的值並加上url_prefix後重新導向
    - type: `jump`
    - content: string  
        按鈕上的文字
    - target: string
        取得重新導向的目標, 應為一個input元素的id, 例如`element-0`
    - url_prefix: string
        加在目標URL前的前綴
