import { A } from "@solidjs/router";
import {
  ErrorBoundary,
  For,
  Suspense,
  createEffect,
  createSignal,
  onCleanup,
  untrack,
} from "solid-js";
import {
  formatDate,
  useQueryBoolean,
  useQueryFilter,
  useQueryNumber,
} from "~/lib/utils";
import { LayoutCenter } from "~/ui/Layout";
import {
  TableBody,
  TableCell,
  TableCellEnd,
  TableHead,
  TableHeadEnd,
  TableHeader,
  TableRoot,
  TableRow,
} from "~/ui/Table";
import { linkVariants } from "~/ui/Link";
import { PageError, PageTitle } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import { RiSystemDeleteBinLine } from "solid-icons/ri";
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
  EventExpandButton,
  EventJSONTableRow,
} from "~/components/Event";
import { api } from "./data";
import { createQuery } from "@tanstack/solid-query";
import { DeviceEvent } from "~/client";
import { DeviceFilterCombobox } from "~/components/Device";
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

export default function EventsLive() {
  const data = createQuery(() => ({
    ...api.devices.list,
    throwOnError: true,
  }));

  const dataQuery = useQueryBoolean("data");

  const deviceQuery = useQueryFilter("devices");
  const codeQuery = useQueryFilter("codes");
  const actionQuery = useQueryFilter("actions");
  const limitQuery = useQueryNumber("limit", 100);

  const [events, setEvents] = createSignal<DeviceEvent[]>([]);

  createEffect(() => {
    const sse = new EventSource(
      "/api/events/sse" +
        getQueryString({
          devices: deviceQuery.values(),
          codes: codeQuery.values(),
          actions: actionQuery.values(),
        }),
    );

    sse.onmessage = (ev: MessageEvent<string>) => {
      const newEvent = JSON.parse(ev.data) as DeviceEvent;
      const end = untrack(limitQuery.value);
      setEvents((prev) =>
        end == 0 ? [newEvent, ...prev] : [newEvent, ...prev.slice(0, end - 1)],
      );
    };

    onCleanup(() => sse.close());
  });

  return (
    <LayoutCenter class="max-w-4xl">
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
                deviceIDs={deviceQuery.values()}
                setDeviceIDs={deviceQuery.setValues}
              />
              <EventCodeFilterCombobox
                codes={codeQuery.values()}
                setCodes={codeQuery.setValues}
              />
              <EventActionFilterCombobox
                actions={actionQuery.values()}
                setActions={actionQuery.setValues}
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
                onChange={(value) => limitQuery.setValue(value.value)}
                value={{ value: limitQuery.value(), name: "" }}
                itemComponent={(props) => (
                  <SelectItem item={props.item}>
                    {props.item.rawValue.name}
                  </SelectItem>
                )}
                class="space-y-2"
              >
                <SelectTrigger>
                  <SelectValue<{ name: string }>>
                    {(state) =>
                      state.selectedOption()?.name ?? limitQuery.value()
                    }
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
                <RiSystemDeleteBinLine class="size-5" />
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
                <TableHeadEnd>
                  <EventExpandButton
                    expanded={dataQuery.value()}
                    onClick={() => dataQuery.setValue(!dataQuery.value())}
                  />
                </TableHeadEnd>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={events()}>
                {(v) => {
                  const [rowDataOpen, setRowDataOpen] = createSignal(
                    dataQuery.value(),
                  );
                  createEffect(() => setRowDataOpen(dataQuery.value()));

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
                        <TableCellEnd>
                          <EventExpandButton
                            expanded={rowDataOpen()}
                            onClick={() => setRowDataOpen(!rowDataOpen())}
                          />
                        </TableCellEnd>
                      </TableRow>
                      <EventJSONTableRow
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
    </LayoutCenter>
  );
}
