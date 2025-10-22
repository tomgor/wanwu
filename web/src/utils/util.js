import router from '@/router/index'
import { menuList } from "@/views/layout/menu"
import { checkPerm, PERMS } from "@/router/permission"
import { i18n } from "@/lang"
import { basePath } from "@/utils/config"

export function guid() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        var r = Math.random() * 16 | 0,
            v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

// 用于登录切组织等找到有权限的第一个菜单路径 (除用模型：用模型为打开的新页面)
export const fetchPermFirPath = (list = menuList) => {
    if (!list.length) return ''

    let path = ''
    for (let i in list) {
        const item = list[i]

        if (item.path && checkPerm(item.perm)) {
            if (item.children && item.children.length) {
                path = fetchPermFirPath(item.children).path
                break
            } else {
                path = item.path || ''
                break
            }
        }
    }
    console.log(path, list,  '----------------------------fetchPermFirPath')

    // 若有权限，跳转左侧菜单第一个有权限的页面；否则跳转 /404
    return {path: path || '/404'}
}

// 找到有权限的第一个菜单的 index
export const fetchCurrentPathIndex = (path, list) => {
    let index = ''
    const findIndex = (list) => {
        for (let i in list) {
            let item = list[i]
            const formatPath = (url) => url + '/'
            if (item.path && formatPath(path).includes(formatPath(item.path))) {
                index = item.index
            } else {
                if (item.children && item.children.length) {
                    findIndex(item.children)
                }
            }
        }
        return index
    }
    return findIndex(list)
}

export const jumpPermUrl = () => {
    const {path} = fetchPermFirPath()

    router.push({path: path || '/404'})
}

export const redirectUrl = () => {
    // 跳到有权限的第一个页面
    jumpPermUrl()
}

export const replaceIcon = (logoPath) => {
    let link = document.querySelector("link[rel*='icon']") || document.createElement("link")
    link.type = "image/x-icon"
    link.rel = "shortcut icon"
    link.href = logoPath ? basePath + '/user/api' + logoPath : basePath + '/aibase/favicon.ico'
    document.getElementsByTagName("head")[0].appendChild(link)
}

export const replaceTitle = (title) => {
    document.title = title || i18n.t('header.title')
}

export const getInitTimeRange = () => {
    const date = new Date()
    const month = date.getMonth() + 1
    const startTime =  date.getFullYear() + "-" + (month < 10 ? "0" : "") + month + "-"  + '01 00:00:00'
    const stamp= new Date().getTime() + 8 * 60 * 60 * 1000
    const endTime = new Date(stamp).toISOString().replace(/T/, ' ').replace(/\..+/, '').substring(0, 19)
    return [startTime, endTime]
}

export function convertLatexSyntax(inputText) {
    // 1. 匹配块级公式，将 `\[` 和 `\]` 替换为 `$$`，支持 `\\[` `\\]` 或单个 `\[` `\]`
    inputText = inputText.replace(/\\\[\s*([\s\S]+?)\s*\\\]/g, (_, formula) => `$$${formula}$$`);
    // 2. 匹配行内公式，将 `\(` 和 `\)` 替换为 `$`，支持 `\\(` `\\)` 或单个 `\(` `\)`
    inputText = inputText.replace(/\\\(\s*([\s\S]+?)\s*\\\)/g, (_, formula) => `$${formula}$`);
    return inputText;
}

export function isSub(data){
    return /\【([0-9]{0,2})\^\】/.test(data)
}

export function parseSub(data){
    return data.replace(/\【([0-9]{0,2})\^\】/g,(item)=>{
        let result = item.match(/\【([0-9]{0,2})\^\】/)[1]
        return `<sup class='citation'>${result}</sup>`
    })
}

/**
 *获取URL参数
 */
export function getQueryString(val, href) {
    const hrefNew = href || window.location.href;
    const search = hrefNew.substring(hrefNew.lastIndexOf('?') + 1, hrefNew.length);
    // 组装?
    const uri = '?' + search;
    const reg = new RegExp('' + val + '=([^&?]*)', 'ig');
    const matchArr = uri.match(reg);
    if (matchArr && matchArr.length) {
        return matchArr[0].substring(val.length + 1);
    }
    return null;
}

// 是否是有效的URL
export function isValidURL(string) {
    const res = string.match(/(https?|ftp|file|ssh):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|]/i);
    return res !== null;
}

export function isExternal(path) {
    return /^(https?:|mailto:|tel:)/.test(path);
}

export const formatTools = (tools) => {
    if (!(tools && tools.length)) return []
    const newTools = tools.map((n,i)=>{
        let params = []
        let properties = n.inputSchema.properties
        for(let key in properties){
            params.push({
                "name": key,
                "requiredBadge": n.inputSchema.required && n.inputSchema.required.includes(key) ? '必填' : '',
                "type": properties[key].type,
                "description": properties[key].description,
            })
        }
        return {
            ...n,
            params
        }
    })
    return newTools
}

export function rawQuery() {
    const [, query = ""] = location.href.split("?");
    if (!query) return {};

    const obj = {};
    query.split("&").forEach((item) => {
      const eq = item.indexOf("=");
      if (eq === -1) {
        // 没有等号，当成空值
        obj[item] = "";
      } else {
        const k = item.slice(0, eq);
        const v = item.slice(eq + 1); // 右边整块，含后续所有 =
        obj[k] = v;
      }
    });
    return obj;
}