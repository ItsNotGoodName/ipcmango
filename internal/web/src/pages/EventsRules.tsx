import Humanize from "humanize-plus";
import { A } from "@solidjs/router";
import {
  createMutation,
  createQuery,
  useQueryClient,
} from "@tanstack/solid-query";
import { RiSystemAddLine, RiSystemDeleteBinLine } from "solid-icons/ri";
import { ErrorBoundary, For, Show, Suspense, createSignal } from "solid-js";
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
  TableHeadBase,
  TableHeader,
  TableRoot,
  TableRow,
} from "~/ui/Table";
import { api } from "./data";
import { deleteApiEventRules } from "~/client";
import { toast } from "~/ui/Toast";
import {
  AlertDialogRoot,
  AlertDialogModal,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogAction,
} from "~/ui/AlertDialog";

export default function () {
  const client = useQueryClient();

  const data = createQuery(() => ({
    ...api.eventRules.list,
    throwOnError: true,
  }));

  const rowSelection = createRowSelection(
    () =>
      data.data?.map((v) => ({ id: v.uuid, disabled: !v.can_delete })) ?? [],
  );

  const [deleteModal, setDeleteModal] = createSignal(false);
  const deleteMutation = createMutation(() => ({
    mutationFn: () =>
      deleteApiEventRules({ requestBody: rowSelection.selections() }),
    onSuccess: () =>
      client
        .invalidateQueries({ queryKey: api.eventRules.list.queryKey })
        .then(() => setDeleteModal(false)),
    onError: (error) => toast.error(error.name, error.message),
  }));

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

      <AlertDialogRoot
        open={deleteModal() && rowSelection.multiple()}
        onOpenChange={setDeleteModal}
      >
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>
              Are you sure you wish to delete {rowSelection.selections().length}{" "}
              event{" "}
              {Humanize.pluralize(
                rowSelection.selections().length,
                "rule",
                "rules",
              )}
              ?
            </AlertDialogTitle>
            <AlertDialogDescription>
              <ul>
                <For each={data.data}>
                  {(e, index) => (
                    <Show when={rowSelection.items[index()]?.checked}>
                      <li>{e.code}</li>
                    </Show>
                  )}
                </For>
              </ul>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              variant="destructive"
              disabled={deleteMutation.isPending}
              onClick={() => deleteMutation.mutate()}
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <div class="flex justify-end gap-2">
        <Button size="icon">
          <RiSystemAddLine class="size-5" />
        </Button>
        <Button
          size="icon"
          variant="destructive"
          disabled={!rowSelection.multiple() || deleteMutation.isPending}
          onClick={() => setDeleteModal(true)}
        >
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
                <TableHeadBase>
                  <button>DB</button>
                </TableHeadBase>
                <TableHeadBase>
                  <button>Live</button>
                </TableHeadBase>
                <TableHeadBase>
                  <button>MQTT</button>
                </TableHeadBase>
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
                        onChange={(value) => rowSelection.set(item.uuid, value)}
                      >
                        <CheckboxControl />
                      </CheckboxRoot>
                    </TableCell>
                    <TableCell class="w-full">{item.code}</TableCell>
                    <TableCell>
                      <div class="flex justify-center">
                        <CheckboxRoot checked={item.allow_db}>
                          <CheckboxControl />
                        </CheckboxRoot>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div class="flex justify-center">
                        <CheckboxRoot checked={item.allow_live}>
                          <CheckboxControl />
                        </CheckboxRoot>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div class="flex justify-center">
                        <CheckboxRoot checked={item.allow_mqtt}>
                          <CheckboxControl />
                        </CheckboxRoot>
                      </div>
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
