<template>
  <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="update">
    <n-form-item :label="$t('pages.users.personal-center.change-password.form.old-password')" path="old">
      <n-input
        v-model:value="payload.old_password"
        type="password"
        :placeholder="$t('pages.users.personal-center.change-password.form.old-password.placeholder')"
      />
    </n-form-item>
    <n-form-item :label="$t('pages.users.personal-center.change-password.form.new-password')" path="new">
      <n-input
        v-model:value="payload.new_password"
        type="password"
        :placeholder="$t('pages.users.personal-center.change-password.form.new-password.placeholder')"
      />
    </n-form-item>
    <n-button type="primary" attr-type="submit" :loading="submitting">Apply</n-button>
  </n-form>
</template>

<script lang="ts" setup>
import { usePrincipal } from "@/stores/account";
import { useMessage, type FormInst, type FormRules } from "naive-ui";
import { reactive, ref } from "vue";
import { http } from "@/utils/http";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $message = useMessage();
const $principal = usePrincipal();

const form = ref<FormInst | null>(null);
const submitting = ref(false);

const rules: FormRules = {
  old_password: {
    required: true,
    validator: (_, v) => v.length >= 6,
    message: t("pages.users.personal-center.change-password.form.old-password.validate"),
    trigger: ["blur", "input"]
  },
  new_password: {
    required: true,
    validator: (_, v) => v.length >= 6,
    message: t("pages.users.personal-center.change-password.form.new-password.validate"),
    trigger: ["blur", "input"]
  }
};

const payload = reactive({
  old_password: "",
  new_password: ""
});

function update() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return;
    }

    try {
      submitting.value = true;
      await http.put("/api/users/password", payload);

      await $principal.fetch();

      $message.success(t('pages.users.personal-center.change-password.feedback.success'));
    } catch (e: any) {
      $message.error(t('common.feedback.unknown-error', [e.response.data ?? e.message]));
    } finally {
      submitting.value = false;
    }
  });
}
</script>
