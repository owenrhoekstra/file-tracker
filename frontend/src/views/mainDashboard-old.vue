<script setup lang="ts">
import { ref } from 'vue'
import Button from 'primevue/button'
import { apiFetch } from '../services/fetch/statusCodeChecks.ts'
import { requestElevation } from '../services/elevation/elevate.ts'
import mainMenuBar from '../components/mainMenuBar.vue'
import { useRouter } from 'vue-router'

const router = useRouter();

function test () {
  apiFetch('/api/protected/test', {})
}

async function testActionElevation() {
  try {
    console.log('Starting action elevation...')
    const ok = await requestElevation('action')
    console.log('Result:', ok)
  } catch (e) {
    console.error('Elevation error:', e)
  }
}

async function testViewElevation() {
  const ok = await requestElevation('view')
  console.log(ok ? 'View elevation granted!' : 'View elevation failed')
}

const ocrResult = ref<Record<string, unknown> | null>(null)
const ocrError = ref<string | null>(null)

async function ocrCall() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.capture = 'environment'

  input.onchange = async () => {
    try {
      const file = input.files?.[0]
      if (!file) return

      ocrResult.value = null
      ocrError.value = null

      const bitmap = await createImageBitmap(file)
      const canvas = document.createElement('canvas')
      canvas.width = bitmap.width
      canvas.height = bitmap.height
      const ctx = canvas.getContext('2d')!
      ctx.drawImage(bitmap, 0, 0)

      const blob = await new Promise<Blob>((resolve, reject) =>
          canvas.toBlob(
              b => b ? resolve(b) : reject(new Error('toBlob failed')),
              'image/webp',
              0.92
          )
      )

      const form = new FormData()
      form.append('image', blob, 'capture.webp')

      const resp = await apiFetch('/api/protected/ocr', {
        method: 'POST',
        body: form,
      })

      if (!resp) return

      if (!resp.ok) {
        ocrError.value = `Error ${resp.status}: ${await resp.text()}`
        return
      }

      ocrResult.value = await resp.json()
    } catch (e) {
      console.error('ocrCall error:', e)
      ocrError.value = `Unexpected error: ${e}`
    }
  }
  input.click()
}
</script>

<template>
  <mainMenuBar />
  <p>Auth complete, dashboard loaded</p>
  <div class="grid grid-cols-1 gap-4 max-w-sm w-full mx-auto px-4 py-3">
    <Button label="Expired Test" @click="test()" />
    <Button label="Test Action Elevation" @click="testActionElevation()" />
    <Button label="Test View Elevation" @click="testViewElevation()" />
    <Button label="OCR Call" @click="ocrCall" />
    <Button label="Print View" @click="router.push('/print')" />

    <pre v-if="ocrError" style="color: red; white-space: pre-wrap;">{{ ocrError }}</pre>
    <pre v-if="ocrResult" style="white-space: pre-wrap;">{{ JSON.stringify(ocrResult, null, 2) }}</pre>
  </div>
</template>