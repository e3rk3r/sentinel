import type { ConnectionState } from '@/types'
import { connectionDotClass, connectionLabel } from '@/lib/connection'
import { TooltipHelper } from '@/components/TooltipHelper'
import { cn } from '@/lib/utils'

type ConnectionBadgeProps = {
  state: ConnectionState
  detail?: string
  onClick?: () => void
}

export default function ConnectionBadge({
  state,
  detail,
  onClick,
}: ConnectionBadgeProps) {
  const label = connectionLabel(state)
  const base = detail && state !== 'connected' ? `${label} — ${detail}` : label
  const tooltip = onClick ? `${base} — click to resync` : base
  return (
    <TooltipHelper content={tooltip}>
      <span
        className={cn(
          'inline-flex h-4 w-4 items-center justify-center rounded-full border border-border-subtle bg-surface-elevated',
          onClick && 'cursor-pointer hover:bg-surface-active',
        )}
        role={onClick ? 'button' : 'status'}
        tabIndex={onClick ? 0 : undefined}
        aria-label={label}
        onClick={onClick}
        onKeyDown={
          onClick
            ? (e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault()
                  onClick()
                }
              }
            : undefined
        }
      >
        <span
          className={`inline-block h-2 w-2 rounded-full ${connectionDotClass(state)}`}
        />
      </span>
    </TooltipHelper>
  )
}
