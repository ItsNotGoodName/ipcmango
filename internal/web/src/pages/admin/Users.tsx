import { CheckboxControl, CheckboxRoot } from "~/ui/Checkbox";
import { action, createAsync, revalidate, useAction, useNavigate, useSearchParams, useSubmission, } from "@solidjs/router";
import { ErrorBoundary, For, Show, Suspense, createSignal, } from "solid-js";
import { RiDesignFocus2Line, RiSystemLockLine, RiUserFacesAdminLine, } from "solid-icons/ri";
import { catchAsToast, createPagePagination, createRowSelection, createToggleSortField, formatDate, parseDate, parseOrder, } from "~/lib/utils";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow, } from "~/ui/Table";
import { Seperator } from "~/ui/Seperator";
import { Skeleton } from "~/ui/Skeleton";
import { PageError } from "~/ui/Page";
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { AdminUsersPageSearchParams, getAdminUsersPage } from "./Users.data";
import { LayoutNormal } from "~/ui/Layout";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot, } from "~/ui/DropdownMenu";
import { getSession } from "~/providers/session";
import { Crud } from "~/components/Crud";
import { useClient } from "~/providers/client";
import { SetUserAdminReq, SetUserDisableReq } from "~/twirp/rpc";
import { AlertDialogAction, AlertDialogCancel, AlertDialogModal, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogRoot, AlertDialogTitle } from "~/ui/AlertDialog";

const actionDelete = action((ids: bigint[]) => useClient()
  .admin.deleteUser({ ids })
  .then(() => revalidate(getAdminUsersPage.key))
  .catch(catchAsToast))

const actionSetDisable = action((input: SetUserDisableReq) => useClient()
  .admin.setUserDisable(input)
  .then(() => revalidate(getAdminUsersPage.key))
  .catch(catchAsToast))

const actionSetAdmin = action((input: SetUserAdminReq) => useClient()
  .admin.setUserAdmin(input)
  .then(() => revalidate(getAdminUsersPage.key))
  .catch(catchAsToast))

