import { Accessor, Setter, batch, createEffect, createMemo, createSignal, onCleanup } from "solid-js";
import { type ClassValue, clsx } from "clsx"
import { createStore } from "solid-js/store";
import { createDateNow, createTimeDifference } from "@solid-primitives/date";

export function cn(...inputs: ClassValue[]) {
  return clsx(inputs)
}

export function formatDate(value: Date): string {
  return value.toLocaleString()
}

export function parseDate(value: Date | string | undefined): Date {
  return new Date(value || "")
}

type RowSelectionItem<T> = {
  id: T,
  checked: boolean,
  disabled: boolean
}

export type CreateRowSelectionReturn<T> = {
  items: Array<RowSelectionItem<T> | undefined>
  multiple: Accessor<boolean>
  all: Accessor<boolean>
  selections: Accessor<Array<T>>
  set: (id: T, value: boolean) => void
  setAll: (value: boolean) => void
}

export function createRowSelection<T>(ids: Accessor<Array<{ id: T, disabled?: boolean }>>): CreateRowSelectionReturn<T> {
  const [items, setItems] = createStore<Array<RowSelectionItem<T>>>(
    ids().map(v => ({ id: v.id, checked: false, disabled: v.disabled || false }))
  )
  createEffect(() =>
    setItems((prev) => ids().map(v => ({ id: v.id, disabled: v.disabled || false, checked: prev.find(p => p.id == v.id)?.checked || false })))
  )

  return {
    items,
    multiple: () => {
      for (let index = 0; index < items.length; index++) {
        if (items[index].checked) return true
      }
      return false
    },
    all: () => {
      let disabled = 0
      for (let index = 0; index < items.length; index++) {
        if (items[index].disabled) disabled++
        else if (!items[index].checked) return false
      }
      if (items.length - disabled == 0) return false
      return true
    },
    selections: () => items.filter(v => v.checked == true).map(v => v.id),
    set: (id, value) => {
      setItems(
        (v) => v.id === id && !v.disabled,
        "checked",
        value,
      );
    },
    setAll: (value) => {
      setItems(
        (v) => !v.disabled,
        "checked",
        value,
      );
    }
  }
}

type CreateValueModalReturn<T> = {
  open: Accessor<boolean>
  value: Accessor<T>
  setClose: () => void
  setValue: Setter<T>
}

export function createModal<T>(value: T): CreateValueModalReturn<T> {
  const [getOpen, setOpen] = createSignal(false)
  const [getValue, setValue] = createSignal(value)
  return {
    open: getOpen,
    value: getValue,
    setClose: () => setOpen(false),
    setValue: (...args) => batch(() => {
      setOpen(true)
      // @ts-ignore
      return setValue(...args)
    })
  }
}

export function relativeWSURL(uri: string): string {
  return `${window.location.protocol === "https:" ? "wss:" : "ws:"}//${window.location.host}${uri}`
}

export function useHiddenScrollbar(): void {
  const html = document.getElementsByTagName("html")[0]
  if (html.style.getPropertyValue("scrollbar-width") == "none") return
  html.style.setProperty("scrollbar-width", "none")
  onCleanup(() => html.style.removeProperty("scrollbar-width"))
}

export function validationState(error?: string | boolean): "invalid" | "valid" {
  return error ? "invalid" : "valid"
}

export function createUptime(date: Accessor<Date>) {
  const [now, update] = createDateNow(() => false);
  const [difference] = createTimeDifference(date, now)
  const timer = setInterval(update, 1000)
  onCleanup(() => clearInterval(timer))

  return createMemo(() => {
    const total = difference() / 1000
    const days = Math.floor(total / 86400)
    const hours = Math.floor((total % 86400) / 3600)
    const minutes = Math.floor((total % 3600) / 60)
    const seconds = Math.floor(total % 60)
    return {
      days,
      hasDays: days > 0,
      hours,
      hasHours: hours > 0 || days > 0,
      minutes,
      hasMinutes: minutes > 0 || hours > 0 || days > 0,
      seconds,
    }
  })
}
