// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/text-field
// https://ui.shadcn.com/docs/components/input
// https://ui.shadcn.com/docs/components/textarea
import { TextField } from "@kobalte/core/text-field"
import { ComponentProps, splitProps } from "solid-js"

import { cn } from "~/lib/utils"
import { labelVariants } from "./Label"

export function TextFieldInput(props: ComponentProps<typeof TextField.Input>) {
  const [_, rest] = splitProps(props, ["class"])
  return <TextField.Input
    class={cn(
      "border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex h-10 w-full rounded-md border px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
      props.class
    )}
    {...rest}
  />
}

export function TextFieldTextArea(props: ComponentProps<typeof TextField.TextArea>) {
  const [_, rest] = splitProps(props, ["class"])
  return <TextField.TextArea
    class={cn(
      "border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex min-h-[80px] w-full rounded-md border px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
      props.class
    )}
    {...rest}
  />
}

export const TextFieldRoot = TextField

export function TextFieldLabel(props: ComponentProps<typeof TextField.Label>) {
  const [_, rest] = splitProps(props, ["class"])
  return <TextField.Label
    class={cn(labelVariants(), props.class)}
    {...rest}
  />
}

export function TextFieldDescription(props: ComponentProps<typeof TextField.Description>) {
  const [_, rest] = splitProps(props, ["class"])
  return <TextField.Description
    class={cn("text-muted-foreground text-sm", props.class)}
    {...rest}
  />
}

export function TextFieldErrorMessage(props: ComponentProps<typeof TextField.ErrorMessage>) {
  const [_, rest] = splitProps(props, ["class"])
  return <TextField.ErrorMessage
    class={cn("text-destructive text-sm font-medium")}
    {...rest}
  />
}
