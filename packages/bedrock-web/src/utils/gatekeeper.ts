import { type RouteLocationNormalizedLoaded } from "vue-router";
import { usePrincipal } from "@/stores/account";
import wildcard from "wildcard-match";

export function hasUserPermissions(...requires: string[]) {
  const $principal = usePrincipal();
  if (!$principal.isSigned || $principal.account == null) {
    return false;
  }

  for (const require of requires) {
    let passed = false;
    for (const [k, v] of Object.entries($principal.account.permissions) ?? []) {
      if (v == true && wildcard(k)(require)) {
        passed = true;
        break;
      }
    }

    if (!passed) {
      return false;
    }
  }

  return true;
}

export function keepGate(to: RouteLocationNormalizedLoaded) {
  const $principal = usePrincipal();
  const meta: any = to?.meta?.gatekeeper ?? {};

  if (meta?.must === true && !$principal.isSigned) {
    return false;
  } else if (meta?.permissions != null && !hasUserPermissions(...meta?.permissions)) {
    return false;
  }

  return true;
}
