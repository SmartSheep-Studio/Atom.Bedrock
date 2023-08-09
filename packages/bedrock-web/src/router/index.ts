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
      name: "auth.sign-out",
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
      path: "/launch/:id",
      name: "framework.sub-app",
      component: () => import("@/views/framework/sub-app.vue"),
    },
  ],
});

export default router;
