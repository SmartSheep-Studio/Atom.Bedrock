<template>
  <div>
    <n-spin :show="requesting">
      <n-list bordered hoverable>
        <n-empty v-if="data.length <= 0" class="py-8" description="There's no data. Why not you add one?" />

        <n-list-item v-for="item in data">
          <n-thing>
            <template #description>
              <div class="flex">
                <div class="font-bold">{{ item.name }}</div>
                <div class="text-gray-400 ml-1">#{{ item.slug }}</div>
              </div>
              <div>{{ item.description }}</div>

              <n-space class="mt-1" size="small">
                <n-button
                  size="tiny"
                  type="warning"
                  @click="
                    () => {
                      refer(item);
                      updating = true;
                    }
                  "
                >
                  <template #icon>
                    <n-icon :component="EditRound" />
                  </template>
                  {{ $t("actions.edit") }}
                </n-button>
                <n-button size="tiny" type="error" @click="destroy(item)">
                  <template #icon>
                    <n-icon :component="DeleteRound" />
                  </template>
                  {{ $t("actions.destroy") }}
                </n-button>
              </n-space>
            </template>
          </n-thing>
        </n-list-item>
      </n-list>
    </n-spin>

    <div class="flex justify-between mt-4">
      <n-button type="primary" size="small" @click="creating = true">
        <template #icon>
          <n-icon :component="PlusRound" />
        </template>
        {{ $t("pages.users.personal-center.oauth-clients.new") }}
      </n-button>
      <n-pagination
        v-model:page="pagination.page"
        :page-count="Math.ceil((rawData.length ?? 0) / pagination.pageSize)"
        :page-slot="pagination.slot"
      />
    </div>

    <n-modal v-model:show="creating">
      <n-card class="w-dialog" :title="$t('pages.users.personal-center.oauth-clients.new')" :bordered="false" size="huge">
        <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="create">
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.slug')" path="slug">
            <n-input
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.slug.placeholder')"
              v-model:value="payload.slug"
            />
          </n-form-item>
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.name')" path="name">
            <n-input
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.name.placeholder')"
              v-model:value="payload.name"
            />
          </n-form-item>
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.secret')" path="secret">
            <n-input
              type="password"
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.secret.placeholder')"
              v-model:value="payload.secret"
            />
          </n-form-item>
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.description')" path="description">
            <n-input
              type="textarea"
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.description.placeholder')"
              v-model:value="payload.description"
            />
          </n-form-item>
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.urls')" path="urls">
            <n-dynamic-input
              v-model:value="payload.urls"
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.urls.placeholder')"
            />
          </n-form-item>
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.callbacks')" path="callbacks">
            <n-dynamic-input
              v-model:value="payload.callbacks"
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.callbacks.placeholder')"
            />
          </n-form-item>

          <n-space class="mt-2">
            <n-button type="primary" attr-type="submit" :loading="submitting">{{ $t("actions.submit") }}</n-button>
          </n-space>
        </n-form>
      </n-card>
    </n-modal>

    <n-modal v-model:show="updating">
      <n-card
        class="w-dialog"
        :title="$t('pages.users.personal-center.oauth-clients.update')"
        :bordered="false"
        size="huge"
        @update:show="(v: boolean) => {if(!v) refer(null)}"
      >
        <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="update">
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.slug')" path="slug">
            <n-input
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.slug.placeholder')"
              v-model:value="payload.slug"
            />
          </n-form-item>
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.name')" path="name">
            <n-input
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.name.placeholder')"
              v-model:value="payload.name"
            />
          </n-form-item>
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.secret')" path="secret">
            <n-input
              type="password"
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.secret.placeholder')"
              v-model:value="payload.secret"
            />
          </n-form-item>
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.description')" path="description">
            <n-input
              type="textarea"
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.description.placeholder')"
              v-model:value="payload.description"
            />
          </n-form-item>
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.urls')" path="urls">
            <n-dynamic-input
              v-model:value="payload.urls"
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.urls.placeholder')"
            />
          </n-form-item>
          <n-form-item :label="$t('pages.users.personal-center.oauth-clients.form.callbacks')" path="callbacks">
            <n-dynamic-input
              v-model:value="payload.callbacks"
              :placeholder="$t('pages.users.personal-center.oauth-clients.form.callbacks.placeholder')"
            />
          </n-form-item>

          <n-space class="mt-2">
            <n-button type="primary" attr-type="submit" :loading="submitting">{{ $t("actions.submit") }}</n-button>
          </n-space>
        </n-form>
      </n-card>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from "vue";
