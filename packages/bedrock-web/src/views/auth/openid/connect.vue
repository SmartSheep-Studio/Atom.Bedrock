<template>
  <div class="container flex justify-center items-center h-max">
    <div class="wrapper">
      <div class="text-center">
        <img src="../../../assets/icon.png" width="64" height="64" />
      </div>

      <n-spin :show="reverting">
        <n-card v-if="error">
          <n-alert :bordered="false" title="Something went wrong..." type="error">
            {{ error }}
          </n-alert>
        </n-card>

        <n-card v-else>
          <div class="text-center">
            <div class="text-lg">{{ $t("pages.auth.openid.connect.title", [client?.name]) }}</div>
            <div class="text-sm text-gray-400">
              {{ $t("pages.auth.openid.connect.desc", [client?.name, $endpoint.info.name ?? "Project Atom"]) }}
            </div>
          </div>

          <div>
            <n-tag type="primary" v-if="client?.is_official">
              <template #icon>
                <n-icon :component="VerifiedRound" />
              </template>
              {{ $t("pages.auth.openid.connect.tags.official") }}
            </n-tag>
            <n-tag type="success" v-if="client?.is_verified">
              <template #icon>
                <n-icon :component="CheckRound" />
              </template>
              {{ $t("pages.auth.openid.connect.tags.verified") }}
            </n-tag>
            <n-tag type="error" v-if="client?.is_danger">
              <template #icon>
                <n-icon :component="WarningRound" />
              </template>
              {{ $t("pages.auth.openid.connect.tags.dangerous") }}
            </n-tag>
          </div>

          <div class="mt-4">
            <div>{{ $t("pages.auth.openid.connect.tips", [client?.name]) }}</div>
            <ol>
              <li v-for="tag in query.scope?.toString().split(' ')">
                <div class="font-bold">{{ tag }}</div>
              </li>
            </ol>
          </div>

          <div class="mt-4">
            {{ $t("pages.auth.openid.connect.tips.redirect", [query?.redirect_uri]) }}
          </div>

          <div class="mt-4">
            <n-form :model="payload" @submit.prevent="approve">
              <n-card size="small" embedded class="unlogin-fields mb-4" v-if="!$principal.isSigned">
                <n-form-item :label="$t('pages.auth.sign-in.form.username')" path="id">
                  <n-input
                    v-model:value="payload.id"
                    :placeholder="$t('pages.auth.sign-in.form.username.placeholder')"
                  />
                </n-form-item>
                <n-form-item :label="$t('pages.auth.sign-in.form.password')" path="password">
                  <n-input
                    v-model:value="payload.password"
                    type="password"
                    :placeholder="$t('pages.auth.sign-in.form.password.placeholder')"
                  />
                </n-form-item>
                <div class="text-center text-gray-400">
                  {{ $t("pages.auth.openid.connect.tips.sign-in", [$endpoint.info.name ?? "Project Atom"]) }}
                </div>
              </n-card>

              <div class="flex justify-center gap-3">
                <n-button type="warning" :loading="submitting" @click="$router.back()">
                  {{ $t("actions.decline") }}
                </n-button>
                <n-button type="primary" :loading="submitting" attr-type="submit">
                  {{ $t("actions.approve") }}
                </n-button>
              </div>
            </n-form>
          </div>
        </n-card>

        <div class="mt-4 flex justify-center items-center text-xs text-center text-gray-400">
          <span>{{ $t("brand.atom-id.powered-by") }}</span>
          <n-icon :component="LockOpenFilled" size="16" class="ml-1" />
        </div>
      </n-spin>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useRouter } from "vue-router";
import { useCookies } from "@vueuse/integrations/useCookies";
import { VerifiedRound, CheckRound, WarningRound, LockOpenFilled } from "@vicons/material";
import { useEndpoint } from "@/stores/connection";
import { useMessage } from "naive-ui";
import { onMounted, reactive, ref } from "vue";
import { usePrincipal } from "@/stores/account";
import { http } from "@/utils/http";
import { useI18n } from "vue-i18n";
import qs from "query-string";

const { t } = useI18n();

const $principal = usePrincipal();
const $endpoint = useEndpoint();
const $message = useMessage();
const $cookies = useCookies(["oauth_mode", "oauth_platform"]);
const $router = useRouter();

const reverting = ref(true);
const submitting = ref(false);

const error = ref<any>(null);
const query = ref<any>({});
const client = ref<any>(null);

const payload = reactive({
  id: "",
  password: "",
});

async function fetch() {
  try {
    reverting.value = true;
    query.value = qs.parse(window.location.search);
    const res = await http.get("/api/auth/openid/connect", {
      params: query.value,
      validateStatus: (status: number) => (status >= 200 && status < 400) || status === 401,
    });
    if (res.data.skip) {
      window.location.href = `${query.value.redirect_uri}?code=${res.data.session.code}&state=${query.value.state}`;
    } else {
      client.value = res.data.client;
    }
  } catch (e: any) {
    if (e.response.status === 400) {
      error.value = t("pages.auth.openid.connect.feedback.failed.fetch-information");
    } else if (e.response.status === 404) {
      error.value = t("pages.auth.openid.connect.feedback.failed.missing-client");
    } else if (e.response.status === 403) {
      error.value = t("pages.auth.openid.connect.feedback.failed.unsafe");
    } else {
      error.value = t("pages.auth.openid.connect.feedback.failed.unknown", [e]);
    }

    $message.error(`Something went wrong... ${e}`);
  } finally {
    reverting.value = false;
  }
}

async function approve() {
  if (!$principal.isSigned && (payload.id.length <= 0 || payload.password.length <= 0)) {
    return;
  }

  try {
    submitting.value = true;
    $cookies.set("oauth_mode", "clients");
    $cookies.set("oauth_platform", "internal");
    const res = await http.post(`/api/auth/openid/connect`, payload, {
      params: {
        client_id: query.value.client_id,
        redirect_uri: encodeURIComponent(query.value.redirect_uri),
        response_type: "code",
        scope: query.value.scope,
      },
    });
    if (query.value.response_type === "code") {
      window.location.href = `${query.value.redirect_uri}?code=${res.data.session.code}&state=${query.value.state}`;
    } else if (query.value.response_type === "token") {
      window.location.href = `${query.value.redirect_uri}?access_token=${res.data.accessToken}&refresh_token=${res.data.refreshToken}&state=${query.value.state}`;
    }
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e]));
  } finally {
    submitting.value = false;
  }
}

onMounted(async () => {
  await fetch();
});
</script>

<style scoped>
.wrapper {
  width: 60vw;
  max-width: 480px;
  min-width: 300px;
}
</style>
