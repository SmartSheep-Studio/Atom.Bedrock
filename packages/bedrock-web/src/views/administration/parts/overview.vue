<template>
  <n-spin :show="reverting">
    <n-card title="Overview">
      <template #header-extra>
        <n-button secondary circle type="info" size="small" :loading="reverting" @click="fetch">
          <template #icon>
            <n-icon :component="RefreshRound" />
          </template>
        </n-button>
      </template>

      <n-grid item-responsive responsive="screen" :x-gap="8" :y-gap="8" :cols="4">
        <n-gi span="4 m:2 l:1">
          <n-card embedded class="w-full">
            <n-statistic label="Users" tabular-nums>
              <n-number-animation :from="0" :to="overview.resources.users" />
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi span="4 m:2 l:1">
          <n-card embedded class="w-full">
            <n-statistic label="Sessions" tabular-nums>
              <n-number-animation :from="0" :to="overview.resources.sessions" />
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi span="4 m:2 l:1">
          <n-card embedded class="w-full">
            <n-statistic label="Contacts" tabular-nums>
              <n-number-animation :from="0" :to="overview.resources.contacts" />
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi span="4 m:2 l:1">
          <n-card embedded class="w-full">
            <n-statistic label="Notifications" tabular-nums>
              <n-number-animation :from="0" :to="overview.resources.notifications" />
            </n-statistic>
          </n-card>
        </n-gi>
      </n-grid>
    </n-card>
  </n-spin>
</template>

<script setup lang="ts">
import { RefreshRound } from "@vicons/material";
import { useMessage } from "naive-ui";
import { onMounted, ref } from "vue";
import { http } from "@/utils/http";

const $message = useMessage();

const reverting = ref(false);

const overview = ref<any>({
  uptime: null,
  resources: {
    users: 0,
    sessions: 0,
    contacts: 0,
    notifications: 0,
  },
});

async function fetch() {
  try {
    reverting.value = true;
    overview.value = (await http.get("/api/metrics")).data;
  } catch (e) {
    $message.error(`Something went wrong... ${e}`);
  } finally {
    reverting.value = false;
  }
}

onMounted(() => {
  fetch();
});
</script>
