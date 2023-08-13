<template>
  <n-form ref="form" :model="payload" @submit.prevent="update">
    <n-form-item v-for="(item, index) in payload">
      <template #label>
        <div class="flex items-center gap-2">
          <div>{{ item.description }}</div>
          <n-badge type="success" v-if="item.verified_at">
            <template #value>
              <n-icon :component="VerifiedRound" />
            </template>
          </n-badge>
        </div>
      </template>
      <n-input
        :placeholder="$t('pages.users.personal-center.contacts.form.name.placeholder')"
        v-model:value="item.description"
      />
      <n-input
        style="margin-left: 8px"
        :placeholder="$t('pages.users.personal-center.contacts.form.content.placeholder')"
        v-model:value="item.content"
      />
      <n-button style="margin-left: 8px" @click="payload.splice(index, 1)">
        <template #icon>
          <n-icon :component="DeleteRound" />
        </template>
      </n-button>
    </n-form-item>
    <n-space size="small">
      <n-button type="primary" attr-type="submit" :loading="submitting">{{ $t("actions.apply") }}</n-button>
      <n-button
        @click="
          payload.push({
            content: '',
            description: $t('pages.users.personal-center.contacts.form.contact', [payload.length]),
            type: 'email'
          })
        "
      >
        {{ $t("actions.add-item") }}
      </n-button>
    </n-space>
  </n-form>
</template>

<script lang="ts" setup>
import { usePrincipal } from "@/stores/account";
import { useMessage, type FormInst } from "naive-ui";
import { VerifiedRound, DeleteRound } from "@vicons/material";
import { ref } from "vue";
import { http } from "@/utils/http";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $message = useMessage();
const $principal = usePrincipal();

const form = ref<FormInst | null>(null);
const submitting = ref(false);

const payload = ref<any[]>($principal.account?.contacts ?? []);

function update() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return;
    }

    try {
      submitting.value = true;
      await http.put("/api/users/contacts", { contacts: payload.value });

      await $principal.fetch();

      $message.success(t("pages.users.personal-center.contacts.feedback.success"));
    } catch (e: any) {
      $message.error(t('common.feedback.unknown-error', [e.response.data ?? e.message]));
    } finally {
      submitting.value = false;
    }
  });
}
</script>
