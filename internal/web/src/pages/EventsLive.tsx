import { A, useSearchParams } from "@solidjs/router";
import {
  ErrorBoundary,
  For,
  Suspense,
  createEffect,
  createSignal,
  onCleanup,
  untrack,
} from "solid-js";
import { formatDate, useQueryFilter } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import {
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRoot,
  TableRow,
} from "~/ui/Table";
import { linkVariants } from "~/ui/Link";
import { PageError, PageTitle } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import { RiArrowsArrowDownSLine, RiSystemDeleteBinLine } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import {
  BreadcrumbsItem,
  BreadcrumbsLink,
  BreadcrumbsRoot,
  BreadcrumbsSeparator,
} from "~/ui/Breadcrumbs";
import { createDate, createTimeAgo } from "@solid-primitives/date";
import {
  TooltipArrow,
  TooltipContent,
  TooltipRoot,
  TooltipTrigger,
} from "~/ui/Tooltip";
import {
  EventActionFilterCombobox,
  EventCodeFilterCombobox,
  JSONTableRow,
} from "./Events";
import { api } from "./data";
import { createQuery } from "@tanstack/solid-query";
import { DeviceEventsOutput } from "~/client";
import { DeviceFilterCombobox } from "~/components/DeviceFilterCombobox";
import { getQueryString } from "~/client/core/request";
import {
  SelectRoot,
  SelectItem,
  SelectTrigger,
  SelectValue,
  SelectPortal,
  SelectContent,
  SelectListbox,
} from "~/ui/Select";

export default function () {
  const [searchParams, setSearchParams] = useSearchParams();

  const data = createQuery(() => api.devices.list);

  const dataOpen = () => Boolean(searchParams.data);
  const setDataOpen = (value: boolean) =>
    setSearchParams({ data: value ? String(value) : "" });

  const deviceFilter = useQueryFilter("devices");
  const codeFilter = useQueryFilter("codes");
  const actionFilter = useQueryFilter("actions");
  const limit = () =>
    searchParams.limit == undefined ? 100 : Number(searchParams.limit);

  const [events, setEvents] = createSignal<DeviceEventsOutput[]>([]);

  createEffect(() => {
    const sse = new EventSource(
      "/api/events" +
        getQueryString({
          devices: deviceFilter.values(),
          codes: codeFilter.values(),
          actions: actionFilter.values(),
        }),
    );

    sse.onmessage = (ev: MessageEvent<string>) => {
      const newEvent = JSON.parse(ev.data) as DeviceEventsOutput;
      const end = untrack(limit);
      setEvents((prev) =>
        end == 0 ? [newEvent, ...prev] : [newEvent, ...prev.slice(0, end - 1)],
      );
    };

    onCleanup(() => sse.close());
  });

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
          <BreadcrumbsItem>Live</BreadcrumbsItem>
        </BreadcrumbsRoot>
      </PageTitle>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex justify-between gap-2">
            <div class="flex flex-wrap gap-2">
              <DeviceFilterCombobox
                deviceIDs={deviceFilter.values()}
                setDeviceIDs={deviceFilter.setValues}
              />
              <EventCodeFilterCombobox
                codes={codeFilter.values()}
                setCodes={codeFilter.setValues}
              />
              <EventActionFilterCombobox
                actions={actionFilter.values()}
                setActions={actionFilter.setValues}
              />
              <SelectRoot
                options={[
                  { value: 0, name: "Unlimited" },
                  { value: 10, name: "10" },
                  { value: 25, name: "25" },
                  { value: 100, name: "100" },
                ]}
                optionTextValue="name"
                optionValue="value"
                onChange={(value) => setSearchParams({ limit: value.value })}
                value={{ value: limit(), name: "" }}
                itemComponent={(props) => (
                  <SelectItem item={props.item}>
                    {props.item.rawValue.name}
                  </SelectItem>
                )}
                class="space-y-2"
              >
                <SelectTrigger>
                  <SelectValue<{ name: string }>>
                    {(state) => state.selectedOption()?.name ?? limit()}
                  </SelectValue>
                </SelectTrigger>
                <SelectPortal>
                  <SelectContent>
                    <SelectListbox />
                  </SelectContent>
                </SelectPortal>
              </SelectRoot>
            </div>
            <div>
              <Button size="icon" onClick={() => setEvents([])}>
                <RiSystemDeleteBinLine class="size-6" />
              </Button>
            </div>
          </div>
          <TableRoot>
            <TableHeader>
              <TableRow>
                <TableHead>Created At</TableHead>
                <TableHead>Device</TableHead>
                <TableHead>Code</TableHead>
                <TableHead>Action</TableHead>
                <TableHead>Index</TableHead>
                <TableHead>
                  <div class="flex items-center justify-end">
                    <Button
                      data-expanded={dataOpen()}
                      onClick={() => setDataOpen(!dataOpen())}
                      title="Data"
                      size="icon"
                      variant="ghost"
                      class="[&[data-expanded=true]>svg]:rotate-180"
                    >
                      <RiArrowsArrowDownSLine class="h-5 w-5 shrink-0 transition-transform duration-200" />
                    </Button>
                  </div>
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={events()}>
                {(v) => {
                  const [rowDataOpen, setRowDataOpen] =
                    createSignal(dataOpen());
                  createEffect(() => setRowDataOpen(dataOpen()));

                  const [createdAt] = createDate(() => v.created_at);
                  const [createdAtAgo] = createTimeAgo(createdAt);

                  return (
                    <>
                      <TableRow class="border-b-0">
                        <TableCell>
                          <TooltipRoot>
                            <TooltipTrigger>{createdAtAgo()}</TooltipTrigger>
                            <TooltipContent>
                              <TooltipArrow />
                              {formatDate(createdAt())}
                            </TooltipContent>
                          </TooltipRoot>
                        </TableCell>
                        <TableCell>
                          <A
                            href={`/devices/${v.device_uuid}`}
                            class={linkVariants()}
                          >
                            {
                              data.data?.find(
                                (d) => d.uuid == String(v.device_uuid),
                              )?.name
                            }
                          </A>
                        </TableCell>
                        <TableCell>{v.code}</TableCell>
                        <TableCell>{v.action}</TableCell>
                        <TableCell>{v.index.toString()}</TableCell>
                        <TableCell>
                          <div class="flex items-center justify-end">
                            <Button
                              data-expanded={rowDataOpen()}
                              onClick={() => setRowDataOpen(!rowDataOpen())}
                              title="Data"
                              size="icon"
                              variant="ghost"
                              class="[&[data-expanded=true]>svg]:rotate-180"
                            >
                              <RiArrowsArrowDownSLine class="h-5 w-5 shrink-0 transition-transform duration-200" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                      <JSONTableRow
                        colspan={6}
                        expanded={rowDataOpen()}
                        data={JSON.stringify(v.data, null, 2)}
                      />
                    </>
                  );
                }}
              </For>
            </TableBody>
          </TableRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  );
}
