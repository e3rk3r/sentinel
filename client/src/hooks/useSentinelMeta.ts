import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'

type MetaResponse = {
  tokenRequired?: boolean
  defaultCwd?: string
  version?: string
  timezone?: string
  locale?: string
  hostname?: string
}

export function useSentinelMeta() {
  const metaQuery = useQuery({
    queryKey: ['meta'],
    queryFn: async ({ signal }) => {
      const response = await fetch('/api/meta', {
        signal,
        headers: { Accept: 'application/json' },
        credentials: 'same-origin',
      })
      if (response.status === 401) {
        return {
          tokenRequired: true,
          defaultCwd: '',
          version: 'dev',
          timezone: 'UTC',
          locale: '',
          hostname: '',
          unauthorized: true,
        }
      }
      if (!response.ok) {
        throw new Error(`meta fetch failed: HTTP ${response.status}`)
      }

      const payload = (await response.json()) as { data?: MetaResponse }
      return {
        tokenRequired: Boolean(payload.data?.tokenRequired),
        defaultCwd: (payload.data?.defaultCwd ?? '').trim(),
        version: (payload.data?.version ?? 'dev').trim() || 'dev',
        timezone: (payload.data?.timezone ?? 'UTC').trim() || 'UTC',
        locale: (payload.data?.locale ?? '').trim(),
        hostname: (payload.data?.hostname ?? '').trim(),
        unauthorized: false,
      }
    },
    retry: 2,
    staleTime: 60_000,
  })

  const value = useMemo(
    () =>
      metaQuery.data ?? {
        tokenRequired: false,
        defaultCwd: '',
        version: 'dev',
        timezone: 'UTC',
        locale: '',
        hostname: '',
        unauthorized: false,
      },
    [metaQuery.data],
  )

  return {
    tokenRequired: value.tokenRequired,
    defaultCwd: value.defaultCwd,
    version: value.version,
    timezone: value.timezone,
    locale: value.locale,
    hostname: value.hostname,
    unauthorized: value.unauthorized,
    loaded: metaQuery.isFetched,
  }
}
