<template>
  <div class="pt-2">
    <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="submit">
      <n-form-item label="Reason" path="reason">
        <n-input
          placeholder="Ban/Lock reason"
          v-model:value="payload.reason"
        />
      </n-form-item>
      <n-form-item label="Expired at" path="expired_at">
        <n-date-picker
          placeholder="Ban/Lock expired at"
          type="datetime"
          class="flex-grow"
          v-model:value="payload.expired_at"
        />
      </n-form-item>
      <n-form-item label="Quiet" path="quiet">
        <n-checkbox
          label="Won't send notification and instant message to user(s)"
          v-model:value="payload.quiet"
        />
      </n-form-item>
      <n-form-item label="Recipients">
        <n-card size="small" embedded>
          <div class="font-mono">
            <b>{{ props.selection.length }}</b> user(s) you selected.
          </div>
        </n-card>
      </n-form-item>

      <div>
        <n-button type="primary" attr-type="submit" :loading="submitting">Submit</n-button>
      </div>
    </n-form>
  </div>
</template>

<script lang="ts" setup>
import { type FormInst, type FormRules, useMessage } from "naive-ui";
import { reactive, ref } from "vue";
import { http } from "@/utils/http";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $message = useMessage();

const props = defineProps<{ selection: any[] }>();
const emits = defineEmits(["close"]);

const form = ref<FormInst | null>(null);
const submitting = ref(false);

const rules: FormRules = {
  reason: {
    required: true,
    message: "Reason is required",
    trigger: ["blur", "input"]
  }
};

const payload = reactive({
  reason: "",
  expired_at: null,
  quiet: false
});

function submit() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return;
    }

    try {
      submitting.value = true;
      for (const item of props.selection) {
        await http.post("/api/administration/locks", {
          reason: payload.reason,
          expired_at: new Date(payload.expired_at),
          user_id: item
        }, {
          params: { quiet: payload.quiet ? "yes" : "no" }
        });
      }

      $message.success("Successfully banned the users you selected!");
      emits("close");
    } catch (e: any) {
      $message.error(t("common.feedback.unknown-error", [e.response.data ?? e.message]));
    } finally {
      submitting.value = false;
    }
  });
}
</script>
