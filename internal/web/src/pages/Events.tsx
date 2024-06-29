import { writeClipboard } from "@solid-primitives/clipboard";
import hljs from "~/lib/hljs";
import { Show } from "solid-js";
import { RiDocumentClipboardLine, RiSystemFilterLine } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { LayoutNormal } from "~/ui/Layout";
import { PageTitle } from "~/ui/Page";
import { createQuery } from "@tanstack/solid-query";
import {
  ComboboxRoot,
  ComboboxItem,
  ComboboxItemLabel,
  ComboboxControl,
  ComboboxTrigger,
  ComboboxIcon,
  ComboboxState,
  ComboboxReset,
  ComboboxContent,
  ComboboxInput,
  ComboboxListbox,
} from "~/ui/Combobox";
import { api } from "./data";

export function JSONTableRow(props: {
  colspan?: number;
  expanded?: boolean;
  data: string;
}) {
  return (
    <tr class="border-b">
      <td colspan={props.colspan} class="p-0">
        <div class="relative overflow-y-hidden">
          <Button
            onClick={() => writeClipboard(props.data)}
            title="Copy"
            size="icon"
            variant="ghost"
            class="absolute right-4 top-2"
          >
            <RiDocumentClipboardLine class="size-5" />
          </Button>
          <pre>
            <Show when={props.expanded}>
              <code
                innerHTML={
                  hljs.highlight(props.data, { language: "json" }).value
                }
                class="hljs"
              />
            </Show>
          </pre>
        </div>
      </td>
    </tr>
  );
}

export function EventCodeFilterCombobox(props: {
  codes: Array<string>;
  setCodes: (value: Array<string>) => void;
}) {
  const data = createQuery(() => ({
    ...api.eventCodes.list,
    throwOnError: true,
  }));

  return (
    <ComboboxRoot<string>
      multiple
      options={data.data || []}
      placeholder="Code"
      value={data.data?.filter((v) => props.codes.includes(v))}
      onChange={(value) => props.setCodes(value)}
      itemComponent={(props) => (
        <ComboboxItem item={props.item}>
          <ComboboxItemLabel>{props.item.rawValue}</ComboboxItemLabel>
        </ComboboxItem>
      )}
    >
      <ComboboxControl<string> aria-label="Code">
        {(state) => (
          <ComboboxTrigger>
            <ComboboxIcon as={RiSystemFilterLine} class="size-4" />
            Code
            <ComboboxState state={state} />
            <ComboboxReset state={state} class="size-4" />
          </ComboboxTrigger>
        )}
      </ComboboxControl>
      <ComboboxContent>
        <ComboboxInput />
        <ComboboxListbox />
      </ComboboxContent>
    </ComboboxRoot>
  );
}

export function EventActionFilterCombobox(props: {
  actions: Array<string>;
  setActions: (value: Array<string>) => void;
}) {
  const data = createQuery(() => ({
    ...api.eventActions.list,
    throwOnError: true,
  }));

  return (
    <ComboboxRoot<string>
      multiple
      options={data.data || []}
      placeholder="Action"
      value={data.data?.filter((v) => props.actions.includes(v))}
      onChange={(value) => props.setActions(value)}
      itemComponent={(props) => (
        <ComboboxItem item={props.item}>
          <ComboboxItemLabel>{props.item.rawValue}</ComboboxItemLabel>
        </ComboboxItem>
      )}
    >
      <ComboboxControl<string> aria-label="Action">
        {(state) => (
          <ComboboxTrigger>
            <ComboboxIcon as={RiSystemFilterLine} class="size-4" />
            Action
            <ComboboxState state={state} />
            <ComboboxReset state={state} class="size-4" />
          </ComboboxTrigger>
        )}
      </ComboboxControl>
      <ComboboxContent>
        <ComboboxInput />
        <ComboboxListbox />
      </ComboboxContent>
    </ComboboxRoot>
  );
}

export default function () {
  return (
    <LayoutNormal class="max-w-4xl">
      <PageTitle>Events</PageTitle>
    </LayoutNormal>
  );
}
