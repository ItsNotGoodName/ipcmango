// # Changes
// - Content overflows in the Y direction
//
// # URLs
// https://kobalte.dev/docs/core/components/sheet
// https://ui.shadcn.com/docs/components/sheet
import { Dialog } from "@kobalte/core/dialog";
import { RiSystemCloseLine } from "solid-icons/ri";
import { cva, type VariantProps } from "class-variance-authority";
import { ComponentProps, JSX, splitProps } from "solid-js";

import { cn } from "~/lib/utils";
import { DialogOverlay, DialogPortal } from "./Dialog";

export const SheetRoot = Dialog;
export const SheetTrigger = Dialog.Trigger;
export const SheetCloseButton = Dialog.CloseButton;

export function SheetOverlay(props: ComponentProps<typeof Dialog.Overlay>) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Dialog.Overlay
      class={cn(
        "fixed inset-0 z-50 bg-black/80 ui-expanded:animate-in ui-expanded:fade-in-0 ui-not-expanded:animate-out ui-not-expanded:fade-out-0",
        props.class,
      )}
      {...rest}
    />
  );
}

const sheetContentVariants = cva(
  "fixed z-50 flex flex-col gap-2 bg-background p-4 shadow-lg transition ease-in-out ui-expanded:duration-500 ui-expanded:animate-in ui-not-expanded:duration-300 ui-not-expanded:animate-out sm:p-6",
  {
    variants: {
      side: {
        top: "inset-x-0 top-0 border-b ui-expanded:slide-in-from-top ui-not-expanded:slide-out-to-top",
        bottom:
          "inset-x-0 bottom-0 border-t ui-expanded:slide-in-from-bottom ui-not-expanded:slide-out-to-bottom",
        left: "inset-y-0 left-0 h-full w-3/4 border-r ui-expanded:slide-in-from-left ui-not-expanded:slide-out-to-left sm:max-w-sm",
        right:
          "inset-y-0 right-0 h-full w-3/4 border-l ui-expanded:slide-in-from-right ui-not-expanded:slide-out-to-right sm:max-w-sm",
      },
    },
    defaultVariants: {
      side: "right",
    },
  },
);

export function SheetContent(
  props: ComponentProps<typeof Dialog.Content> &
    VariantProps<typeof sheetContentVariants>,
) {
  const [_, rest] = splitProps(props, ["class", "side", "children"]);
  return (
    <DialogPortal>
      <DialogOverlay />
      <Dialog.Content
        class={cn(sheetContentVariants({ side: props.side }), props.class)}
        {...rest}
      >
        {props.children}
        <Dialog.CloseButton class="absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none ui-expanded:bg-secondary">
          <RiSystemCloseLine class="h-4 w-4" />
          <span class="sr-only">Close</span>
        </Dialog.CloseButton>
      </Dialog.Content>
    </DialogPortal>
  );
}

export function SheetHeader(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <div
      class={cn("flex flex-col text-center sm:text-left", props.class)}
      {...rest}
    />
  );
}

export function SheetOverflow(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"]);
  return <div class={cn("flex-grow overflow-y-auto", props.class)} {...rest} />;
}

export function SheetFooter(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <div
      class={cn(
        "flex flex-col-reverse gap-2 sm:flex-row sm:justify-end",
        props.class,
      )}
      {...rest}
    />
  );
}

export function SheetTitle(props: ComponentProps<typeof Dialog.Title>) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Dialog.Title
      class={cn("text-lg font-semibold text-foreground", props.class)}
      {...rest}
    />
  );
}

export function SheetDescription(
  props: ComponentProps<typeof Dialog.Description>,
) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Dialog.Description
      class={cn("text-sm text-muted-foreground", props.class)}
      {...rest}
    />
  );
}
