import { RiArrowsArrowDownSLine } from "solid-icons/ri";
import { ParentProps } from "solid-js";
import { PagePagination } from "~/client";

import { cn } from "~/lib/utils";

enum Sort {
  ASC = 1,
  DESC = 2,
}

export function SortButton(
  props: ParentProps<{
    onClick?: () => void;
    sort?: Sort;
  }>,
) {
  return (
    <button
      onClick={props.onClick}
      class={cn(
        "flex items-center whitespace-nowrap text-nowrap",
        props.sort && "text-blue-500",
      )}
    >
      {props.children}
      <RiArrowsArrowDownSLine
        data-selected={props.sort == Sort.ASC}
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
