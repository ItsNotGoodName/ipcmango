import { A, useSearchParams } from "@solidjs/router";
import {
  ErrorBoundary,
  For,
  Suspense,
  createEffect,
  createSignal,
  onCleanup,
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
import { RiArrowsArrowDownSLine } from "solid-icons/ri";
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

export default function () {
  const [searchParams, setSearchParams] = useSearchParams();

  const data = createQuery(() => api.devices.list);

  const dataOpen = () => Boolean(searchParams.data);
  const setDataOpen = (value: boolean) =>
    setSearchParams({ data: value ? String(value) : "" });

  const deviceFilter = useQueryFilter("device");
  const codeFilter = useQueryFilter("code");
  const actionFilter = useQueryFilter("action");

  const [events, setEvents] = createSignal<DeviceEventsOutput[]>([]);

  createEffect(() => {
    const sse = new EventSource(
      "/api/events" +
        getQueryString({
          "device-uuids": deviceFilter.values(),
          code: codeFilter.values(),
        }),
    );

    sse.onmessage = (ev: MessageEvent<string>) => {
      setEvents((prev) => [JSON.parse(ev.data) as DeviceEventsOutput, ...prev]);
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
          <div class="flex gap-2">
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
