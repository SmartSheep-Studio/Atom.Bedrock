<template>
  <div class="container h-[100vh] flex justify-center items-center">
    <div>
      <div class="text-center">
        <img src="../../assets/icon.png" width="64" height="64" />
        <div class="text-2xl font-bold">{{ $t("actions.sign-in") }}</div>
      </div>
      <div class="w-96 mt-4">
        <n-card>
          <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="submit">
            <n-form-item :label="$t('pages.auth.sign-in.form.username')" path="id">
              <n-input v-model:value="payload.id" :placeholder="$t('pages.auth.sign-in.form.username.placeholder')" />
            </n-form-item>
            <n-form-item :label="$t('pages.auth.sign-in.form.password')" path="password">
              <n-input
                v-model:value="payload.password"
                type="password"
                :placeholder="$t('pages.auth.sign-in.form.password.placeholder')"
              />
            </n-form-item>

            <n-button class="w-full" type="primary" attr-type="submit" :loading="submitting">
              {{ $t("actions.sign-in") }}
            </n-button>

            <n-alert class="mt-4" title="Sign In Required" type="warning" v-if="$history.message">
              {{ $history.message }}
              {{ $t("common.feedback.redirect-to", [$route.query.redirect_uri]) }}
            </n-alert>

            <div class="mt-4 flex justify-center items-center text-xs text-center text-gray-400">
              <span>{{ $t("brand.atom-id.powered-by") }}</span>
              <n-icon :component="LockOpenFilled" size="16" class="ml-1" />
            </div>
          </n-form>
        </n-card>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useMessage, useDialog, type FormInst, type FormRules } from "naive-ui";
import { usePrincipal } from "@/stores/account";
import { parseRedirect } from "@/utils/callback";
import { LockOpenFilled } from "@vicons/material";
import { useRoute, useRouter } from "vue-router";
import { reactive, ref, h } from "vue";
import { http } from "@/utils/http";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $route = useRoute();
const $router = useRouter();
const $dialog = useDialog();
const $message = useMessage();
const $principal = usePrincipal();
const $history: any = window.history.state;

const form = ref<FormInst | null>(null);
const submitting = ref(false);

const rules: FormRules = {
  id: {
    required: true,
    message: t('pages.auth.sign-in.form.username.validate'),
    trigger: ["blur", "input"]
  },
  password: {
    required: true,
    validator: (_, v) => v.length >= 6,
    message: t('pages.auth.sign-in.form.password.validate'),
    trigger: ["blur", "input"]
  }
};

const payload = reactive({
  id: "",
  password: ""
});

function submit() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return;
    }

    try {
      submitting.value = true;
      const res = await http.post("/api/auth/sign-in", payload);

      await $principal.fetch();

      $message.success(t('pages.auth.sign-in.feedback.success', [res.data.user.nickname]));
      $router.push(await parseRedirect($route.query));
    } catch (e: any) {
      if(e.response.status === 403) {
        $dialog.error({
          title: "Your account has been locked.",
          content: () => {
            const messages = [];
            messages.push(h("div", `Your account has been locked by the administrator because ${e.response.data.reason}.`))
            if(e.response.data.expired_at != null) {
              messages.push(h("div", `Your account lock will be unlocked on ${new Date(e.response.data.expired_at).toLocaleString()}`))
            } else {
              messages.push(h("div", "Your account lock will not be automatically unlocked."))
            }
            messages.push(h("div", "If you have any questions, you can contact our members to appeal."))
            messages.push(h("div", "Thanks for your understanding."))
            return h("div", messages)
          },
        })
      } else {
        $message.error(t('common.feedback.unknown-error', [e.response.data ?? e.message]));
      }
      submitting.value = false;
    }
  });
}
</script>
