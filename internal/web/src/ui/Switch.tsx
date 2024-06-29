// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/switch
// https://ui.shadcn.com/docs/components/switch
import { Switch } from "@kobalte/core/switch";
import { ComponentProps, splitProps } from "solid-js";

import { cn } from "~/lib/utils";
import { labelVariants } from "./Label";

export const SwitchRoot = Switch;

export function SwitchLabel(props: ComponentProps<typeof Switch.Label>) {
  const [_, rest] = splitProps(props, ["class"]);
  return <Switch.Label class={cn(labelVariants(), props.class)} {...rest} />;
}

export function SwitchDescription(
  props: ComponentProps<typeof Switch.Description>,
) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Switch.Description
      class={cn("text-sm text-muted-foreground", props.class)}
      {...rest}
    />
  );
}

export function SwitchErrorMessage(
  props: ComponentProps<typeof Switch.ErrorMessage>,
) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Switch.ErrorMessage
      class={cn("text-sm font-medium text-destructive", props.class)}
      {...rest}
    />
  );
}

export function SwitchControl(
  props: Omit<ComponentProps<typeof Switch.Control>, "children"> & {
    inputProps?: Omit<ComponentProps<typeof Switch.Input>, "class">;
  },
) {
  const [_, rest] = splitProps(props, ["class", "inputProps"]);
  return (
    <>
      <Switch.Input class="peer" {...props.inputProps} />
      <Switch.Control
        class={cn(
          "inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors peer-focus-visible:outline-none peer-focus-visible:ring-2 peer-focus-visible:ring-ring peer-focus-visible:ring-offset-2 peer-focus-visible:ring-offset-background ui-disabled:cursor-not-allowed ui-disabled:opacity-50 ui-checked:bg-primary ui-not-checked:bg-input",
          props.class,
        )}
        {...rest}
      >
        <Switch.Thumb
          class={cn(
            "pointer-events-none block h-5 w-5 rounded-full bg-background shadow-lg ring-0 transition-transform ui-checked:translate-x-5 ui-not-checked:translate-x-0",
          )}
        />
      </Switch.Control>
    </>
  );
}
