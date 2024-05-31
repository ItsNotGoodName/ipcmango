import { cva } from "class-variance-authority"
import { BiRegularCctv } from "solid-icons/bi";
import { ErrorBoundary, Suspense, batch, createSignal, } from "solid-js";
import { A, RouteSectionProps } from "@solidjs/router";
import { RiDocumentFileLine, RiBuildingsHomeLine, RiSystemMenuLine, RiWeatherFlashlightLine, RiBusinessMailLine, RiDocumentBookLine, } from "solid-icons/ri";
import { Portal } from "solid-js/web";
import { makePersisted } from "@solid-primitives/storage";

import { ThemeIcon } from "~/ui/ThemeIcon";
import { toggleTheme, useThemeTitle } from "~/ui/theme";
import { ToastList, ToastRegion } from "~/ui/Toast";
import { cn } from "~/lib/utils";
import { PageError, PageLoading } from "~/ui/Page";
import { SheetContent, SheetHeader, SheetOverflow, SheetRoot, SheetTitle } from "./ui/Sheet";

const menuLinkVariants = cva("ui-disabled:pointer-events-none ui-disabled:opacity-50 relative flex cursor-pointer select-none items-center gap-1 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors", {
  variants: {
    size: {
      default: "h-10 px-4 py-2",
      icon: "h-10 w-10",
    },
    variant: {
      default: "hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground",
      active: "bg-primary text-primary-foreground hover:bg-primary/90 focus:bg-primary/90",
    }
  },
  defaultVariants: {
    variant: "default"
  }
})

function MenuLinks(props: { onClick?: () => void }) {
  return (
    <div class="flex flex-col">
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/" noScroll end>
        <RiBuildingsHomeLine class="size-5" />Home
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/devices" noScroll>
        <BiRegularCctv class="size-5" />Devices
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/emails" noScroll>
        <RiBusinessMailLine class="size-5" />Emails
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/events" noScroll>
        <RiWeatherFlashlightLine class="size-5" />Events
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/files" noScroll>
        <RiDocumentFileLine class="size-5" />Files<div class="flex flex-1 justify-end"><div>🚧</div></div>
      </A>
      <a class={menuLinkVariants()} onClick={props.onClick}
        href="/docs" noScroll>
        <RiDocumentBookLine class="size-5" />Docs
      </a>
    </div>
  )
}

type HeaderProps = {
  onMenuClick: () => void
  onMobileMenuClick: () => void
}

function Header(props: HeaderProps) {
  return (
    <div class="overflow-x-hidden z-10 w-full h-12 border-b bg-background text-foreground border-b-border">
      <div class="flex gap-1 items-center px-1 h-full">
        <div onClick={props.onMobileMenuClick} title="Menu" class={cn(menuLinkVariants(), "md:hidden")}>
          <RiSystemMenuLine class="size-6" />
        </div>
        <button onClick={props.onMenuClick} title="Menu" class={cn(menuLinkVariants(), "hidden md:inline-flex")}>
          <RiSystemMenuLine class="size-6" />
        </button>
        <div class="flex flex-1 gap-2 items-baseline truncate">
          <A href="/" class="flex items-center text-xl">
            IPCManView
          </A>
        </div>
        <div class="flex gap-1">
          <button class={menuLinkVariants({ size: "icon" })} onClick={toggleTheme} title={useThemeTitle()}>
            <ThemeIcon class="size-6" />
          </button>
        </div>
      </div>
    </div>
  )
}

function createMenu() {
  const [mobileOpen, setMobileOpen] = createSignal(false)
  const toggleMobileOpen = () => setMobileOpen(!mobileOpen())
  const closeMobile = () => setMobileOpen(false)

  const [open, setOpen] = makePersisted(createSignal(true), { "name": "menu-open" })
  const toggleOpen = () => {
    if (open()) {
      batch(() => {
        setOpen(false)
        setMobileOpen(false)
      })
    } else {
      setOpen(true)
    }
  }

  return {
    mobileOpen,
    toggleMobileOpen,
    closeMobile,
    open,
    toggleOpen,
  }
}

export function Root(props: RouteSectionProps) {
  const menu = createMenu()

  return (
    <ErrorBoundary fallback={(e) =>
      <div class="p-4">
        <PageError error={e} />
      </div>
    }>
      <Suspense fallback={<PageLoading class="pt-10" />}>
        <Portal>
          <ToastRegion class="top-12 sm:top-12">
            <ToastList class="top-12 sm:top-12" />
          </ToastRegion>
        </Portal>
        <SheetRoot open={menu.mobileOpen()} onOpenChange={menu.toggleMobileOpen}>
          <SheetContent side="left" class="p-2">
            <SheetHeader class="px-2 sm:pt-2">
              <SheetTitle>IPCManView</SheetTitle>
            </SheetHeader>
            <SheetOverflow class="pb-2">
              <MenuLinks onClick={menu.closeMobile} />
            </SheetOverflow>
          </SheetContent>
        </SheetRoot>
        <Header
          onMenuClick={menu.toggleOpen}
          onMobileMenuClick={menu.toggleMobileOpen}
        />
        <div class="flex">
          <div data-open={menu.open()} class="border-border w-0 shrink-0 border-r-0 transition-all duration-300 md:data-[open=true]:w-48 md:data-[open=true]:border-r">
            <div class="overflow-y-auto sticky top-0 max-h-screen overflow-x-clip">
              <div class="p-2">
                <MenuLinks />
              </div>
            </div>
          </div>
          <div class="min-h-[calc(100vh-3rem)] w-full overflow-x-clip">
            {props.children}
          </div>
        </div>
      </Suspense>
    </ErrorBoundary>
  )
}

