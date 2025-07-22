import service from "@/utils/request";
//查询知识库关键词列表
export const getKeyWord = (data)=>{
    return service({
        url: '/konwledgeServe/v1/knowledge/keywords',
        method: 'get',
        params: data
    })
};

//新增知识库关键词列表
export const addKeyWord = (data)=>{
    return service({
        url: '/konwledgeServe/v1/knowledge/keywords',
        method: 'post',
        data
    })
};
//编辑知识库关键词列表
export const editKeyWord = (data)=>{
    return service({
        url: '/konwledgeServe/v1/knowledge/keywords',
        method: 'put',
        data
    })
};
//删除知识库关键词列表
export const delKeyWord = (data)=>{
    return service({
        url: '/konwledgeServe/v1/knowledge/keywords',
        method: 'delete',
        data
    })
};
//知识库关键词详情
export const keyWordDetail = (data)=>{
    return service({
        url: '/konwledgeServe/v1/knowledge/keywords/detail',
        method: 'get',
        params: data
    })
};