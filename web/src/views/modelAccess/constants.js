export const LLM = 'llm'

export const MODEL_TYPE_OBJ = {
    [LLM]: 'LLM',
    rerank: 'Rerank',
    embedding: 'Embedding',
}

export const MODEL_TYPE = Object.keys(MODEL_TYPE_OBJ).map(key => ({key, name: MODEL_TYPE_OBJ[key]}))

export const YUAN_JING = 'YuanJing'
export const OPENAI_API = 'OpenAI-API-compatible'
export const OLLAMA = 'Ollama'
export const QWEN = 'Qwen'
export const HUOSHAN = 'Huoshan'
export const PROVIDER_OBJ = {
    [OPENAI_API]: 'OpenAI-API-compatible',
    [YUAN_JING]: '联通元景',
    [OLLAMA]: 'Ollama',
    [QWEN]: '通义千问',
    // [HUOSHAN]: '火山引擎',
}

export const PROVIDER_IMG_OBJ = {
    [OPENAI_API]: require('@/assets/imgs/openAI.png'),
    [YUAN_JING]: require('@/assets/imgs/yuanjing.png'),
    [OLLAMA]: require('@/assets/imgs/ollama.png'),
    [QWEN]: require('@/assets/imgs/qwen.png'),
    [HUOSHAN]: require('@/assets/imgs/volcano.png'),
}

export const PROVIDER_TYPE = Object.keys(PROVIDER_OBJ)
    .map(key => ({key, name: PROVIDER_OBJ[key], children: MODEL_TYPE}))

export const DEFAULT_CALLING = 'noSupport'
export const FUNC_CALLING = [
    {key: 'noSupport', name: '不支持'},
    {key: 'toolCall', name: 'Tool call'},
    {key: 'functionCall', name: 'Function call'},
]

export const TYPE_OBJ = {
    apiKey: {
        [YUAN_JING]: 'sk-abc********************xyz',
        [OPENAI_API]: 'sk_7e4*************4s-BpI1l',
        [OLLAMA]: '',
        [QWEN]: 'sk-b************c70d',
        [HUOSHAN]: 'd8008ac0-****-****-****-**************'
    },
    inferUrl: {
        [YUAN_JING]: 'https://maas.ai-yuanjing.com/openapi/compatible-mode/v1',
        [OPENAI_API]: 'https://api.siliconflow.cn/v1',
        [OLLAMA]: 'https://192.168.21.100:11434',
        [QWEN]: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
        [HUOSHAN]: 'https://ark.cn-beijing.volces.com/api/v3'
    },
}