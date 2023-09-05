<template>
  <div class="pt-2">
    <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="submit">
      <n-form-item label="Title" path="title">
        <n-input placeholder="Notification title" v-model:value="payload.title" />
      </n-form-item>
      <n-form-item label="Description" path="description">
        <n-input placeholder="Notification description" v-model:value="payload.description" />
      </n-form-item>
      <n-form-item label="Level" path="level">
        <n-select :options="levels" v-model:value="payload.level" />
      </n-form-item>
      <n-form-item label="Content" path="content">
        <n-input type="textarea" placeholder="Notification content" v-model:value="payload.content" />
      </n-form-item>
      <n-form-item label="Recipients">
        <n-card size="small" embedded>
          <div class="font-mono">
            <b>{{ props.selection.length }}</b> users you selected.
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
import { useMessage, type FormInst, type FormRules } from "naive-ui";
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
  title: {
    required: true,
    message: "Notification title is required",
    trigger: ["blur", "input"],
  },
  description: {
    required: true,
    message: "Notification description is required",
    trigger: ["blur", "input"],
  },
  content: {
    required: true,
    message: "Notification description is required",
    trigger: ["blur", "input"],
  },
};

const levels = [
  { label: "Tips —— Won't send instant message", value: "tips" },
  { label: "Info —— Will send instant message", value: "info" },
  { label: "Warning —— Will send instant message", value: "warning" },
  { label: "Alert —— Will send instant message", value: "alert" },
];

const payload = reactive({
  title: "",
  description: "",
  level: "tips",
  content: "",
});

function submit() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return;
    }

    try {
      submitting.value = true;
      for (const item of props.selection) {
        await http.post("/api/administration/notifications", {
          title: payload.title,
          description: payload.description,
          level: payload.level,
          content: payload.content,
          recipient_id: item,
        });
      }

      $message.success("Successfully sent notifications to users you selected!");
      emits("close")
    } catch (e: any) {
      $message.error(t("common.feedback.unknown-error", [e.response.data ?? e.message]));
    } finally {
      submitting.value = false;
    }
  });
}
</script>
