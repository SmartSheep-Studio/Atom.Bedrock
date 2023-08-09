import defaultBanner from "@/assets/placeholders/default-banner.jpg"
import defaultAvatar from "@/assets/placeholders/default-avatar.png"

export const placeholders: { [id: string]: any } = {
  description: "No description yet",
  banner: defaultBanner,
  avatar: defaultAvatar,
}

export function usePlaceholder(id: string, val?: string): any {
  if (val == null || val.length <= 0) {
    return placeholders[id]
  } else {
    return val
  }
}
