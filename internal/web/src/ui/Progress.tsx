// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/progress
// https://ui.shadcn.com/docs/components/progress
import { Progress } from "@kobalte/core/progress";
import { ComponentProps, splitProps } from "solid-js";

import { cn } from "~/lib/utils";

export const ProgressRoot = Progress;
export const ProgressLabel = Progress.Label;
export const ProgressValueLabel = Progress.ValueLabel;

export function ProgressTrack(props: ComponentProps<typeof Progress.Track>) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Progress.Track
      class={cn(
        "relative h-4 w-full overflow-hidden rounded-full bg-secondary",
        props.class,
      )}
      {...rest}
    />
  );
}

export function ProgressFill(props: ComponentProps<typeof Progress.Fill>) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Progress.Fill
      class={cn(
        "h-full w-[var(--kb-progress-fill-width)] flex-1 bg-primary transition-all",
        props.class,
      )}
      {...rest}
    />
  );
}
