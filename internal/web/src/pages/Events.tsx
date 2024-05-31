import { writeClipboard } from "@solid-primitives/clipboard";
import hljs from "~/lib/hljs"
import { Show, } from "solid-js";
import { RiDocumentClipboardLine, } from "solid-icons/ri";
import { Button, } from "~/ui/Button";

export function JSONTableRow(props: { colspan?: number, expanded?: boolean, data: string }) {
  return (
    <tr class="border-b">
      <td colspan={props.colspan} class="p-0">
        <div class="overflow-y-hidden relative">
          <Button onClick={() => writeClipboard(props.data)} title="Copy" size="icon" variant="ghost" class="absolute top-2 right-4">
            <RiDocumentClipboardLine class="size-5" />
          </Button>
          <pre>
            <Show when={props.expanded}>
              <code innerHTML={hljs.highlight(props.data, { language: "json" }).value} class="hljs" />
            </Show>
          </pre>
        </div>
      </td>
    </tr>
  )
}
