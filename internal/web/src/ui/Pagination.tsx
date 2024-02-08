// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/pagination
// https://ui.shadcn.com/docs/components/pagination
import { ComponentProps, JSX, mergeProps, splitProps } from "solid-js"
import { Pagination } from "@kobalte/core"
import { RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiSystemMoreLine } from "solid-icons/ri"

import { cn } from "~/lib/utils"
import { ButtonProps, buttonVariants } from "./Button"

export function PaginationRoot(props: ComponentProps<typeof Pagination.Root>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Pagination.Root
    class={cn("mx-auto flex w-full justify-center", props.class)}
    {...rest}
  />
}

export function PaginationContent(props: JSX.HTMLAttributes<HTMLUListElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <ul
    class={cn("flex flex-row items-center gap-1", props.class)}
    {...rest}
  />
}

export const PaginationItem = Pagination.Item
export const PaginationItems = Pagination.Items

type PaginationLinkProps = {
  isActive?: boolean
} & Pick<ButtonProps, "size"> &
  JSX.ButtonHTMLAttributes<HTMLButtonElement>

export function PaginationLink(props: PaginationLinkProps) {
  const [_, rest] = splitProps(mergeProps({ size: "icon" }, props), ["class", "isActive", "size"])
  return <button
    class={cn(
      buttonVariants({
        variant: props.isActive ? "outline" : "ghost",
        size: props.size as any,
      }), props.class)}
    {...rest}
  />
}

type PaginationNextPreviousProps = Pick<ButtonProps, "size"> & JSX.ButtonHTMLAttributes<HTMLButtonElement>

export function PaginationPrevious(props: PaginationNextPreviousProps) {
  const [_, rest] = splitProps(mergeProps({ size: "default" }, props), ["class", "size"])
  return <Pagination.Previous
    aria-label="Go to previous page"
    class={cn(
      buttonVariants({
        variant: "ghost",
        size: props.size as any,
      }),
      "gap-1 pl-2.5",
      props.class)}
    {...rest}
  >
    <RiArrowsArrowLeftSLine class="h-4 w-4" />
    <span>Previous</span>
  </Pagination.Previous>
}

export function PaginationNext(props: PaginationNextPreviousProps) {
  const [_, rest] = splitProps(mergeProps({ size: "default" }, props), ["class", "size"])
  return <Pagination.Next
    aria-label="Go to next page"
    class={cn(
      buttonVariants({
        variant: "ghost",
        size: props.size as any,
      }),
      "gap-1 pl-2.5",
      props.class)}
    {...rest}
  >
    <span>Next</span>
    <RiArrowsArrowRightSLine class="h-4 w-4" />
  </Pagination.Next>
}

export function PaginationEllipsis(props: ComponentProps<typeof Pagination.Ellipsis>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Pagination.Ellipsis
    aria-hidden
    class={cn("flex h-9 w-9 items-center justify-center", props.class)}
    {...rest}
  >
    <RiSystemMoreLine class="h-4 w-4" />
    <span class="sr-only">More pages</span>
  </Pagination.Ellipsis>
}
