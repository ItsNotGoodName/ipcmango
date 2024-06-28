// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/breadcrumbs
import { Breadcrumbs } from "@kobalte/core/breadcrumbs";
import { ComponentProps, ParentProps, splitProps } from "solid-js";

import { cn } from "~/lib/utils";

export function BreadcrumbsRoot(props: ComponentProps<typeof Breadcrumbs>) {
  const [_, rest] = splitProps(props, ["children"]);
  return (
    <Breadcrumbs {...rest}>
      <ol class="inline-flex items-center">{props.children}</ol>
    </Breadcrumbs>
  );
}

export function BreadcrumbsItem(props: ParentProps) {
  return <li class="inline-flex items-center">{props.children}</li>;
}

export function BreadcrumbsLink(
  props: ComponentProps<typeof Breadcrumbs.Link>,
) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Breadcrumbs.Link class={cn("hover:text-sky-500", props.class)} {...rest} />
  );
}

export function BreadcrumbsSeparator(
  props: ComponentProps<typeof Breadcrumbs.Separator>,
) {
  const [_, rest] = splitProps(props, ["class"]);
  return <Breadcrumbs.Separator class={cn("px-2", props.class)} {...rest} />;
}
