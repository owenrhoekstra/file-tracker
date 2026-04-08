<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import bwipjs from 'bwip-js'

interface Props {
  labelWidthMm?: number
  labelHeightMm?: number
  testValue?: string
}

const props = withDefaults(defineProps<Props>(), {
  labelWidthMm: 12,
  labelHeightMm: 12,
  testValue: '123456789'
})

const canvas = ref<HTMLCanvasElement | null>(null)
const error = ref<string | null>(null)

const RENDER_DPI = 300
const MM_TO_INCH = 25.4
const PX_PER_MM = RENDER_DPI / MM_TO_INCH

const canvasWidthPx = computed(() => Math.round(props.labelWidthMm * PX_PER_MM))
const canvasHeightPx = computed(() => Math.round(props.labelHeightMm * PX_PER_MM))

function renderDataMatrix() {
  if (!canvas.value) return
  error.value = null

  try {
    bwipjs.toCanvas(canvas.value, {
      bcid: 'datamatrix',
      text: props.testValue,
      scale: 2,
    })
  } catch (e: unknown) {
    error.value = `Failed to render barcode: ${e instanceof Error ? e.message : String(e)}`
  }
}

onMounted(renderDataMatrix)
watch(() => [props.labelWidthMm, props.labelHeightMm, props.testValue], renderDataMatrix)

function printLabel() {
  if (!canvas.value) return

  const dataUrl = canvas.value.toDataURL('image/png')
  const win = window.open('', '_blank')
  if (!win) {
    error.value = 'Popup blocked — allow popups for this site.'
    return
  }

  win.document.title = 'Label Print'
  win.document.documentElement.lang = 'en'
  win.document.body.innerHTML = `<img src="${dataUrl}" alt="Label Barcode" />`

  const style = win.document.createElement('style')
  style.textContent = `
    @page {
      size: ${props.labelWidthMm}mm ${props.labelHeightMm}mm;
      margin: 0;
    }
    * {
      margin: 0;
      padding: 2;
      box-sizing: border-box;
    }
    body {
      width: ${props.labelWidthMm}mm;
      height: ${props.labelHeightMm}mm;
      overflow: hidden;
    }
    img {
      width: ${props.labelWidthMm}mm;
      height: ${props.labelHeightMm}mm;
      display: block;
      image-rendering: pixelated;
    }
  `
  win.document.head.appendChild(style)

  const script = win.document.createElement('script')
  script.textContent = `
    const runPrint = () => {
      window.print();
      window.onafterprint = () => window.close();
    };
    if (document.readyState === 'complete') {
      runPrint();
    } else {
      window.addEventListener('load', runPrint);
    }
  `
  win.document.body.appendChild(script)
}

</script>

<template>
  <div class="min-h-screen bg-zinc-950 text-zinc-100 font-mono flex items-center justify-center p-8">
    <div class="w-full max-w-md space-y-8">

      <!-- Header -->
      <div class="border-b border-zinc-700 pb-4">
        <p class="text-xs text-zinc-500 uppercase tracking-widest">Label Print Test</p>
        <h1 class="text-2xl font-bold text-zinc-100 mt-1">Data Matrix Preview</h1>
      </div>

      <!-- Meta -->
      <div class="grid grid-cols-2 gap-4 text-sm">
        <div class="bg-zinc-900 border border-zinc-800 rounded p-3">
          <p class="text-zinc-500 text-xs uppercase tracking-wider mb-1">Size</p>
          <p class="text-zinc-100">{{ labelWidthMm }} × {{ labelHeightMm }} mm</p>
        </div>
        <div class="bg-zinc-900 border border-zinc-800 rounded p-3">
          <p class="text-zinc-500 text-xs uppercase tracking-wider mb-1">Test Value</p>
          <p class="text-zinc-100 truncate">{{ testValue }}</p>
        </div>
      </div>

      <!-- Canvas Preview -->
      <div class="bg-white rounded-lg flex items-center justify-center p-4">
        <canvas
            ref="canvas"
            :width="canvasWidthPx"
            :height="canvasHeightPx"
            class="block"
            style="image-rendering: pixelated;"
        />
      </div>

      <!-- Error -->
      <p v-if="error" class="text-red-400 text-sm bg-red-950 border border-red-800 rounded p-3">
        {{ error }}
      </p>

      <!-- Print Button -->
      <button
          @click="printLabel"
          class="w-full bg-zinc-100 hover:bg-white text-zinc-950 font-bold py-3 px-6 rounded transition-colors duration-150 uppercase tracking-widest text-sm"
      >
        Print Label
      </button>

      <p class="text-zinc-600 text-xs text-center">
        Verify scan output matches <span class="text-zinc-400">{{ testValue }}</span> before deploying labels.
      </p>

    </div>
  </div>
</template>