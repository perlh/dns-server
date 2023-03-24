import {createRouter, createWebHashHistory} from 'vue-router'

const routes = [
    {
        path: '/domain/forward',
        component: () => import("@/pages/domain/Forward"),
        name: "DomainForward"
    },
    {
        path: '/domain/reverse',
        component: () => import("@/pages/domain/Reverse"),
        name: "DomainReverse"
    },
    {
        path: '/settings/person_center',
        component: () => import("@/pages/settings/PersonCenter"),
        name: "PersonCenter"
    },
    {
        path: '/settings/other',
        component: () => import("@/pages/settings/Other"),
        name: "Other"
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export default router