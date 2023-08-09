<template>
  <div class="gatekeeper-container h-full w-full">
    <slot />
  </div>
</template>

<script lang="ts" setup>
import { onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { usePrincipal } from "@/stores/account";
import { keepGate } from "@/utils/gatekeeper";

const $route = useRoute();
const $router = useRouter();
const $principal = usePrincipal();

onMounted(() => {
  watch($principal, () => {
    if (!keepGate($route)) {
      $router.push({
        name: "auth.sign-in",
        query: { redirect_uri: $route.fullPath },
        state: { message: `You need sign in before access that page!` },
      });
    }
  });

  watch(
    $route,
    (v) => {
      if (!keepGate(v)) {
        $router.push({
          name: "auth.sign-in",
          query: { redirect_uri: $route.fullPath },
          state: { message: `You need sign in before access that page!` },
        });
      }
    },
    { deep: true }
  );
});
</script>