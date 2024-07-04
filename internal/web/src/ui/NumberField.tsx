// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/number-field
// https://ui.shadcn.com/docs/components/number-field
import { ComponentProps, splitProps } from "solid-js";
import { NumberField } from "@kobalte/core/number-field";

import { cn } from "~/lib/utils";
import { labelVariants } from "./Label";

export function NumberFieldInput(
  props: ComponentProps<typeof NumberField.Input>,
) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <NumberField.Input
      class={cn(
        "flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
        props.class,
      )}
      {...rest}
    />
  );
}

export const NumberFieldRoot = NumberField;

export function NumberFieldLabel(
  props: ComponentProps<typeof NumberField.Label>,
) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <NumberField.Label class={cn(labelVariants(), props.class)} {...rest} />
  );
}

export function NumberFieldDescription(
  props: ComponentProps<typeof NumberField.Description>,
) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <NumberField.Description
      class={cn("text-sm text-muted-foreground", props.class)}
      {...rest}
    />
  );
}

export function NumberFieldErrorMessage(
  props: ComponentProps<typeof NumberField.ErrorMessage>,
) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <NumberField.ErrorMessage
      class={cn("text-sm font-medium text-destructive")}
      {...rest}
    />
  );
}
