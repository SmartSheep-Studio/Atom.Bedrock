import {createI18n} from "vue-i18n"
import zhCN from "@/i18n/zh-cn.json"
import zhHant from "@/i18n/zh-hant.json"
import en from "@/i18n/en.json"

const i18n = createI18n({
    locale: localStorage.getItem("locale") || navigator.language,
    fallbackLocale: "en-US",
    messages: {
        "zh-CN": zhCN,
        "zh-Hant": zhHant,
        "en-US": en,
        "en-UK": en,
    }
})

export default i18n