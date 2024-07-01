import {
  ErrorBoundary,
  For,
  Suspense,
  createEffect,
  createSignal,
} from "solid-js";
import { RiSystemDeleteBinLine } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { LayoutCenter } from "~/ui/Layout";
import { PageError, PageTitle } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import {
  createMutation,
  createQuery,
  useQueryClient,
} from "@tanstack/solid-query";
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
  useQuerySort,
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
  TableHeadEnd,
  TableCellEnd,
} from "~/ui/Table";
import { PageMetadata, SortButton } from "~/components/Utils";
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
import {
  EventActionFilterCombobox,
  EventCodeFilterCombobox,
  ExpandButton,
  JSONTableRow,
} from "~/components/Events";

export default function Events() {
  const client = useQueryClient();

  const pageQuery = useQueryNumber("page", 1);
  const perPageQuery = useQueryNumber("perPage", 10);

  const devicesQuery = useQueryFilter("devices");
  const codesQuery = useQueryFilter("codes");
  const actionsQuery = useQueryFilter("actions");
  const sortQuery = useQuerySort("sort");

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
      order: sortQuery.value().order,
    }),
    throwOnError: true,
  }));

  const queryData = useQueryBoolean("data");

  const [deleteDialog, setDeleteDialog] = createSignal(false);
  const deleteMutation = createMutation(() => ({
    mutationFn: deleteApiEvents,
    onSuccess: () =>
      client
        .invalidateQueries({ queryKey: api.events.list._def })
        .then(() => setDeleteDialog(false)),
    onError: (error) => toast.error(error.name, error.message),
  }));

  return (
    <LayoutCenter class="max-w-4xl">
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
                  title="Delete Events"
                >
                  <RiSystemDeleteBinLine class="size-5" />
                </AlertDialogTrigger>
                <AlertDialogModal>
                  <AlertDialogHeader>
                    <AlertDialogTitle>
                      Are you sure you wish to delete all events?
                    </AlertDialogTitle>
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
                    field="createdAt"
                    sort={sortQuery.value()}
                    onToggle={sortQuery.setValue}
                  >
                    Created At
                  </SortButton>
                </TableHead>
                <TableHead>Device</TableHead>
                <TableHead>Code</TableHead>
                <TableHead>Action</TableHead>
                <TableHead>Index</TableHead>
                <TableHeadEnd>
                  <ExpandButton
                    expanded={queryData.value()}
                    onClick={() => queryData.setValue(!queryData.value())}
                  />
                </TableHeadEnd>
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
                        <TableCellEnd>
                          <ExpandButton
                            expanded={rowDataOpen()}
                            onClick={() => setRowDataOpen(!rowDataOpen())}
                          />
                        </TableCellEnd>
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
    </LayoutCenter>
  );
}
