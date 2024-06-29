import { writeClipboard } from "@solid-primitives/clipboard";
import hljs from "~/lib/hljs";
import {
  ErrorBoundary,
  For,
  Show,
  Suspense,
  batch,
  createEffect,
  createSignal,
} from "solid-js";
import {
  RiArrowsArrowDownSLine,
  RiDocumentClipboardLine,
  RiSystemDeleteBinLine,
  RiSystemFilterLine,
} from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { LayoutNormal } from "~/ui/Layout";
import { PageError, PageTitle } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import {
  createMutation,
  createQuery,
  useQueryClient,
} from "@tanstack/solid-query";
import {
  ComboboxRoot,
  ComboboxItem,
  ComboboxItemLabel,
  ComboboxControl,
  ComboboxTrigger,
  ComboboxIcon,
  ComboboxState,
  ComboboxReset,
  ComboboxContent,
  ComboboxInput,
  ComboboxListbox,
} from "~/ui/Combobox";
import { api } from "./data";
import {
  PaginationEllipsis,
  PaginationEnd,
  PaginationItem,
  PaginationItems,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
  PaginationRoot,
  PaginationStart,
} from "~/ui/Pagination";
import { A } from "@solidjs/router";
import {
  formatDate,
  parseDate,
  useQueryBoolean,
  useQueryFilter,
  useQueryNumber,
  useQueryString,
} from "~/lib/utils";
import { linkVariants } from "~/ui/Link";
import {
  TableRoot,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
  TableCaption,
} from "~/ui/Table";
import { PageMetadata, PositionEnd, SortButton } from "~/components/Utils";
import { DeviceFilterCombobox } from "~/components/DeviceFilterCombobox";
import {
  SelectRoot,
  SelectItem,
  SelectTrigger,
  SelectValue,
  SelectPortal,
  SelectContent,
  SelectListbox,
} from "~/ui/Select";
import { deleteApiEvents } from "~/client";
import { toast } from "~/ui/Toast";
import {
  AlertDialogRoot,
  AlertDialogTrigger,
  AlertDialogModal,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogAction,
} from "~/ui/AlertDialog";

export function JSONTableRow(props: {
  colspan?: number;
  expanded?: boolean;
  data: string;
}) {
  return (
    <tr class="border-b">
      <td colspan={props.colspan} class="p-0">
        <div class="relative overflow-y-hidden">
          <Button
            onClick={() => writeClipboard(props.data)}
            title="Copy"
            size="icon"
            variant="ghost"
            class="absolute right-4 top-2"
          >
            <RiDocumentClipboardLine class="size-5" />
          </Button>
          <pre>
            <Show when={props.expanded}>
              <code
                innerHTML={
                  hljs.highlight(props.data, { language: "json" }).value
                }
                class="hljs"
              />
            </Show>
          </pre>
        </div>
      </td>
    </tr>
  );
}

export function EventCodeFilterCombobox(props: {
  codes: Array<string>;
  setCodes: (value: Array<string>) => void;
}) {
  const data = createQuery(() => ({
    ...api.eventCodes.list,
    throwOnError: true,
  }));

  return (
    <ComboboxRoot<string>
      multiple
      options={data.data || []}
      placeholder="Code"
      value={data.data?.filter((v) => props.codes.includes(v))}
      onChange={(value) => props.setCodes(value)}
      itemComponent={(props) => (
        <ComboboxItem item={props.item}>
          <ComboboxItemLabel>{props.item.rawValue}</ComboboxItemLabel>
        </ComboboxItem>
      )}
    >
      <ComboboxControl<string> aria-label="Code">
        {(state) => (
          <ComboboxTrigger>
            <ComboboxIcon as={RiSystemFilterLine} class="size-4" />
            Code
            <ComboboxState state={state} />
            <ComboboxReset state={state} class="size-4" />
          </ComboboxTrigger>
        )}
      </ComboboxControl>
      <ComboboxContent>
        <ComboboxInput />
        <ComboboxListbox />
      </ComboboxContent>
    </ComboboxRoot>
  );
}

export function EventActionFilterCombobox(props: {
  actions: Array<string>;
  setActions: (value: Array<string>) => void;
}) {
  const data = createQuery(() => ({
    ...api.eventActions.list,
    throwOnError: true,
  }));

  return (
    <ComboboxRoot<string>
      multiple
      options={data.data || []}
      placeholder="Action"
      value={data.data?.filter((v) => props.actions.includes(v))}
      onChange={(value) => props.setActions(value)}
      itemComponent={(props) => (
        <ComboboxItem item={props.item}>
          <ComboboxItemLabel>{props.item.rawValue}</ComboboxItemLabel>
        </ComboboxItem>
      )}
    >
      <ComboboxControl<string> aria-label="Action">
        {(state) => (
          <ComboboxTrigger>
            <ComboboxIcon as={RiSystemFilterLine} class="size-4" />
            Action
            <ComboboxState state={state} />
            <ComboboxReset state={state} class="size-4" />
          </ComboboxTrigger>
        )}
      </ComboboxControl>
      <ComboboxContent>
        <ComboboxInput />
        <ComboboxListbox />
      </ComboboxContent>
    </ComboboxRoot>
  );
}

export function ExpandButton(props: {
  expanded?: boolean;
  onClick: () => void;
}) {
  return (
    <Button
      data-expanded={props.expanded}
      onClick={props.onClick}
      title="Expand"
      size="icon"
      variant="ghost"
      class="[&[data-expanded=true]>svg]:rotate-180"
    >
      <RiArrowsArrowDownSLine class="h-5 w-5 shrink-0 transition-transform duration-200" />
    </Button>
  );
}

