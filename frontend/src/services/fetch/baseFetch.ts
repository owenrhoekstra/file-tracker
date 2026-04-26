const BASE_URL = import.meta.env.VITE_API_URL

if (!BASE_URL) throw new Error('VITE_API_URL environment variable is not set')
if (!BASE_URL.startsWith('https://')) throw new Error('VITE_API_URL must be HTTPS')

export async function baseFetch(url: string, options: RequestInit = {}): Promise<Response> {
    const fullUrl = new URL(url, BASE_URL).toString()
    return fetch(fullUrl, { ...options, credentials: 'include' })
}

export { BASE_URL }