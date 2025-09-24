import request from "@/utils/request"
const BASE_URL = '/use/model/api/v1'

/*----自定义工具------*/
export const getCustomList = (data)=>{
    return request({
        url: `${BASE_URL}/tool/custom/list`,
        method: 'get',
        params: data
    })
};

export const getCustom = (data)=>{
    return request({
        url: `${BASE_URL}/tool/custom`,
        method: 'get',
        params: data
    })
};

export const editCustom = (data)=>{
    return request({
        url: `${BASE_URL}/tool/custom`,
        method: 'put',
        data
    })
};

export const addCustom = (data)=>{
    return request({
        url: `${BASE_URL}/tool/custom`,
        method: 'post',
        data
    })
};

export const deleteCustom = (data)=>{
    return request({
        url: `${BASE_URL}/tool/custom`,
        method: 'delete',
        data
    })
};

export const getSchema = (data)=>{
    return request({
        url: `${BASE_URL}/tool/custom/schema`,
        method: 'post',
        data
    })
};


/*---mcp列表---*/
export const getList = (data)=>{
    return request({
        url: `${BASE_URL}/mcp/list`,
        method: 'get',
        params: data
    })
};

export const getDetail = (data)=>{
    return request({
        url: `${BASE_URL}/mcp`,
        method: 'get',
        params: data
    })
};

export const setDelete = (data)=>{
    return request({
        url: `${BASE_URL}/mcp`,
        method: 'delete',
        data: data
    })
};

export const getTools = (data)=>{
    return request({
        url: `${BASE_URL}/mcp/tool/list`,
        method: 'get',
        params: data
    })
};

export const setCreate = (data)=>{
    return request({
        url: `${BASE_URL}/mcp`,
        method: 'post',
        data: data
    })
};

/*---第三方工具广场---*/
export const getPublicMcpList = (data)=>{
    return request({
        url: `${BASE_URL}/mcp/square/list`,
        method: 'get',
        params: data
    })
};
export const getPublicMcpInfo = (data)=>{
    return request({
        url: `${BASE_URL}/mcp/square`,
        method: 'get',
        params: data
    })
};
export const getRecommendsList = (data)=>{
    return request({
        url: `${BASE_URL}/mcp/square/recommend`,
        method: 'get',
        params: data
    })
};

/*----内置工具------*/
export const getBuiltInList = (data)=>{
    return request({
        url: `${BASE_URL}/tool/square/list`,
        method: 'get',
        params: data
    })
};

export const getToolDetail = (data)=>{
    return request({
        url: `${BASE_URL}/tool/square`,
        method: 'get',
        params: data
    })
};

export const changeApiKey = (data)=>{
    return request({
        url: `${BASE_URL}/tool/builtin`,
        method: 'post',
        data
    })
};