export function AdminUsers() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams<AdminUsersPageSearchParams>()
  const data = createAsync(() => getAdminUsersPage({
    page: {
      page: Number(searchParams.page) || 1,
      perPage: Number(searchParams.perPage) || 10
    },
    sort: {
      field: searchParams.sort || "",
      order: parseOrder(searchParams.order)
    },
  }))
  const rowSelection = createRowSelection(() => data()?.items.map(v => v.id) || [])

  // List
  const pagination = createPagePagination(() => data()?.pageResult)
  const toggleSort = createToggleSortField(() => data()?.sort)

  // Delete
  const deleteSubmission = useSubmission(actionDelete)
  const deleteAction = useAction(actionDelete)
  // Single
  const [openDeleteConfirm, setOpenDeleteConfirm] = createSignal<{ username: string, id: bigint } | undefined>()
  const deleteSubmit = () => deleteAction([openDeleteConfirm()?.id || BigInt(0)])
    .then(() => setOpenDeleteConfirm(undefined))
  // Multiple
  const [openDeleteMultipleConfirm, setOpenDeleteMultipleConfirm] = createSignal(false)
  const deleteMultipleSubmit = () => deleteAction(rowSelection.selections())
    .then(() => setOpenDeleteMultipleConfirm(false))

  // Disable
  const setDisableSubmission = useSubmission(actionSetDisable)
  const setDisableAction = useAction(actionSetDisable)
  const setDisableMultipleSubmit = (disable: boolean) => setDisableAction({ items: rowSelection.selections().map(v => ({ id: v, disable })) })
    .then(() => rowSelection.setAll(false))

  // Admin
  const setAdminSubmission = useSubmission(actionSetAdmin)
  const setAdminAction = useAction(actionSetAdmin)

  const session = createAsync(getSession)

  return (
    <LayoutNormal class="max-w-4xl">

      <AlertDialogRoot open={openDeleteConfirm() != undefined} onOpenChange={() => setOpenDeleteConfirm(undefined)}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {openDeleteConfirm()?.username}?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteSubmission.pending} onClick={deleteSubmit}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <AlertDialogRoot open={openDeleteMultipleConfirm()} onOpenChange={setOpenDeleteMultipleConfirm}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {rowSelection.selections().length} users?</AlertDialogTitle>
            <AlertDialogDescription class="max-h-32 overflow-y-auto">
              <For each={data()?.items}>
                {(e, index) =>
                  <Show when={rowSelection.rows[index()].checked}>
                    <div>
                      {e.username}
                    </div>
                  </Show>
                }
              </For>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteSubmission.pending} onClick={deleteMultipleSubmit}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <div class="text-xl">Users</div>
      <Seperator />

      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex justify-between gap-2">
            <Crud.PerPageSelect
              class="w-20"
              perPage={data()?.pageResult?.perPage}
              onChange={pagination.setPerPage}
            />
            <Crud.PageButtons
              previousPageDisabled={pagination.previousPageDisabled()}
              previousPage={pagination.previousPage}
              nextPageDisabled={pagination.nextPageDisabled()}
              nextPage={pagination.nextPage}
            />
          </div>
          <TableRoot>
            <TableHeader>
              <TableRow>
                <TableHead>
                  <CheckboxRoot
                    checked={rowSelection.multiple()}
                    indeterminate={rowSelection.indeterminate()}
                    onChange={(v) => rowSelection.setAll(v)}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableHead>
                <TableHead>
                  <Crud.SortButton
                    name="username"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Username
                  </Crud.SortButton>
                </TableHead>
                <TableHead>
                  <Crud.SortButton
                    name="email"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Email
                  </Crud.SortButton>
                </TableHead>
                <TableHead>
                  <Crud.SortButton
                    name="createdAt"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Created At
                  </Crud.SortButton>
                </TableHead>
                <Crud.LastTableHead>
                  <DropdownMenuRoot placement="bottom-end">
                    <Crud.MoreDropdownMenuTrigger />
                    <DropdownMenuPortal>
                      <DropdownMenuContent>
                        <DropdownMenuItem>
                          Create
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={rowSelection.selections().length == 0 || setDisableSubmission.pending}
                          onClick={() => setDisableMultipleSubmit(false)}
                        >
                          Enable
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={rowSelection.selections().length == 0 || setDisableSubmission.pending}
                          onClick={() => setDisableMultipleSubmit(true)}
                        >
                          Disable
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={rowSelection.selections().length == 0 || deleteSubmission.pending}
                          onClick={() => setOpenDeleteMultipleConfirm(true)}
                        >
                          Delete
                        </DropdownMenuItem>
                        <DropdownMenuArrow />
                      </DropdownMenuContent>
                    </DropdownMenuPortal>
                  </DropdownMenuRoot>
                </Crud.LastTableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data()?.items}>
                {(item, index) => {
                  const onClick = () => navigate(`./${item.id}`)
                  const toggleDisable = () => setDisableAction({ items: [{ id: item.id, disable: !item.disabled }] })
                  const toggleAdmin = () => setAdminAction({ id: item.id, admin: !item.admin })

                  return (
                    <TableRow>
                      <TableHead>
                        <CheckboxRoot checked={rowSelection.rows[index()]?.checked} onChange={(v) => rowSelection.set(item.id, v)}>
                          <CheckboxControl />
                        </CheckboxRoot>
                      </TableHead>
                      <TableCell class="cursor-pointer select-none" onClick={onClick}>{item.username}</TableCell>
                      <TableCell class="cursor-pointer select-none" onClick={onClick}>{item.email}</TableCell>
                      <TableCell class="cursor-pointer select-none" onClick={onClick}>{formatDate(parseDate(item.createdAtTime))}</TableCell>
                      <Crud.LastTableCell>
                        <Show when={item.id == BigInt(session()?.user_id || 0)}>
                          <TooltipRoot>
                            <TooltipTrigger>
                              <RiDesignFocus2Line class="h-5 w-5" />
                            </TooltipTrigger>
                            <TooltipContent>
                              You
                            </TooltipContent>
                          </TooltipRoot>
                        </Show>
                        <Show when={item.admin}>
                          <TooltipRoot>
                            <TooltipTrigger>
                              <RiUserFacesAdminLine class="h-5 w-5" />
                            </TooltipTrigger>
                            <TooltipContent>
                              Admin
                            </TooltipContent>
                          </TooltipRoot>
                        </Show>
                        <Show when={item.disabled}>
                          <TooltipRoot>
                            <TooltipTrigger>
                              <RiSystemLockLine class="h-5 w-5" />
                            </TooltipTrigger>
                            <TooltipContent>
                              Disabled since {formatDate(parseDate(item.disabledAtTime))}
                            </TooltipContent>
                          </TooltipRoot>
                        </Show>
                        <DropdownMenuRoot placement="bottom-end">
                          <Crud.MoreDropdownMenuTrigger />
                          <DropdownMenuPortal>
                            <DropdownMenuContent>
                              <DropdownMenuItem>
                                Edit
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                Reset password
                              </DropdownMenuItem>
                              <Show when={item.id != BigInt(session()?.user_id || 0)}>
                                <DropdownMenuItem disabled={setDisableSubmission.pending} onSelect={toggleDisable}>
                                  <Show when={item.disabled} fallback={<>Disable</>}>
                                    Enable
                                  </Show>
                                </DropdownMenuItem>
                                <DropdownMenuItem
                                  disabled={setAdminSubmission.pending}
                                  onClick={toggleAdmin}
                                >
                                  <Show when={!item.admin} fallback={<>Demote</>}>
                                    Promote
                                  </Show>
                                </DropdownMenuItem>
                                <DropdownMenuItem
                                  disabled={deleteSubmission.pending}
                                  onClick={() => setOpenDeleteConfirm(item)}
                                >
                                  Delete
                                </DropdownMenuItem>
                              </Show>
                              <DropdownMenuArrow />
                            </DropdownMenuContent>
                          </DropdownMenuPortal>
                        </DropdownMenuRoot>
                      </Crud.LastTableCell>
                    </TableRow>
                  )
                }}
              </For>
            </TableBody>
            <TableCaption>
              <Crud.Metadata pageResult={data()?.pageResult} />
            </TableCaption>
          </TableRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

