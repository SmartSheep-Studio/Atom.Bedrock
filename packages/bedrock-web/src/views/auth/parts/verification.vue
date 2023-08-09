<template>
  <div>
    <n-steps class="px-1" v-model:current="step">
      <n-step title="Make a request" description="Select a contact to make request." />
      <n-step title="Enter the one time passcode" description="Find out the passcode from your inbox and enter it!" />
    </n-steps>

    <div class="mt-5">
      <div v-if="step === 1">
        <div class="flex gap-3">
          <n-select :options="options" v-model:value="payload.id" placeholder="Contact you want to verify" />
          <n-button type="primary" class="max-h-[34px]" :loading="submitting" @click="makeRequest">Next</n-button>
        </div>
      </div>
      <div v-if="step === 2">
        <div class="flex gap-3">
          <n-input v-model:value="payload.code" placeholder="The code in your inbox" />
          <n-button type="primary" class="max-h-[34px]" :loading="submitting" @click="apply">Apply</n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, reactive, ref } from "vue";
import { usePrincipal } from "@/stores/account";
import { useMessage } from "naive-ui";
import { useI18n } from "vue-i18n";
import { http } from "@/utils/http";

const { t } = useI18n();

const $principal = usePrincipal();
const $message = useMessage();

const options = computed(() => $principal.account?.contacts.map((v) => ({ label: v.content, value: v.id })) ?? []);
const submitting = ref(false);
const step = ref(1);

const payload = reactive({
  id: null,
  code: "",
});

async function makeRequest() {
  if (payload.id == null) {
    return;
  }
  try {
    submitting.value = true;
    await http.get("/api/users/self/verify", { params: { id: payload.id } });
    $message.success("Successfully sent one time passcode to your inbox!");
    step.value++;
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e]));
  } finally {
    submitting.value = false;
  }
}

async function apply() {
  if (payload.code.length <= 0) {
    return;
  }
  try {
    submitting.value = true;
    await http.get("/api/users/self/verify", { params: { code: payload.code } });
    await $principal.fetch();
    $message.success("Successfully verified your contact and account!");
    step.value = 1;
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e]));
  } finally {
    submitting.value = false;
  }
}
</script>
