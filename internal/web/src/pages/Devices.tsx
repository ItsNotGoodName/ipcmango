import { A, useSearchParams } from "@solidjs/router"
import { ErrorBoundary, For, Show, Suspense, createSignal, } from "solid-js"
import { PageError } from "~/ui/Page"
import { LayoutNormal } from "~/ui/Layout"
import { TabsContent, TabsList, TabsRoot, TabsTrigger } from "~/ui/Tabs"
import { TableBody, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table"
import { Skeleton } from "~/ui/Skeleton"
import { createUptime, formatDate, parseDate, useQueryFilter, } from "~/lib/utils"
import { linkVariants } from "~/ui/Link"
import { createMutation, createQuery, useQueryClient } from "@tanstack/solid-query"
import { GetApiDevicesResponse, postApiDevicesByUuidVideoInModeSync } from "~/client"
import { Button } from "~/ui/Button"
import { RiMediaImageLine, RiSystemRefreshLine } from "solid-icons/ri"
import { Image } from "@kobalte/core/image"
import { ToggleButton } from "@kobalte/core/toggle-button"
import Humanize from "humanize-plus"
import { TooltipArrow, TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip"
import { createDate, createTimeAgo } from "@solid-primitives/date"
import { api } from "./data"
import { DeviceFilterCombobox } from "~/components/DeviceFilterCombobox"
import { toast } from "~/ui/Toast"


function EmptyTableCell(props: { colspan: number }) {
  return <TableCell colspan={props.colspan}>N/A</TableCell>
}

function LoadingTableCell(props: { colspan: number }) {
  return (
    <TableCell colspan={props.colspan} class="py-0">
      <Skeleton class="h-8" />
    </TableCell>
  )
}

function ErrorTableCell(props: { colspan: number, error: Error }) {
  return (
    <TableCell colspan={props.colspan} class="py-0">
      <div class="bg-destructive text-destructive-foreground rounded p-2">
        {props.error.message}
      </div>
    </TableCell>
  )
}

function DeviceNameCell(props: { device: { uuid: string, name: string } }) {
  return (
    <TableCell>
      <A class={linkVariants()} href={`./${props.device.uuid}`}>{props.device.name}</A>
    </TableCell>
  )
}

export function Devices() {
  const [searchParams, setSearchParams] = useSearchParams()

  const data = createQuery(() => api.devices.list)

  const queryFilter = useQueryFilter("device")
  const devices = () => queryFilter.values().length == 0 ? data.data : data.data?.filter(v => queryFilter.values().includes(v.uuid))

  return (
    <LayoutNormal>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <h1>Devices</h1>
        <div class="flex">
          <DeviceFilterCombobox
            deviceIDs={queryFilter.values()}
            setDeviceIDs={queryFilter.setValues}
          />
        </div>
        <TabsRoot value={searchParams.tab || "device"} onChange={(value) => setSearchParams({ tab: value })}>
          <div class="flex flex-col gap-2">
            <div class="overflow-x-auto">
              <TabsList>
                <TabsTrigger value="device">Device</TabsTrigger>
                <TabsTrigger value="status">Status</TabsTrigger>
                <TabsTrigger value="uptime">Uptime</TabsTrigger>
                <TabsTrigger value="snapshot">Snapshot</TabsTrigger>
                <TabsTrigger value="detail">Detail</TabsTrigger>
                <TabsTrigger value="software-version">Software Version</TabsTrigger>
                <TabsTrigger value="license">License</TabsTrigger>
                <TabsTrigger value="storage">Storage</TabsTrigger>
                <TabsTrigger value="videoinmode">VideoInMode</TabsTrigger>
              </TabsList>
            </div>
          </div>
          <TabsContent value="device">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <DeviceTable devices={devices()} />
            </Suspense>
          </TabsContent>
          <TabsContent value="status">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <StatusTable devices={devices()} />
            </Suspense>
          </TabsContent>
          <TabsContent value="uptime">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <UptimeTable devices={devices()} />
            </Suspense>
          </TabsContent>
          <TabsContent value="snapshot">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <SnapshotGrid devices={devices()} />
            </Suspense>
          </TabsContent>
          <TabsContent value="detail">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <DetailTable devices={devices()} />
            </Suspense>
          </TabsContent>
          <TabsContent value="software-version">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <SoftwareVersionTable devices={devices()} />
            </Suspense>
          </TabsContent>
          <TabsContent value="license">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <LicenseTable devices={devices()} />
            </Suspense>
          </TabsContent>
          <TabsContent value="storage">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <StorageTable devices={devices()} />
            </Suspense>
          </TabsContent>
          <TabsContent value="videoinmode">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <VideoInMode devices={devices()} />
            </Suspense>
          </TabsContent>
        </TabsRoot>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

function DeviceTable(props: { devices?: GetApiDevicesResponse }) {
  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>IP</TableHead>
          <TableHead>Username</TableHead>
          <TableHead>Created At</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => (
            <TableRow>
              <DeviceNameCell device={item} />
              <TableCell><a class={linkVariants()} href={'http://' + item.ip} target="_blank">{item.ip}</a></TableCell>
              <TableCell>{item.username}</TableCell>
              <TableCell>{formatDate(parseDate(item.created_at))}</TableCell>
            </TableRow>
          )}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function StatusTable(props: { devices?: GetApiDevicesResponse }) {
  const colspan = 8

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>State</TableHead>
          <TableHead>Last Login</TableHead>
          <TableHead>Error</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createQuery(() => api.devices.status(item.uuid))

            return (
              <TableRow>
                <DeviceNameCell device={item} />
                <ErrorBoundary fallback={e => <ErrorTableCell colspan={colspan} error={e} />}>
                  <Suspense fallback={<LoadingTableCell colspan={colspan} />}>
                    <TableCell>{data.data?.state}</TableCell>
                    <TableCell>{formatDate(parseDate(data?.data?.last_login))}</TableCell>
                    <TableCell>{data.data?.error}</TableCell>
                  </Suspense>
                </ErrorBoundary>
              </TableRow>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function UptimeTable(props: { devices?: GetApiDevicesResponse }) {
  const colspan = 2

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>Last</TableHead>
          <TableHead>Total</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createQuery(() => api.devices.uptime(item.uuid))

            return (
              <TableRow>
                <DeviceNameCell device={item} />
                <ErrorBoundary fallback={e => <ErrorTableCell colspan={colspan} error={e} />}>
                  <Suspense fallback={<LoadingTableCell colspan={colspan} />}>
                    <Show when={data.data?.supported} fallback={
                      <EmptyTableCell colspan={colspan} />
                    }>
                      <UptimeTableCell date={parseDate(data.data?.last)} />
                      <UptimeTableCell date={parseDate(data.data?.total)} />
                    </Show>
                  </Suspense>
                </ErrorBoundary>
              </TableRow>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function UptimeTableCell(props: { date: Date }) {
  const uptime = createUptime(() => props.date)

  return (
    <TableCell>
      <Show when={uptime().hasDays}>
        {uptime().days} days &nbsp
      </Show>
      <Show when={uptime().hasHours}>
        {uptime().hours} hours &nbsp
      </Show>
      <Show when={uptime().hasMinutes}>
        {uptime().minutes} minutes &nbsp
      </Show>
      {uptime().seconds} seconds
    </TableCell>
  )
}

function SnapshotGrid(props: { devices?: GetApiDevicesResponse }) {
  return (
    <div class="grid grid-cols-1 gap-4 lg:grid-cols-2 2xl:grid-cols-3">
      <For each={props.devices}>
        {item => {
          const [t, setT] = createSignal(new Date().getTime())
          const refresh = () => setT(new Date().getTime())
          const src = () => `/api/devices/${item.uuid}/snapshot?t=${t()}`

          return (
            <div>
              <div class="flex flex-col rounded-t border">
                <div class="flex items-center justify-between gap-2 border-b p-2">
                  <Button size="icon" variant="ghost">
                    <RiSystemRefreshLine class="size-5" onClick={refresh} />
                  </Button>
                  <div class="px-2">
                    <A href={`/devices/${item.uuid}`}>{item.name}</A>
                  </div>
                </div>
                <Image>
                  <Image.Img src={src()} />
                  <Image.Fallback>
                    <RiMediaImageLine class="h-full w-full" />
                  </Image.Fallback>
                </Image>
              </div>
            </div>
          )
        }}
      </For>
    </div>
  )
}

function DetailTable(props: { devices?: GetApiDevicesResponse }) {
  const colspan = 9

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>SN</TableHead>
          <TableHead>Device Class</TableHead>
          <TableHead>Device Type</TableHead>
          <TableHead>Hardware Version</TableHead>
          <TableHead>Market Area</TableHead>
          <TableHead>Process Info</TableHead>
          <TableHead>Vendor</TableHead>
          <TableHead>Onvif Version</TableHead>
          <TableHead>Algorithm Version</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createQuery(() => api.devices.detail(item.uuid))

            return (
              <TableRow>
                <DeviceNameCell device={item} />
                <ErrorBoundary fallback={e => <ErrorTableCell colspan={colspan} error={e} />}>
                  <Suspense fallback={<LoadingTableCell colspan={colspan} />}>
                    <TableCell>
                      <ToggleButton>
                        {state => (
                          <Show when={state.pressed()} fallback={<>***************</>}>
                            {data.data?.sn}
                          </Show>
                        )}
                      </ToggleButton>
                    </TableCell>
                    <TableCell>{data.data?.device_class}</TableCell>
                    <TableCell>{data.data?.device_type}</TableCell>
                    <TableCell>{data.data?.hardware_version}</TableCell>
                    <TableCell>{data.data?.market_area}</TableCell>
                    <TableCell>{data.data?.process_info}</TableCell>
                    <TableCell>{data.data?.vendor}</TableCell>
                    <TableCell>{data.data?.onvif_version}</TableCell>
                    <TableCell>{data.data?.algorithm_version}</TableCell>
                  </Suspense>
                </ErrorBoundary>
              </TableRow>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function SoftwareVersionTable(props: { devices?: GetApiDevicesResponse }) {
  const colspan = 9

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>Build</TableHead>
          <TableHead>Build Date</TableHead>
          <TableHead>Security Base Line Version</TableHead>
          <TableHead>Version</TableHead>
          <TableHead>Web Version</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createQuery(() => api.devices.software(item.uuid))

            return (
              <TableRow>
                <DeviceNameCell device={item} />
                <ErrorBoundary fallback={e => <ErrorTableCell colspan={colspan} error={e} />}>
                  <Suspense fallback={<LoadingTableCell colspan={colspan} />}>
                    <TableCell>{data.data?.build}</TableCell>
                    <TableCell>{data.data?.build_date}</TableCell>
                    <TableCell>{data.data?.security_base_line_version}</TableCell>
                    <TableCell>{data.data?.version}</TableCell>
                    <TableCell>{data.data?.web_version}</TableCell>
                  </Suspense>
                </ErrorBoundary>
              </TableRow>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function LicenseTable(props: { devices?: GetApiDevicesResponse }) {
  const colspan = 9

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>Abroad Info</TableHead>
          <TableHead>All Type</TableHead>
          <TableHead>Digit Channel</TableHead>
          <TableHead>Effective Days</TableHead>
          <TableHead>Effective Time</TableHead>
          <TableHead>License ID</TableHead>
          <TableHead>Product Type</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Username</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createQuery(() => api.devices.licenses(item.uuid))

            return (
              <ErrorBoundary fallback={e =>
                <TableRow>
                  <DeviceNameCell device={item} />
                  <ErrorTableCell colspan={colspan} error={e} />
                </TableRow>
              }>
                <Suspense fallback={
                  <TableRow>
                    <DeviceNameCell device={item} />
                    <LoadingTableCell colspan={colspan} />
                  </TableRow>
                }>
                  <For each={data.data} fallback={
                    <TableRow>
                      <DeviceNameCell device={item} />
                      <EmptyTableCell colspan={colspan} />
                    </TableRow>
                  }>
                    {v => {
                      const [effectiveTime] = createDate(() => parseDate(v.effective_time));
                      const [effectiveTimeAgo] = createTimeAgo(effectiveTime, { interval: 0 });

                      return (
                        <TableRow>
                          <DeviceNameCell device={item} />
                          <TableCell>{v.abroad_info}</TableCell>
                          <TableCell>{v.all_type}</TableCell>
                          <TableCell>{v.digit_channel}</TableCell>
                          <TableCell>{v.effective_days}</TableCell>
                          <TableCell>
                            <TooltipRoot>
                              <TooltipTrigger>{formatDate(effectiveTime())}</TooltipTrigger>
                              <TooltipContent>
                                <TooltipArrow />
                                {effectiveTimeAgo()}
                              </TooltipContent>
                            </TooltipRoot>
                          </TableCell>
                          <TableCell>{v.license_id}</TableCell>
                          <TableCell>{v.product_type}</TableCell>
                          <TableCell>{v.status}</TableCell>
                          <TableCell>{v.username}</TableCell>
                        </TableRow>
                      )
                    }
                    }
                  </For>
                </Suspense>
              </ErrorBoundary>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function StorageTable(props: { devices?: GetApiDevicesResponse }) {
  const colspan = 6;

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>Name</TableHead>
          <TableHead>State</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Used</TableHead>
          <TableHead>Total</TableHead>
          <TableHead>Is Error</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createQuery(() => api.devices.storage(item.uuid))

            return (
              <ErrorBoundary fallback={e =>
                <TableRow>
                  <DeviceNameCell device={item} />
                  <ErrorTableCell colspan={colspan} error={e} />
                </TableRow>
              }>
                <Suspense fallback={
                  <TableRow>
                    <DeviceNameCell device={item} />
                    <LoadingTableCell colspan={colspan} />
                  </TableRow>
                }>
                  <For each={data.data} fallback={
                    <TableRow>
                      <DeviceNameCell device={item} />
                      <EmptyTableCell colspan={colspan} />
                    </TableRow>
                  }>
                    {v => (
                      <TableRow>
                        <DeviceNameCell device={item} />
                        <TableCell>{v.name}</TableCell>
                        <TableCell>{v.state}</TableCell>
                        <TableCell>{v.type}</TableCell>
                        <TableCell>{Humanize.fileSize(Number(v.used_bytes))}</TableCell>
                        <TableCell>{Humanize.fileSize(Number(v.total_bytes))}</TableCell>
                        <TableCell>{v.is_error}</TableCell>
                      </TableRow>
                    )
                    }
                  </For>
                </Suspense>
              </ErrorBoundary>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function VideoInMode(props: { devices?: GetApiDevicesResponse }) {
  const colspan = 4;

  const client = useQueryClient()

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>Switch Mode</TableHead>
          <TableHead>Time Section</TableHead>
          <TableHead></TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createQuery(() => api.devices.video_in_mode(item.uuid))
            const sync = createMutation(() => ({
              mutationFn: () => postApiDevicesByUuidVideoInModeSync({ uuid: item.uuid, requestBody: {} }),
              onSuccess: (data) => client.setQueryData(api.devices.video_in_mode(item.uuid).queryKey, data),
              onError: (error) => toast.error(error.name, error.message)
            }))

            return (
              <ErrorBoundary fallback={e =>
                <TableRow>
                  <DeviceNameCell device={item} />
                  <ErrorTableCell colspan={colspan} error={e} />
                </TableRow>
              }>
                <Suspense fallback={
                  <TableRow>
                    <DeviceNameCell device={item} />
                    <LoadingTableCell colspan={colspan} />
                  </TableRow>
                }>
                  <TableRow>
                    <DeviceNameCell device={item} />
                    <TableCell>{data.data?.switch_mode}</TableCell>
                    <TableCell>{data.data?.time_section}</TableCell>
                    <TableCell><Button size="xs" disabled={sync.isPending} onClick={() => sync.mutate()}>Sync</Button></TableCell>
                  </TableRow>
                </Suspense>
              </ErrorBoundary>
            )
          }}
        </For>
      </TableBody>
    </TableRoot >
  )
}

export default Devices
