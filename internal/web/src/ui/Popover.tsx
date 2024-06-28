// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/popover
// https://ui.shadcn.com/docs/components/popover
import { Popover } from "@kobalte/core/popover";
import { RiSystemCloseLine } from "solid-icons/ri";
import { ComponentProps, splitProps } from "solid-js";

import { cn } from "~/lib/utils";

export const PopoverRoot = Popover;
export const PopoverTrigger = Popover.Trigger;
export const PopoverAnchor = Popover.Anchor;
export const PopoverPortal = Popover.Portal;

export function PopoverContent(props: ComponentProps<typeof Popover.Content>) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Popover.Content
      class={cn(
        "z-50 w-72 max-w-[var(--kb-popper-content-available-width)] origin-[var(--kb-menu-content-transform-origin)] rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-none ui-expanded:animate-in ui-expanded:fade-in-0 ui-expanded:zoom-in-95 ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-not-expanded:zoom-out-95",
        props.class,
      )}
      {...rest}
    />
  );
}

export const PopoverArrow = Popover.Arrow;
export const PopoverCloseButton = Popover.CloseButton;
export const PopoverCloseIcon = RiSystemCloseLine;
export const PopoverTitle = Popover.Title;
export const PopoverDescription = Popover.Description;
