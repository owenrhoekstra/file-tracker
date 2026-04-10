import { base64ToUint8Array, uint8ArrayToBase64 } from './utilFunctions'

type PublicKeyCredentialRequestOptionsJSON = {
    challenge: string
    rpId?: string
    allowCredentials?: Array<{
        id: string
        type: PublicKeyCredentialType
        transports?: AuthenticatorTransport[]
    }>
    userVerification?: UserVerificationRequirement
    timeout?: number
    [key: string]: any
}

type PasskeyLoginOptions = {
    options: {
        publicKey: PublicKeyCredentialRequestOptionsJSON
    }
    sessionId: string
}

export async function passkeyLogin(email: string): Promise<void> {
    // 1. Get challenge
    const res = await fetch('/api/auth/passkey/login-challenge', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
    })

    if (!res.ok) throw new Error('Backend not available')

    const { options, sessionId }: PasskeyLoginOptions = await res.json()

    // 2. Build WebAuthn request
    const pk: PublicKeyCredentialRequestOptions = {
        challenge: base64ToUint8Array(options.publicKey.challenge),
        rpId: options.publicKey.rpId,
        userVerification: options.publicKey.userVerification,
        timeout: options.publicKey.timeout
    }

    if (options.publicKey.allowCredentials?.length) {
        pk.allowCredentials = options.publicKey.allowCredentials.map(
            (cred: { id: string; type: PublicKeyCredentialType; transports?: AuthenticatorTransport[] }) => ({
                type: cred.type,
                id: base64ToUint8Array(cred.id),
                transports: cred.transports
            })
        )
    }

    // 3. Get credential
    const credential = await navigator.credentials.get({
        publicKey: pk
    }) as PublicKeyCredential | null

    if (!credential) throw new Error('No credential returned')

    const response = credential.response as AuthenticatorAssertionResponse

    // 4. Send to backend (NOW includes sessionId)
    const authRes = await fetch('/api/auth/passkey/login-verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            email,
            sessionId, // 🔥 REQUIRED NOW
            id: uint8ArrayToBase64(credential.rawId),
            rawId: uint8ArrayToBase64(credential.rawId),
            type: credential.type,
            response: {
                authenticatorData: uint8ArrayToBase64(response.authenticatorData),
                clientDataJSON: uint8ArrayToBase64(response.clientDataJSON),
                signature: uint8ArrayToBase64(response.signature),
                userHandle: response.userHandle
                    ? uint8ArrayToBase64(response.userHandle)
                    : null
            }
        })
    })

    if (!authRes.ok) throw new Error('Authentication failed')

    // 5. Redirect
    window.location.href = '/dashboard'
}