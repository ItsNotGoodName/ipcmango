import { RiArrowsArrowDownSLine } from "solid-icons/ri";
import { ParentProps } from "solid-js";
import { PagePagination } from "~/client";

import { cn } from "~/lib/utils";

enum Order {
  Unknown = "",
  Ascending = "ascending",
  Descending = "descending",
}

export function SortButton(
  props: ParentProps<{
    onToggle: (order: Order) => void;
    order?: Order | string;
  }>,
) {
  return (
    <button
      onClick={() => {
        if (props.order == Order.Ascending) {
          props.onToggle(Order.Descending);
        } else if (props.order == Order.Descending) {
          props.onToggle(Order.Unknown);
        } else {
          props.onToggle(Order.Ascending);
        }
      }}
      class={cn(
        "flex items-center whitespace-nowrap text-nowrap",
        (props.order == Order.Ascending || props.order == Order.Descending) &&
          "text-blue-500",
      )}
    >
      {props.children}
      <RiArrowsArrowDownSLine
        data-selected={props.order == Order.Ascending}
        class="size-5 transition-all data-[selected=true]:rotate-180"
      />
    </button>
  );
}

export function PositionEnd(props: ParentProps) {
  return <div class="flex items-center justify-end" {...props} />;
}

export function PageMetadata(props: { pageResult?: PagePagination }) {
  return (
    <div class="flex justify-between">
      <div>
        Seen {props.pageResult?.seen_items.toString() || 0} of{" "}
        {props.pageResult?.total_items.toString() || 0}
      </div>
      <div>
        Page {props.pageResult?.page || 0} of{" "}
        {props.pageResult?.total_pages || 0}
      </div>
    </div>
  );
}
