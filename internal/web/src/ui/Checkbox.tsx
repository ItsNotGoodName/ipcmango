// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/checkbox
// https://ui.shadcn.com/docs/components/checkbox
import { Checkbox } from "@kobalte/core";
import { RiSystemCheckLine } from "solid-icons/ri";
import { ComponentProps, splitProps } from "solid-js";

import { cn } from "~/lib/utils"

export function CheckboxRoot(props: Checkbox.CheckboxRootProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Checkbox.Root
    class={cn("flex flex-wrap items-center space-x-2", props.class)}
    {...rest}
  />
}

export function CheckboxControl(props: Omit<Checkbox.CheckboxControlProps, "children"> & { inputProps?: Omit<Checkbox.CheckboxInputProps, "class"> }) {
  const [_, rest] = splitProps(props, ["class", "inputProps"])
  return <>
    <Checkbox.Input class="peer" {...props.inputProps} />
    <Checkbox.Control
      class={cn(
        "border-primary peer-focus-visible:ring-ring ui-checked:bg-primary ui-checked:text-primary-foreground peer h-4 w-4 shrink-0 cursor-pointer rounded-sm border shadow ui-disabled:cursor-not-allowed ui-disabled:opacity-50 peer-focus-visible:outline-none peer-focus-visible:ring-1",
        props.class
      )}
      {...rest}
    >
      <Checkbox.Indicator class="flex items-center justify-center text-current">
        <RiSystemCheckLine class="h-4 w-4" />
      </Checkbox.Indicator>
    </Checkbox.Control>
  </>
}

export const CheckboxLabel = Checkbox.Label

export function CheckboxDescription(props: ComponentProps<typeof Checkbox.Description>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Checkbox.Description
    class={cn("w-full text-sm font-medium")}
    {...rest}
  />
}

export function CheckboxErrorMessage(props: ComponentProps<typeof Checkbox.ErrorMessage>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Checkbox.ErrorMessage
    class={cn("text-destructive w-full text-sm font-medium")}
    {...rest}
  />
}
