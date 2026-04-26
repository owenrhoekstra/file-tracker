<script setup lang="ts">
import Menubar from 'primevue/menubar';
import { ref } from "vue";
import 'primeicons/primeicons.css'

const items = ref([
  {
    label: 'Options',
    icon: 'pi pi-bars',
    items: [
      {
        label: 'New',
        icon: 'pi pi-fw pi-plus',
      },
    ]
  },
  {
    label: 'Search',
    icon: 'pi pi-search',
  },

])
</script>

<template>
  <div class="card">
    <Menubar :model="items" class="bg-black text-white" />
  </div>
</template>

<style scoped>

</style>