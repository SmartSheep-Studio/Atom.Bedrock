<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-dialog-provider>
      <n-message-provider>
        <div class="w-full h-screen relative">
          <n-layout position="absolute">
            <n-layout-header bordered>
              <div class="px-8 md:px-16 h-full flex items-center">
                <div class="cursor-pointer" @click="$router.push({ name: 'landing' })">
                  {{ $endpoint.info?.name ?? $t("brand.fullname") }}
                </div>

                <n-menu
                  class="mx-8 lg:mx-16 xl:mx-36 flex-grow"
                  mode="horizontal"
                  :options="menuOptions"
                  v-model:value="menuKey"
                />

                <div class="flex gap-3" v-if="!$principal.isSigned">
                  <n-button @click="$router.push({ name: 'auth.sign-in' })">{{ $t("actions.sign-in") }}</n-button>
                  <n-button type="primary" @click="$router.push({ name: 'auth.sign-out' })">{{
                    $t("actions.sign-up")
                  }}</n-button>
                </div>
                <div class="flex gap-3" v-else>
                  <n-dropdown placement="bottom-end" show-arrow :options="dropdownOptions" @select="dropdownHandler">
                    <n-avatar
                      size="medium"
                      color="transparent"
                      :src="usePlaceholder('avatar', $principal.account?.avatar_url)"
                    ></n-avatar>
                  </n-dropdown>
                </div>
              </div>
            </n-layout-header>

            <n-layout class="w-full h-max">
              <data-provider>
                <gatekeeper>
                  <router-view />
                </gatekeeper>
              </data-provider>
            </n-layout>
          </n-layout>
        </div>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script lang="ts" setup>
import DataProvider from "@/data-provider.vue";
import Gatekeeper from "@/components/global/gatekeeper.vue";
import { useEndpoint } from "@/stores/connection";
import { usePrincipal } from "@/stores/account";
import { h, type Component, computed, type Ref, ref, watch } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { NIcon, type MenuOption, type DropdownOption } from "naive-ui";
import { usePlaceholder } from "@/utils/placeholders";
import { AccountCircleRound, LogOutRound } from "@vicons/material";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $route = useRoute();
const $router = useRouter();
const $endpoint = useEndpoint();
const $principal = usePrincipal();

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

const themeOverrides = {
  common: {
    primaryColor: "#ca4d4dFF",
    primaryColorHover: "#DF5656FF",
    primaryColorPressed: "#C04747FF",
    primaryColorSuppl: "#A84141FF",
  },
};

const menuKey = ref($route.name);
const menuOptions: Ref<MenuOption[]> = computed(() =>
  $principal.isSigned
    ? [
        ...($endpoint.nav?.map((v) => {
          return {
            label: () =>
              h(
                RouterLink,
                {
                  to: {
                    name: "framework.sub-app",
                    params: { id: v.name },
                  },
                },
                { default: () => v.title }
              ),
            // @ts-ignore
            icon: () => h(NIcon, null, { default: () => h("span", { class: `mdi ${v.icon}` }, null) }),
            key: v.name,
          };
        }) ?? []),
      ]
    : []
);

watch($route, (v) => {
  menuKey.value = v.name;
});

const dropdownOptions: DropdownOption[] = [
  { label: t("nav.users.personal-center"), key: "users.personal-center", icon: renderIcon(AccountCircleRound) },
  { label: t("actions.sign-out"), key: "auth.sign-out", icon: renderIcon(LogOutRound) },
];

function dropdownHandler(key: string) {
  $router.push({ name: key });
}
</script>

<style>
.n-layout-header {
  height: 72px;
}

.n-layout-footer {
  padding: 24px;
}

.h-max {
  height: calc(100vh - 72px);
}

.w-dialog {
  width: 520px;
}
</style>
