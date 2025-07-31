import request from "@/utils/request";
//编辑敏感词表
export const editSensitive = (data)=>{
    return request({
        url: '/user/api/v1/safe/sensitive/table',
        method: 'put',
        data
    })
};
//创建敏感词表
export const createSensitive = (data)=>{
    return request({
        url: '/user/api/v1/safe/sensitive/table',
        method: 'post',
        data
    })
};
//删除敏感词表
export const delSensitive = (data)=>{
    return request({
        url: '/user/api/v1/safe/sensitive/table',
        method: 'delete',
        data
    })
};
//查看敏感词表列表
export const getSensitiveList = ()=>{
    return request({
        url: '/user/api/v1/safe/sensitive/table/list',
        method: 'get',
    })
};
//编辑回复设置
export const setReply = (data)=>{
    return request({
        url: '/user/api/v1/safe/sensitive/table/reply',
        method: 'put',
        data
    })
};
//获取敏感词表下拉列表
export const sensitiveSelect = ()=>{
    return request({
        url: '/user/api/v1/safe/sensitive/table/select',
        method: 'get',
    })
};
//删除敏感词
export const delSensitiveWord = (data)=>{
    return request({
        url: '/user/api/v1/safe/sensitive/word',
        method: 'delete',
        data
    })
};
//查询词表数据列表
export const getSensitiveWord = (data)=>{
    return request({
        url: '/user/api/v1/safe/sensitive/word/list',
        method: 'get',
        params:data
    })
};
//上传敏感词
export const uploadSensitiveWord = (data)=>{
    return request({
        url: '/user/api/v1/safe/sensitive/word',
        method: 'post',
        data
    })
};