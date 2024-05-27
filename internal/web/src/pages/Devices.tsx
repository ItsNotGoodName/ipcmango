import { A, useSearchParams } from "@solidjs/router"
import { ErrorBoundary, For, Show, Suspense, } from "solid-js"
import { PageError } from "~/ui/Page"
import { LayoutNormal } from "~/ui/Layout"
import { TabsContent, TabsList, TabsRoot, TabsTrigger } from "~/ui/Tabs"
import { TableBody, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table"
import { Skeleton } from "~/ui/Skeleton"
import { formatDate, } from "~/lib/utils"
import { linkVariants } from "~/ui/Link"
import { TooltipArrow, TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip"
import { RiSystemLockLine, } from "solid-icons/ri"
import { createQuery } from "@tanstack/solid-query"

// function EmptyTableCell(props: { colspan: number }) {
//   return <TableCell colspan={props.colspan}>N/A</TableCell>
// }
//
// function LoadingTableCell(props: { colspan: number }) {
//   return (
//     <TableCell colspan={props.colspan} class="py-0">
//       <Skeleton class="h-8" />
//     </TableCell>
//   )
// }
//
// function ErrorTableCell(props: { colspan: number, error: Error }) {
//   return (
//     <TableCell colspan={props.colspan} class="py-0">
//       <div class="bg-destructive text-destructive-foreground rounded p-2">
//         {props.error.message}
//       </div>
//     </TableCell>
//   )
// }

function DeviceNameCell(props: { device: { uuid: string, name: string } }) {
  return (
    <TableCell>
      <A class={linkVariants()} href={`./${props.device.uuid}`}>{props.device.name}</A>
    </TableCell>
  )
}

export function Devices() {
  const [searchParams, setSearchParams] = useSearchParams()

  const data = createQuery(() => ({
    queryKey: ["devices"],
    queryFn: async () => {
      const result = await fetch('/api/devices')
      if (!result.ok) throw new Error('Failed to fetch data')
      return result.json()
    },
    throwOnError: true, // Throw an error if the query fails
  }))

  return (
    <LayoutNormal>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <h1>Devices</h1>
        <TabsRoot value={searchParams.tab || "device"} onChange={(value) => setSearchParams({ tab: value })}>
          <div class="flex flex-col gap-2">
            <div class="overflow-x-auto">
              <TabsList>
                <TabsTrigger value="device">Device</TabsTrigger>
              </TabsList>
            </div>
          </div>
          <TabsContent value="device">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <DeviceTable devices={data.data} />
            </Suspense>
          </TabsContent>
        </TabsRoot>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

function DeviceTable(props: { devices?: any[] }) {
  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>IP</TableHead>
          <TableHead>Username</TableHead>
          <TableHead>Created At</TableHead>
          <TableHead></TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => (
            <TableRow>
              <DeviceNameCell device={item} />
              <TableCell><a class={linkVariants()} href={item.url}>{item.ip}</a></TableCell>
              <TableCell>{item.username}</TableCell>
              <TableCell>{formatDate(item.createdAtTime)}</TableCell>
              <TableCell>
                <Show when={item.disabled}>
                  <TooltipRoot>
                    <TooltipTrigger>
                      <RiSystemLockLine class="size-5" />
                    </TooltipTrigger>
                    <TooltipContent>
                      <TooltipArrow />
                      Disabled
                    </TooltipContent>
                  </TooltipRoot>
                </Show>
              </TableCell>
            </TableRow>
          )}
        </For>
      </TableBody>
    </TableRoot>
  )
}

export default Devices
