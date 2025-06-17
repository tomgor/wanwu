import { PERMS } from "@/router/permission"
import { i18n } from '@/lang'

export const menuList = [
    {
        name: i18n.t('menu.modelAccess'),
        key: 'modelAccess',
        img: require('@/assets/imgs/model.png'),
        imgActive: require('@/assets/imgs/model_active.png'),
        path: '/modelAccess',
        perm: PERMS.MODEL,
    },
    {
        name: i18n.t('menu.knowledge'),
        key: 'knowledge',
        img: require('@/assets/imgs/knowledge.png'),
        imgActive: require('@/assets/imgs/knowledge_active.png'),
        path: '/knowledge',
        perm: PERMS.KNOWLEDGE,
    },
    {
        name: i18n.t('menu.mcp'),
        key: 'mcpManage',
        img: require('@/assets/imgs/knowledge.png'),
        imgActive: require('@/assets/imgs/knowledge_active.png'),
        path: '/mcp',
        perm: PERMS.MCP,
    },
    {
        key: 'line'
    },
    {
        name: i18n.t('menu.app.rag'),
        key: 'rag',
        img: require('@/assets/imgs/task.png'),
        imgActive: require('@/assets/imgs/task_active.png'),
        path: '/appSpace/rag',
        perm: PERMS.RAG
    },
    {
        name: i18n.t('menu.app.workflow'),
        key: 'workflow',
        img: require('@/assets/imgs/task.png'),
        imgActive: require('@/assets/imgs/task_active.png'),
        path: '/appSpace/workflow',
        perm: PERMS.WORKFLOW
    },
    {
        name: i18n.t('menu.app.agent'),
        key: 'agent',
        img: require('@/assets/imgs/task.png'),
        imgActive: require('@/assets/imgs/task_active.png'),
        path: '/appSpace/agent',
        perm: PERMS.AGENT
    },
    {
        key: 'line'
    },
    {
        name: i18n.t('menu.explore'),
        key: 'explore',
        img: require('@/assets/imgs/explore.png'),
        imgActive: require('@/assets/imgs/explore_active.png'),
        path: '/explore',
        perm: PERMS.EXPLORE
    },
]