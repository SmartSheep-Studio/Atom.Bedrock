<template>
  <div class="pt-8 px-12">
    <!-- Header -->
    <div class="px-4">
      <div class="text-lg font-bold">Hello, {{ $principal.account.nickname }}!</div>
      <div class="text-md">Welcome to Atom administration area!</div>
    </div>
    
    <overview class="mt-6" />
  </div>
</template>

<script lang="ts" setup>
import { usePrincipal } from "@/stores/account";
import Overview from "@/views/administration/parts/overview.vue";

const $principal = usePrincipal();
</script>
