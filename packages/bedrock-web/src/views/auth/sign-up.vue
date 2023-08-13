<template>
  <div class="container h-max flex justify-center items-center">
    <div>
      <div class="text-center">
        <img src="../../assets/icon.png" width="64" height="64" />
        <div class="text-2xl font-bold">{{ $t('actions.sign-up') }}</div>
      </div>
      <div class="w-96 mt-4">
        <n-card>
          <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="submit">
            <n-form-item :label="$t('pages.auth.sign-up.form.username')" path="name">
              <n-input
                v-model:value="payload.name"
                :placeholder="$t('pages.auth.sign-up.form.username.placeholder')"
              />
            </n-form-item>
            <n-form-item :label="$t('pages.auth.sign-up.form.nickname')" path="nickname">
              <n-input
                v-model:value="payload.nickname"
                :placeholder="$t('pages.auth.sign-up.form.nickname.placeholder')"
              />
            </n-form-item>
            <n-form-item :label="$t('pages.auth.sign-up.form.email')" path="contact">
              <n-input v-model:value="payload.contact" :placeholder="$t('pages.auth.sign-up.form.email.placeholder')" />
            </n-form-item>
            <n-form-item :label="$t('pages.auth.sign-up.form.password')" path="password">
              <n-input
                v-model:value="payload.password"
                type="password"
                :placeholder="$t('pages.auth.sign-up.form.password.placeholder')"
              />
            </n-form-item>

            <n-button class="w-full" type="primary" attr-type="submit" :loading="submitting">
              {{$t("actions.sign-up") }}
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
import { usePrincipal } from "@/stores/account";
import { parseRedirect } from "@/utils/callback";
import { http } from "@/utils/http";
import { LockOpenFilled } from "@vicons/material";
import { useMessage, type FormRules, type FormInst } from "naive-ui";
import { reactive, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";

const { t } = useI18n();

const $route = useRoute();
const $router = useRouter();
const $message = useMessage();
const $principal = usePrincipal();
const $history: any = window.history.state;

const form = ref<FormInst | null>(null);
const submitting = ref(false);

const rules: FormRules = {
  name: {
    required: true,
    validator: (_, v) => new RegExp(/^\w+$/).test(v),
    message: t('pages.auth.sign-up.form.username.validate'),
    trigger: ["blur", "input"]
  },
  nickname: {
    required: true,
    validator: (_, v) => v.length >= 4,
    message: t('pages.auth.sign-up.form.nickname.validate'),
    trigger: ["blur", "input"]
  },
  contact: {
    required: true,
    validator: (_, v) => new RegExp(/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/).test(v),
    message: t('pages.auth.sign-up.form.email.validate'),
    trigger: ["blur", "input"]
  },
  password: {
    required: true,
    validator: (_, v) => v.length >= 6,
    message: t('pages.auth.sign-up.form.password.validate'),
    trigger: ["blur", "input"]
  }
};

const payload = reactive({
  name: "",
  nickname: "",
  contact: "",
  password: ""
});

function submit() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return;
    }

    try {
      submitting.value = true;
      await http.post("/api/auth/sign-up", payload);
      await http.post("/api/auth/sign-in", { id: payload.name, password: payload.password });

      await $principal.fetch();

      $message.success(t('pages.auth.sign-up.feedback.success'));
      $router.push(await parseRedirect($route.query));
    } catch (e: any) {
      $message.error(t('common.feedback.unknown-error', [e.response.data ?? e.message]));
    } finally {
      submitting.value = false;
    }
  });
}
</script>
