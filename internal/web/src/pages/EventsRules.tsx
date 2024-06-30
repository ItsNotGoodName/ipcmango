import { A } from "@solidjs/router";
import { createQuery } from "@tanstack/solid-query";
import { RiSystemAddLine, RiSystemDeleteBinLine } from "solid-icons/ri";
import { ErrorBoundary, For, Show, Suspense } from "solid-js";
import { createRowSelection } from "~/lib/utils";
import {
  BreadcrumbsRoot,
  BreadcrumbsItem,
  BreadcrumbsLink,
  BreadcrumbsSeparator,
} from "~/ui/Breadcrumbs";
import { Button } from "~/ui/Button";
import { CheckboxRoot, CheckboxControl } from "~/ui/Checkbox";
import { LayoutNormal } from "~/ui/Layout";
import { PageError, PageTitle } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import {
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRoot,
  TableRow,
} from "~/ui/Table";
import { TextFieldRoot, TextFieldInput } from "~/ui/TextField";
import { api } from "./data";

export default function () {
  const data = createQuery(() => ({
    ...api.eventRules.list,
    throwOnError: true,
  }));

  const rowSelection = createRowSelection(
    () => data.data?.map((v) => ({ id: v.code, disabled: !v.deletable })) ?? [],
  );

  return (
    <LayoutNormal class="max-w-4xl">
      <PageTitle>
        <BreadcrumbsRoot>
          <BreadcrumbsItem>
            <BreadcrumbsLink as={A} href="/events">
              Events
            </BreadcrumbsLink>
            <BreadcrumbsSeparator />
          </BreadcrumbsItem>
          <BreadcrumbsItem>Rules</BreadcrumbsItem>
        </BreadcrumbsRoot>
      </PageTitle>

      <div class="flex justify-end gap-2">
        <Button size="icon">
          <RiSystemAddLine class="size-5" />
        </Button>
        <Button size="icon" variant="destructive">
          <RiSystemDeleteBinLine class="size-5" />
        </Button>
      </div>

      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <TableRoot>
            <TableHeader>
              <TableRow>
                <TableHead>
                  <CheckboxRoot
                    indeterminate={rowSelection.multiple()}
                    checked={rowSelection.all()}
                    onChange={rowSelection.setAll}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableHead>
                <TableHead class="w-full">Code</TableHead>
                <TableHead class="w-0">
                  <button>DB</button>
                </TableHead>
                <TableHead class="w-0">
                  <button>Live</button>
                </TableHead>
                <TableHead class="w-0">
                  <button>MQTT</button>
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data.data}>
                {(item, index) => (
                  <TableRow>
                    <TableCell>
                      <CheckboxRoot
                        disabled={rowSelection.items[index()]?.disabled}
                        checked={rowSelection.items[index()]?.checked}
                        onChange={(value) => rowSelection.set(item.code, value)}
                      >
                        <CheckboxControl />
                      </CheckboxRoot>
                    </TableCell>
                    <Show
                      when={item.deletable}
                      fallback={
                        <TableCell class="w-full">{item.code}</TableCell>
                      }
                    >
                      <td class="w-full min-w-32 py-0 align-middle">
                        <TextFieldRoot>
                          <TextFieldInput />
                        </TextFieldRoot>
                      </td>
                    </Show>
                    <TableCell>
                      <CheckboxRoot checked={!item.ignore_db}>
                        <CheckboxControl />
                      </CheckboxRoot>
                    </TableCell>
                    <TableCell>
                      <CheckboxRoot checked={!item.ignore_live}>
                        <CheckboxControl />
                      </CheckboxRoot>
                    </TableCell>
                    <TableCell>
                      <CheckboxRoot checked={!item.ignore_mqtt}>
                        <CheckboxControl />
                      </CheckboxRoot>
                    </TableCell>
                  </TableRow>
                )}
              </For>
            </TableBody>
          </TableRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  );
}
