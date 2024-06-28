// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/hover-card
// https://ui.shadcn.com/docs/components/hover-card
import { ComponentProps, splitProps } from "solid-js";
import { HoverCard } from "@kobalte/core/hover-card";

import { cn } from "~/lib/utils";

export const HoverCardRoot = HoverCard;
export const HoverCardTrigger = HoverCard.Trigger;
export const HoverCardArrow = HoverCard.Arrow;

export function HoverCardContent(
  props: ComponentProps<typeof HoverCard.Content>,
) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <HoverCard.Portal>
      <HoverCard.Content
        class={cn(
          "z-50 w-64 max-w-[var(--kb-popper-content-available-width)] origin-[var(--kb-hovercard-content-transform-origin)] rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-none ui-expanded:animate-in ui-expanded:fade-in-0 ui-expanded:zoom-in-95 ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-not-expanded:zoom-out-95",
          props.class,
        )}
        {...rest}
      />
    </HoverCard.Portal>
  );
}
