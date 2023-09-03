<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-dialog-provider>
      <n-message-provider>
        <div class="w-full h-screen relative">
          <n-layout position="absolute" :has-sider="menuMode === 'width'">
            <n-layout-sider :collapse-mode="menuMode" :collapsed-width="menuMode === 'transform' ? 0 : 64"
                            :show-trigger="menuMode === 'transform' ? 'bar' : 'arrow-circle'" :width="280"
                            :collapsed="menuCollapsed"
                            :native-scrollbar="false"
                            :class="menuMode === 'transform' ? 'fixed z-100 shadow-xl' : undefined" bordered
                            @collapse="menuCollapsed = true" @expand="menuCollapsed = false">

              <div class="flex flex-col h-[100vh]">
                <div :class="menuCollapsed ? 'nav-item-collapsed' : 'nav-item-expand'"
                     class="nav-item pt-[8px] h-[42px] flex gap-2 items-center cursor-pointer"
                     @click="$router.push({ name: 'landing' })">

                  <img src="./assets/icon.png" alt="Logo" class="block brand-item-icon" />

                  <div v-if="!menuCollapsed">Atom</div>

                </div>

                <n-menu mode="vertical" :collapsed="menuCollapsed" :collapsed-width="64" :collapsed-icon-size="22"
                        :options="menuOptions" v-model:value="menuKey" />

                <div class="grow"></div>

                <div :class="menuCollapsed ? 'nav-item-collapsed' : 'nav-item-expand'"
                     class="nav-item ml-[-2px] pb-[8px] h-[42px] flex gap-2 items-center" v-if="!$principal.isSigned">
                  <n-dropdown placement="right-end" show-arrow :options="unsignedDropdownOptions"
                              @select="dropdownHandler">
                    <n-avatar size="medium" color="transparent">
                      <n-icon color="black" :component="SupervisorAccountRound" />
                    </n-avatar>
                  </n-dropdown>

                  <div v-if="!menuCollapsed">
                    <div class="ml-[2px]">Unsigned in</div>
                    <div class="text-xs text-gray-500 mt-[-4px]">401 Unauthorized</div>
                  </div>
                </div>

                <div :class="menuCollapsed ? 'nav-item-collapsed' : 'nav-item-expand'"
                     class="nav-item ml-[-2px] pb-[8px] h-[42px] flex gap-2 items-center" v-else>
                  <div class="flex gap-3">
                    <n-dropdown placement="right-end" show-arrow :options="signedDropdownOptions"
                                @select="dropdownHandler">
                      <n-avatar size="medium" color="transparent"
                                :src="usePlaceholder('avatar', $principal.account?.avatar_url)"></n-avatar>
                    </n-dropdown>

                    <div v-if="!menuCollapsed">
                      <div class="ml-[2px]">{{ $principal.account?.nickname }}</div>
                      <div class="text-xs text-gray-500 mt-[-4px]">@{{ $principal.account?.name }}</div>
                    </div>
                  </div>
                </div>
              </div>

            </n-layout-sider>

            <n-layout class="w-full h-full" :native-scrollbar="false">
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
import { usePlaceholder } from "@/utils/placeholders";
import { type Component, computed, h, onMounted, type Ref, ref, watch } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { type DropdownOption, type MenuOption, NIcon } from "naive-ui";
import {
  AccountCircleRound,
  EmailRound,
  LogInRound,
  LogOutRound,
  NewLabelRound,
  SupervisorAccountRound
} from "@vicons/material";
import { hasUserPermissions } from "@/utils/gatekeeper";
import { useI18n } from "vue-i18n";
import { useLocalStorage } from "@vueuse/core";

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
    primaryColorSuppl: "#A84141FF"
  }
};

const menuCollapsed = useLocalStorage("layout-nav-collapsed", false);
const menuKey = ref($route.name);
const menuMode = ref(window.innerWidth >= 768 ? "width" : "transform");
const menuOptions: Ref<MenuOption[]> = computed(() => {
  const items: MenuOption[] = [
    {
      label: () =>
        h(RouterLink, { to: { name: "landing" } }, { default: () => "Landing" }),
      icon: () =>
        h(NIcon,
          null,
          {
            default: () => h("span", { class: `mdi mdi-earth` }, null)
          }
        ),
      key: "landing"
    }
  ];

  const build = (list: any[]): any[] => {
    if (list.length === 0) {
      return [];
    }

    return list?.filter((v) => {
      const page = $endpoint.pages.filter((i) => {
        return i.name === v.name;
      })[0];

      if (v.children != null && v.children.length > 0) {
        return true;
      } else if (page == null) {
        return false;
      }

      if (page.meta != null) {
        if (page.meta.gatekeeper != null) {
          if (page.meta.gatekeeper?.must === true && !$principal.isSigned) {
            return false;
          } else if (page.meta.gatekeeper?.permissions != null && !hasUserPermissions(...page.meta.gatekeeper?.permissions)) {
            return false;
          }
        }
      }

      return true;
    })?.map((v) => {
      const page = $endpoint.pages.filter((i) => {
        return i.name === v.name;
      })[0];

      return {
        label: (
          page?.to == null
            ? () =>
              h("span", v.title)
            : () =>
              h(
                RouterLink,
                {
                  to: {
                    name: "framework.subapp",
                    params: { id: v.name }
                  }
                },
                { default: () => v.title }
              )
        ),
        icon: () =>
          h(NIcon,
            null,
            {
              default: () => h("span", { class: `mdi ${v.icon}` }, null)
            }
          ),
        children: v.children ? build(v.children) : undefined,
        key: v.name
      };
    }) ?? [];
  };

  items.push(...build($endpoint.nav));

  return items;
});

onMounted(() => {
  addEventListener("resize", () => {
    menuMode.value = window.innerWidth >= 768 ? "width" : "transform";
  });
});

watch($route, (v) => {
  menuKey.value = v.name;
});

const unsignedDropdownOptions: DropdownOption[] = [
  { label: t("actions.sign-in"), key: "auth.sign-in", icon: renderIcon(LogInRound) },
  { label: t("actions.sign-up"), key: "auth.sign-up", icon: renderIcon(NewLabelRound) }
];

const signedDropdownOptions: DropdownOption[] = [
  { label: t("nav.users.personal-center"), key: "users.personal-center", icon: renderIcon(AccountCircleRound) },
  { label: t("actions.sign-out"), key: "auth.sign-out", icon: renderIcon(LogOutRound) },
  { label: "Notifications", key: "users.notifications", icon: renderIcon(EmailRound) }
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

.w-dialog {
  width: 520px;
}

.nav-item,
.nav-item-icon,
.brand-item-icon {
  @apply transition-all ease-in-out delay-[.05s]
}

.nav-item-collapsed {
  @apply pl-[18px] pr-[18px]
}

.nav-item-expand {
  @apply pl-[32px] pr-[18px]
}

.nav-item-collapsed .nav-item-icon {
  @apply w-[28px] h-[28px]
}

.nav-item-expand .nav-item-icon {
  @apply w-[22px] h-[22px]
}

.nav-item-collapsed .brand-item-icon {
  @apply w-[28px] h-[28px]
}

.nav-item-expand .brand-item-icon {
  @apply w-[24px] h-[24px]
}
</style>
