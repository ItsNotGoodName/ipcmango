import { CheckboxControl, CheckboxRoot } from "~/ui/Checkbox";
import { action, createAsync, revalidate, useAction, useNavigate, useSearchParams, useSubmission, } from "@solidjs/router";
import { ErrorBoundary, For, Show, Suspense, } from "solid-js";
import { RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiDesignFocus2Line, RiSystemLockLine, RiSystemMore2Line, RiUserFacesAdminLine, } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { catchAsToast, createPagePagination, createRowSelection, createToggleSortField, formatDate, parseDate, parseOrder, } from "~/lib/utils";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow, } from "~/ui/Table";
import { Seperator } from "~/ui/Seperator";
import { Skeleton } from "~/ui/Skeleton";
import { PageError } from "~/ui/Page";
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { AdminUsersPageSearchParams, getAdminUsersPage } from "./Users.data";
import { LayoutNormal } from "~/ui/Layout";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { getSession } from "~/providers/session";
import { Crud } from "~/components/Crud";
import { useClient } from "~/providers/client";
import { SetUserAdminReq, SetUserDisableReq } from "~/twirp/rpc";

const actionSetUserDisable = action((input: SetUserDisableReq) => useClient()
  .admin.setUserDisable(input)
  .then(() => revalidate(getAdminUsersPage.key))
  .catch(catchAsToast)
)

const actionSetUserAdmin = action((input: SetUserAdminReq) => useClient()
  .admin.setUserAdmin(input)
  .then(() => revalidate(getAdminUsersPage.key))
  .catch(catchAsToast)
)

export function AdminUsers() {
  const session = createAsync(getSession)

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

  // Disable/Enable
  const setUserDisableSubmission = useSubmission(actionSetUserDisable)
  const setUserDisable = useAction(actionSetUserDisable)
  const setUserDisableByRowSelector = (disable: boolean) => setUserDisable({ items: rowSelection.selections().map(v => ({ id: v, disable })) })
    .then(() => rowSelection.setAll(false))

  const setUserAdminSubmission = useSubmission(actionSetUserAdmin)
  const setUserAdmin = useAction(actionSetUserAdmin)

  return (
    <LayoutNormal>
      <div class="text-xl">Users</div>
      <Seperator />

      <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex justify-between gap-2">
            <Crud.PerPageSelect
              class="w-20"
              perPage={data()?.pageResult?.perPage}
              onChange={pagination.setPerPage}
            />
            <div class="flex gap-2">
              <Button
                title="Previous"
                size="icon"
                disabled={pagination.previousPageDisabled()}
                onClick={pagination.previousPage}
              >
                <RiArrowsArrowLeftSLine class="h-6 w-6" />
              </Button>
              <Button
                title="Next"
                size="icon"
                disabled={pagination.nextPageDisabled()}
                onClick={pagination.nextPage}
              >
                <RiArrowsArrowRightSLine class="h-6 w-6" />
              </Button>
            </div>
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
                <TableHead>
                  <div class="flex items-center justify-end">
                    <DropdownMenuRoot placement="bottom-end">
                      <DropdownMenuTrigger class="hover:bg-accent hover:text-accent-foreground rounded p-1" title="Actions">
                        <RiSystemMore2Line class="h-5 w-5" />
                      </DropdownMenuTrigger>
                      <DropdownMenuPortal>
                        <DropdownMenuContent>
                          <DropdownMenuItem>
                            Create
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            disabled={rowSelection.selections().length == 0}
                            onClick={() => setUserDisableByRowSelector(false)}
                          >
                            Enable
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            disabled={rowSelection.selections().length == 0}
                            onClick={() => setUserDisableByRowSelector(true)}
                          >
                            Disable
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            disabled={rowSelection.selections().length == 0}
                          >
                            Delete
                          </DropdownMenuItem>
                          <DropdownMenuArrow />
                        </DropdownMenuContent>
                      </DropdownMenuPortal>
                    </DropdownMenuRoot>
                  </div>
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data()?.items}>
                {(item, index) => {
                  const onClick = () => navigate(`./${item.id}`)
                  const toggleUserDisable = () => setUserDisable({ items: [{ id: item.id, disable: !item.disabled }] })
                  const toggleUserAdmin = () => setUserAdmin({ id: item.id, admin: !item.admin })

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
                          <DropdownMenuTrigger class="hover:bg-accent hover:text-accent-foreground rounded p-1" title="Actions">
                            <RiSystemMore2Line class="h-5 w-5" />
                          </DropdownMenuTrigger>
                          <DropdownMenuPortal>
                            <DropdownMenuContent>
                              <DropdownMenuItem>
                                Edit
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                Reset password
                              </DropdownMenuItem>
                              <Show when={item.id != BigInt(session()?.user_id || 0)}>
                                <DropdownMenuItem disabled={setUserDisableSubmission.pending} onSelect={toggleUserDisable}>
                                  <Show when={item.disabled} fallback={<>Disable</>}>
                                    Enable
                                  </Show>
                                </DropdownMenuItem>
                                <DropdownMenuItem
                                  disabled={setUserAdminSubmission.pending}
                                  onClick={toggleUserAdmin}
                                >
                                  <Show when={!item.admin} fallback={<>Demote</>}>
                                    Promote
                                  </Show>
                                </DropdownMenuItem>
                                <DropdownMenuItem

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

