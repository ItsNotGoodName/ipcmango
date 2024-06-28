// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/tabs
// https://ui.shadcn.com/docs/components/tabs
import { Tabs } from "@kobalte/core/tabs";
import { ComponentProps, splitProps } from "solid-js";

import { cn } from "~/lib/utils";

export const TabsRoot = Tabs;

export function TabsList(props: ComponentProps<typeof Tabs.List>) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Tabs.List
      class={cn(
        "inline-flex h-10 items-center rounded-md bg-muted p-1 text-muted-foreground",
        props.class,
      )}
      {...rest}
    />
  );
}

export function TabsTrigger(props: ComponentProps<typeof Tabs.Trigger>) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Tabs.Trigger
      class={cn(
        "ui-selected::shadow-sm inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1.5 text-sm font-medium ring-offset-background transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 ui-selected:bg-background ui-selected:text-foreground",
        props.class,
      )}
      {...rest}
    />
  );
}

export function TabsContent(props: ComponentProps<typeof Tabs.Content>) {
  const [_, rest] = splitProps(props, ["class"]);
  return (
    <Tabs.Content
      class={cn(
        "mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        props.class,
      )}
      {...rest}
    />
  );
}
