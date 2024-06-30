import { JSX, Show, splitProps } from "solid-js";
import { FormStore } from "@modular-forms/solid";

import { cn } from "~/lib/utils";

export function FormMessage(
  props: JSX.HTMLAttributes<HTMLParagraphElement> & {
    form: FormStore<any, any>;
  },
) {
  const [_, rest] = splitProps(props, ["class", "form", "children"]);
  const body = () =>
    props.form.response.message ? props.form.response.message : props.children;

  return (
    <Show when={body()}>
      <p
        class={cn("text-sm font-medium text-destructive", props.class)}
        {...rest}
      >
        {body()}
      </p>
    </Show>
  );
}
