<template>
  <div class="w-full relative h-[100vh]">
    <n-layout position="absolute">
      <n-layout-header bordered class="px-8 md:px-16 h-full flex items-center">
        <n-menu mode="horizontal" v-model:value="menuKey" :options="menuOptions" />
      </n-layout-header>

      <n-layout class="w-full h-[100vh-72px]">
        <router-view />
      </n-layout>
    </n-layout>
  </div>
</template>

<script lang="ts" setup>
import { DashboardRound } from "@vicons/material";
import { RouterLink, useRoute } from "vue-router";
import { NIcon } from "naive-ui";
import { h } from "vue";

const $route = useRoute();

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

const menuKey = ref($route.name);
const menuOptions: Ref<MenuOption[]> = computed(() => [
  {
    label: () => h(RouterLink, { to: { name: "administration.dashboard" } }, { default: () => "Dashboard" }),
    icon: renderIcon(DashboardRound),
    key: "administration.dashboard",
  },
]);
</script>
