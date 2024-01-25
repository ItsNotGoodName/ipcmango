import { PartialMessage } from "@protobuf-ts/runtime";
import { Accessor, createSignal } from "solid-js";
import { Timestamp } from "~/twirp/google/protobuf/timestamp";
import { type ClassValue, clsx } from "clsx"
import { toast } from "~/ui/Toast";
import { RpcError } from "@protobuf-ts/runtime-rpc";
import { FormError } from "@modular-forms/solid";
import { Order, Sort } from "~/twirp/rpc";

export function cn(...inputs: ClassValue[]) {
  return clsx(inputs)
}

export function createLoading(fn: () => Promise<void>): [Accessor<boolean>, () => Promise<void>] {
  const [loading, setLoading] = createSignal(false)
  return [loading, () => {
    if (loading()) {
      return Promise.resolve()
    }
    setLoading(true)
    return fn().finally(() => setLoading(false))
  }]
}

export function parseDate(value: PartialMessage<Timestamp> | undefined): Date {
  return Timestamp.toDate(Timestamp.create(value))
}

export function formatDate(value: Date): string {
  return value.toLocaleString()
}

export function catchAsToast(e: Error) {
  toast.error("Error", e.message)
}

export function throwAsFormError(e: unknown) {
  if (e instanceof RpcError)
    // @ts-ignore
    throw new FormError(e.message, e.meta ?? {})
  if (e instanceof Error)
    throw new FormError(e.message)
  throw new FormError("Unknown error has occured.")
}

export type PageProps<T> = {
  params: Partial<T>
}

export const paginateOptions = [10, 25, 50, 100]

export function parseOrder(s?: string): Order {
  if (s == "desc")
    return Order.DESC
  if (s == "asc")
    return Order.ASC
  return Order.ORDER_UNSPECIFIED
}

export function encodeOrder(o: Order): string {
  if (o == Order.DESC)
    return "desc"
  if (o == Order.ASC)
    return "asc"
  return ""
}

export function nextSort(sort?: Sort, field?: string): { field?: string, order: Order } {
  if (field == sort?.field) {
    const order = ((sort?.order ?? Order.ORDER_UNSPECIFIED) + 1) % 3

    if (order == Order.ORDER_UNSPECIFIED) {
      return { field: undefined, order: Order.ORDER_UNSPECIFIED }
    }

    return { field: field, order: order }
  }

  return { field: field, order: Order.DESC }
}
