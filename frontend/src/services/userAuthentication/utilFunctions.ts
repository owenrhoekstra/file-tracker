export function base64ToUint8Array(
    input: string | ArrayBuffer | ArrayBufferView
): Uint8Array {
    if (input instanceof Uint8Array) return input
    if (input instanceof ArrayBuffer) return new Uint8Array(input)
    if (ArrayBuffer.isView(input)) return new Uint8Array(input.buffer)

    const pad = '='.repeat((4 - (input.length % 4)) % 4)
    const b64 = (input + pad)
        .replace(/-/g, '+')
        .replace(/_/g, '/')

    const binary = atob(b64)
    return Uint8Array.from(binary, c => c.charCodeAt(0))
}

export function uint8ArrayToBase64(buffer: ArrayBuffer | Uint8Array): string {
    const bytes = buffer instanceof Uint8Array ? buffer : new Uint8Array(buffer)
    return btoa(String.fromCharCode(...bytes))
}