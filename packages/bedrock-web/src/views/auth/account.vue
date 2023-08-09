<template>
  <div>
    <img
      class="w-full object-cover"
      style="height: 180px"
      :src="usePlaceholder('banner', $principal.account?.banner_url)"
    />

    <n-tabs type="line" justify-content="center" animated>
      <n-tab-pane
        name="personal-information"
        :tab="$t('pages.auth.principal.personal-information')"
        display-directive="show:lazy"
      >
        <div class="container pt-6">
          <div class="md:px-8 lg:px-32 xl:px-48">
            <div class="mb-8">
              <div class="text-2xl font-bold">{{ greetings }}, {{ $principal.account?.nickname }}.</div>
              <div class="text-lg">{{ $principal.account?.description }}</div>
              <div class="text-gray-400">#{{ $principal.account?.name }}</div>
            </div>

            <n-alert class="mt-2" type="warning" title="Unverified account" v-if="$principal.account?.verified_at == null">
              Currently your account wasn't verified, some features isn't available. You can follow the steps below to
              verify your account!
            </n-alert>
            <n-card
              class="mt-2"
              :title="$t('pages.auth.principal.verification')"
              v-if="$principal.account?.verified_at == null"
            >
              <verification />
            </n-card>
            <n-card class="mt-2" :title="$t('pages.auth.principal.personal-information')">
              <personal-information />
            </n-card>
          </div>
        </div>
      </n-tab-pane>
      <n-tab-pane name="contacts" :tab="$t('pages.auth.principal.contacts')" display-directive="show:lazy">
        <div class="container pt-4">
          <div class="md:px-8 lg:px-32 xl:px-48">
            <n-card class="mt-2" :title="$t('pages.auth.principal.contacts')">
              <contacts />
            </n-card>
          </div>
        </div>
      </n-tab-pane>
      <n-tab-pane name="security" :tab="$t('pages.auth.principal.security')" display-directive="show:lazy">
        <div class="container pt-4">
          <div class="md:px-8 lg:px-32 xl:px-48">
            <n-grid :cols="2" item-responsive responsive="screen" :x-gap="8" :y-gap="8">
              <n-gi span="2 l:1">
                <n-card class="mt-2 h-full" :title="$t('pages.auth.principal.sessions')">
                  <sessions />
                </n-card>
              </n-gi>
              <n-gi span="2 l:1">
                <n-card class="mt-2" :title="$t('pages.auth.principal.change-password')">
                  <new-password />
                </n-card>
              </n-gi>
              <n-gi span="2 l:1">
                <n-card class="mt-2" :title="$t('pages.auth.principal.api-tokens')">
                  <api-tokens />
                </n-card>
              </n-gi>
            </n-grid>
          </div>
        </div>
      </n-tab-pane>
      <n-tab-pane name="applications" :tab="$t('pages.auth.principal.applications')" display-directive="show:lazy">
        <div class="container pt-4">
          <div class="md:px-12 lg:px-48 xl:px-72">
            <n-card :title="$t('pages.auth.principal.oauth-clients')">
              <oauth-clients />
            </n-card>
          </div>
        </div>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script lang="ts" setup>
import { usePrincipal } from "@/stores/account";
import { usePlaceholder } from "@/utils/placeholders";
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import Verification from "@/views/auth/parts/verification.vue";
import PersonalInformation from "@/views/auth/parts/personal-information.vue";
import ApiTokens from "@/views/auth/parts/api-tokens.vue";
import OauthClients from "@/views/auth/parts/oauth-clients.vue";
import Sessions from "@/views/auth/parts/sessions.vue";
import Contacts from "@/views/auth/parts/contacts.vue";
import NewPassword from "@/views/auth/parts/new-password.vue";

const { t } = useI18n();

const $principal = usePrincipal();

const greetings = computed(() =>
  new Date().getHours() < 12
    ? t("common.greetings.morning")
    : new Date().getHours() <= 18 && new Date().getHours() >= 12
    ? t("common.greetings.afternoon")
    : t("common.greetings.night")
);
</script>