export default function () {
  const client = useQueryClient();

  const pageQuery = useQueryNumber("page", 1);
  const perPageQuery = useQueryNumber("perPage", 10);

  const devicesQuery = useQueryFilter("devices");
  const codesQuery = useQueryFilter("codes");
  const actionsQuery = useQueryFilter("actions");
  const orderQuery = useQueryString("order");

  const devices = createQuery(() => ({
    ...api.devices.list,
    throwOnError: true,
  }));

  const data = createQuery(() => ({
    ...api.events.list({
      page: pageQuery.value(),
      perPage: perPageQuery.value(),
      device: devicesQuery.values(),
      codes: codesQuery.values(),
      actions: actionsQuery.values(),
      order: orderQuery.value(),
    }),
    throwOnError: true,
  }));

  const queryData = useQueryBoolean("data");

  const [deleteDialog, setDeleteDialog] = createSignal(false);
  const deleteMutation = createMutation(() => ({
    mutationFn: () => deleteApiEvents(),
    onSuccess: () =>
      batch(() => {
        client.invalidateQueries({ queryKey: api.events.list._def });
        setDeleteDialog(false);
      }),
    onError: (error) => toast.error(error.name, error.message),
  }));

  return (
    <LayoutNormal class="max-w-4xl">
      <PageTitle>Events</PageTitle>

      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex justify-between gap-2">
            <div class="flex flex-wrap gap-2">
              <DeviceFilterCombobox
                deviceIDs={devicesQuery.values()}
                setDeviceIDs={devicesQuery.setValues}
              />
              <EventCodeFilterCombobox
                codes={codesQuery.values()}
                setCodes={codesQuery.setValues}
              />
              <EventActionFilterCombobox
                actions={actionsQuery.values()}
                setActions={actionsQuery.setValues}
              />
              <SelectRoot
                options={[
                  { value: 10, name: "10" },
                  { value: 25, name: "25" },
                  { value: 100, name: "100" },
                ]}
                optionTextValue="name"
                optionValue="value"
                onChange={(value) => perPageQuery.setValue(value?.value)}
                value={{ value: perPageQuery.value(), name: "" }}
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
                      state.selectedOption()?.name ?? perPageQuery.value()
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
              <AlertDialogRoot
                open={deleteDialog()}
                onOpenChange={setDeleteDialog}
              >
                <AlertDialogTrigger
                  as={Button}
                  disabled={deleteMutation.isPending}
                  variant="destructive"
                  size="icon"
                >
                  <RiSystemDeleteBinLine class="size-5" />
                </AlertDialogTrigger>
                <AlertDialogModal>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Delete events?</AlertDialogTitle>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                      disabled={deleteMutation.isPending}
                      onClick={() => deleteMutation.mutate()}
                      variant="destructive"
                    >
                      Delete
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogModal>
              </AlertDialogRoot>
            </div>
          </div>

          <PaginationRoot
            page={data.data?.pagination.page}
            count={data.data?.pagination.total_pages || 0}
            onPageChange={(page) => pageQuery.setValue(page)}
            itemComponent={(props) => (
              <PaginationItem page={props.page}>
                <PaginationLink
                  isActive={props.page == data.data?.pagination.page}
                >
                  {props.page}
                </PaginationLink>
              </PaginationItem>
            )}
            ellipsisComponent={() => <PaginationEllipsis />}
          >
            <PaginationStart>
              <PaginationItems />
            </PaginationStart>
            <PaginationEnd>
              <PaginationPrevious />
              <PaginationNext />
            </PaginationEnd>
          </PaginationRoot>

          <TableRoot>
            <TableHeader>
              <TableRow>
                <TableHead>
                  <SortButton
                    order={orderQuery.value()}
                    onToggle={(order) => orderQuery.setValue(order)}
                  >
                    Created At
                  </SortButton>
                </TableHead>
                <TableHead>Device</TableHead>
                <TableHead>Code</TableHead>
                <TableHead>Action</TableHead>
                <TableHead>Index</TableHead>
                <TableHead class="p-0">
                  <PositionEnd>
                    <ExpandButton
                      expanded={queryData.value()}
                      onClick={() => queryData.setValue(!queryData.value())}
                    />
                  </PositionEnd>
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data.data?.data}>
                {(v) => {
                  const [rowDataOpen, setRowDataOpen] = createSignal(
                    queryData.value(),
                  );
                  createEffect(() => setRowDataOpen(queryData.value()));

                  return (
                    <>
                      <TableRow class="border-b-0">
                        <TableCell>
                          {formatDate(parseDate(v.created_at))}
                        </TableCell>
                        <TableCell>
                          <A
                            href={`/devices/${v.device_uuid}`}
                            class={linkVariants()}
                          >
                            {
                              devices.data?.find(
                                (d) => d.uuid == String(v.device_uuid),
                              )?.name
                            }
                          </A>
                        </TableCell>
                        <TableCell>{v.code}</TableCell>
                        <TableCell>{v.action}</TableCell>
                        <TableCell>{v.index.toString()}</TableCell>
                        <TableCell class="py-0">
                          <PositionEnd>
                            <ExpandButton
                              expanded={rowDataOpen()}
                              onClick={() => setRowDataOpen(!rowDataOpen())}
                            />
                          </PositionEnd>
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
            <TableCaption>
              <PageMetadata pageResult={data.data?.pagination} />
            </TableCaption>
          </TableRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  );
}
