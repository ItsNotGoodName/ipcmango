import { createVirtualizer } from "@tanstack/solid-virtual";
import { createSignal } from "solid-js";
import { formatDate } from "~/lib/utils";
import { Button } from "~/ui/Button";

type Item = {
  color: string;
  start_time: Date;
  duration: number;
};

export default function () {
  // seconds, must be even
  const [range, setRange] = createSignal(60 * 60);
  // current cursor position
  const [cursor, setCursor] = createSignal(new Date(Date.now()));
  // items to render
  const [items, setItems] = createSignal<Array<Item>>();
  const [offset, setOffset] = createSignal(0);

  let down = false;
  let startX = 0;

  let ref: HTMLDivElement | null;
  const virtualizer = createVirtualizer({
    count: 10000,
    getScrollElement: () => ref,
    estimateSize: () => 32,
    overscan: 5,
  });

  return (
    <div class="flex h-full flex-col">
      <div class="flex-1">{formatDate(cursor())}</div>
      <div class="flex h-20 flex-col border-t">
        <div
          class="relative h-14"
          ref={ref!}
          style={{
            width: `${virtualizer.getTotalSize()}px`,
          }}
        >
          <div class="flex-1"></div>
          <div
            class="h-2 bg-red-500"
            style={{ translate: offset() + "px" }}
          ></div>
        </div>
        <div class="flex h-10 justify-center gap-2 p-2">
          <Button size="xs" onClick={() => setOffset((prev) => prev - 1)}>
            Back
          </Button>
          <Button size="xs" onClick={() => setOffset((prev) => prev + 1)}>
            Forward
          </Button>
        </div>
      </div>
    </div>
  );
}
