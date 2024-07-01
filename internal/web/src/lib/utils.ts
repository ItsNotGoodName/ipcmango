import {
  Accessor,
  Setter,
  batch,
  createEffect,
  createMemo,
  createSignal,
  onCleanup,
} from "solid-js";
import { type ClassValue, clsx } from "clsx";
import { createStore } from "solid-js/store";
import { createDateNow, createTimeDifference } from "@solid-primitives/date";
import { Location, Params, useSearchParams } from "@solidjs/router";

export function cn(...inputs: ClassValue[]) {
  return clsx(inputs);
}

export function formatDate(value: Date): string {
  return value.toLocaleString();
}

export function parseDate(value: Date | string | undefined): Date {
  return new Date(value || "");
}

type RowSelectionItem<T> = {
  id: T;
  checked: boolean;
  disabled: boolean;
};

export type CreateRowSelectionReturn<T> = {
  items: Array<RowSelectionItem<T> | undefined>;
  multiple: Accessor<boolean>;
  all: Accessor<boolean>;
  selections: Accessor<Array<T>>;
  set: (id: T, value: boolean) => void;
  setAll: (value: boolean) => void;
};

export function createRowSelection<T>(
  ids: Accessor<Array<{ id: T; disabled?: boolean }>>,
): CreateRowSelectionReturn<T> {
  const [items, setItems] = createStore<Array<RowSelectionItem<T>>>(
    ids().map((v) => ({
      id: v.id,
      checked: false,
      disabled: v.disabled || false,
    })),
  );
  createEffect(() =>
    setItems((prev) =>
      ids().map((v) => ({
        id: v.id,
        disabled: v.disabled || false,
        checked: prev.find((p) => p.id == v.id)?.checked || false,
      })),
    ),
  );

  return {
    items,
    multiple: () => {
      for (let index = 0; index < items.length; index++) {
        if (items[index].checked) return true;
      }
      return false;
    },
    all: () => {
      let disabled = 0;
      for (let index = 0; index < items.length; index++) {
        if (items[index].disabled) disabled++;
        else if (!items[index].checked) return false;
      }
      if (items.length - disabled == 0) return false;
      return true;
    },
    selections: () => items.filter((v) => v.checked == true).map((v) => v.id),
    set: (id, value) => {
      setItems((v) => v.id === id && !v.disabled, "checked", value);
    },
    setAll: (value) => {
      setItems((v) => !v.disabled, "checked", value);
    },
  };
}

type CreateValueDialogReturn<T> = {
  open: Accessor<boolean>;
  value: Accessor<T>;
  setClose: () => void;
  setValue: Setter<T>;
};

export function createValueDialog<T>(
  defaultValue: T,
): CreateValueDialogReturn<T> {
  const [open, setOpen] = createSignal(false);
  const [value, setValue] = createSignal(defaultValue);
  return {
    open,
    value,
    setClose: () => setOpen(false),
    setValue: (...args) =>
      batch(() => {
        setOpen(true);
        // @ts-ignore
        return setValue(...args);
      }),
  };
}

export function createFormToggle(formOpen: boolean, onOpen: () => void) {
  const [open, setOpen] = createSignal(formOpen);
  return {
    open,
    setOpen: () => {
      if (open() == true) return;
      batch(() => {
        setOpen(true);
        onOpen();
      });
    },
    setClose: () => setOpen(false),
  };
}

export function relativeWSURL(uri: string): string {
  return `${window.location.protocol === "https:" ? "wss:" : "ws:"}//${window.location.host}${uri}`;
}

export function useHiddenScrollbar(): void {
  const html = document.getElementsByTagName("html")[0];
  if (html.style.getPropertyValue("scrollbar-width") == "none") return;
  html.style.setProperty("scrollbar-width", "none");
  onCleanup(() => html.style.removeProperty("scrollbar-width"));
}

export function validationState(
  error?: string | boolean | null | Error,
): "invalid" | "valid" {
  return error ? "invalid" : "valid";
}

export function createUptime(date: Accessor<Date>) {
  const [now, update] = createDateNow(() => false);
  const [difference] = createTimeDifference(date, now);
  const timer = setInterval(update, 1000);
  onCleanup(() => clearInterval(timer));

  return createMemo(() => {
    const total = difference() / 1000;
    const days = Math.floor(total / 86400);
    const hours = Math.floor((total % 86400) / 3600);
    const minutes = Math.floor((total % 3600) / 60);
    const seconds = Math.floor(total % 60);
    return {
      days,
      hasDays: days > 0,
      hours,
      hasHours: hours > 0 || days > 0,
      minutes,
      hasMinutes: minutes > 0 || hours > 0 || days > 0,
      seconds,
    };
  });
}

export function useQueryFilter(key: string) {
  const [searchParams, setSearchParams] = useSearchParams();

  const values: Accessor<string[]> = createMemo(
    () => searchParams[key]?.split(".") || [],
  );
  const setValues = (value: string[]) =>
    setSearchParams({ [key]: value.join(".") });

  return {
    values,
    setValues,
  };
}

export function useQueryBoolean(key: string, defaultValue?: boolean) {
  const [searchParams, setSearchParams] = useSearchParams();

  const value = () =>
    searchParams[key] == undefined ? defaultValue : searchParams[key] == "true";
  const setValue = (value: boolean) =>
    setSearchParams({
      [key]: value != undefined ? String(value) : defaultValue,
    });

  return {
    value,
    setValue,
  };
}

export function useQueryNumber(key: string, defaultValue?: number) {
  const [searchParams, setSearchParams] = useSearchParams();

  const value = () => Number(searchParams[key] ?? defaultValue);
  const setValue = (value: number) => {
    setSearchParams({
      [key]: value != undefined ? String(value) : defaultValue,
    });
  };

  return {
    value,
    setValue,
  };
}

export type Sort = {
  field?: string;
  order?: "ascending" | "descending";
};

type QuerySortReturn = {
  value: Accessor<Sort>;
  setValue: (value: Sort) => void;
};

export function useQuerySort(
  key: string,
  defaultValue?: Sort,
): QuerySortReturn {
  const [searchParams, setSearchParams] = useSearchParams();

  const value = (): Sort => {
    const value = searchParams[key];

    if (value == undefined || value.length == 0) return defaultValue || {};
    if (value[0] == "-") {
      return { order: "descending", field: value.slice(1) };
    }
    return { order: "ascending", field: value };
  };
  const setValue = (sort: Sort) => {
    const prefix = sort.order == "descending" ? "-" : "";
    const suffix = sort.field ?? "";
    setSearchParams({ [key]: prefix + suffix });
  };

  return {
    value,
    setValue,
  };
}

export function useQueryString(key: string, defaultValue?: string) {
  const [searchParams, setSearchParams] = useSearchParams();

  const value = () =>
    searchParams[key] == undefined ? defaultValue : searchParams[key];
  const setValue = (value: string) =>
    setSearchParams({
      [key]: value != undefined ? value : defaultValue,
    });

  return {
    value,
    setValue,
  };
}

export type RouteLoadFuncArgs<T extends Params = Params, S = unknown> = {
  params: T;
  location: Location<S>;
  intent: "initial" | "navigate" | "native" | "preload";
};
