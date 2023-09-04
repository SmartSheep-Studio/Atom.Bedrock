import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "landing",
      component: () => import("@/views/landing.vue"),
    },

    {
      path: "/auth/sign-in",
      name: "auth.sign-in",
      component: () => import("@/views/auth/sign-in.vue"),
    },
    {
      path: "/auth/sign-up",
      name: "auth.sign-up",
      component: () => import("@/views/auth/sign-up.vue"),
    },
    {
      path: "/auth/sign-out",
      name: "auth.sign-out",
      component: () => import("@/views/auth/sign-out.vue"),
      meta: { gatekeeper: { must: true } },
    },
    {
      path: "/auth/openid/connect",
      name: "auth.openid.connect",
      component: () => import("@/views/auth/openid/connect.vue"),
    },

    {
      path: "/users/self",
      name: "users.personal-center",
      component: () => import("@/views/users/account.vue"),
      meta: { gatekeeper: { must: true } },
    },
    {
      path: "/users/self/notifications",
      name: "users.notifications",
      component: () => import("@/views/users/notifications.vue"),
      meta: { gatekeeper: { must: true } },
    },

    {
      path: "/administration",
      name: "administration",
      component: () => import("@/views/administration/layout.vue"),
      redirect: { name: "administration.dashboard" },
      meta: { gatekeeper: { must: true, permissions: ["bedrock.admin.view"] } },
      children: [
        {
          path: "/administration",
          name: "administration.dashboard",
          component: () => import("@/views/administration/dashboard.vue"),
          meta: { gatekeeper: { must: true, permissions: ["bedrock.admin.view"] } },
        },
      ],
    },

    {
      path: "/svm/subapps/:id/:matchAll*",
      name: "framework.subapp",
      component: () => import("@/views/framework/subapp.vue"),
    },
  ],
});

export default router;
