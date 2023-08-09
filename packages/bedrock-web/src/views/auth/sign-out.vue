<template>
  <div class="container h-max flex justify-center items-center">
    <div>
      <div class="text-center">
        <img src="../../assets/icon.png" width="64" height="64" />
        <div class="text-2xl font-bold">{{ $t("actions.sign-out") }}</div>
      </div>
      <div class="w-96 mt-4">
        <n-card>
          <div>
            {{ $t("pages.auth.sign-out.notice") }}
          </div>

          <n-button class="w-full mt-4" type="error" @click="submit" :loading="submitting">{{
            $t("actions.sign-out")
          }}</n-button>
        </n-card>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { usePrincipal } from "@/stores/account";
import { http } from "@/utils/http";
import { useMessage } from "naive-ui";
import { useRouter } from "vue-router";
import { ref } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $router = useRouter();
const $principal = usePrincipal();
const $message = useMessage();

const submitting = ref(false);

async function submit() {
  submitting.value = true;

  try {
    await http.delete("/api/auth/sessions");
    $router.push({ name: "landing" }).then(() => {
      $principal.logout();
    });
  } catch (e: any) {
    $message.error(t('common.feedback.unknown-error', [e]));
  } finally {
    submitting.value = false;
  }
}
</script>
