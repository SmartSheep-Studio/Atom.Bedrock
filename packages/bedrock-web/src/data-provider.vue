<template>
  <n-layout-content class="content-placeholder w-full h-max" v-if="reverting">
    <div class="w-full h-max flex justify-center items-center">
      <n-spin show>
        <template #description>{{ $t("common.feedback.connecting") }}</template>
      </n-spin>
    </div>
  </n-layout-content>

  <n-layout-content class="content" v-else>
    <slot />
  </n-layout-content>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { usePrincipal } from "@/stores/account";
import { useEndpoint } from "@/stores/connection";
import { useDialog } from "naive-ui";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $dialog = useDialog();
const $endpoint = useEndpoint();
const $principal = usePrincipal();

const reverting = ref(true);

async function fetch() {
  reverting.value = true;

  try {
    await Promise.all([$endpoint.fetch(), $principal.fetch()]);
  } catch (e: any) {
    $dialog.error({
      closable: false,
      closeOnEsc: false,
      title: t("common.feedback.network-error"),
      content: e.toString(),
      positiveText: t("actions.retry"),
      onPositiveClick: () => {
        fetch();
      },
    });
  } finally {
    reverting.value = false;
  }
}

onMounted(() => {
  fetch();
});
</script>
