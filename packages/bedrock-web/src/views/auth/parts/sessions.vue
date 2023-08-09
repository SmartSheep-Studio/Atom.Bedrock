<template>
  <div>
    <n-spin :show="requesting">
      <n-list bordered hoverable>
        <n-empty v-if="data.length <= 0" class="py-8" :description="$t('common.feedback.no-data')" />

        <n-list-item v-for="item in data">
          <n-thing :title="item.location">
            <template #header-extra>
              <n-tooltip trigger="hover" placement="right">
                <template #trigger>
                  <n-button
                    circle
                    secondary
                    size="tiny"
                    type="error"
                    :disabled="item.id === $principal.session.id"
                    @click="terminate(item)"
                  >
                    <template #icon>
                      <n-icon size="12" :component="LogOutRound" />
                    </template>
                  </n-button>
                </template>
                {{ $t("pages.auth.principal.sessions.terminate") }}
              </n-tooltip>
            </template>
            <template #description>
              <n-space class="mb-2">
                <n-tag :bordered="false" size="small" type="info" v-if="item.client">{{ item.client.name }}</n-tag>
                <n-tag :bordered="false" size="small" type="info" v-else>{{ $t("brand.nucleus.name") }}</n-tag>
              </n-space>
              <div class="flex">
                <div>{{ item.ip }}</div>
                <div class="text-gray-400 ml-1">#{{ item.id }}</div>
              </div>
            </template>
          </n-thing>
        </n-list-item>
      </n-list>
    </n-spin>

    <div class="flex justify-center mt-4">
      <n-pagination
        v-model:page="pagination.page"
        :page-count="Math.ceil((rawData.length ?? 0) / pagination.pageSize)"
        :page-slot="pagination.slot"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { usePrincipal } from "@/stores/account";
import { computed, onMounted, reactive, ref } from "vue";
import { LogOutRound } from "@vicons/material";
import { http } from "@/utils/http";
import { useMessage } from "naive-ui";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $message = useMessage();
const $principal = usePrincipal();

const rawData = ref<any[]>([]);
const data = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize;
  return (
    rawData.value
      ?.reverse()
      .filter((v) => v.type !== 2) // Remove api tokens
      .slice(start, start + pagination.pageSize) ?? []
  );
});

const requesting = ref(true);

const pagination = reactive({
  page: 1,
  pageSize: 5,
  slot: 5
});

async function fetch() {
  try {
    requesting.value = true;
    const res = await http.get("/api/auth/sessions");

    rawData.value = res.data.sessions;
    rawData.value.map((v: any) => {
      if (v.client_id != null) {
        v.client = res.data["clients"][v.client_id];
      } else {
        v.client = null;
      }

      return v;
    });
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e]));
  } finally {
    requesting.value = false;
  }
}

async function terminate(item: any) {
  try {
    requesting.value = true;
    await http.delete("/api/auth/sessions", { params: { id: item.id } });

    await Promise.all([fetch(), $principal.fetch()]);

    $message.success(t('pages.auth.principal.sessions.feedback.success'));
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e]));
  } finally {
    requesting.value = false;
  }
}

onMounted(() => {
  fetch();
});
</script>
