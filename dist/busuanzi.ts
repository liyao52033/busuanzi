(function (){
    // 在此处设置您的后端地址 如 https://example.com/api
    let url: string = `${__BSZ_API_DOMAIN__}/api`,
        tags: string[] = ["site_pv", "site_uv", "page_pv", "page_uv", "today_pv", "today_uv", "yesterday_pv", "yesterday_uv", "month_pv", "month_uv"],
        current: HTMLOrSVGScriptElement = document.currentScript,
        pjax: boolean = current.hasAttribute("pjax"),                          // 是否启用 pjax
        api: string = current.getAttribute("data-api") || url,                 // 自定义后端地址
        prefix: string = current.getAttribute("data-prefix") || "busuanzi",    // 自定义标签ID前缀
        style: string = current.getAttribute("data-style") || "default",       // 数字显示风格 default | comma | short
        storageName: string = "bsz-id";                                        // 本地存储名称

    let format = (num: number, style: string = 'default'): string => {
        if (style === "comma") return num.toLocaleString();
        if (style === "short") {
            if (num === 0) return "0";  // 处理 0 的特殊情况
            const units = ["", "K", "M", "B", "T"];
            let index = Math.floor(Math.log10(num) / 3);
            num /= Math.pow(1000, index);
            return `${Math.round(num * 100) / 100}${units[index]}`;
        }
        return num.toString();
    };

    let bsz_send = () => {
        let xhr: XMLHttpRequest = new XMLHttpRequest();
        xhr.open("POST", api, true);

        // set user identity
        let token: string | null = localStorage.getItem(storageName);
        if (token != null) xhr.setRequestHeader("Authorization", "Bearer " + token);
        xhr.setRequestHeader("x-bsz-referer", window.location.href);
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
                if (xhr.status === 200) {
                    let res: any = JSON.parse(xhr.responseText);
                    if (res.success === true) {
                        tags.map((tag: string) => {
                            let element = document.getElementById(`${prefix}_${tag}`);
                            if (element != null) element.innerHTML = format(res['data'][tag], style);

                            let container = document.getElementById(`${prefix}_container_${tag}`);
                            if (container != null) container.style.display = "inline";
                        })

                        let setIdentity = xhr.getResponseHeader("Set-Bsz-Identity")
                        if (setIdentity != null && setIdentity != "") localStorage.setItem(storageName, setIdentity);
                    }
                }
            }
        }
        xhr.send();
    };
    bsz_send();

    if (!!pjax) {
        let history_pushState: Function = window.history.pushState;
        window.history.pushState = function () {
            history_pushState.apply(this, arguments);
            bsz_send();
        };

        window.addEventListener("popstate", function (_e: PopStateEvent) {
            bsz_send();
        }, false);
    }
})()