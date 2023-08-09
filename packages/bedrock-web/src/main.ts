import "./assets/style.css";

import { createApp } from "vue";
import { createPinia } from "pinia";

import WujieVue from "wujie-vue3";

import wrapper from "./app.vue";
import router from "./router";
import i18n from "./i18n";

import "vfonts/Lato.css";
import "vfonts/FiraCode.css";

import "@mdi/font/css/materialdesignicons.min.css";

const app = createApp(wrapper);

app.use(createPinia());
app.use(router);
app.use(i18n);
app.use(WujieVue);

app.mount("#app");
