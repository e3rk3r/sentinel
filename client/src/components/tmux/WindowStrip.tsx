import { useMemo } from 'react'
import { ChevronDown, Plus, X } from 'lucide-react'
import { getSessionIcon } from '@/components/sidebar/sessionIcons'
import type { TmuxLauncher, WindowInfo } from '@/types'
import { Button } from '@/components/ui/button'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from '@/components/ui/context-menu'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { TooltipHelper } from '@/components/TooltipHelper'
import { cn } from '@/lib/utils'
import { useIsMobileLayout } from '@/hooks/useIsMobileLayout'

type WindowStripProps = {
  hasActiveSession: boolean
  inspectorLoading: boolean
  inspectorError: string
  windows: Array<WindowInfo>
  activeWindowIndex: number | null
  launchers: Array<TmuxLauncher>
  recentLauncher: TmuxLauncher | null
  onSelectWindow: (windowIndex: number) => void
  onCloseWindow: (windowIndex: number) => void
  onRenameWindow: (windowInfo: WindowInfo) => void
  onCreateWindow: () => void
  onLaunchLauncher: (launcherID: string) => void
  onOpenLaunchers: () => void
}

export default function WindowStrip({
  hasActiveSession,
  inspectorLoading,
  inspectorError,
  windows,
  activeWindowIndex,
  launchers,
  recentLauncher,
  onSelectWindow,
  onCloseWindow,
  onRenameWindow,
  onCreateWindow,
  onLaunchLauncher,
  onOpenLaunchers,
}: WindowStripProps) {
  const isMobile = useIsMobileLayout()
  const sortedWindows = useMemo(
    () => [...windows].sort((left, right) => left.index - right.index),
    [windows],
  )
  const secondaryLaunchers = useMemo(
    () =>
      recentLauncher === null
        ? launchers
        : launchers.filter((launcher) => launcher.id !== recentLauncher.id),
    [launchers, recentLauncher],
  )
  const stripClass = 'flex min-h-[24px] items-center gap-1.5 overflow-x-auto'

  if (!hasActiveSession) {
    return (
      <div className={stripClass}>
        <span className="truncate text-[11px] text-secondary-foreground">
          Select and attach a session.
        </span>
      </div>
    )
  }
  if (inspectorLoading) {
    return (
      <div className={stripClass} aria-busy="true" aria-live="polite">
        <div className="h-6 w-6 shrink-0 rounded border border-border-subtle bg-surface-elevated motion-safe:animate-pulse" />
        <div className="h-5 w-20 shrink-0 rounded border border-border-subtle bg-surface-elevated motion-safe:animate-pulse" />
        <div className="h-5 w-24 shrink-0 rounded border border-border-subtle bg-surface-elevated motion-safe:animate-pulse" />
        <span className="sr-only">Loading windows</span>
      </div>
    )
  }
  if (inspectorError) {
    return (
      <div className={stripClass}>
        <span className="truncate text-[11px] text-destructive-foreground">
          {inspectorError}
        </span>
      </div>
    )
  }

  return (
    <div className={stripClass}>
      <div className="flex shrink-0 items-center">
        <TooltipHelper content="Create blank window">
          <Button
            variant="outline"
            size="icon-sm"
            className="rounded-r-none border-r-0"
            onClick={onCreateWindow}
            aria-label="Create blank window"
          >
            <Plus className="h-4 w-4" />
          </Button>
        </TooltipHelper>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="outline"
              size="icon-sm"
              className="rounded-l-none px-1.5"
              aria-label="Open launcher menu"
            >
              <ChevronDown className="h-3.5 w-3.5" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" className="w-56">
            <DropdownMenuItem onSelect={onCreateWindow}>
              <Plus className="h-3.5 w-3.5" />
              New blank window
            </DropdownMenuItem>
            {recentLauncher !== null && (
              <>
                <DropdownMenuSeparator />
                <DropdownMenuLabel>Last used</DropdownMenuLabel>
                <DropdownMenuItem
                  onSelect={() => onLaunchLauncher(recentLauncher.id)}
                >
                  {(() => {
                    const Icon = getSessionIcon(recentLauncher.icon)
                    return <Icon className="h-3.5 w-3.5" />
                  })()}
                  <span className="flex min-w-0 flex-1 items-center gap-2">
                    <span className="truncate">{recentLauncher.name}</span>
                    <span className="truncate text-[10px] text-muted-foreground">
                      {recentLauncher.command}
                    </span>
                  </span>
                </DropdownMenuItem>
              </>
            )}
            {secondaryLaunchers.length > 0 && (
              <>
                <DropdownMenuSeparator />
                <DropdownMenuLabel>Launchers</DropdownMenuLabel>
                {secondaryLaunchers.map((launcher) => {
                  const Icon = getSessionIcon(launcher.icon)
                  return (
                    <DropdownMenuItem
                      key={launcher.id}
                      onSelect={() => onLaunchLauncher(launcher.id)}
                    >
                      <Icon className="h-3.5 w-3.5" />
                      <span className="flex min-w-0 flex-1 items-center gap-2">
                        <span className="truncate">{launcher.name}</span>
                        <span className="truncate text-[10px] text-muted-foreground">
                          {launcher.command}
                        </span>
                      </span>
                    </DropdownMenuItem>
                  )
                })}
              </>
            )}
            <DropdownMenuSeparator />
            <DropdownMenuItem onSelect={onOpenLaunchers}>
              Manage launchers...
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {sortedWindows.length === 0 && (
        <span className="truncate">No windows found for this session.</span>
      )}
      {sortedWindows.map((windowInfo) => {
        const isActive = activeWindowIndex === windowInfo.index
        const unreadPanes = windowInfo.unreadPanes ?? 0
        const hasUnread = windowInfo.hasUnread ?? unreadPanes > 0
        return (
          <ContextMenu key={`${windowInfo.session}:${windowInfo.index}`}>
            <ContextMenuTrigger asChild>
              <div
                className={cn(
                  'inline-flex shrink-0 items-center overflow-hidden rounded border text-[11px]',
                  isActive
                    ? 'border-primary/50 text-primary-text'
                    : hasUnread
                      ? 'border-amber-400/60 text-amber-100'
                      : 'border-border text-secondary-foreground',
                )}
              >
                <button
                  className="inline-flex cursor-pointer items-center gap-1 px-1.5 py-0.5 whitespace-nowrap hover:text-foreground"
                  type="button"
                  onClick={() => onSelectWindow(windowInfo.index)}
                  aria-label={
                    isMobile ? `Select window ${windowInfo.name}` : undefined
                  }
                >
                  {isMobile ? windowInfo.index : windowInfo.name}
                </button>
                {!isMobile && (
                  <button
                    className="grid h-5 w-5 cursor-pointer place-items-center border-l border-border-subtle text-secondary-foreground hover:bg-surface-close-hover hover:text-destructive-foreground"
                    type="button"
                    onClick={() => onCloseWindow(windowInfo.index)}
                    aria-label={`Close window #${windowInfo.index}`}
                  >
                    <X className="h-3 w-3" />
                  </button>
                )}
              </div>
            </ContextMenuTrigger>
            <ContextMenuContent className="w-44">
              <ContextMenuItem onSelect={() => onRenameWindow(windowInfo)}>
                Rename window
              </ContextMenuItem>
              <ContextMenuItem
                className="text-destructive-foreground focus:text-destructive-foreground"
                onSelect={() => onCloseWindow(windowInfo.index)}
              >
                Close window
              </ContextMenuItem>
            </ContextMenuContent>
          </ContextMenu>
        )
      })}
    </div>
  )
}
