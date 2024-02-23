// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/switch
// https://ui.shadcn.com/docs/components/switch
import { Switch } from "@kobalte/core"
import { ComponentProps, splitProps } from "solid-js"

import { cn } from "~/lib/utils"
import { labelVariants } from "./Label"

export const SwitchRoot = Switch.Root

export function SwitchLabel(props: Switch.SwitchLabelProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Switch.Label
    class={cn(labelVariants(), props.class)}
    {...rest}
  />
}
export function SwitchDescription(props: Switch.SwitchDescriptionProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Switch.Description
    class={cn("text-muted-foreground w-full text-sm", props.class)}
    {...rest}
  />
}

export function SwitchErrorMessage(props: Switch.SwitchErrorMessageProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Switch.ErrorMessage
    class={cn("text-destructive w-full text-sm font-medium", props.class)}
    {...rest}
  />
}

export function SwitchControl(props: ComponentProps<typeof Switch.Control>) {
  const [_, rest] = splitProps(props, ["class"])
  return <>
    <Switch.Input class="peer" />
    <Switch.Control
      class={cn(
        "peer-focus-visible:ring-ring peer-focus-visible:ring-offset-background ui-checked:bg-primary ui-not-checked:bg-input inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors peer-focus-visible:outline-none peer-focus-visible:ring-2 peer-focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
        props.class
      )}
      {...rest}
    >
      <Switch.Thumb
        class={cn(
          "bg-background ui-checked:translate-x-5 ui-not-checked:translate-x-0 pointer-events-none block h-5 w-5 rounded-full shadow-lg ring-0 transition-transform"
        )}
      />
    </Switch.Control>
  </>
}

