<template>
  <div class="pt-8 px-12">
    <!-- Header -->
    <div class="px-4">
      <div class="text-lg font-bold">Users</div>
      <div class="text-md">Manage Atom IDs on this instance.</div>
    </div>

    <!-- Data Table -->
    <div class="mt-6">
      <n-spin :show="reverting">
        <n-card>
          <n-card size="small" embedded>
            <div class="flex justify-between items-center pr-2">
              <n-space>
                <n-button type="info" :disabled="selection.length <= 0" @click="popups.notifications = true">
                  <template #icon>
                    <n-icon :component="SendRound" />
                  </template>
                  Send Notifications
                </n-button>
                <n-button type="error" :disabled="selection.length <= 0" @click="popups.locks = true">
                  <template #icon>
                    <n-icon :component="BlockRound" />
                  </template>
                  Ban Account
                </n-button>
              </n-space>

              <div class="flex gap-[5px] font-mono">
                <span>Selected</span>
                <span class="font-bold">{{ selection.length }}</span>
                <span>user(s)</span>
              </div>
            </div>
          </n-card>

          <n-data-table
            class="mt-4 h-[480px]"
            flex-height
            :row-key="(v: any) => v.id"
            :columns="tableOptions.columns"
            :data="users"
            v-model:checked-row-keys="selection"
          />
        </n-card>
      </n-spin>
    </div>

    <!-- Send Notifications -->
    <n-drawer v-model:show="popups.notifications" :default-width="520" placement="right" resizable>
      <n-drawer-content title="Send Notifications">
        <send-notifications :selection="selection" @close="popups.notifications = false" />
      </n-drawer-content>
    </n-drawer>

    <!-- Ban Users -->
    <n-drawer v-model:show="popups.locks" :default-width="520" placement="right" resizable>
      <n-drawer-content title="Ban Account">
        <ban-account :selection="selection" @close="popups.locks = false" />
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script lang="ts" setup>
import { CheckRound, CloseRound, BlockRound, SendRound } from "@vicons/material";
import { h, onMounted, reactive, ref } from "vue";
import { NCode, NIcon, NTime, useMessage } from "naive-ui";
import { http } from "@/utils/http";
import SendNotifications from "@/views/administration/parts/send-notifications.vue";
import BanAccount from "@/views/administration/parts/ban-account.vue";
import hljs from "highlight.js/lib/core";
import json from "highlight.js/lib/languages/json";

hljs.registerLanguage("json", json);

const $message = useMessage();

const reverting = ref(true);

const users = ref<any[]>([]);
const selection = ref<any[]>([]);

const popups = reactive({
  notifications: false,
  locks: false
});

const tableOptions = {
  columns: [
    { type: "selection", width: 60 },
    { title: "UID", key: "id", sorter: "default", width: 80 },
    { title: "Name", key: "name", width: 180, ellipsis: { tooltip: true } },
    { title: "Nickname", key: "nickname", width: 180, ellipsis: { tooltip: true } },
    { title: "Description", key: "description", width: 360, ellipsis: { tooltip: true } },
    {
      title: "Permissions",
      key: "permissions",
      width: 240,
      ellipsis: { tooltip: true },
      render(row: any) {
        return h(NCode, { hljs, code: JSON.stringify(row.permissions) });
      }
    },
    {
      title: "Is Verified",
      key: "is_verified",
      render(row: any) {
        return h(NIcon, { component: row.verified_at != null ? CheckRound : CloseRound });
      }
    },
    {
      title: "Created At",
      key: "created_at",
      sorter: "default",
      render(row: any) {
        return h(NTime, { time: new Date(row.created_at), format: "yyyy-MM-dd hh:mm" });
      }
    }
  ]
};

async function fetch() {
  try {
    reverting.value = true;
    users.value = (await http.get("/api/administration/users")).data;
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
