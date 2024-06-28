// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/seperator
// https://ui.shadcn.com/docs/components/seperator
import { Separator } from "@kobalte/core/separator";
import { JSX, splitProps } from "solid-js";

export function Seperator(
  props: JSX.HTMLAttributes<HTMLDivElement> & {
    orientation?: "horizontal" | "vertical";
  },
) {
  const [_, rest] = splitProps(props, ["orientation"]);
  return (
    <div {...rest}>
      <Separator
        class={
          "shrink-0 bg-border ui-horizontal:h-[1px] ui-horizontal:w-full ui-vertical:h-full ui-vertical:w-[1px]"
        }
        orientation={props.orientation}
      />
    </div>
  );
}
