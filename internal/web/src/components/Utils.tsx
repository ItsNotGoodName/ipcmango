import { RiArrowsArrowDownSLine } from "solid-icons/ri";
import { ParentProps } from "solid-js";
import { PagePagination } from "~/client";

import { Sort, cn } from "~/lib/utils";

export function SortButton(
  props: ParentProps<{
    field: string;
    onToggle: (value: Sort) => void;
    sort: Sort;
  }>,
) {
  const active = () => props.field == props.sort.field;
  return (
    <button
      onClick={() => {
        if (!active() || props.sort.order == undefined)
          return props.onToggle({ field: props.field, order: "descending" });
        if (props.sort.order == "descending")
          return props.onToggle({ field: props.field, order: "ascending" });
        return props.onToggle({ order: undefined });
      }}
      class={cn(
        "flex items-center whitespace-nowrap text-nowrap",
        active() &&
          (props.sort.order == "ascending" ||
            props.sort.order == "descending") &&
          "text-blue-500",
      )}
    >
      {props.children}
      <RiArrowsArrowDownSLine
        data-selected={active() && props.sort.order == "ascending"}
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
