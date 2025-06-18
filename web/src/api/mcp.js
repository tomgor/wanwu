import request from "@/utils/request";

/*----元景------*/
//mcp列表
export const getList = (data)=>{
    return request({
        url: '/use/model/api/v1/mcp/list',
        method: 'get',
        params: data
    })
};

export const getDetail = (data)=>{
    return request({
        url: '/use/model/api/v1/mcp/info',
        method: 'get',
        params: data
    })
};

export const setDelete = (data)=>{
    return request({
        url: '/use/model/api/v1/mcp/delete',
        method: 'delete',
        data: data
    })
};

export const setUpdate = (data)=>{
    return request({
        url: '/use/model/api/v1/mcp/update',
        method: 'put',
        data: data
    })
};

export const getTools = (data)=>{
    return request({
        url: '/use/model/api/v1/mcp/getTools',
        method: 'get',
        params: data
    })
};

export const setCreate = (data)=>{
    return request({
        url: '/use/model/api/v1/mcp/create',
        method: 'post',
        data: data
    })
};

export const cubmCreate = (data)=>{
    return request({
        url: '/use/model/api/v1/assistant/mcp/create',
        method: 'post',
        data: data
    })
};

export const cubmDelete = (data)=>{
    return request({
        url: '/use/model/api/v1/assistant/mcp/delete',
        method: 'delete',
        data: data
    })
};
/*---第三方mcp广场---*/
export const getMcpCategoryList = (data)=>{
    return request({
        url: '/use/model/api/v1/mcpSquare/categories',
        method: 'get',
        params: data
    })
};
export const getPublicMcpList = (data)=>{
    return request({
        url: '/use/model/api/v1/mcpSquare/list',
        method: 'get',
        params: data
    })
};
export const getPublicMcpInfo = (data)=>{
    return request({
        url: '/use/model/api/v1/mcpSquare/info',
        method: 'get',
        params: data
    })
};

export const getMarkDownContent = (path)=>{
    return request({
        //url: '/file/api/' + path,
        url: 'https://obs-nmhhht6.cucloud.cn/' + path,
        method: 'get',
        responseType: "blob"
    })
};
export const getRecommendsList = (data)=>{
    return request({
        url: '/use/model/api/v1/mcpSquare/recommends',
        method: 'get',
        params: data
    })
};