<template>
  <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="update">
    <n-form-item :label="$t('pages.users.personal-center.personal-information.avatar')" path="avatar">
      <n-spin :show="submitting">
        <div class="flex items-center">
          <n-avatar
            class="cursor-pointer"
            color="transparent"
            :size="48"
            :src="usePlaceholder('avatar', $principal.account?.avatar_url)"
            @click="$refs['personalize.avatar'].click()"
          />
          <div class="flex text-sm text-gray-400">
            <n-icon size="20" :component="ArrowLeftRound" />
            {{ $t("pages.users.personal-center.personal-information.avatar.hint") }}
          </div>
        </div>

        <input
          ref="personalize.avatar"
          type="file"
          accept="image/*"
          @change="(e: any) => personalise('avatar')(e)"
          style="display: none"
        />
      </n-spin>
    </n-form-item>
    <n-form-item :label="$t('pages.users.personal-center.personal-information.banner')" path="banner">
      <n-spin :show="submitting">
        <n-image
          class="w-full rounded cursor-pointer"
          color="transparent"
          :src="usePlaceholder('banner', $principal.account?.banner_url)"
          object-fit="cover"
          height="120"
          width="340"
          preview-disabled
          @click="$refs['personalize.banner'].click()"
        />
        <div class="flex text-sm text-gray-400">
          <n-icon size="20" :component="ArrowDropUpRound" />
          {{ $t("pages.users.personal-center.personal-information.banner.hint") }}
        </div>

        <input
          ref="personalize.banner"
          type="file"
          accept="image/*"
          @change="(e: any) => personalise('banner')(e)"
          style="display: none"
        />
      </n-spin>
    </n-form-item>
    <n-form-item :label="$t('pages.users.personal-center.personal-information.form.username')" path="name">
      <n-tooltip trigger="hover" placement="top">
        <template #trigger>
          <n-input
            v-model:value="payload.name"
            :placeholder="$t('pages.users.personal-center.personal-information.form.username.placeholder')"
            disabled
          />
        </template>
        {{ $t('pages.users.personal-center.personal-information.form.username.hint') }}
      </n-tooltip>
    </n-form-item>
    <n-form-item :label="$t('pages.users.personal-center.personal-information.form.nickname')" path="nickname">
      <n-input
        v-model:value="payload.nickname"
        :placeholder="$t('pages.users.personal-center.personal-information.form.nickname.placeholder')"
      />
    </n-form-item>
    <n-form-item :label="$t('pages.users.personal-center.personal-information.form.description')" path="description">
      <n-input
        type="textarea"
        v-model:value="payload.description"
        :placeholder="$t('pages.users.personal-center.personal-information.form.description.placeholder')"
      />
    </n-form-item>
    <n-button type="primary" attr-type="submit" :loading="submitting">{{ $t("actions.apply") }}</n-button>
  </n-form>
</template>

<script lang="ts" setup>
import { usePrincipal } from "@/stores/account";
import { useMessage, type FormInst, type FormRules } from "naive-ui";
import { usePlaceholder } from "@/utils/placeholders";
import { reactive, ref } from "vue";
import { http } from "@/utils/http";
import { ArrowDropUpRound, ArrowLeftRound } from "@vicons/material";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $message = useMessage();
const $principal = usePrincipal();

const form = ref<FormInst | null>(null);
const submitting = ref(false);

const rules: FormRules = {
  name: {
    required: true,
    validator: (_, v) => new RegExp(/^\w+$/).test(v),
    message: t('pages.users.personal-center.personal-information.form.username.validate'),
    trigger: ["blur", "input"]
  },
  nickname: {
    required: true,
    validator: (_, v) => v.length >= 4,
    message: t('pages.users.personal-center.personal-information.form.nickname.validate'),
    trigger: ["blur", "input"]
  }
};

const payload = reactive({
  name: $principal.account?.name,
  nickname: $principal.account?.nickname,
  description: $principal.account?.description
});

function update() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return;
    }

    try {
      submitting.value = true;
      await http.put("/api/users/self", payload);

      await $principal.fetch();

      $message.success(t('pages.users.personal-center.personal-information.feedback.success'));
    } catch (e: any) {
      $message.error(t('common.feedback.unknown-error', [e.response.data ?? e.message]));
    } finally {
      submitting.value = false;
    }
  });
}

function personalise(mode: "avatar" | "banner") {
  return async (e: any) => {
    const file = e.target.files[0];
    if (file == null) {
      return;
    }

    try {
      submitting.value = true;
      const payload = new FormData();

      switch (mode) {
        case "avatar":
          payload.set("avatar", file);
          await http.put("/api/users/self/personalize?field=avatar", payload);
          break;
        case "banner":
          payload.set("banner", file);
          await http.put("/api/users/self/personalize?field=banner", payload);
          break;
      }

      await $principal.fetch();

      $message.success(t('pages.users.personal-center.personal-information.feedback.success.personalise', [t(`pages.users.personal-center.personal-information.${mode}`)]));
    } catch (e: any) {
      $message.error(t('common.feedback.unknown-error', [e.response.data ?? e.message]));
    } finally {
      submitting.value = false;
    }
  };
}
</script>
