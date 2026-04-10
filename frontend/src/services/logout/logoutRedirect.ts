export async function apiFetch(url: string, options: RequestInit = {}) {
    const res = await fetch(url, {
        ...options,
        credentials: 'include',
    })

    if (res.status === 401) {
        window.location.href = '/login'
        return
    }

    return res
}