import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        component: () => import('../views/userAuthentication.vue'),
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router