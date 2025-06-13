import { PERMS } from "@/router/permission"
import { i18n } from '@/lang'

export const menuList = [
    {
        name: i18n.t('menu.explore'),
        key: 'explore',
        img: require('@/assets/imgs/explore.png'),
        imgActive: require('@/assets/imgs/explore_active.png'),
        path: '/explore',
        perm: PERMS.EXPLORE
    },
    {
        name: i18n.t('menu.workspace'),
        key: 'workspace',
        img: require('@/assets/imgs/workspace.png'),
        imgActive: require('@/assets/imgs/workspace_active.png'),
        perm: PERMS.WORKSPACE,
        children: [
            {
                name: i18n.t('menu.app.index'),
                img: require('@/assets/imgs/task.png'),
                imgActive: require('@/assets/imgs/task_active.png'),
                index: 'workspace-1',
                perm: PERMS.WORKSPACE_APP,
                children: [
                    {
                        name: i18n.t('menu.app.all'),
                        routeName: 'all',
                        path: '/appSpace/all',
                        index: 'workspace-1-1',
                    },
                    {
                        name: i18n.t('menu.app.agent'),
                        routeName: 'agent',
                        path: '/appSpace/agent',
                        index: 'workspace-1-2',
                    },
                    {
                        name: i18n.t('menu.app.rag'),
                        routeName: 'rag',
                        path: '/appSpace/rag',
                        index: 'workspace-1-3',
                    },
                    {
                        name: i18n.t('menu.app.workflow'),
                        routeName: 'workflow',
                        path: '/appSpace/workflow',
                        index: 'workspace-1-4',
                    }
                ]
            },
            {
                name: i18n.t('menu.knowledge'),
                routeName: 'knowledge',
                path: '/knowledge',
                img: require('@/assets/imgs/knowledge.png'),
                imgActive: require('@/assets/imgs/knowledge_active.png'),
                index: 'workspace-2',
                perm: PERMS.WORKSPACE_KNOWLEDGE,
            },
        ]
    },
    {
        name: i18n.t('menu.account'),
        key: 'account',
        img: require('@/assets/imgs/account.png'),
        imgActive: require('@/assets/imgs/account_active.png'),
        children: [
            {
                name: i18n.t('menu.userInfo'),
                routeName: 'userInfo',
                path: '/userInfo',
                img: require('@/assets/imgs/userInfo.png'),
                imgActive: require('@/assets/imgs/userInfo_active.png'),
                index: 'account-1',
            },
            {
                name: i18n.t('menu.org'),
                routeName: 'org',
                path: '/permission',
                index: 'account-2',
                img: require('@/assets/imgs/org.png'),
                imgActive: require('@/assets/imgs/org_active.png'),
                perm: PERMS.PERMISSION
            },
            {
                name: i18n.t('menu.modelAccess'),
                routeName: 'modelAccess',
                path: '/modelAccess',
                img: require('@/assets/imgs/model.png'),
                imgActive: require('@/assets/imgs/model_active.png'),
                index: 'account-3',
            },
        ]
    }
]