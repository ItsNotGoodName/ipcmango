import { CheckboxControl, CheckboxRoot } from "~/ui/Checkbox";
import { createAsync, useNavigate, useSearchParams, } from "@solidjs/router";
import { ErrorBoundary, For, Show, Suspense, } from "solid-js";
import { RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiSystemLockLine, RiSystemMore2Line, RiUserFacesAdminLine, } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { SelectContent, SelectItem, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select";
import { createRowSelection, formatDate, parseDate, } from "~/lib/utils";
import { encodeOrder, toggleSortField, parseOrder } from "~/lib/utils";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableMetadata, TableRoot, TableRow, TableSortButton } from "~/ui/Table";
import { Seperator } from "~/ui/Seperator";
import { Skeleton } from "~/ui/Skeleton";
import { PageError } from "~/ui/Page";
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { AdminUsersPageSearchParams, getAdminUsersPage } from "./Users.data";
import { defaultPerPageOptions } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";

export function AdminUsers() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams<AdminUsersPageSearchParams>()
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
  const previousPageDisabled = () => data()?.pageResult?.previousPage == data()?.pageResult?.page
  const previousPage = () => !previousPageDisabled() && setSearchParams({ page: data()?.pageResult?.previousPage.toString() } as AdminUsersPageSearchParams)
  const nextPageDisabled = () => data()?.pageResult?.nextPage == data()?.pageResult?.page
  const nextPage = () => !nextPageDisabled() && setSearchParams({ page: data()?.pageResult?.nextPage.toString() } as AdminUsersPageSearchParams)
  const toggleSort = (field: string) => {
    const sort = toggleSortField(data()?.sort, field)
    return setSearchParams({ sort: sort.field, order: encodeOrder(sort.order) } as AdminUsersPageSearchParams)
  }
  const setPerPage = (value: number) => value && setSearchParams({ page: 1, perPage: value })

  return (
    <LayoutNormal>
      <div class="text-xl">Users</div>
      <Seperator />

      <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex justify-between gap-2">
            <SelectRoot
              class="w-20"
              value={data()?.pageResult?.perPage}
              onChange={setPerPage}
              options={defaultPerPageOptions}
              itemComponent={props => (
                <SelectItem item={props.item}>
                  {props.item.rawValue}
                </SelectItem>
              )}
            >
              <SelectTrigger aria-label="Per page">
                <SelectValue<number>>
                  {state => state.selectedOption()}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                <SelectListbox />
              </SelectContent>
            </SelectRoot>
            <div class="flex gap-2">
              <Button
                title="Previous"
                size="icon"
                disabled={previousPageDisabled()}
                onClick={previousPage}
              >
                <RiArrowsArrowLeftSLine class="h-6 w-6" />
              </Button>
              <Button
                title="Next"
                size="icon"
                disabled={nextPageDisabled()}
                onClick={nextPage}
              >
                <RiArrowsArrowRightSLine class="h-6 w-6" />
              </Button>
            </div>
          </div>
          <TableRoot>
            <TableHeader>
              <tr class="border-b">
                <TableHead>
                  <CheckboxRoot
                    checked={rowSelection.multiple()}
                    indeterminate={rowSelection.indeterminate()}
                    onChange={(v) => rowSelection.checkAll(v)}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableHead>
                <TableHead>
                  <TableSortButton
                    name="username"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Username
                  </TableSortButton>
                </TableHead>
                <TableHead class="w-full">
                  <TableSortButton
                    name="email"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Email
                  </TableSortButton>
                </TableHead>
                <TableHead>
                  <TableSortButton
                    name="createdAt"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Created At
                  </TableSortButton>
                </TableHead>
                <TableHead>
                  <div class="flex items-center justify-end">
                    <DropdownMenuRoot placement="bottom-end">
                      <DropdownMenuTrigger class="hover:bg-accent hover:text-accent-foreground rounded p-1" title="Actions">
                        <RiSystemMore2Line class="h-5 w-5" />
                      </DropdownMenuTrigger>
                      <DropdownMenuPortal>
                        <DropdownMenuContent>
                          <DropdownMenuArrow />
                        </DropdownMenuContent>
                      </DropdownMenuPortal>
                    </DropdownMenuRoot>
                  </div>
                </TableHead>
              </tr>
            </TableHeader>
            <TableBody>
              <For each={data()?.items}>
                {(item, index) => {
                  const onClick = () => navigate(`./${item.id}`)

                  return (
                    <TableRow>
                      <TableHead>
                        <CheckboxRoot checked={rowSelection.rows[index()]?.checked} onChange={(v) => rowSelection.check(item.id, v)}>
                          <CheckboxControl />
                        </CheckboxRoot>
                      </TableHead>
                      <TableCell onClick={onClick} class="cursor-pointer select-none">{item.username}</TableCell>
                      <TableCell onClick={onClick} class="cursor-pointer select-none">{item.email}</TableCell>
                      <TableCell onClick={onClick} class="text-nowrap cursor-pointer select-none whitespace-nowrap">{formatDate(parseDate(item.createdAtTime))}</TableCell>
                      <TableCell class="py-0">
                        <div class="flex justify-end gap-2">
                          <Show when={item.admin}>
                            <TooltipRoot>
                              <TooltipTrigger class="p-1">
                                <RiUserFacesAdminLine class="h-5 w-5" />
                              </TooltipTrigger>
                              <TooltipContent>
                                Admin
                              </TooltipContent>
                            </TooltipRoot>
                          </Show>
                          <Show when={item.disabled}>
                            <TooltipRoot>
                              <TooltipTrigger class="p-1">
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
                                <DropdownMenuArrow />
                              </DropdownMenuContent>
                            </DropdownMenuPortal>
                          </DropdownMenuRoot>
                        </div>
                      </TableCell>
                    </TableRow>
                  )
                }}
              </For>
            </TableBody>
            <TableCaption>
              <TableMetadata pageResult={data()?.pageResult} />
            </TableCaption>
          </TableRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

