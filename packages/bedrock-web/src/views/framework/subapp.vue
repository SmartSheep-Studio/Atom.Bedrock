<template>
  <wujie-vue
    sync
    id="subapp-container"
    height="100vh"
    width="100%"
    name="q"
    :url="page.to"
    :prefix="pagePrefix"
  />
</template>

<script lang="ts" setup>
import { computed } from "vue";
import { useRoute } from "vue-router";
import * as Base64 from "js-base64";
import { useEndpoint } from "@/stores/connection";

const $route = useRoute();
const $endpoint = useEndpoint();

const page = computed(() => $endpoint.pages.filter((v) => v.name === $route.params.id)[0]);
const host = computed(() => Base64.decode($route.params.host as string));

const pagePrefix = computed(() => {
  return { q: `/srv/subapps/${page.value.name}` };
});

function updateState() {
  window.__BEDROCK = {
    isUnderShadow: true,
    layout: {
      isDisplayHeader: true,
      contentHeight: "calc(100vh - 72px)"
    }
  };
}

updateState();
</script>