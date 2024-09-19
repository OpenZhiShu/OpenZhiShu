# OpenZhiShu

## List

`list.json`

- freshmen: []int
- seniors: []int

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
    當視窗比例不符合時，"黑邊"的顏色
- ratio: string  
    框架的比例，背景及所有元素都會在這個框架中
- elements: []element  
    一個元素的陣列


### 元素

element

- type: string  
    指定這個元素的內容應該以什麼形式解析
- content: string  
    這個元素的內容，見該類型的說明
- style: string  
    附加在這個元素上的css屬性, 
    - size: `width` and `height`  
        推薦使用百分比, 能在不同的視窗大小下保持比例
    - position: (`top` or `bottom`) and (`left` or `right`)  
        推薦使用百分比, 能在不同的視窗大小下保持比例
    - others: 例如圖片的`object-fit: cover;`, 文字的`color: lightgray;`
- appear: int  
    元素出現的時間, 以ms為單位, 數值應>=0
- hide: int  
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
    能取的變數的值並顯示在一個"text"中
    - type: `variable`
    - content: string  
        變數的名稱

