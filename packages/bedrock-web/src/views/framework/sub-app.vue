<template>
  <wujie-vue
    sync
    id="sub-app-container"
    height="calc(100vh - 72px)"
    width="100%"
    name="p"
    :url="nav.to"
  />
</template>

<script lang="ts" setup>
import { computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import * as Base64 from "js-base64";
import { useEndpoint } from "@/stores/connection";

const $route = useRoute();
const $endpoint = useEndpoint();

const nav = computed(() => $endpoint.nav.filter((v) => v.name === $route.params.id)[0]);
const host = computed(() => Base64.decode($route.params.host as string));
const url = computed(() => `//${host.value}${$route.params.path}`);
</script>
