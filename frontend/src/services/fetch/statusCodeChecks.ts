import router from "../../router/index.ts"
import { baseFetch } from './baseFetch.ts'
import { requestElevation } from '../elevation/elevate.ts'

export async function apiFetch(url: string, options: RequestInit = {}) {
    const res = await baseFetch(url, options)
    console.log('apiFetch response status:', res.status, url)

    if (res.status === 401) {
        await router.push('/?logout=true')
        return
    }

    if (res.status === 403) {
        console.log('403 received')
        const elevationType = res.headers.get('X-Require-Elevation')
        console.log('elevation header:', elevationType)
        if (elevationType === 'action' || elevationType === 'view') {
            const ok = await requestElevation(elevationType)
            if (ok) return baseFetch(url, options)
        }
        return
    }

    return res
}