import { DeleteRound, EditRound, PlusRound } from "@vicons/material";
import { useMessage, type FormRules, type FormInst, useDialog } from "naive-ui";
import { http } from "@/utils/http";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $dialog = useDialog();
const $message = useMessage();

const rawData = ref<any[]>([]);
const data = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize;
  return rawData.value.reverse().slice(start, start + pagination.pageSize) ?? [];
});

const requesting = ref(true);
const submitting = ref(false);
const creating = ref(false);
const updating = ref(false);

const form = ref<FormInst | null>(null);
const rules: FormRules = {
  slug: {
    required: true,
    validator: (_, v) => new RegExp(/^[A-Za-z0-9-_]+$/).test(v),
    message: t("pages.users.personal-center.oauth-clients.form.slug.validate"),
    trigger: ["blur", "input"],
  },
  name: {
    required: true,
    validator: (_, v) => v.length >= 4,
    message: t("pages.users.personal-center.oauth-clients.form.name.validate"),
    trigger: ["blur", "input"],
  },
  description: {
    required: true,
    message: t("pages.users.personal-center.oauth-clients.form.description.validate"),
    trigger: ["blur", "input"],
  },
  secret: {
    required: true,
    validator: (_, v) => v.length >= 6,
    message: t("pages.users.personal-center.oauth-clients.form.secret.validate"),
    trigger: ["blur", "input"],
  },
};

const pagination = reactive({
  page: 1,
  pageSize: 5,
  slot: 5,
});

const payload = ref<any>({
  slug: "",
  name: "",
  description: "",
  secret: "",
  urls: [],
  callbacks: [],
});

async function fetch() {
  try {
    requesting.value = true;
    const res = await http.get("/api/users/oauth");
    rawData.value = res.data;
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e]));
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

      await http.post("/api/users/oauth", payload.value);
      await fetch();

      creating.value = false;
      $message.success(t("pages.users.personal-center.oauth-clients.feedback.success.create"));
    } catch (e: any) {
      $message.error(t("common.feedback.unknown-error", [e]));
    } finally {
      submitting.value = false;
    }
  });
}

function update() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return;
    }

    try {
      submitting.value = true;

      await http.put(`/api/users/oauth/${payload.value.target}`, payload.value);
      await fetch();

      updating.value = false;
      $message.success(t("pages.users.personal-center.oauth-clients.feedback.success.update"));
    } catch (e: any) {
      $message.error(t("common.feedback.unknown-error", [e]));
    } finally {
      submitting.value = false;
    }
  });
}

function destroy(item: any) {
  $dialog.warning({
    title: "Warning",
    content: "This operation cannot be undo. Are you confirm?",
    positiveText: "Yes",
    negativeText: "Not really",
    onPositiveClick: async () => {
      try {
        requesting.value = true;
        await http.delete(`/api/users/oauth/${item.slug}`);

        await fetch();

        $message.success(t("pages.users.personal-center.oauth-clients.feedback.success.destroy"));
      } catch (e: any) {
        $message.error(t("common.feedback.unknown-error", [e]));
      } finally {
        requesting.value = false;
      }
    },
  });
}

function refer(item: any | null) {
  if (!item) {
    payload.value = {
      slug: "",
      name: "",
      description: "",
      secret: "",
      urls: [],
      callbacks: [],
    };
  } else {
    payload.value = JSON.parse(JSON.stringify(item));
    payload.value.target = item.slug;
  }
}

onMounted(() => {
  fetch();
});
</script>
