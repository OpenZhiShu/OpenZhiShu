# OpenZhiShu

## Config

`config.json`

### 常見屬性

- type: string  
    指定這個元素的內容應該以什麼形式解析
- content: string  
    這個元素的內容，見該類型的說明
- layout: string  
    決定這個元素在畫面上的布局
    - size: `width` and `height`  
        推薦使用百分比, 能在不同的視窗大小下保持比例
    - position: (`top` or `bottom`) and (`left` or `right`)  
        推薦使用百分比, 能在不同的視窗大小下保持比例
- style: string  
    附加在這個元素上的css屬性, 例如圖片的`object-fit: cover;`
- link: string (optional)  
    設定連結網址, 例如: `https://github.com/Shiphan`, 或是在網站內移動, 例如回到首頁: `/`

### 各欄位的說明
    
- homepage
    - body_color: string  
        當視窗比例不符合時，"黑邊"的顏色
    - ratio: string  
        框架的比例，背景及所有元素都會在這個框架中
    - background
        - image
            - type: `image`
            - content: string  
                圖片的網址, 例如: `https://path.to.img`, 或在`./assets/static/`中的圖片, 例如: `/static/images/image.png`
            - style: string  
                附加在這張圖片上的css屬性, 例如`width: 10%; top: 10%; right: 15%`
                - object-fit: `object-fit`  
                    推薦使用`cover`, 可用的值請見[說明](https://developer.mozilla.org/en-US/docs/Web/CSS/object-fit)
        - video
            - type: `video`
            - content: string  
                影片的網址, 例如: `https://path.to.video`, 或在`./assets/static/`中的影片, 例如: `/static/images/home_icon.png`
            - style: string  
            - autoplay: bool
            - loop: bool
            - muted: bool

    - elements: []element
        - image
            - type: `image`
            - content: string  
                圖片的網址, 例如: `https://path.to.img`, 或在`./assets/static/`中的圖片, 例如: `/static/images/image.png`
            - layout: string  
            - style: string  
                附加在這張圖片上的css屬性, 例如`width: 10%; top: 10%; right: 15%`
                - object-fit: `object-fit`  
                    推薦使用`cover`, 可用的值請見[說明](https://developer.mozilla.org/en-US/docs/Web/CSS/object-fit)
            - link: string (optional)  
                設定連結網址, 例如: `https://github.com/Shiphan`, 或是在網站內移動, 例如回到首頁: `/`
- drawing

-  
