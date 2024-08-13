<script setup>
import { RouterView } from 'vue-router'
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
</script>

<template>
  <header>
    <Toolbar>
      <template #start>
        <IconField>
          <InputIcon class="pi pi-amazon" />
          <InputText v-model="a" placeholder="Pablo" />
        </IconField>
        asdf
      </template>
    </Toolbar>
  </header>

  <RouterView />
</template>

<style scoped>
</style>
