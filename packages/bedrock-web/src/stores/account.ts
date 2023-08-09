import { defineStore } from "pinia"
import { useCookies } from "@vueuse/integrations/useCookies"
import { http } from "@/utils/http"
import { computed, ref } from "vue"
import { useLocalStorage } from "@vueuse/core"

export const usePrincipal = defineStore("principal", () => {
  const isSigned = ref(false)
  const cookies = useCookies(["authorization"])
  const token = computed(() => cookies.get("authorization"))

  const session = ref<any>({})
  const account = useLocalStorage<any>("account-data", null, {
    deep: true,
    listenToStorageChanges: true,
    serializer: {
      read(v: any) {
        try {
          return JSON.parse(v)
        } catch {
          return null
        }
      },
      write(v: any) {
        if (v != null) {
          return JSON.stringify(v)
        } else {
          return "null"
        }
      }
    }
  })

  async function fetch() {
    if (cookies.get("authorization") != null) {
      try {
        const res = await http.get("/api/users/self")
        account.value = res.data.user
        session.value = res.data.session

        isSigned.value = true
      } catch {
        account.value = null
        isSigned.value = false
      }
    }
  }

  function logout() {
    cookies.remove("authorization")
    account.value = null
    window.location.reload()
  }

  return { account, session, isSigned, token, fetch, logout }
})
