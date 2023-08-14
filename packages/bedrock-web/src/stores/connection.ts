import { defineStore } from "pinia"
import { http } from "@/utils/http"
import { reactive, ref } from "vue"

export const useEndpoint = defineStore("endpoint", () => {
  const isPrepared = ref(true)
  const firmware = reactive({ name: "Project Atom", version: "0" })
  const external = ref<any[]>([])
  const limit = ref<any>({})
  const info = ref<any>({})
  const pages = ref<any[]>([])
  const nav = ref<any[]>([])

    async function fetch() {
    try {
      const res = await http.get("/api/info")
      firmware.name = res.data.firmware
      firmware.version = res.data.firmware_version
      limit.value = res.data.limit
      info.value = res.data.manifest
      pages.value = res.data.pages
      nav.value = res.data.nav
    } catch (e: any) {
      isPrepared.value = e.response.status !== 503;

      throw e
    }

    document.title = info.value.name ?? "Project Atom"
  }

  return { isPrepared, pages, nav, external, firmware, limit, info, fetch }
})
