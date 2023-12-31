<template>
  <div>
    <div class="container pt-10">
      <div class="md:px-6 lg:px-16 xl:px-20">
        <div>
          <div class="text-xl font-bold">Notifications</div>
          <div class="text-md">
            You have <b class="font-mono">{{ $principal.account?.notification_count }}</b> unread messages.
          </div>
        </div>

        <div class="pt-4">
          <n-space>
            <n-checkbox label="Only unread" v-model:checked="options.onlyUnread" />
          </n-space>
        </div>

        <div class="pt-6 mx-[-8px]">
          <n-list bordered hoverable v-if="data.length > 0">
            <n-list-item v-for="item in data">
              <n-thing :title="item.title" content-style="margin-top: 10px;">
                <template #header-extra>
                  <n-space>
                    <n-button
                      circle
                      quaternary
                      type="info"
                      size="tiny"
                      :disabled="item.read_at != null"
                      @click="read(item)"
                    >
                      <template #icon>
                        <n-icon :component="CheckRound" />
                      </template>
                    </n-button>
                  </n-space>
                </template>

                <template #description>
                  <n-space vertical>
                    <div class="text-xs text-gray-600 mt-[-4px]">{{ item.description }}</div>
                    <n-space size="small" class="ml-[-2px]">
                      <n-tag type="error" size="small" :bordered="false" v-if="item.read_at == null"> Unread</n-tag>
                      <n-tag type="warning" size="small" :bordered="false" class="capitalize">
                        {{ item.level }}
                      </n-tag>
                      <n-tag type="info" size="small" :bordered="false">
                        {{ new Date(item.created_at).toLocaleString() }}
                      </n-tag>
                    </n-space>
                  </n-space>
                </template>

                <div class="mb-[-6px]">
                  <vue-markdown :source="item.content" />
                </div>
              </n-thing>
            </n-list-item>
          </n-list>

          <n-list bordered v-else>
            <n-list-item>
              <n-empty description="No notification for you." />
            </n-list-item>
          </n-list>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { CheckRound } from "@vicons/material";
import { usePrincipal } from "@/stores/account";
import { reactive, ref, watch } from "vue";
import { useMessage } from "naive-ui";
import { http } from "@/utils/http";
import { useI18n } from "vue-i18n";
import VueMarkdown from "vue-markdown-render";

const { t } = useI18n();

const $principal = usePrincipal();
const $message = useMessage();

const reverting = ref(false);

const data = ref<any[]>([]);

const options = reactive({
  onlyUnread: true
});

async function fetch() {
  try {
    reverting.value = true;

    const res = await http.get("/api/users/self/notifications", {
      params: {
        only_unread: options.onlyUnread ? "yes" : "no"
      }
    });

    data.value = res.data;
  } catch (e) {
    $message.error(t("common.feedback.unknown-error", [e]));
  } finally {
    reverting.value = false;
  }
}

async function read(item: any) {
  try {
    reverting.value = true;

    await http.post(`/api/notifications/${item.id}/read`);
    await fetch();
  } catch (e) {
    $message.error(t("common.feedback.unknown-error", [e]));
  } finally {
    reverting.value = false;
  }
}

watch(
  options,
  () => {
    fetch();
  },
  { immediate: true, deep: true }
);
</script>
