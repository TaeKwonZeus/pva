<script setup>
import { RouterView } from 'vue-router'
import { ref } from "vue";

const items = ref([
  {
    label: "bruh",
    icon: "pi pi-home"
  },
  {
    label: "bruh",
    icon: "pi pi-home"
  }
]);
</script>

<template>
  <header>
    <Menubar :model="items" />
  </header>

  <RouterView />
</template>

