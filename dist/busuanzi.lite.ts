(function (){
    // 在此处设置您的后端地址 如 https://example.com/api
    let url: string = `${__BSZ_API_DOMAIN__}/api`,
        tags: string[] = ["site_pv", "site_uv", "page_pv", "page_uv", "today_pv", "today_uv", "yesterday_pv", "yesterday_uv", "month_pv", "month_uv"],
        prefix: string = "busuanzi",    // 自定义标签ID前缀
        storageName: string = "bsz-id";                                        // 本地存储名称

    let bsz_send = () => {
        let xhr: XMLHttpRequest = new XMLHttpRequest();
        xhr.open("POST", url, true);

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
                            if (element != null) element.innerHTML = res['data'][tag];

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
})()