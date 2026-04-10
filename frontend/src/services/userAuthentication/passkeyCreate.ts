import { base64ToUint8Array, uint8ArrayToBase64 } from './utilFunctions'

type PublicKeyCredentialCreationOptionsJSON = {
    rp: PublicKeyCredentialRpEntity
    user: {
        id: string
        name: string
        displayName: string
    }
    challenge: string
    pubKeyCredParams: PublicKeyCredentialParameters[]
    excludeCredentials?: Array<{
        id: string
        type: PublicKeyCredentialType
        transports?: AuthenticatorTransport[]
    }>
    [key: string]: any
}

type PasskeyCreateOptions = {
    options: PublicKeyCredentialCreationOptionsJSON
    sessionId: string
}

export async function passkeyCreate(email: string): Promise<void> {
    // 1. Get registration challenge + sessionId
    const res = await fetch('/api/auth/passkey/register-challenge', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
    })

    if (!res.ok) throw new Error('Backend not available')

    const { options, sessionId }: PasskeyCreateOptions = await res.json()

    // 2. Prepare WebAuthn options
    const pk: PublicKeyCredentialCreationOptions = {
        rp: options.publicKey.rp,
        user: {
            ...options.publicKey.user,
            id: base64ToUint8Array(options.publicKey.user.id)
        },
        challenge: base64ToUint8Array(options.publicKey.challenge).buffer as ArrayBuffer,
        pubKeyCredParams: options.publicKey.pubKeyCredParams,
    }

    if (options.publicKey.excludeCredentials) {
        pk.excludeCredentials = options.publicKey.excludeCredentials.map(
            (cred: PublicKeyCredentialDescriptor) => ({
                ...cred,
                id: new Uint8Array(base64ToUint8Array(cred.id)).buffer
            })
        )
    }

    // 3. Create credential
    const credential = await navigator.credentials.create({
        publicKey: pk
    }) as PublicKeyCredential | null

    if (!credential) throw new Error('No credential created')

    const response = credential.response as AuthenticatorAttestationResponse

    // 4. Send to backend WITH sessionId
    const regRes = await fetch('/api/auth/passkey/register-verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            email,
            sessionId, // 🔥 REQUIRED NOW
            id: uint8ArrayToBase64(credential.rawId),
            rawId: uint8ArrayToBase64(credential.rawId),
            type: credential.type,
            response: {
                clientDataJSON: uint8ArrayToBase64(response.clientDataJSON),
                attestationObject: uint8ArrayToBase64(response.attestationObject)
            }
        })
    })

    if (!regRes.ok) throw new Error('Registration failed')

    // 5. Redirect
    window.location.href = '/dashboard'
}