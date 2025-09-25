import _ from "lodash"

declare const window: any

export function getViteUrl(key: string) {
    return window?.PX_BASE_URL?.[key] ?? import.meta.env[key]
}