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
                  <n-button circle secondary size="tiny" type="error" @click="terminate(item)">
                    <template #icon>
                      <n-icon size="12" :component="LogOutRound" />
                    </template>
                  </n-button>
                </template>
                {{ $t("pages.auth.principal.api-tokens.remove") }}
              </n-tooltip>
            </template>
            <template #description>
              <div class="flex">
                <div>{{ item.description }}</div>
                <div class="text-gray-400 ml-1">#{{ item.id }}</div>
              </div>
            </template>
          </n-thing>
        </n-list-item>
      </n-list>
      <div class="mt-2">
        <n-button class="w-full" type="primary" @click="creating = true">
          <template #icon>
            <n-icon :component="KeyRound" />
          </template>
          {{ $t("pages.auth.principal.api-tokens.new") }}
        </n-button>
      </div>
    </n-spin>

    <div class="flex justify-center mt-4">
      <n-pagination
        v-model:page="pagination.page"
        :page-count="Math.ceil((rawData.length ?? 0) / pagination.pageSize)"
        :page-slot="pagination.slot"
      />
    </div>

    <n-modal v-model:show="creating">
      <n-card class="w-dialog" :title="$t('pages.auth.principal.api-tokens.new')" :bordered="false" size="huge">
        <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="create">
          <n-form-item :label="$t('pages.auth.principal.api-tokens.form.description')" path="description">
            <n-input
              type="textarea"
              :placeholder="$t('pages.auth.principal.api-tokens.form.description.placeholder')"
              v-model:value="payload.description"
            />
          </n-form-item>
          <n-form-item
            v-for="(_, index) in payload.scope"
            :label="$t('pages.auth.principal.api-tokens.form.scope', [index + 1])"
          >
            <n-auto-complete
              :placeholder="$t('pages.auth.principal.api-tokens.form.scope.placeholder')"
              v-model:value="payload.scope[index]"
              :options="scopeOptions"
            />
            <n-button style="margin-left: 8px" @click="payload.scope.splice(index, 1)">
              <template #icon>
                <n-icon :component="DeleteRound" />
              </template>
            </n-button>
          </n-form-item>

          <n-space class="mt-2">
            <n-button type="primary" attr-type="submit" :loading="submitting">{{ $t("actions.submit") }}</n-button>
            <n-button @click="payload.scope.push('')">{{ $t("actions.add-item") }}</n-button>
          </n-space>
        </n-form>
      </n-card>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { usePrincipal } from "@/stores/account";
import { computed, onMounted, reactive, ref } from "vue";
import { LogOutRound, KeyRound, DeleteRound } from "@vicons/material";
import { http } from "@/utils/http";
import { useMessage, type FormRules, type FormInst, useDialog } from "naive-ui";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $dialog = useDialog();
const $message = useMessage();
const $principal = usePrincipal();

const rawData = ref<any[]>([]);
const data = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize;
  return rawData.value?.reverse().slice(start, start + pagination.pageSize) ?? [];
});

const requesting = ref(true);
const submitting = ref(false);
const creating = ref(false);

const form = ref<FormInst | null>(null);
const rules: FormRules = {
  description: {
    required: true,
    message: t("pages.auth.principal.api-tokens.form.description.validate"),
    trigger: ["blur", "input"]
  }
};

const scopeOptions = [
  "principal",
  "sessions.read",
  "sessions.delete",
  "oauth.read",
  "oauth.create",
  "oauth.update",
  "oauth.delete",
  "auth.openid",
  "auth.openid.approve",
  "users.update",
  "users.update.avatar",
  "users.update.banner"
];

const pagination = reactive({
  page: 1,
  pageSize: 5,
  slot: 5
});

const payload = reactive({
  description: "",
  scope: [""]
});

async function fetch() {
  try {
    requesting.value = true;
    const res = await http.get("/api/auth/tokens");
    rawData.value = res.data;
  } catch (e: any) {
    $message.error(t('common.feedback.unknown-error', [e]));
  } finally {
    requesting.value = false;
  }
}

function create() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return;
    }

    try {
      submitting.value = true;

      const res = await http.post("/api/auth/tokens", payload);
      await Promise.all([fetch(), $principal.fetch()]);

      creating.value = false;
      $dialog.success({
        title: t('pages.auth.principal.api-tokens.feedback.success.create'),
        content: t('pages.auth.principal.api-tokens.feedback.success.create.desc', [res.data.token]),
        positiveText: t('actions.ok')
      });
    } catch (e: any) {
      $message.error(t('common.feedback.unknown-error', [e]));
    } finally {
      submitting.value = false;
    }
  });
}

async function terminate(item: any) {
  try {
    requesting.value = true;
    await http.delete("/api/auth/sessions", { params: { id: item.id } });

    await Promise.all([fetch(), $principal.fetch()]);

    $message.success(t('pages.auth.principal.api-tokens.feedback.success.remove'));
  } catch (e: any) {
    $message.error(t('common.feedback.unknown-error', [e]));
  } finally {
    requesting.value = false;
  }
}

onMounted(() => {
  fetch();
});
</script>
